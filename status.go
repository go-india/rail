package rail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// LiveTrainStatusReq parameters
type LiveTrainStatusReq struct {
	// Specifies the train number.
	TrainNumber uint32 `validate:"required"`
	// Specifies the date for which result is required.
	Date time.Time `validate:"required"`
}

// Request encodes LiveTrainStatusReq parameters returning a new http.Request
func (r LiveTrainStatusReq) Request() (*http.Request, error) {
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

// LiveTrainStatusResp of the request
type LiveTrainStatusResp struct {
	Train          *Train     `json:"train,omitempty"`
	CurrentStation *Station   `json:"current_station,omitempty"`
	Route          []Route    `json:"route,omitempty"`
	StartDate      *time.Time // `json:"start_date,omitempty"`
	PositionRemark *string    `json:"position,omitempty"`

	*Response
}

// UnmarshalJSON convert JSON data to struct
func (s *LiveTrainStatusResp) UnmarshalJSON(data []byte) error {
	type Alias LiveTrainStatusResp
	t := struct {
		Alias
		Start string `json:"start_date"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}
	*s = LiveTrainStatusResp(t.Alias)

	if t.Start != "" {
		start, err := time.Parse("2 Jan 2006", t.Start)
		if err != nil {
			return errors.Wrap(err, "parse StartDate failed")
		}
		s.StartDate = &start
	}

	return nil
}

// LiveTrainStatus gets live running status of a Train.
func (c Client) LiveTrainStatus(ctx context.Context,
	TrainNumber uint32,
	Date time.Time,
) (LiveTrainStatusResp, error) {
	if c.Auth == nil {
		return LiveTrainStatusResp{}, ErrNoAuth
	}

	var r LiveTrainStatusResp
	err := c.Do(c.Auth(WithCtx(ctx, LiveTrainStatusReq{
		TrainNumber: TrainNumber,
		Date:        Date,
	})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// TrainRouteReq parameters
type TrainRouteReq struct {
	TrainNumber uint32 `validate:"required"` // Specifies the train number.
}

// Request encodes TrainRouteReq parameters returning a new http.Request
func (r TrainRouteReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/route"
	urlStr += fmt.Sprintf("/train/%d", r.TrainNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainRouteResp holds route information of a train
type TrainRouteResp struct {
	Train *Train  `json:"train,omitempty"`
	Route []Route `json:"route,omitempty"`

	*Response
}

// TrainRoute gets details about all the stations in the train’s route.
func (c Client) TrainRoute(ctx context.Context, TrainNumber uint32) (TrainRouteResp, error) {
	if c.Auth == nil {
		return TrainRouteResp{}, ErrNoAuth
	}

	var r TrainRouteResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainRouteReq{TrainNumber})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
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

// CheckSeat gets train seat availability.
func (c Client) CheckSeat(ctx context.Context,
	TrainNumber uint32,
	FromStationCode string,
	ToStationCode string,
	Class string,
	Quota string,
	Date time.Time,
) (CheckSeatResp, error) {
	if c.Auth == nil {
		return CheckSeatResp{}, ErrNoAuth
	}

	var r CheckSeatResp
	err := c.Do(c.Auth(WithCtx(ctx, CheckSeatReq{
		TrainNumber:     TrainNumber,
		FromStationCode: FromStationCode,
		ToStationCode:   ToStationCode,
		Class:           Class,
		Quota:           Quota,
		Date:            Date,
	})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// PNRStatusReq parameters
type PNRStatusReq struct {
	PNRNumber uint64 `validate:"required"` // Specifies the pnr number.
}

// Request encodes PNRStatusReq parameters returning a new http.Request
func (r PNRStatusReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/pnr-status"
	urlStr += fmt.Sprintf("/pnr/%d", r.PNRNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// PNRStatusResp is the response for a PNRReq
type PNRStatusResp struct {
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
func (p *PNRStatusResp) UnmarshalJSON(data []byte) error {
	type Alias PNRStatusResp
	t := struct {
		Alias
		DOJ string `json:"doj"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}
	*p = PNRStatusResp(t.Alias)

	if t.DOJ != "" {
		doj, err := time.Parse("02-01-2006", t.DOJ)
		if err != nil {
			return errors.Wrap(err, "parse DateOfJourney failed")
		}
		p.DateOfJourney = &doj
	}

	return nil
}

// PNRStatus gets PNR status details.
func (c Client) PNRStatus(ctx context.Context, PNRNumber uint64) (PNRStatusResp, error) {
	if c.Auth == nil {
		return PNRStatusResp{}, ErrNoAuth
	}

	var r PNRStatusResp
	err := c.Do(c.Auth(WithCtx(ctx, PNRStatusReq{PNRNumber})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}

// TrainFareReq parameters
type TrainFareReq struct {
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

// Request encodes TrainFareReq parameters returning a new http.Request
func (r TrainFareReq) Request() (*http.Request, error) {
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

// TrainFareResp holds fare details for a train journey
type TrainFareResp struct {
	FromStation  *Station    `json:"from_station,omitempty"`
	ToStation    *Station    `json:"to_station,omitempty"`
	Quota        *Quota      `json:"quota,omitempty"`
	Train        *Train      `json:"train,omitempty"`
	Fare         *float64    `json:"fare,omitempty"`
	JourneyClass *Class      `json:"journey_class,omitempty"`
	Availability []Available `json:"availability,omitempty"`

	*Response
}

// TrainFare gets fares of a train.
func (c Client) TrainFare(ctx context.Context,
	TrainNumber uint32,
	FromStationCode string,
	ToStationCode string,
	Age uint8,
	Class string,
	Quota string,
	Date time.Time,
) (TrainFareResp, error) {
	if c.Auth == nil {
		return TrainFareResp{}, ErrNoAuth
	}

	var r TrainFareResp
	err := c.Do(c.Auth(WithCtx(ctx, TrainFareReq{
		TrainNumber:     TrainNumber,
		FromStationCode: FromStationCode,
		ToStationCode:   ToStationCode,
		Age:             Age,
		Class:           Class,
		Quota:           Quota,
		Date:            Date,
	})), &r)
	return r, errors.Wrap(err, "Client.Do failed")
}
