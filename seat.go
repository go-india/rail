package rail

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

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
