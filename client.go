package cidirectinvesting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	// APIEndpoint is the production API offered by CI Direct Investing
	APIEndpoint = "https://app.cidirectinvesting.com/api"
	// UserAgent is a nice gesture so that the CI Direct Investing folks know who's hitting their API
	UserAgent = "cidirectinvesting/go-v0"
)

// Client holds the required state to perform API calls with the CI Direct Investing API
type Client struct {
	APIEndpoint string

	apiKey    string
	apiSecret string

	client *http.Client
}

// Options allow the configuration of an API Client
type Options func(c *Client) error

// WithThirdPartyKey sets the required authentication parameters
func WithThirdPartyKey(key, secret string) Options {
	return func(c *Client) error {
		c.apiSecret = secret
		c.apiKey = key
		return nil
	}
}

// WithAPIEndpoint sets the APIEndpoint to be used. By default the real production endpoint is used
func WithAPIEndpoint(endpoint string) Options {
	return func(c *Client) error {
		c.APIEndpoint = endpoint
		return nil
	}
}

// New creates a new API Client
func New(opts ...Options) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c := &Client{
		APIEndpoint: APIEndpoint,
		client: &http.Client{
			Jar: jar,
		},
	}
	for _, option := range opts {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", c.APIEndpoint, req.URL)
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	req.URL = u

	//log.Println(req.URL.String())

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Encoding", "UTF8")
	req.Header.Set("Accept-Language", "en-us")

	return c.client.Do(req)
}

func (c *Client) doWithProcessing(req *http.Request, decodingTarget interface{}, responseCodes []int) error {
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	hasMatchingResponseCode := false
	for _, acceptableCode := range responseCodes {
		if acceptableCode == resp.StatusCode {
			hasMatchingResponseCode = true
			break
		}
	}
	if !hasMatchingResponseCode {
		return fmt.Errorf("no matching HTTP status code found: expected %q, got %d", responseCodes, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(decodingTarget); err != nil {
		return err
	}

	return nil
}
