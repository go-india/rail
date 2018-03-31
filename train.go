package rail

import (
	"encoding/json"
	"fmt"
	"net/http"
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
