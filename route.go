package rail

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

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
