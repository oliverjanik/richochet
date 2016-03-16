package ricochet

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var suites = make(map[string]*Suite)

// TestFunc is signature for tests
type TestFunc func(r *Ricochet)

// Suite contains multiple tests
type Suite struct {
	name  string
	tests map[string]TestFunc
}

// NewSuite creates new test suite
func NewSuite(name string) *Suite {
	s := &Suite{
		name:  name,
		tests: make(map[string]TestFunc),
	}
	suites[name] = s
	return s
}

// OAuth sets up credential
func (s *Suite) OAuth(endpoint, client, secret, username, password string) *Suite {
	params := url.Values{}

	params.Add("grant_type", "password")
	params.Add("client_id", client)
	params.Add("client_secret", secret)
	params.Add("username", username)
	params.Add("password", password)

	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		fmt.Println("OAuth error:", err)
		return nil
	}

	d, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(d))

	return s
}

// Test defines a test in a suit
func (s *Suite) Test(name string, test TestFunc) *Suite {
	s.tests[name] = test
	return s
}

// Run test suit
func (s *Suite) Run() {
	fmt.Println("Running", s.name)
	for n, t := range s.tests {
		fmt.Println("\t", "...", n)
		t(&Ricochet{})
	}
}
