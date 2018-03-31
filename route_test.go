package rail_test

import (
	"testing"

	"github.com/go-india/rail"
)

func TestRoute(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.RouteReq{
		TrainNumber: 14311,
	}

	var resp rail.RouteResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Route) < 1 {
		t.Fatal("invalid routes length")
	}
}
