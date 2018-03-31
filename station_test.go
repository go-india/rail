package rail_test

import (
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestTrainBetweenStations(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.TrainBetweenStationsReq{
		StationFrom: "BE",
		StationTo:   "ADI",
		Date:        time.Now(),
	}

	var resp rail.TrainBetweenStationsResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Trains) < 1 {
		t.Fatal("invalid trains length")
	}
}

func TestTrainArrivals(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.TrainArrivalsReq{
		Station: "BE",
		Hours:   rail.WindowHour2,
	}

	var resp rail.TrainArrivalsResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Trains) < 1 {
		t.Fatal("invalid trains length")
	}
}

func TestStationCode(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.StationCodeReq{
		StationName: "bareilly",
	}

	var resp rail.Stations
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Stations) < 1 {
		t.Fatal("invalid Stations length")
	}
}

func TestStationName(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.StationNameReq{
		StationCode: "BE",
	}

	var resp rail.Stations
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Stations) < 1 {
		t.Fatal("invalid Stations length")
	}
}

func TestStationAutoComplete(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.StationAutoCompleteReq{
		StationName: "bareilly",
	}

	var resp rail.Stations
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Stations) < 1 {
		t.Fatal("invalid Stations length")
	}
}
