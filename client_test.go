package rail_test

import (
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
			inputRequester: mockRequester{},
			inputIntoPtr:   temp,
			expected:       "generate HTTP request failed",
		},
		{
			inputRequester: mockRequester{true},
			inputIntoPtr:   temp,
			expected:       "HTTP request failed",

			setup: func() { c.HTTPClient = nil },
		},
		{
			inputRequester: mockRequester{true},
			inputIntoPtr:   temp,
			expected:       "HTTP request failed",

			setup: func() { c.BaseURL = url },
		},
		{
			inputRequester: mockRequester{true},
			inputIntoPtr:   temp,
			expected:       "invalid character 'I' looking for beginning of value",

			setup: func() {
				c.HTTPClient = &http.Client{}
				c.HTTPClient.Transport = loaderTransport{
					testServer.String() + "/Invalid.txt",
				}
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
type mockRequester struct{ r bool }

func (r mockRequester) Request() (*http.Request, error) {
	if r.r {
		return http.NewRequest(http.MethodGet, "http://0.0.0.0", nil)
	}
	return nil, errors.New("error")
}
