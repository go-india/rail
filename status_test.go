package rail_test

import (
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestLiveStatus(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.LiveStatusReq{
		Date:        time.Now(),
		TrainNumber: 14311,
	}

	var resp rail.LiveStatusResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Route) < 1 {
		t.Fatal("invalid routes length")
	}
}
