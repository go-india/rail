package rail

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// TrainByNumberReq parameters
type TrainByNumberReq struct {
	TrainNumber int `validate:"required"` // Specifies the train number.
}

// Request encodes TrainByNumberReq parameters returning a new http.Request
func (r TrainByNumberReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/name-number"
	urlStr += fmt.Sprintf("/train/%d", r.TrainNumber)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}

// TrainResp holds train details
type TrainResp struct {
	Train Train `json:"train"`
	Response
}

// TrainByNameReq parameters
type TrainByNameReq struct {
	TrainName string `validate:"required"` // Specifies the train name.
}

// Request encodes TrainByNameReq parameters returning a new http.Request
func (r TrainByNameReq) Request() (*http.Request, error) {
	err := validate.Struct(r)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}

	urlStr := DefaultBaseURL + "/v2/name-number"
	urlStr += fmt.Sprintf("/train/%s", r.TrainName)

	return http.NewRequest(http.MethodGet, urlStr, nil)
}
