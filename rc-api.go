package rc_api

import (
	"golang.org/x/oauth2"
	"log"
	"math/rand"
)

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

func getStateString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randString)
}

func (c *RCConfig) IsStateString(state string) bool {
	return state == c.StateString
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
