package ricochet

import (
	"log"
	"net/http"
	"net/url"
	"path"
)

// Ricochet executes http calls
type Ricochet struct {
	baseURL *url.URL
	token   string
}

// Get performs http GET request
func (r *Ricochet) Get(url string) *Response {
	url = combineURL(r.baseURL, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Get:", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// todo: deal with error
		// this should fail the test and short-circuit it
		log.Println("Get:", err)
	}

	return &Response{
		resp,
	}
}

func combineURL(base *url.URL, relative string) string {
	if base == nil {
		return relative
	}

	if relParsed, err := url.Parse(relative); err == nil && relParsed.IsAbs() {
		return relative
	}

	copy := *base

	copy.Path = path.Join(copy.Path, relative)

	return copy.String()
}
