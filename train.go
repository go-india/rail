package rail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TrainByNumberReq parameters
type TrainByNumberReq struct {
	TrainNumber int `validate:"required"` // Specifies the train number.
}

// Request encodes TrainByNumberReq parameters returning a new http.Request
func (r TrainByNumberReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/name-number"
	urlStr += fmt.Sprintf("/train/%d", r.TrainNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainResp holds train details
type TrainResp struct {
	Train Train `json:"train"`
	Response
}

// TrainByNameReq parameters
type TrainByNameReq struct {
	TrainName string `validate:"required"` // Specifies the train name.
}

// Request encodes TrainByNameReq parameters returning a new http.Request
func (r TrainByNameReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/name-number"
	urlStr += fmt.Sprintf("/train/%s", r.TrainName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// CancelledTrainsReq parameters
type CancelledTrainsReq struct {
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
}

// Request encodes CancelledTrainsReq parameters returning a new http.Request
func (r CancelledTrainsReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/cancelled"
	urlStr += fmt.Sprintf("/date/%s", date(r.Date))

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainSemi holds semi train information
type TrainSemi struct {
	Source      Station   `json:"source"`
	Destination Station   `json:"dest"`
	Type        string    `json:"type"`
	StartDate   time.Time `json:"start_time"`

	Train
}

// UnmarshalJSON convert JSON data to struct
func (s *TrainSemi) UnmarshalJSON(data []byte) error {
	type Alias TrainSemi
	t := struct {
		Alias
		Start string `json:"start_time"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	start, err := time.Parse("2 Jan 2006", t.Start)
	if err != nil {
		return errors.Wrap(err, "parse StartTime failed")
	}

	*s = TrainSemi(t.Alias)
	s.StartDate = start
	return nil
}

// CancelledTrainsResp holds cancelled trains details
type CancelledTrainsResp struct {
	Trains []TrainSemi `json:"trains"`
	Total  int         `json:"total"`

	Response
}

// RescheduledTrainsReq parameters
type RescheduledTrainsReq struct {
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
}

// Request encodes RescheduledTrainsReq parameters returning a new http.Request
func (r RescheduledTrainsReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/rescheduled"
	urlStr += fmt.Sprintf("/date/%s", date(r.Date))

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// RescheduledTrain holds rescheduled train detail
type RescheduledTrain struct {
	StationFrom Station `json:"from_station"`
	StationTo   Station `json:"to_station"`

	TimeDifference  time.Duration // `json:"time_diff"`
	RescheduledDate time.Time     // `json:"rescheduled_date"`
	RescheduledTime time.Time     // `json:"rescheduled_time"`

	Train
}

// UnmarshalJSON convert JSON data to struct
func (s *RescheduledTrain) UnmarshalJSON(data []byte) error {
	type Alias RescheduledTrain
	t := struct {
		Alias
		TimeDifference  string `json:"time_diff"`
		RescheduledDate string `json:"rescheduled_date"`
		RescheduledTime string `json:"rescheduled_time"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*s = RescheduledTrain(t.Alias)

	if t.RescheduledDate != "" {
		rd, err := time.Parse("02-01-2006", t.RescheduledDate)
		if err != nil {
			return errors.Wrap(err, "parse RescheduledDate failed")
		}
		s.RescheduledDate = rd
	}

	if len(t.RescheduledTime) == 5 {
		rt, err := time.Parse("15:04", t.RescheduledTime)
		if err != nil {
			return errors.Wrap(err, "parse RescheduledTime failed")
		}
		s.RescheduledTime = rt
	}

	if len(t.TimeDifference) == 5 {
		d := func(s string) string { return strings.Replace(s, ":", "h", -1) + "m" }
		dur, err := time.ParseDuration(d(t.TimeDifference))
		if err != nil {
			return errors.Wrap(err, "parse TimeDifference failed")
		}
		s.TimeDifference = dur
	}
	return nil
}

// RescheduledTrainsResp holds rescheduled trains
type RescheduledTrainsResp struct {
	Trains []RescheduledTrain `json:"trains"`
	Response
}
