package rail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"
)

const (
	// DefaultBaseURL is the default base server URL of the API.
	DefaultBaseURL = "https://api.railwayapi.com"
	// DefaultUserAgent is the default user agent used by client.
	DefaultUserAgent = "go-india/rail"
)

var (
	// ErrNoAuth is returned when auth is not defined for a client.
	ErrNoAuth = errors.New("rail: no authenticator in client")
)

// Requester is implemented by any value that has a Request method.
type Requester interface {
	// Request should generate an HTTP request from parameters.
	Request() (*http.Request, error)
}

// RequesterFunc implements Requester
type RequesterFunc func() (*http.Request, error)

// Request invokes 'f'
func (f RequesterFunc) Request() (*http.Request, error) {
	return f()
}

// WithCtx applies 'ctx' to the the http.Request and returns a new Requester
func WithCtx(ctx context.Context, r Requester) Requester {
	if ctx == nil {
		return r
	}
	return RequesterFunc(func() (*http.Request, error) {
		req, err := r.Request()
		return req.WithContext(ctx), err
	})
}

// Client is an railwayapi's HTTP REST API client instance.
//
// Its zero value is usable client that uses http.DefaultTransport.
// Client is safe for use by multiple go routines.
type Client struct {
	// BaseURL is the base URL of the API server.
	BaseURL *url.URL
	// User agent used when communicating with the API.
	UserAgent string
	// HTTPClient is a reusable http client instance.
	HTTPClient *http.Client

	// Auth holds authenticator function used to authenticate requests.
	//
	// Client methods uses Auth to add APIKey to requests.
	// Use NewAuth(apikey) to generate a new authenticator.
	Auth func(Requester) Requester
}

// Do sends the http.Request and unmarshalls the JSON response into 'intoPtr'.
func (c Client) Do(r Requester, intoPtr interface{}) error {
	if r == nil {
		return errors.New("requester is nil")
	}

	req, err := r.Request()
	if err != nil {
		return errors.Wrap(err, "generate HTTP request failed")
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
		client.Transport = http.DefaultTransport
		client.Timeout = 15 * time.Second
	}

	if c.BaseURL != nil {
		req.URL.Scheme = c.BaseURL.Scheme
		req.URL.Host = c.BaseURL.Host
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "HTTP request failed")
	}

	defer func() {
		// Read the body if small so underlying TCP connection will be re-used.
		// No need to check for errors: if it fails, Transport won't reuse it anyway.
		if rsp.Body != nil {
			const maxBodySlurpSize = 2 << 10
			if rsp.ContentLength == -1 || rsp.ContentLength <= maxBodySlurpSize {
				io.CopyN(ioutil.Discard, rsp.Body, maxBodySlurpSize)
			}
			rsp.Body.Close()
		}
	}()

	if rsp.StatusCode != http.StatusOK {
		return ErrAPI{rsp}
	}

	return errors.Wrap(json.NewDecoder(rsp.Body).Decode(intoPtr), "UnmarshalJSON failed")
}

// ErrAPI is returned by API calls when the response status code isn't 200.
type ErrAPI struct {
	// Response from the request which returned error.
	Response *http.Response
}

// Error implements the error interface.
func (err ErrAPI) Error() (errStr string) {
	if err.Response != nil {
		errStr += fmt.Sprintf(
			"request to %s returned %d (%s)",
			err.Response.Request.URL,
			err.Response.StatusCode,
			http.StatusText(err.Response.StatusCode),
		)
	}
	return errStr
}

// NewAuth returns a new authenticator function.
//
// Assign to client.Auth field to make client methods use it for requests.
func NewAuth(APIKey string) func(Requester) Requester {
	return func(r Requester) Requester {
		return RequesterFunc(func() (*http.Request, error) {
			req, err := r.Request()
			if err != nil {
				return req, errors.Wrap(err, "generate HTTP request failed")
			}

			req.URL.Path = path.Join(req.URL.Path, fmt.Sprintf("/apikey/%s/", APIKey))
			return req, nil
		})
	}
}

// NewClient returns a new RailwayAPI.com authenticated API client.
//
// Use returned client's methods to access various API functions.
func NewClient(APIKey string) Client {
	return Client{Auth: NewAuth(APIKey)}
}
