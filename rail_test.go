package rail_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-india/rail"
	"github.com/pkg/errors"
)

var (
	testServer     *url.URL
	testDataDir    = "./testdata/"
	updateTestData = flag.Bool("update", false, "if True run integration tests; if False run internal tests")
)

// returns APIKey from the environment
func getAPIKey() string {
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	return env("RAILWAYAPI_TEST_API_KEY", "API_KEY")
}

func TestMain(m *testing.M) {
	flag.Parse()

	// Run testServer for unit tests
	if !*updateTestData {
		server := httptest.NewServer(http.FileServer(http.Dir(testDataDir)))

		surl, err := url.Parse(server.URL)
		if err != nil {
			fmt.Println("testServer URL parse failed:", err)
			os.Exit(1)
		}
		testServer = surl

		defer server.Close()
	}

	os.Exit(m.Run())
	return
}

func testClient(c *rail.Client, t *testing.T) {
	c.HTTPClient = &http.Client{}

	if *updateTestData {
		c.HTTPClient.Transport = &saverTransport{t}
		return
	}

	c.HTTPClient.Transport = &loaderTransport{
		filepath: testServer.String() + "/" + filename(t),
	}
	return
}

// saverTransport saves response body to testdata file
type saverTransport struct{ t *testing.T }

func (st saverTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return resp, errors.Wrap(err, "request failed")
	}

	if resp.StatusCode != http.StatusOK {
		return resp, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, errors.Wrap(err, "read body failed")
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = ioutil.WriteFile(testDataDir+filename(st.t), body, 0644)
	return resp, errors.Wrap(err, "write file failed")
}

// loaderTransport loads response from testdata file
type loaderTransport struct{ filepath string }

func (lt loaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return http.Get(lt.filepath)
}

func filename(t *testing.T) string {
	name := t.Name()
	if strings.Contains(name, "/") { // If a subtest
		name = name[strings.LastIndex(t.Name(), "/")+1:]
	}
	name = strings.TrimPrefix(name, "Test")
	return name + ".json"
}
