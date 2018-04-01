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

// TrainBetweenStationsReq parameters
type TrainBetweenStationsReq struct {
	// Specifies the source station code.
	FromStationCode string `validate:"required"`
	// Specifies the destination station code.
	ToStationCode string `validate:"required"`
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
		r.FromStationCode,
		r.ToStationCode,
		date(r.Date),
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// ExtendedTrain holds extended train details
type ExtendedTrain struct {
	*Train

	ToStation              *Station       `json:"to_station,omitempty"`
	FromStation            *Station       `json:"from_station,omitempty"`
	SourceDepartureTime    *time.Time     // `json:"src_departure_time,omitempty"`
	DestinationArrivalTime *time.Time     // `json:"dest_arrival_time,omitempty"`
	TravelDuration         *time.Duration // `json:"travel_time,omitempty"`
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
		et.SourceDepartureTime = &sdt
	}

	if len(t.DestinationArrivalTime) == 5 {
		dat, err := time.Parse("15:04", t.DestinationArrivalTime)
		if err != nil {
			return errors.Wrap(err, "parse DestinationArrivalTime failed")
		}
		et.DestinationArrivalTime = &dat
	}

	if len(t.TravelTime) == 5 {
		d := func(s string) string { return strings.Replace(s, ":", "h", -1) + "m" }
		dur, err := time.ParseDuration(d(t.TravelTime))
		if err != nil {
			return errors.Wrap(err, "parse travelTime failed")
		}
		et.TravelDuration = &dur
	}

	return nil
}

// TrainBetweenStationsResp holds trains between stations
type TrainBetweenStationsResp struct {
	Trains []ExtendedTrain `json:"trains,omitempty"`
	Total  *int            `json:"total,omitempty"`

	*Response
}

// TrainBetweenStations gets trains running between stations.
func (c Client) TrainBetweenStations(ctx context.Context,
	FromStationCode string,
	ToStationCode string,
	Date time.Time,
) (TrainBetweenStationsResp, error) {
	if c.Auth == nil {
		return TrainBetweenStationsResp{}, ErrNoAuth
	}

	var r TrainBetweenStationsResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainBetweenStationsReq{
		FromStationCode: FromStationCode,
		ToStationCode:   ToStationCode,
		Date:            Date,
	})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// TrainArrivalsReq parameters
type TrainArrivalsReq struct {
	// Specifies the source station code.
	StationCode string `validate:"required"`

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
	urlStr += fmt.Sprintf("/station/%s/hours/%d", r.StationCode, hours)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainWithTimings holds train timings
type TrainWithTimings struct {
	*Train

	DelayArrivalTime       *time.Time //`json:"delayarr,omitempty"`
	DelayDepartureTime     *time.Time //`json:"delaydep,omitempty"`
	ScheduledArrivalTime   *time.Time //`json:"scharr,omitempty"`
	ScheduledDepartureTime *time.Time //`json:"schdep,omitempty"`
	ActualDepartureTime    *time.Time //`json:"actdep,omitempty"`
	ActualArrivalTime      *time.Time //`json:"actarr,omitempty"`
}

// UnmarshalJSON convert JSON data to struct
func (r *TrainWithTimings) UnmarshalJSON(data []byte) error {
	t := struct {
		ScheduledArrivalTime   string `json:"scharr"`
		ScheduledDepartureTime string `json:"schdep"`
		ActualDepartureTime    string `json:"actdep"`
		ActualArrivalTime      string `json:"actarr"`

		DelayArrivalTime   string `json:"delayarr"`
		DelayDepartureTime string `json:"delaydep"`

		*Train
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	r.Train = t.Train

	if len(t.ScheduledArrivalTime) == 5 {
		sa, err := time.Parse("15:04", t.ScheduledArrivalTime)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledArrival failed")
		}
		r.ScheduledArrivalTime = &sa
	}

	if len(t.ScheduledDepartureTime) == 5 {
		sd, err := time.Parse("15:04", t.ScheduledDepartureTime)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledDeparture failed")
		}
		r.ScheduledDepartureTime = &sd
	}

	if len(t.ActualDepartureTime) == 5 {
		ad, err := time.Parse("15:04", t.ActualDepartureTime)
		if err != nil {
			return errors.Wrap(err, "parse ActualDeparture failed")
		}
		r.ActualDepartureTime = &ad
	}

	if len(t.ActualArrivalTime) == 5 {
		aa, err := time.Parse("15:04", t.ActualArrivalTime)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrival failed")
		}
		r.ActualArrivalTime = &aa
	}

	if len(t.DelayArrivalTime) == 5 {
		da, err := time.Parse("15:04", t.DelayArrivalTime)
		if err != nil {
			return errors.Wrap(err, "parse DelayArrival failed")
		}
		r.DelayArrivalTime = &da
	}

	if len(t.DelayDepartureTime) == 5 {
		dd, err := time.Parse("15:04", t.DelayDepartureTime)
		if err != nil {
			return errors.Wrap(err, "parse DelayDeparture failed")
		}
		r.DelayDepartureTime = &dd
	}

	return nil
}

// TrainArrivalsResp holds train arrivals details
type TrainArrivalsResp struct {
	Trains []TrainWithTimings `json:"trains,omitempty"`
	Total  *int               `json:"total,omitempty"`

	*Response
}

// TrainArrivals get list of trains arriving at a station within
// a window period along with their live status.
//
// Window time in hours to search, valid values are 2 or 4.
func (c Client) TrainArrivals(ctx context.Context,
	StationCode string,
	Hours WindowHour,
) (TrainArrivalsResp, error) {
	if c.Auth == nil {
		return TrainArrivalsResp{}, ErrNoAuth
	}

	var r TrainArrivalsResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainArrivalsReq{
		StationCode: StationCode,
		Hours:       Hours,
	})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// StationNameToCodeReq parameters
type StationNameToCodeReq struct {
	// Specifies the source station name.
	StationName string `validate:"required"`
}

// Request encodes StationNameToCodeReq parameters returning a new http.Request
func (r StationNameToCodeReq) Request() (*http.Request, error) {
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
	*Response
}

// StationNameToCode gets station details of the given station and
// its nearby stations using partial station name.
// Station’s name is autocompleted.
func (c Client) StationNameToCode(ctx context.Context, name string) (Stations, error) {
	if c.Auth == nil {
		return Stations{}, ErrNoAuth
	}

	var r Stations
	err := c.Do(c.Auth(WithCtx(ctx, StationNameToCodeReq{name})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// StationCodeToNameReq parameters
type StationCodeToNameReq struct {
	// Specifies the source station code.
	StationCode string `validate:"required"`
}

// Request encodes StationCodeToNameReq parameters returning a new http.Request
func (r StationCodeToNameReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/code-to-name"
	urlStr += fmt.Sprintf("/code/%s", r.StationCode)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// StationCodeToName gets station details of the given station and
// its nearby stations using partial station name.
// Station’s name is autocompleted.
func (c Client) StationCodeToName(ctx context.Context, code string) (Stations, error) {
	if c.Auth == nil {
		return Stations{}, ErrNoAuth
	}

	var r Stations
	err := c.Do(c.Auth(WithCtx(ctx, StationCodeToNameReq{code})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// SuggestStationReq parameters
type SuggestStationReq struct {
	// Specifies the source station name.
	StationName string `validate:"required"`
}

// Request encodes SuggestStationReq parameters returning a new http.Request
func (r SuggestStationReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/suggest-station"
	urlStr += fmt.Sprintf("/name/%s", r.StationName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// SuggestStation suggests full station names given a partial station name.
func (c Client) SuggestStation(ctx context.Context, name string) (Stations, error) {
	if c.Auth == nil {
		return Stations{}, ErrNoAuth
	}

	var r Stations
	err := c.Do(c.Auth(WithCtx(ctx, SuggestStationReq{name})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}
