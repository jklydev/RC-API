package RC_Oauth

import (
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"io/ioutil"
	"log"
)

///////////////////////////////////////////////////////
// AUTH:
///////////////////////////////////////////////////////

var	oauthStateString = getStateString(10)

var RCOauthConfig *oauth2.Config;

var RCOauthToken *oauth2.Token;

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

func SetToken(code string) {
	RCOauthToken = GetToken(code)
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

func IsStateString(state string) bool {
	return state == oauthStateString
}

func IsTokenValid() bool {
	return RCOauthToken.Valid()
}


///////////////////////////////////////////////////////
// Request Utilities
///////////////////////////////////////////////////////

var baseUrl = "https://www.recurse.com/api/v1/"

func getAccessToken() string {
	token := RCOauthToken.AccessToken
	param := "?access_token=" + token
	return param
}

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


///////////////////////////////////////////////////////
// Recurser
///////////////////////////////////////////////////////

var recurserUrl = baseUrl +  "people/"

func GetMe() string {
	me := GetRecurser("me")
	return me
}

func GetRecurser(id string) string {
	url :=  recurserUrl + id + getAccessToken()
	res := makeRequest(url)
	return res
}


///////////////////////////////////////////////////////
// Batch
///////////////////////////////////////////////////////

var batchUrl = baseUrl + "batches/"

func GetBatchList() string {
	url :=  batchUrl + getAccessToken()
	res := makeRequest(url)
	return res
}

func GetBatch(id string) string {
	url := batchUrl + id + getAccessToken()
	res := makeRequest(url)
	return res
}

func GetBatchMembers(id string) string {
	url := batchUrl + id + "/people" + getAccessToken()
	res := makeRequest(url)
	return res
}
