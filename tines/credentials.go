package tines

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"

	"github.com/tines/go-sdk/internal/paginate"
)

type CredentialType string

const (
	CredentialTypeAws   CredentialType = "AWS"
	CredentialTypeHttp  CredentialType = "HTTP_REQUEST_AGENT"
	CredentialTypeJwt   CredentialType = "JWT"
	CredentialTypeMtls  CredentialType = "MTLS"
	CredentialTypeMulti CredentialType = "MULTI_REQUEST"
	CredentialTypeOauth CredentialType = "OAUTH"
	CredentialTypeText  CredentialType = "TEXT"
)

type Credential struct {
	Id           int            `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Mode         CredentialType `json:"mode,omitempty"`
	TeamId       int            `json:"team_id,omitempty"`
	FolderId     int            `json:"folder_id,omitempty"`
	ReadAccess   string         `json:"read_access,omitempty"`
	SharedTeams  []string       `json:"shared_team_slugs,omitempty"`
	Description  string         `json:"description,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	AllowedHosts map[string]any `json:"allowed_hosts,omitempty"`
	TestCred     bool           `json:"test_credential_enabled,omitempty"`
	IsTest       bool           `json:"is_test,omitempty"`
	CredentialPayload
}

type CredentialList struct {
	UserCredentials []Credential  `json:"user_credentials,omitempty"`
	Meta            paginate.Meta `json:"meta,omitempty"`
}

// The credential payload values to set when creating or updating a credential
// depend on the CredentialType set for the Mode value. The only exception is
// for CredentialTypeMulti payloads, which are a form of the HTTP credential
// type and will accept the HttpReqTokenLoc and HttpReqTtl values as well.
type CredentialPayload struct {
	// CredentialTypeAws
	AwsAuthType    string `json:"aws_authentication_type,omitempty"`
	AwsAccessKey   string `json:"aws_access_key,omitempty"`
	AwsSecretKey   string `json:"aws_secret_key,omitempty"`
	AwsAssumedRole string `json:"aws_assumed_role_arn,omitempty"`

	// CredentialTypeHttp
	HttpReqOpts     string `json:"http_request_options,omitempty"`
	HttpReqTokenLoc string `json:"http_request_location_of_token,omitempty"`
	HttpReqSecret   string `json:"http_request_secret,omitempty"`
	HttpReqTtl      int    `json:"http_request_ttl,omitempty"`

	// CredentialTypeJwt
	JwtAlgo       string         `json:"jwt_algorithm,omitempty"`
	JwtPayload    map[string]any `json:"jwt_payload,omitempty"`
	JwtAutoClaims bool           `json:"jwt_auto_generate_time_claims,omitempty"`
	JwtPrivKey    string         `json:"jwt_private_key,omitempty"`

	// CredentialTypeMtls
	MtlsCliCert    string `json:"mtls_client_certificate,omitempty"`
	MtlsCliPrivKey string `json:"mtls_client_private_key,omitempty"`
	MtlsRootCert   string `json:"mtls_root_certificate,omitempty"`

	// CredentialTypeMulti
	MultiCredReqs []CredentialMultiReq `json:"credential_requests,omitempty"`

	// CredentialTypeOauth
	OauthUrl          string `json:"oauth_url,omitempty"`
	OauthClientId     string `json:"oauth_client_id,omitempty"`
	OauthClientSecret string `json:"oauth_client_secret,omitempty"`
	OauthScope        string `json:"oauth_scope,omitempty"`
	OauthGrant        string `json:"oauth_grant_type,omitempty"`
	OauthPkce         string `json:"oauthPkceCodeChallengeMethod,omitempty"`

	// CredentialTypeText
	TextValue string `json:"value,omitempty"`
}

type CredentialMultiReq struct {
	Options struct {
		Url         string         `json:"url,omitempty"`
		ContentType string         `json:"content_type,omitempty"`
		Method      string         `json:"method,omitempty"`
		Payload     map[string]any `json:"payload,omitempty"`
		Headers     map[string]any `json:"headers,omitempty"`
	} `json:"options,omitempty"`
	Secret string `json:"http_request_secret,omitempty"`
}

func (c *Client) CreateCredential(ctx context.Context, cred *Credential) (*Credential, error) {
	return &Credential{}, nil
}

func (c *Client) GetCredential(ctx context.Context, id int) (*Credential, error) {
	return &Credential{}, nil
}

func (c *Client) UpdateCredential(ctx context.Context, cred *Credential) (*Credential, error) {
	return &Credential{}, nil
}

func (c *Client) ListCredentials(ctx context.Context, f ListFilter) iter.Seq2[Credential, error] {
	var credentialList, resultList CredentialList
	resource := "/api/v1/user_credentials"
	params := f.ToParamMap()
	page := paginate.Cursor{
		TotalRequested: f.MaxResults(),
	}
	return func(yield func(Credential, error) bool) {

		for !page.MaxResultsReturned() {
			res, err := c.doRequest(ctx, http.MethodGet, resource, params, nil)
			if err != nil {
				yield(Credential{}, err)
				return
			}

			err = json.Unmarshal(res, &resultList)
			if err != nil {
				yield(Credential{}, err)
				return
			}

			page.UpdatePagination(resultList.Meta)
			params = page.GetNextPageParams()

			for _, v := range resultList.UserCredentials {
				credentialList.UserCredentials = append(credentialList.UserCredentials, v)
				page.IncrementCounter()
				if page.MaxResultsReturned() {
					c.logger.Debug("hit the limit of results to return")
					break
				}
			}

			// Clear the temporary result buffer
			resultList = CredentialList{}

			if !page.ReturnMoreResults() {
				c.logger.Debug("no more results to return")
				break
			}
		}

		for _, v := range credentialList.UserCredentials {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func (c *Client) DeleteCredential(ctx context.Context, id int) error {
	return nil
}
