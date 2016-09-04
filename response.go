package ricochet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliverjanik/ricochet/j"
)

// Response contains assertion methods
type Response struct {
	*http.Response
}

// AsJSON reads response as json and wraps it in j.Node for assertions
func (r *Response) AsJSON() *j.Node {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var result interface{}

	err := dec.Decode(&result)
	if err != nil {
		fmt.Println("\t\tError:", "Could not decode body to json")
	}

	return j.New(result)
}
