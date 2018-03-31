package rail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

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
