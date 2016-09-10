package RC-API

import (
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

///////////////////////////////////////////////////////
// Config:
///////////////////////////////////////////////////////

type RCConfig struct {
	*oauth2.Config
	StateString string
}

func MakeConfig(url, id, secret string) *RCConfig {
	c := &oauth2.Config{
		Scopes: []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://recurse.com/oauth/authorize",
			TokenURL: "https://recurse.com/oauth/token",
		},
		RedirectURL:  url,
		ClientID:     id,
		ClientSecret: secret,
	}
	rcConfig := &RCConfig{
		c,
		getStateString(20),
	}
	return rcConfig
}

// Generates a random 20 char string, as per the protocol
func getStateString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randString)
}

// Generates the URL on which use user can give consent for the app to use their RC data
func (c *RCConfig) GetUrl() string {
	url := c.AuthCodeURL(c.StateString)
	return url
}

func (c *RCConfig) MakeAuth(code string) RCAuth {
	token, err := c.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(token)
	}
	t := RCAuth{token,
		"https://www.recurse.com/api/v1/",
		"people/",
		"batches/",
		genAccessParam(token.AccessToken)}
	return t
}

///////////////////////////////////////////////////////
// Auth:
///////////////////////////////////////////////////////

type RCAuth struct {
	*oauth2.Token
	BaseUrl      string
	RecurserPath string
	BatchPath    string
	TokenParam   string
}

func genAccessParam(token string) string {
	param := "?access_token=" + token
	return param
}

///////////////////////////////////////////////////////
// Request Utilities
///////////////////////////////////////////////////////

// Makes a request and returns the result
// Should probably be JSON instead of a string
func makeRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyS := string(body)
	return bodyS
}

// Checks if the state string passed back matches the one the user sent
func (c *RCConfig) IsStateString(state string) bool {
	return state == c.StateString
}

///////////////////////////////////////////////////////
// Recurser
///////////////////////////////////////////////////////

// Get the details of the authed Recurser
func (t *RCAuth) Me() string {
	me := t.Recurser("me")
	return me
}

// Get any given Recurser
// Takes ether a user ID or an email
func (t *RCAuth) Recurser(id string) string {
	url := t.BaseUrl + t.RecurserPath + id + t.TokenParam
	res := makeRequest(url)
	return res
}

///////////////////////////////////////////////////////
// Batch
///////////////////////////////////////////////////////

// Returns a list of every batch
func (t *RCAuth) GetBatchList() string {
	url := t.BaseUrl + t.BatchPath + t.TokenParam
	res := makeRequest(url)
	return res
}

// Returns a particular batch
// Takes a batch ID
func (t *RCAuth) GetBatch(id string) string {
	url := t.BaseUrl + t.BatchPath + id + t.TokenParam
	res := makeRequest(url)
	return res
}

// Returns the details of every member of a batch
// Takes a batch ID
func (t *RCAuth) GetBatchMembers(id string) string {
	url := t.BaseUrl + t.BatchPath + id + "/people" + t.TokenParam
	res := makeRequest(url)
	return res
}
