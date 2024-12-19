package tines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/tines/go-sdk/internal/utils"
	"go.uber.org/zap"
)

type Client struct {
	tenantUrl  string
	apiKey     string
	userAgent  string
	httpClient *http.Client
	logger     *zap.Logger
}

// Create a new Tines API client. The Tenant URL and Tines API Key
// must be specified when creating a new client. Adding a custom User
// Agent string is optional, but recommended for identifying the particular
// application making requests to the Tines API.
//
// Example Usage:
//
//	client, err := tines.NewClient(
//	  tines.SetTenantUrl("https://example.tines.com/"),
//	  tines.SetApiKey("foobar"),
//	)
func NewClient(opts ...func(*Client)) (*Client, error) {
	c := Client{
		httpClient: &http.Client{},
		logger:     zap.NewNop(),
	}
	errs := Error{Type: ErrorTypeRequest}

	for _, o := range opts {
		o(&c)
	}

	if c.tenantUrl == "" {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: "host error",
			Details: errEmptyTenant,
		})
	}

	if c.apiKey == "" {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: "credential error",
			Details: errEmptyApiKey,
		})
	}
	// We do additional error checking when crafting the HTTP request to
	// ensure it goes to a valid URL, but we should at least make sure
	// the identified tenant URL starts with a valid protocol since that
	// aspect is not checked by the url.Parse() function.
	ok, err := regexp.Match("^https:\\/\\/", []byte(c.tenantUrl))
	if err != nil {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: errParseError,
			Details: err.Error(),
		})
	}

	if !ok {
		errs.Errors = append(errs.Errors, ErrorMessage{
			Message: "host error",
			Details: errMalformedTenant,
		})
	}

	if errs.HasErrors() {
		return nil, errs
	}

	return &c, nil
}

func SetTenantUrl(s string) func(*Client) {
	return func(c *Client) {
		c.tenantUrl = s
	}
}

func SetApiKey(s string) func(*Client) {
	return func(c *Client) {
		c.apiKey = s
	}
}

func SetUserAgent(s string) func(*Client) {
	return func(c *Client) {
		c.userAgent = s
	}
}

func SetLogger(l *zap.Logger) func(*Client) {
	return func(c *Client) {
		c.logger = l
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, params map[string]any, data []byte) ([]byte, error) {
	tenant, err := url.Parse(c.tenantUrl)
	if err != nil {
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errParseError,
					Details: err.Error(),
				},
			},
		}
	}

	for k, v := range params {
		c.logger.Debug("params", zap.Any(k, v))
	}

	q := tenant.Query()
	for k, v := range params {
		switch reflect.TypeOf(v) {
		case reflect.TypeFor[int]():
			v = strconv.Itoa(v.(int))
		case reflect.TypeFor[float64]():
			v = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		case reflect.TypeFor[bool]():
			v = strconv.FormatBool(v.(bool))
		}
		c.logger.Debug(fmt.Sprintf("Setting %s as %s", k, v))
		q.Add(k, v.(string))
	}

	c.logger.Debug(fmt.Sprintf("query string: %s", q.Encode()))

	fullUrl := tenant.JoinPath(path)
	fullUrl.RawQuery = q.Encode()

	c.logger.Debug(fmt.Sprintf("sending request to url %s", fullUrl.String()))

	req, err := http.NewRequestWithContext(ctx, method, fullUrl.String(), bytes.NewBuffer(data))
	if err != nil {
		c.logger.Debug(err.Error())
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: err.Error(),
				},
			},
		}
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", utils.SetUserAgent(c.userAgent))
	req.Header.Set("x-tines-client-version", fmt.Sprintf("tines-go-sdk-%s", utils.SetClientVersion()))
	req.Header.Set("x-user-token", c.apiKey)

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		return nil, Error{
			Type: ErrorTypeRequest,
			Errors: []ErrorMessage{
				{
					Message: errDoRequestError,
					Details: respErr.Error(),
				},
			},
		}
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		c.logger.Debug(readErr.Error())
		return nil, Error{
			Type:       ErrorTypeServer,
			StatusCode: resp.StatusCode,
			Errors: []ErrorMessage{
				{
					Message: errReadBodyError,
					Details: readErr.Error(),
				},
			},
		}
	}

	// Return a server error for 5XX responses
	if resp.StatusCode >= http.StatusInternalServerError {
		errMsgs := c.getErrorMessages(body)

		c.logger.Debug(fmt.Sprintf("received a %d status code from the server", resp.StatusCode))
		return nil, Error{
			Type:       ErrorTypeServer,
			StatusCode: resp.StatusCode,
			Errors:     errMsgs,
		}
	}

	// Return a request error for 4XX responses
	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		errMsgs := c.getErrorMessages(body)

		c.logger.Debug(fmt.Sprintf("received a %d status code from the server", resp.StatusCode))
		return nil, Error{
			Type:       ErrorTypeRequest,
			StatusCode: resp.StatusCode,
			Errors:     errMsgs,
		}
	}

	return body, nil
}

func (c *Client) getErrorMessages(body []byte) []ErrorMessage {
	var errorInfo Error
	var errorMsgs []ErrorMessage

	// The structure of an error response body can be inconsistent between API endpoints,
	// so we try a couple techniques to capture the error messages.
	jsonErr := json.Unmarshal(body, &errorInfo)
	if jsonErr != nil {
		jsonErr := json.Unmarshal(body, &errorMsgs)
		if jsonErr != nil && body != nil {
			errorMsgs = []ErrorMessage{{Message: "message", Details: string(body)}}
		}
		return errorMsgs
	}

	return errorInfo.Errors
}
