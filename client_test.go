package rail_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-india/rail"
	"github.com/pkg/errors"
)

func TestClientDo(t *testing.T) {
	c := rail.Client{}

	var (
		temp   interface{}
		url, _ = url.Parse("http://0.0.0.0")
		req, _ = http.NewRequest(http.MethodGet, url.String(), nil)
	)

	tests := []struct {
		inputRequester rail.Requester
		inputIntoPtr   interface{}

		setup func()
		clear func()

		expected string
	}{
		{
			inputRequester: nil,
			inputIntoPtr:   temp,
			expected:       "requester is nil",
		},
		{
			inputRequester: mockRequester(
				func() (*http.Request, error) { return nil, errors.New("error") },
			),
			inputIntoPtr: temp,
			expected:     "generate HTTP request failed",
		},
		{
			inputRequester: mockRequester(
				func() (*http.Request, error) { return req, nil },
			),
			inputIntoPtr: temp,
			expected:     "HTTP request failed",

			setup: func() { c.HTTPClient = nil },
		},
		{
			inputRequester: mockRequester(
				func() (*http.Request, error) { return req, nil },
			),
			inputIntoPtr: temp,
			expected:     "HTTP request failed",

			setup: func() { c.BaseURL = url },
		},
		{
			inputRequester: mockRequester(
				func() (*http.Request, error) { return req, nil },
			),
			inputIntoPtr: temp,
			expected:     "UnmarshalJSON failed",

			setup: func() {
				c.HTTPClient = &http.Client{}
				c.HTTPClient.Transport = mockTransport(
					func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(bytes.NewBufferString("Boom")),
							Request:    req,
						}, nil
					},
				)
			},
			clear: func() { c.HTTPClient = nil },
		},
		{
			inputRequester: mockRequester(
				func() (*http.Request, error) { return req, nil },
			),
			inputIntoPtr: temp,
			expected:     "request to",

			setup: func() {
				c.HTTPClient = &http.Client{}
				c.HTTPClient.Transport = mockTransport(
					func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Request:    req,
						}, nil
					},
				)
			},
			clear: func() { c.HTTPClient = nil },
		},
	}

	for _, tt := range tests {
		if tt.setup != nil {
			tt.setup()
		}

		output := c.Do(tt.inputRequester, &tt.inputIntoPtr)

		if output != nil && !strings.Contains(output.Error(), tt.expected) {
			t.Fatalf("expected: `%s`, actual `%s`", tt.expected, output.Error())
		}

		if tt.clear != nil {
			tt.clear()
		}
	}
}

// mockRequester mocks Requester and helps in testing.
type mockRequester func() (*http.Request, error)

func (mr mockRequester) Request() (*http.Request, error) { return mr() }

// loaderTransport loads response passed
type mockTransport func(*http.Request) (*http.Response, error)

func (mt mockTransport) RoundTrip(r *http.Request) (*http.Response, error) { return mt(r) }
