package ricochet

import (
	"fmt"
	"net/http"
)

// RicochetResponse contains assertion methods
type Response struct {
	*http.Response
}

func (r *Response) IsSuccess() *Response {
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		fmt.Println("\t\tError:", "Expected Success, was", r.StatusCode)
	}

	return r
}

func (r *Response) Is(code int) *Response {
	if r.StatusCode != code {
		fmt.Println("\t\tError:", "Expected status", code, "was", r.StatusCode)
	}

	return r
}
