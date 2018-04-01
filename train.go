package rail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TrainByNumberReq parameters
type TrainByNumberReq struct {
	TrainNumber uint32 `validate:"required"` // Specifies the train number.
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
	Train *Train `json:"train,omitempty"`
	*Response
}

// TrainByNumber gets train details by its number.
func (c Client) TrainByNumber(ctx context.Context, TrainNumber uint32) (TrainResp, error) {
	if c.Auth == nil {
		return TrainResp{}, ErrNoAuth
	}

	var r TrainResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainByNumberReq{TrainNumber})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
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

// TrainByName gets train details by its number.
func (c Client) TrainByName(ctx context.Context, TrainName string) (TrainResp, error) {
	if c.Auth == nil {
		return TrainResp{}, ErrNoAuth
	}

	var r TrainResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainByNameReq{TrainName})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
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
	Source      *Station   `json:"source,omitempty"`
	Destination *Station   `json:"dest,omitempty"`
	Type        *string    `json:"type,omitempty"`
	StartDate   *time.Time //`json:"start_time,omitempty"`

	*Train
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
	*s = TrainSemi(t.Alias)

	if t.Start != "" {
		start, err := time.Parse("2 Jan 2006", t.Start)
		if err != nil {
			return errors.Wrap(err, "parse StartDate failed")
		}
		s.StartDate = &start
	}

	return nil
}

// CancelledTrainsResp holds cancelled trains details
type CancelledTrainsResp struct {
	Trains []TrainSemi `json:"trains,omitempty"`
	Total  *int        `json:"total,omitempty"`

	*Response
}

// CancelledTrains gets list of all cancelled trains on a particular day.
func (c Client) CancelledTrains(ctx context.Context, Date time.Time) (CancelledTrainsResp, error) {
	if c.Auth == nil {
		return CancelledTrainsResp{}, ErrNoAuth
	}

	var r CancelledTrainsResp
	err := c.Do(c.Auth(WithCtx(ctx, CancelledTrainsReq{Date})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
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
	FromStation *Station `json:"from_station,omitempty"`
	ToStation   *Station `json:"to_station,omitempty"`

	TimeDifference  *time.Duration // `json:"time_diff,omitempty"`
	RescheduledDate *time.Time     // `json:"rescheduled_date,omitempty"`
	RescheduledTime *time.Time     // `json:"rescheduled_time,omitempty"`

	*Train
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
		s.RescheduledDate = &rd
	}

	if len(t.RescheduledTime) == 5 {
		rt, err := time.Parse("15:04", t.RescheduledTime)
		if err != nil {
			return errors.Wrap(err, "parse RescheduledTime failed")
		}
		s.RescheduledTime = &rt
	}

	if len(t.TimeDifference) == 5 {
		d := func(s string) string { return strings.Replace(s, ":", "h", -1) + "m" }
		dur, err := time.ParseDuration(d(t.TimeDifference))
		if err != nil {
			return errors.Wrap(err, "parse TimeDifference failed")
		}
		s.TimeDifference = &dur
	}
	return nil
}

// RescheduledTrainsResp holds rescheduled trains
type RescheduledTrainsResp struct {
	Trains []RescheduledTrain `json:"trains,omitempty"`
	*Response
}

// RescheduledTrains gets list of all rescheduled trains on a particular date.
func (c Client) RescheduledTrains(ctx context.Context, Date time.Time) (RescheduledTrainsResp, error) {
	if c.Auth == nil {
		return RescheduledTrainsResp{}, ErrNoAuth
	}

	var r RescheduledTrainsResp
	err := c.Do(c.Auth(WithCtx(ctx, RescheduledTrainsReq{Date})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// SuggestTrainByNameReq parameters
type SuggestTrainByNameReq struct {
	// Specifies the Train name.
	TrainName string `validate:"required"`
}

// Request encodes SuggestTrainByNameReq parameters returning a new http.Request
func (r SuggestTrainByNameReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/suggest-train"
	urlStr += fmt.Sprintf("/train/%s", r.TrainName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// Trains holds trains details
type Trains struct {
	Trains []Train `json:"trains,omitempty"`
	*Response
}

// SuggestTrainByName suggests full train names or numbers given a partial train name.
func (c Client) SuggestTrainByName(ctx context.Context, TrainName string) (Trains, error) {
	if c.Auth == nil {
		return Trains{}, ErrNoAuth
	}

	var r Trains
	err := c.Do(c.Auth(WithCtx(ctx, SuggestTrainByNameReq{TrainName})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// SuggestTrainByCodeReq parameters
type SuggestTrainByCodeReq struct {
	// Specifies the Train code.
	TrainCode uint32 `validate:"required"`
}

// Request encodes SuggestTrainByCodeReq parameters returning a new http.Request
func (r SuggestTrainByCodeReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/suggest-train"
	urlStr += fmt.Sprintf("/train/%d", r.TrainCode)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// SuggestTrainByCode suggests full train names or numbers given a partial train code.
func (c Client) SuggestTrainByCode(ctx context.Context, TrainCode uint32) (Trains, error) {
	if c.Auth == nil {
		return Trains{}, ErrNoAuth
	}

	var r Trains
	err := c.Do(c.Auth(WithCtx(ctx, SuggestTrainByCodeReq{TrainCode})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}
