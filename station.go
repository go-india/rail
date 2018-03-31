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
