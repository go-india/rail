package rail_test

import (
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestCheckSeat(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.CheckSeatReq{
		TrainNumber: 14311,
		StationFrom: "BE",
		StationTo:   "ADI",
		Class:       "SL",
		Quota:       "GN",
		Date:        time.Date(2018, time.April, 05, 0, 0, 0, 0, time.UTC),
	}

	var resp rail.CheckSeatResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Availability) < 1 {
		t.Fatalf("invalid Availability length")
	}
}
