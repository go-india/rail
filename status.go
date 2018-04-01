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
	TrainNumber uint32 `validate:"required"`
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
	Train          *Train     `json:"train,omitempty"`
	CurrentStation *Station   `json:"current_station,omitempty"`
	Route          []Route    `json:"route,omitempty"`
	StartDate      *time.Time // `json:"start_date,omitempty"`
	PositionRemark *string    `json:"position,omitempty"`

	*Response
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
	*s = LiveStatusResp(t.Alias)

	if t.Start != "" {
		start, err := time.Parse("2 Jan 2006", t.Start)
		if err != nil {
			return errors.Wrap(err, "parse StartDate failed")
		}
		s.StartDate = &start
	}

	return nil
}

// RouteReq parameters
type RouteReq struct {
	TrainNumber uint32 `validate:"required"` // Specifies the train number.
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
	Train *Train  `json:"train,omitempty"`
	Route []Route `json:"route,omitempty"`

	*Response
}

// CheckSeatReq parameters
type CheckSeatReq struct {
	// Specifies the train number.
	TrainNumber uint32 `validate:"required"`
	// Specifies the source station code.
	FromStationCode string `validate:"required"`
	// Specifies the destination station code.
	ToStationCode string `validate:"required"`
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
		r.FromStationCode,
		r.ToStationCode,
		date(r.Date),
		r.Class,
		r.Quota,
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// CheckSeatResp holds seat availability response
type CheckSeatResp struct {
	Train        *Train      `json:"train,omitempty"`
	FromStation  *Station    `json:"from_station,omitempty"`
	ToStation    *Station    `json:"to_station,omitempty"`
	Quota        *Quota      `json:"quota,omitempty"`
	JourneyClass *Class      `json:"journey_class,omitempty"`
	Availability []Available `json:"availability,omitempty"`

	*Response
}

// PNRReq parameters
type PNRReq struct {
	PNRNumber uint64 `validate:"required"` // Specifies the pnr number.
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
	ChartPrepared   *bool       `json:"chart_prepared,omitempty"`
	DateOfJourney   *time.Time  // `json:"doj,omitempty"`
	BoardingPoint   *Station    `json:"boarding_point,omitempty"`
	FromStation     *Station    `json:"from_station,omitempty"`
	ToStation       *Station    `json:"to_station,omitempty"`
	TotalPassengers *int        `json:"total_passengers,omitempty"`
	JourneyClass    *Class      `json:"journey_class,omitempty"`
	Train           *Train      `json:"train,omitempty"`
	Passengers      []Passenger `json:"passengers,omitempty"`
	PNR             *uint64     `json:"pnr,string,omitempty"`
	ReservationUpto *Station    `json:"reservation_upto,omitempty"`

	*Response
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
	*p = PNRResp(t.Alias)

	if t.DOJ != "" {
		doj, err := time.Parse("02-01-2006", t.DOJ)
		if err != nil {
			return errors.Wrap(err, "parse DateOfJourney failed")
		}
		p.DateOfJourney = &doj
	}

	return nil
}

// FareReq parameters
type FareReq struct {
	// Specifies the train number.
	TrainNumber uint32 `validate:"required"`
	// Specifies the source station code.
	FromStationCode string `validate:"required"`
	// Specifies the destination station code.
	ToStationCode string `validate:"required"`
	// Specifies the age code of passenger
	Age uint8 `url:"age" validate:"required"`
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
		r.FromStationCode,
		r.ToStationCode,
		r.Age,
		r.Class,
		r.Quota,
		date(r.Date),
	)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// FareResp holds fare details for a train journey
type FareResp struct {
	FromStation  *Station    `json:"from_station,omitempty"`
	ToStation    *Station    `json:"to_station,omitempty"`
	Quota        *Quota      `json:"quota,omitempty"`
	Train        *Train      `json:"train,omitempty"`
	Fare         *float64    `json:"fare,omitempty"`
	JourneyClass *Class      `json:"journey_class,omitempty"`
	Availability []Available `json:"availability,omitempty"`

	*Response
}
