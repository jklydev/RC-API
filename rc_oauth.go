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
	"os"
)

var (
	RCOauthConfig = &oauth2.Config{
		RedirectURL:    "urn:ietf:wg:oauth:2.0:oob",
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Scopes:       []string{"public"},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  "https://recurse.com/oauth/authorize",
			TokenURL: "https://recurse.com/oauth/token",
		},
	}
	oauthStateString = getRandString(10)
)

func GetUrl() string {
	url := RCOauthConfig.AuthCodeURL(oauthStateString)
	return url
}

//func GetToken() string {	
//}


func getRandString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randString)
}

