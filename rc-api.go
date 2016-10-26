package rcAPI

import (
	"fmt"

	"golang.org/x/oauth2"
)

// Config is a wrapper around the oauth2 config struct
type Config struct {
	*oauth2.Config
}

// MakeConfig generates the config object that is used to authenticate the user
func MakeConfig(url, id, secret string) *Config {
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
	rcConfig := &Config{
		c,
	}
	return rcConfig
}

// GetURL generates the URL on which use user can give consent for the app to use their RC data
func (c *Config) GetURL(state string) string {
	url := c.AuthCodeURL(state)
	return url
}

// MakeAuth generates the auth object that is used to make requests to the api
func (c *Config) MakeAuth(code string) Auth {
	token, err := c.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println(err)
	}
	t := Auth{token,
		"https://www.recurse.com/api/v1/",
		"people/",
		"batches/",
		genAccessParam(token.AccessToken)}
	return t
}

// Auth wraps the oauth2 token struct and adds aditional feilds to make requesting easier
type Auth struct {
	*oauth2.Token
	BaseURL      string
	RecurserPath string
	BatchPath    string
	TokenParam   string
}

func genAccessParam(token string) string {
	param := "?access_token=" + token
	return param
}
