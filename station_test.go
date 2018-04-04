package rail_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestTrainBetweenStations(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	d := time.Date(2018, time.April, 5, 0, 0, 0, 0, time.UTC)
	resp, err := c.TrainBetweenStations(context.Background(), "BE", "ADI", d)
	if err != nil {
		t.Fatal("TrainBetweenStations failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestTrainArrivals(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.TrainArrivals(context.Background(), "BE", rail.WindowHour2)
	if err != nil {
		t.Fatal("TrainArrivals failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestStationNameToCode(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.StationNameToCode(context.Background(), "bareilly")
	if err != nil {
		t.Fatal("StationNameToCode failed:", err)
	}

	if len(resp.Stations) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestStationCodeToName(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.StationCodeToName(context.Background(), "BE")
	if err != nil {
		t.Fatal("StationCodeToName failed:", err)
	}

	if len(resp.Stations) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestSuggestStation(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.SuggestStation(context.Background(), "bareilly")
	if err != nil {
		t.Fatal("SuggestStation failed:", err)
	}

	if len(resp.Stations) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}
