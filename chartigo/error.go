package chartigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	ErrMsg401 string = "Unauthorized. The provided api_key and secret combination is invalid."
	ErrMsg403 string = "Forbidden. The api_key is not allowed to access the content provided by the requested URL."
)

var (
	// ErrMissingAPIKey is an error that is returned when the required
	// environment variable CHARTIO_API_KEY was not found.
	ErrMissingAPIKey = errors.New("CHARTIO_API_KEY not found in environment, required for authentication.")

	// ErrMissingAPIPassword is an error that is returned when the required
	// environment variable CHARTIO_API_PASSWORD was not found.
	ErrMissingAPIPassword = errors.New("CHARTIO_API_PASSWORD not found in environment, required for authentication.")
)

// HTTPError wraps non-2XX responses from the Chartio API.
type HTTPError struct {
	StatusCode int
	Detail     string `json:"detail"`
}

// NewHTTPError creates a new HTTP error from the given code.
func NewHTTPError(resp *http.Response) *HTTPError {
	var e HTTPError
	e.StatusCode = resp.StatusCode

	if resp.Body == nil {
		return &e
	}
	switch e.StatusCode {
	case 401:
		e.Detail = ErrMsg401
	case 403:
		e.Detail = ErrMsg403
	default:
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			panic(err)
		}
	}

	return &e
}

// Error implements the error interface and returns the string representing the
// error text that includes the status code and the corresponding status text.
func (e *HTTPError) Error() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%d - %s", e.StatusCode, http.StatusText(e.StatusCode))
	if e.Detail != "" {
		fmt.Fprintf(&b, "\n    Detail:     %s", e.Detail)
	}

	return b.String()
}

// String implements the stringer interface and returns the string representing
// the string text that includes the status code and corresponding status text.
func (e *HTTPError) String() string {
	return e.Error()
}

// wrapHTTPResp surfaces non-2XX error codes from an HTTP Request.
func wrapHTTPResp(resp *http.Response, err error) (*http.Response, error) {
	// Return the error if one was already there
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, NewHTTPError(resp)
	}
}
