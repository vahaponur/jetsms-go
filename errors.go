package jetsms

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// HTTPError describes a non-2xx HTTP response from the API.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("jetsms: http error %d: %s", e.StatusCode, e.Body)
}

func httpError(resp *resty.Response) error {
	return &HTTPError{
		StatusCode: resp.StatusCode(),
		Body:       string(resp.Body()),
	}
}

