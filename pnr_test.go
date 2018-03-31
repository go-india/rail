package rail_test

import (
	"testing"

	"github.com/go-india/rail"
)

func TestPNR(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.PNRReq{
		PNRNumber: 2144287856,
	}

	var resp rail.PNRResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Passengers) < 1 {
		t.Fatal("invalid Passengers length")
	}
}
