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

	// If no custom UserAgent has been set during client initialization,
	// apply the default one.
	if c.userAgent == "" {
		ua := SetUserAgent("")
		ua(&c)
	}

	// We do additional error checking when crafting the HTTP request to
	// ensure it goes to a valid URL, but we should at least make sure
	// the identified tenant URL starts with a valid protocol since that
	// aspect is not checked by the url.Parse() function.
	ok, err := regexp.Match("^https?:\\/\\/", []byte(c.tenantUrl))
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
	c.logger.Debug("Tines client", zap.String("version", utils.SetClientVersion()))
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
		c.userAgent = utils.SetUserAgent(s)
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
		c.logger.Debug("found param", zap.Any(k, v))
	}

	q := tenant.Query()
	for k, v := range params {
		switch reflect.TypeOf(v) {
		case reflect.TypeFor[int]():
			i, ok := v.(int)
			if ok {
				v = strconv.Itoa(i)
			} else {
				c.logger.Debug("unable to convert int value to string, skipping", zap.Any(k, v))
			}
		case reflect.TypeFor[float64]():
			f, ok := v.(float64)
			if ok {
				v = strconv.FormatFloat(f, 'f', 0, 64)
			} else {
				c.logger.Debug("unable to convert float64 value to string, skipping", zap.Any(k, v))
			}
		case reflect.TypeFor[bool]():
			b, ok := v.(bool)
			if ok {
				v = strconv.FormatBool(b)
			} else {
				c.logger.Debug("unable to convert bool value to string, skipping", zap.Any(k, v))
			}
		}

		s, ok := v.(string)
		if ok {
			c.logger.Debug(fmt.Sprintf("setting query param %s to value %s", k, v))
			q.Add(k, s)
		} else {
			c.logger.Debug("invalid string value, skipping", zap.Any(k, v))
		}
	}

	c.logger.Debug(fmt.Sprintf("final query string: %s", q.Encode()))

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
