package rail_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestLiveTrainStatus(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.LiveTrainStatus(context.Background(), 12138, time.Now())
	if err != nil {
		t.Fatal("LiveTrainStatus failed:", err)
	}

	if len(resp.Route) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestTrainRoute(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.TrainRoute(context.Background(), 14311)
	if err != nil {
		t.Fatal("TrainRoute failed:", err)
	}

	if len(resp.Route) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestCheckSeat(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	d := time.Date(2018, time.April, 5, 0, 0, 0, 0, time.UTC)

	resp, err := c.CheckSeat(context.Background(), 14311, "BE", "ADI", "SL", "GN", d)
	if err != nil {
		t.Fatal("CheckSeat failed:", err)
	}

	if len(resp.Availability) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestPNRStatus(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.PNRStatus(context.Background(), 2144287856)
	if err != nil {
		t.Fatal("PNRStatus failed:", err)
	}

	if len(resp.Passengers) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestTrainFare(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	d := time.Date(2018, time.April, 5, 0, 0, 0, 0, time.UTC)

	resp, err := c.TrainFare(context.Background(), 14311, "BE", "ADI", 24, "SL", "GN", d)
	if err != nil {
		t.Fatal("TrainFare failed:", err)
	}

	if resp.Train == nil || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}
