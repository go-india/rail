package rail_test

import (
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestFare(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.FareReq{
		TrainNumber: 14311,
		StationFrom: "BE",
		StationTo:   "ADI",
		Age:         24,
		Class:       "SL",
		Quota:       "GN",
		Date:        time.Date(2018, time.April, 05, 0, 0, 0, 0, time.UTC),
	}

	var resp rail.FareResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if resp.ResponseCode != 200 {
		t.Fatalf("invalid ResponseCode")
	}
}
