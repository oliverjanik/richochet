package ricochet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Response contains assertion methods
type Response struct {
	*http.Response
}

// Any represents any json value
type any interface{}

// JSONNode represents json response body
type JSONNode struct {
	any
}

// IsSuccess tests for successful status code between 200 and 300 exclusive
func (r *Response) IsSuccess() *Response {
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		fmt.Println("\t\tError:", "Expected Success, was", r.StatusCode)
	}

	return r
}

// Is tests for specific status code
func (r *Response) Is(code int) *Response {
	if r.StatusCode != code {
		fmt.Println("\t\tError:", "Expected status", code, "was", r.StatusCode)
	}

	return r
}

// JSON reads response as json and wraps it in JSONNode for assertions
func (r *Response) JSON() *JSONNode {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var result interface{}

	err := dec.Decode(&result)
	if err != nil {
		fmt.Println("\t\tError:", "Could not decode body to json")
	}

	return &JSONNode{result}
}

// Dump pretty prints the parsed Json
func (n *JSONNode) Dump() {
	spew.Dump(n.any)
}

// IsArray checks if response is an array
func (n *JSONNode) IsArray() *JSONNode {
	if _, ok := n.any.([]interface{}); !ok {
		fmt.Println("\t\tError:", "Expected Array")
	}

	return n
}

// Path traverses the json and returns a descendent node
func (n *JSONNode) Path(parts ...interface{}) (result *JSONNode) {
	result = &JSONNode{nil}
	r := n.any

	for _, part := range parts {
		switch v := part.(type) {
		case int:
			if array, ok := r.([]interface{}); ok {
				r = array[v]
			} else {
				return
			}
		case string:
			if obj, ok := r.(map[string]interface{}); ok {
				r = obj[v]
			} else {
				return
			}
		default:
			panic("Unknown type")
		}
	}

	result.any = r
	return
}

// Equals asserts equality
func (n *JSONNode) Equals(val interface{}) (r *JSONNode) {
	r = n

	switch expected := val.(type) {
	case float64:
		if actual, ok := n.any.(float64); ok && actual == float64(expected) {
			return
		}
	case int:
		if actual, ok := n.any.(float64); ok && actual == float64(expected) {
			return
		}
	case string:
		if actual, ok := n.any.(string); ok && actual == expected {
			return
		}
	case bool:
		if actual, ok := n.any.(bool); ok && actual == expected {
			return
		}
	case nil:
		if val == nil {
			return
		}
	}

	spew.Printf("\t\tError: Expected %#+v but found %#+v\n", val, n.any)
	return
}

// func path(path ...interface{]}) (interface{}, err) {

// }
