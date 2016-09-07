package RC_Oauth

/*
Auth API:
Auth: Takes in the secrets and callback uri and returns an auth object

Create Auth Url: Uses the Auth object and returns an authorisation url

Get Token: makes a call to the access token endpoint with the given url

Recurser API:

Batch API:

*/


import (
	"golang.org/x/oauth2"
	"math/rand"
)

var	oauthStateString = getStateString(10)

var RCOauthConfig *oauth2.Config;

func MakeConfig(url, id, secret string) {
	RCOauthConfig = &oauth2.Config{
		RedirectURL:   url,
		ClientID:     id,
		ClientSecret: secret,
		Scopes:       []string{"public"},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  "https://recurse.com/oauth/authorize",
			TokenURL: "https://recurse.com/oauth/token",
		},
	}
	
}


func GetUrl() string {
	url := RCOauthConfig.AuthCodeURL(oauthStateString)
	return url
}

func IsStateString(state string) bool {
	return state == oauthStateString
}

func GetToken(code string) *oauth2.Token {
	token, err := RCOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(token)
	}
	return token
}


func getStateString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randString)
}

