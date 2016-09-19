package rc_api

import (
	"golang.org/x/oauth2"
	"fmt"
)

type Config struct {
	*oauth2.Config
}

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

// Generates the URL on which use user can give consent for the app to use their RC data
func (c *Config) GetUrl(state string) string {
	url := c.AuthCodeURL(state)
	return url
}

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

type Auth struct {
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
