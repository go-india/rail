package rail

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

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
