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
