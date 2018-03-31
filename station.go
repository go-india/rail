package rail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// ExtendedTrain holds extended train details
type ExtendedTrain struct {
	Train

	StationTo          Station       `json:"to_station"`
	StationFrom        Station       `json:"from_station"`
	SourceDeparture    time.Time     // `json:"src_departure_time"`
	DestinationArrival time.Time     // `json:"dest_arrival_time"`
	TravelDuration     time.Duration // `json:"travel_time"`
}

// UnmarshalJSON convert JSON data to struct
func (et *ExtendedTrain) UnmarshalJSON(data []byte) error {
	type Alias ExtendedTrain
	t := struct {
		Alias
		SourceDepartureTime    string `json:"src_departure_time"`
		DestinationArrivalTime string `json:"dest_arrival_time"`
		TravelTime             string `json:"travel_time"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*et = ExtendedTrain(t.Alias)

	if len(t.SourceDepartureTime) == 5 {
		sdt, err := time.Parse("15:04", t.SourceDepartureTime)
		if err != nil {
			return errors.Wrap(err, "parse SourceDepartureTime failed")
		}
		et.SourceDeparture = sdt
	}

	if len(t.DestinationArrivalTime) == 5 {
		dat, err := time.Parse("15:04", t.DestinationArrivalTime)
		if err != nil {
			return errors.Wrap(err, "parse DestinationArrivalTime failed")
		}
		et.DestinationArrival = dat
	}

	if len(t.TravelTime) == 5 {
		d := func(s string) string { return strings.Replace(s, ":", "h", -1) + "m" }
		dur, err := time.ParseDuration(d(t.TravelTime))
		if err != nil {
			return errors.Wrap(err, "parse travelTime failed")
		}
		et.TravelDuration = dur
	}

	return nil
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

// TrainWithTimings holds train timings
type TrainWithTimings struct {
	Train

	DelayArrival   time.Time //`json:"delayarr"`
	DelayDeparture time.Time //`json:"delaydep"`

	ScheduledArrival   time.Time //`json:"scharr"`
	ScheduledDeparture time.Time //`json:"schdep"`
	ActualDeparture    time.Time //`json:"actdep"`
	ActualArrival      time.Time //`json:"actarr"`
}

// UnmarshalJSON convert JSON data to struct
func (r *TrainWithTimings) UnmarshalJSON(data []byte) error {
	t := struct {
		ScheduledArrival   string `json:"scharr"`
		ScheduledDeparture string `json:"schdep"`
		ActualDeparture    string `json:"actdep"`
		ActualArrival      string `json:"actarr"`

		DelayArrival   string `json:"delayarr"`
		DelayDeparture string `json:"delaydep"`

		Train
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	r.Train = t.Train

	if len(t.ScheduledArrival) == 5 {
		sa, err := time.Parse("15:04", t.ScheduledArrival)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledArrival failed")
		}
		r.ScheduledArrival = sa
	}

	if len(t.ScheduledDeparture) == 5 {
		sd, err := time.Parse("15:04", t.ScheduledDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledDeparture failed")
		}
		r.ScheduledDeparture = sd
	}

	if len(t.ActualDeparture) == 5 {
		ad, err := time.Parse("15:04", t.ActualDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ActualDeparture failed")
		}
		r.ActualDeparture = ad
	}

	if len(t.ActualArrival) == 5 {
		aa, err := time.Parse("15:04", t.ActualArrival)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrival failed")
		}
		r.ActualArrival = aa
	}

	if len(t.DelayArrival) == 5 {
		da, err := time.Parse("15:04", t.DelayArrival)
		if err != nil {
			return errors.Wrap(err, "parse DelayArrival failed")
		}
		r.DelayArrival = da
	}

	if len(t.DelayDeparture) == 5 {
		dd, err := time.Parse("15:04", t.DelayDeparture)
		if err != nil {
			return errors.Wrap(err, "parse DelayDeparture failed")
		}
		r.DelayDeparture = dd
	}

	return nil
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

// StationAutoCompleteReq parameters
type StationAutoCompleteReq struct {
	// Specifies the source station name.
	StationName string `validate:"required"`
}

// Request encodes StationAutoCompleteReq parameters returning a new http.Request
func (r StationAutoCompleteReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/suggest-station"
	urlStr += fmt.Sprintf("/name/%s", r.StationName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}
