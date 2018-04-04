package rail_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-india/rail"
)

func TestTrainByNumber(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.TrainByNumber(context.Background(), 14311)
	if err != nil {
		t.Fatal("TrainByNumber failed:", err)
	}

	if resp.Train == nil || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestTrainByName(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.TrainByName(context.Background(), "duranto")
	if err != nil {
		t.Fatal("TrainByName failed:", err)
	}

	if resp.Train == nil || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestCancelledTrains(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.CancelledTrains(context.Background(), time.Now())
	if err != nil {
		t.Fatal("CancelledTrains failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestRescheduledTrains(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.RescheduledTrains(context.Background(), time.Now())
	if err != nil {
		t.Fatal("RescheduledTrains failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestSuggestTrainByName(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.SuggestTrainByName(context.Background(), "duranto")
	if err != nil {
		t.Fatal("SuggestTrainByName failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}

func TestSuggestTrainByCode(t *testing.T) {
	c := rail.NewClient(getAPIKey())
	testClient(&c, t)

	resp, err := c.SuggestTrainByCode(context.Background(), 14311)
	if err != nil {
		t.Fatal("SuggestTrainByCode failed:", err)
	}

	if len(resp.Trains) < 1 || resp.ResponseCode != 200 {
		t.Fatal("invalid response")
	}
}
