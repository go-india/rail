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
