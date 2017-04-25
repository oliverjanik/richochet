package ricochet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TestFunc is signature for tests
type TestFunc func(r *R)

// SuiteFunc is signature for setup/teardown
type SuiteFunc func(s *Suite)

// TestGroup represents tests the will run sequentially in parallel to other groups
type TestGroup struct {
	name   string
	tests  []test
	failed bool
	indent string
}

// Suite contains multiple tests
type Suite struct {
	TestGroup
	groups   []*TestGroup
	baseURL  *url.URL
	header   http.Header
	setUp    SuiteFunc
	tearDown SuiteFunc
}

type test struct {
	name string
	f    TestFunc
}

// NewSuite creates new test suite
func NewSuite(name string) *Suite {
	return &Suite{
		TestGroup: TestGroup{
			name: name,
		},
	}
}

// NewGroup creates a group of tests
func NewGroup(name string) *TestGroup {
	return &TestGroup{
		name:   name,
		indent: "\t",
	}
}

// SetUp records a setup funcion
func (s *Suite) SetUp(f SuiteFunc) *Suite {
	s.setUp = f
	return s
}

// TearDown records a teardown function
func (s *Suite) TearDown(f SuiteFunc) *Suite {
	s.tearDown = f
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

// SetHeader sets header for every future request
func (s *Suite) SetHeader(key string, value string) *Suite {
	if s.header == nil {
		s.header = make(http.Header)
	}

	s.header.Set(key, value)

	return s
}

// CreateR creates an instace of R
func (s *Suite) CreateR() *R {
	return &R{
		baseURL: s.baseURL,
		header:  s.header,
	}
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

	r, _ := http.NewRequest("POST", endpoint, strings.NewReader(params.Encode()))
	for k, v := range s.header {
		r.Header[k] = v
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
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

	s.SetHeader("Authorization", "Bearer "+msg.AccessToken)

	return s
}

// Test defines a test in a suit
func (s *Suite) Test(name string, testFunc TestFunc) *Suite {
	s.tests = append(s.tests, test{name, testFunc})
	return s
}

// Group adds a group of tests to the suite
func (s *Suite) Group(group *TestGroup) *Suite {
	s.groups = append(s.groups, group)
	return s
}

// Test adds a test to a group
func (g *TestGroup) Test(name string, testFunc TestFunc) *TestGroup {
	g.tests = append(g.tests, test{name, testFunc})
	return g
}
