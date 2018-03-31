package rail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// LiveStatusReq parameters
type LiveStatusReq struct {
	// Specifies the train number.
	TrainNumber int `validate:"required"`
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
}

// Request encodes LiveStatusReq parameters returning a new http.Request
func (r LiveStatusReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/live"

	urlStr += fmt.Sprintf(
		"/train/%d/date/%s",
		r.TrainNumber,
		date(r.Date),
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// LiveStatusResp of the request
type LiveStatusResp struct {
	Train          Train     `json:"train"`
	CurrentStation Station   `json:"current_station"`
	Route          []Route   `json:"route"`
	StartDate      time.Time // `json:"start_date"`
	PositionRemark string    `json:"position"`

	Response
}

// UnmarshalJSON convert JSON data to struct
func (s *LiveStatusResp) UnmarshalJSON(data []byte) error {
	type Alias LiveStatusResp
	t := struct {
		Alias
		Start string `json:"start_date"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	start, err := time.Parse("2 Jan 2006", t.Start)
	if err != nil {
		return errors.Wrap(err, "parse StartDate failed")
	}

	*s = LiveStatusResp(t.Alias)
	s.StartDate = start
	return nil
}

// RouteReq parameters
type RouteReq struct {
	TrainNumber int `validate:"required"` // Specifies the train number.
}

// Request encodes RouteReq parameters returning a new http.Request
func (r RouteReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/route"
	urlStr += fmt.Sprintf("/train/%d", r.TrainNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// RouteResp holds route information of a train
type RouteResp struct {
	Train Train   `json:"train"`
	Route []Route `json:"route"`

	Response
}

// CheckSeatReq parameters
type CheckSeatReq struct {
	// Specifies the train number.
	TrainNumber int `validate:"required"`
	// Specifies the source station code.
	StationFrom string `validate:"required"`
	// Specifies the destination station code.
	StationTo string `validate:"required"`
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
	// Specifies the class code. Ex: SL/AC/2S
	Class string `validate:"required"`
	// Specifies the quota code. Ex: GN etc
	Quota string `validate:"required"`
}

// Request encodes CheckSeat parameters returning a new http.Request
func (r CheckSeatReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/check-seat"
	urlStr += fmt.Sprintf(
		"/train/%d/source/%s/dest/%s/date/%s/pref/%s/quota/%s",
		r.TrainNumber,
		r.StationFrom,
		r.StationTo,
		date(r.Date),
		r.Class,
		r.Quota,
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// CheckSeatResp holds seat availability response
type CheckSeatResp struct {
	Train        Train       `json:"train"`
	StationFrom  Station     `json:"from_station"`
	StationTo    Station     `json:"to_station"`
	Quota        Quota       `json:"quota"`
	JourneyClass Class       `json:"journey_class"`
	Availability []Available `json:"availability"`

	Response
}

// PNRReq parameters
type PNRReq struct {
	PNRNumber int `validate:"required"` // Specifies the pnr number.
}

// Request encodes PNRReq parameters returning a new http.Request
func (r PNRReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/pnr-status"
	urlStr += fmt.Sprintf("/pnr/%d", r.PNRNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// PNRResp is the response for a PNRReq
type PNRResp struct {
	ChartPrepared   bool        `json:"chart_prepared"`
	DateOfJourney   time.Time   // `json:"doj"`
	BoardingPoint   Station     `json:"boarding_point"`
	StationFrom     Station     `json:"from_station"`
	StationTo       Station     `json:"to_station"`
	TotalPassengers int         `json:"total_passengers"`
	JourneyClass    Class       `json:"journey_class"`
	Train           Train       `json:"train"`
	Passengers      []Passenger `json:"passengers"`
	PNR             int64       `json:"pnr,string"`
	ReservationUpto Station     `json:"reservation_upto"`

	Response
}

// UnmarshalJSON convert JSON data to struct
func (p *PNRResp) UnmarshalJSON(data []byte) error {
	type Alias PNRResp
	t := struct {
		Alias
		DOJ string `json:"doj"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	doj, err := time.Parse("02-01-2006", t.DOJ)
	if err != nil {
		return errors.Wrap(err, "parse DateOfJourney failed")
	}

	*p = PNRResp(t.Alias)
	p.DateOfJourney = doj
	return nil
}

// FareReq parameters
type FareReq struct {
	// Specifies the train number.
	TrainNumber int `validate:"required"`
	// Specifies the source station code.
	StationFrom string `validate:"required"`
	// Specifies the destination station code.
	StationTo string `validate:"required"`
	// Specifies the age code of passenger
	Age int `url:"age" validate:"required"`
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
	// Specifies the class code. Ex: SL/AC/2S
	Class string `validate:"required"`
	// Specifies the quota code. Ex: GN etc
	Quota string `validate:"required"`
}

// Request encodes FareReq parameters returning a new http.Request
func (r FareReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/fare"
	urlStr += fmt.Sprintf(
		"/train/%d/source/%s/dest/%s/age/%d/pref/%s/quota/%s/date/%s",
		r.TrainNumber,
		r.StationFrom,
		r.StationTo,
		r.Age,
		r.Class,
		r.Quota,
		date(r.Date),
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// FareResp holds fare details for a train journey
type FareResp struct {
	StationTo    Station     `json:"to_station"`
	Quota        Quota       `json:"quota"`
	Train        Train       `json:"train"`
	StationFrom  Station     `json:"from_station"`
	Fare         float64     `json:"fare"`
	JourneyClass Class       `json:"journey_class"`
	Availability []Available `json:"availability"`

	Response
}
