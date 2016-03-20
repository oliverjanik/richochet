package ricochet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var suites = make(map[string]*Suite)

// TestFunc is signature for tests
type TestFunc func(r *Ricochet)

// Suite contains multiple tests
type Suite struct {
	name    string
	tests   map[string]TestFunc
	baseURL *url.URL
	token   string
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

// BaseURL sets base URL for following operations
func (s *Suite) BaseURL(baseURL string) *Suite {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic("Error parsing base URL" + err.Error())
	}

	s.baseURL = u
	return s
}

type oauthResult struct {
	AccessToken string `json:"access_token"`
}

// OAuth sets up credential
func (s *Suite) OAuth(endpoint, client, secret, username, password string) *Suite {
	params := url.Values{}

	params.Add("grant_type", "password")
	params.Add("client_id", client)
	params.Add("client_secret", secret)
	params.Add("username", username)
	params.Add("password", password)

	endpoint = combineURL(s.baseURL, endpoint)
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		fmt.Println("OAuth error:", err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("OAuth did returned " + resp.Status)
	}

	d := json.NewDecoder(resp.Body)

	var msg oauthResult
	err = d.Decode(&msg)
	if err != nil {
		panic("Error decoding OAuth response " + err.Error())
	}

	s.token = msg.AccessToken

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
		t(&Ricochet{
			baseURL: s.baseURL,
			token:   s.token,
		})
	}
}
