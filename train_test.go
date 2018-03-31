package rail_test

import (
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestTrainByNumber(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.TrainByNumberReq{
		TrainNumber: 14311,
	}

	var resp rail.TrainResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Train.Classes) < 1 {
		t.Fatal("invalid classes length")
	}
}

func TestTrainByName(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.TrainByNameReq{
		TrainName: "bhopal",
	}

	var resp rail.TrainResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Train.Classes) < 1 {
		t.Fatal("invalid classes length")
	}
}

func TestCancelledTrains(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.CancelledTrainsReq{
		Date: time.Now(),
	}

	var resp rail.CancelledTrainsResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Trains) < 1 {
		t.Fatal("invalid trains length")
	}
}

func TestRescheduledTrains(t *testing.T) {
	c := &rail.Client{
		Auth: rail.NewAuth(getAPIKey()),
	}
	testClient(c, t)

	req := rail.RescheduledTrainsReq{
		Date: time.Now(),
	}

	var resp rail.RescheduledTrainsResp
	err := c.Do(c.Auth(req), &resp)
	if err != nil {
		t.Fatalf("client Do failed: %+v", err)
	}

	if len(resp.Trains) < 1 {
		t.Fatal("invalid trains length")
	}
}
