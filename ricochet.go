package ricochet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/oliverjanik/ricochet/j"
)

// R executes http calls and provides asserts
type R struct {
	baseURL *url.URL
	header  http.Header
}

// Get performs http GET request
func (r *R) Get(url string) *Response {
	return r.send("GET", url, nil)
}

// Post posts data to specified url
func (r *R) Post(url string, data interface{}) *Response {
	return r.send("POST", url, data)
}

// Put performs a Put request
func (r *R) Put(url string, data interface{}) *Response {
	return r.send("PUT", url, data)
}

// Delete sends delete request to given url
func (r *R) Delete(url string) *Response {
	return r.send("DELETE", url, nil)
}

func (r *R) send(method string, url string, data interface{}) *Response {
	url = combineURL(r.baseURL, url)

	var body io.Reader
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			r.Fail("Error serializing data:", err)
		}

		body = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		r.Fail("Error preparing request:", err)
	}

	req.Header.Set("Accept", "application/json")

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// apply headers from parent
	for k, v := range r.header {
		req.Header[k] = v
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.Fail(method, url, "failed:", err)
	}

	return &Response{
		resp,
	}
}

// Fail stops exection of current test or suite
func (r *R) Fail(a ...interface{}) {
	msg := fmt.Sprintln(a)
	panic(msg)
}

func combineURL(base *url.URL, relative string) string {
	if base == nil {
		return relative
	}

	return base.String() + "/" + relative
}

// AssertSuccess tests for successful status code between 200 and 300 exclusive
func (r *R) AssertSuccess(resp *Response) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		r.Fail("Expected Success, was", resp.StatusCode)
	}
}

// AssertStatus tests for specific status code
func (r *R) AssertStatus(resp *Response, code int) {
	if resp.StatusCode != code {
		r.Fail("Expected status", code, "was", resp.StatusCode)
	}
}

// AssertNotEmpty checks if the array is not empty
func (r *R) AssertNotEmpty(a *j.Array) {
	if a == nil {
		r.Fail("Array is nil")
	}

	if a.Len() == 0 {
		r.Fail("Array is empty")
	}
}

// AssertUndefinedOrNull checks for missing property or null
func (r *R) AssertUndefinedOrNull(n *j.Node) {
	if n == nil {
		return
	}

	if n.Raw() == nil {
		return
	}

	r.Fail(spew.Sprintf("Expected undefiend or nil, found %#+v", n.Raw()))
}

// AssertEquals asserts equality
func (r *R) AssertEquals(n *j.Node, val interface{}) {
	if n == nil {
		r.Fail("Node undefined")
	}

	if n.Equals(val) {
		return
	}

	r.Fail(spew.Sprintf("Expected %#+v but found %#+v", val, n))
}

// AssertNotNil test if value is not nil
func (r *R) AssertNotNil(val interface{}) {
	if val == nil {
		r.Fail("Unexpected nil")
	}
}
