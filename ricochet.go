package ricochet

import (
	"log"
	"net/http"
)

// Ricochet executes http calls
type Ricochet struct {
}

// Get performs http GET request
func (r *Ricochet) Get(url string) *Response {
	resp, err := http.Get(url)
	if err != nil {
		// todo: deal with error
		// this should fail the test and short-circuit it
		log.Println("Get:", err)
	}

	return &Response{
		resp,
	}
}
