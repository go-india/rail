package rail

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// TrainBetweenStationsReq parameters
type TrainBetweenStationsReq struct {
	// Specifies the source station code.
	StationFrom string `validate:"required"`
	// Specifies the destination station code.
	StationTo string `validate:"required"`
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
}

// Request encodes TrainBetweenStationsReq parameters returning a new http.Request
func (r TrainBetweenStationsReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/between"
	urlStr += fmt.Sprintf(
		"/source/%s/dest/%s/date/%s",
		r.StationFrom,
		r.StationTo,
		date(r.Date),
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainBetweenStationsResp holds trains between stations
type TrainBetweenStationsResp struct {
	Trains []ExtendedTrain `json:"trains"`
	Total  int             `json:"total"`

	Response
}

// TrainArrivalsReq parameters
type TrainArrivalsReq struct {
	// Specifies the source station code.
	Station string `validate:"required"`

	// Specifies the windows hours to search.
	//
	// Window time in hours to search, valid values are 2 or 4.
	Hours WindowHour `validate:"required"`
}

// Request encodes TrainArrivalsReq parameters returning a new http.Request
func (r TrainArrivalsReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	var hours uint8
	switch r.Hours {
	case WindowHour2:
		hours = 2
	case WindowHour4:
		hours = 4
	default:
		return nil, errors.New("invalid WindowHour")
	}

	urlStr := DefaultBaseURL + "/v2/arrivals"
	urlStr += fmt.Sprintf("/station/%s/hours/%d", r.Station, hours)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainArrivalsResp holds train arrivals details
type TrainArrivalsResp struct {
	Trains []TrainWithTimings `json:"trains"`
	Total  int                `json:"total"`

	Response
}

// StationCodeReq parameters
type StationCodeReq struct {
	// Specifies the source station name.
	StationName string `validate:"required"`
}

// Request encodes StationCodeReq parameters returning a new http.Request
func (r StationCodeReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/name-to-code"
	urlStr += fmt.Sprintf("/station/%s", r.StationName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// Stations holds stations
type Stations struct {
	Stations []Station `json:"stations"`
	Response
}

// StationNameReq parameters
type StationNameReq struct {
	// Specifies the source station code.
	StationCode string `validate:"required"`
}

// Request encodes StationNameReq parameters returning a new http.Request
func (r StationNameReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/code-to-name"
	urlStr += fmt.Sprintf("/code/%s", r.StationCode)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}
