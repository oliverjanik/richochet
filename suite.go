package ricochet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// TestFunc is signature for tests
type TestFunc func(r *R)

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
	groups  []*TestGroup
	baseURL *url.URL
	token   string
	oauth   *oauth
}

type oauth struct {
	endpoint     string
	clientID     string
	clientSecret string
	username     string
	password     string
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
	s.oauth = &oauth{
		endpoint:     endpoint,
		clientID:     client,
		clientSecret: secret,
		username:     username,
		password:     password,
	}

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

func (s *Suite) authenticate() {
	if s.oauth == nil {
		return
	}

	params := url.Values{}

	params.Add("grant_type", "password")
	params.Add("client_id", s.oauth.clientID)
	params.Add("client_secret", s.oauth.clientSecret)
	params.Add("username", s.oauth.username)
	params.Add("password", s.oauth.password)

	endpoint := combineURL(s.baseURL, s.oauth.endpoint)
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		fmt.Println("OAuth error:", err)
		return
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
}
