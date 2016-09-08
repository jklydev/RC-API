package RC_Oauth

import (
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
)

///////////////////////////////////////////////////////
// AUTH:
///////////////////////////////////////////////////////

var	oauthStateString = getStateString(10)

var RCOauthConfig *oauth2.Config;

var RCOauthToken *oauth2.Token;

// Takes the applications details and generates the Config object
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

// Generates the URL on which use user can give consent for the app to use their RC data
func GetUrl() string {
	url := RCOauthConfig.AuthCodeURL(oauthStateString)
	return url
}

// Set the token for use internally
func SetToken(code string) {
	token, err := RCOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(token)
	}
	RCOauthToken = token
}

// Returns Token object for use in the app
func GetToken() *oauth2.Token {
	return RCOauthToken
}

// Returns access token string for use in the app
func GetAccessToken() string {
	return RCOauthToken.AccessToken
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

// Checks if the state string passed back matches the one the user sent
func IsStateString(state string) bool {
	return state == oauthStateString
}

// Checks that the token is non-nil and has not expired
func IsTokenValid() bool {
	return RCOauthToken.Valid()
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
    if ! IsStateString(state) {
        fmt.Printf("invalid oauth state")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
	}
	code := r.FormValue("code")
	SetToken(code)
	if ! IsTokenValid() {
		fmt.Printf("invalid token")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
	}
	http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
}


///////////////////////////////////////////////////////
// Request Utilities
///////////////////////////////////////////////////////

var baseUrl = "https://www.recurse.com/api/v1/"

// Generates the access token param to append to the end of a url
func getAccessToken() string {
	token := RCOauthToken.AccessToken
	param := "?access_token=" + token
	return param
}

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


///////////////////////////////////////////////////////
// Recurser
///////////////////////////////////////////////////////

var recurserUrl = baseUrl +  "people/"

// Get the details of the current Recurser
func GetMe() string {
	me := GetRecurser("me")
	return me
}

// Get any given Recurser
// Takes ether a user ID or an email
func GetRecurser(id string) string {
	url :=  recurserUrl + id + getAccessToken()
	res := makeRequest(url)
	return res
}


///////////////////////////////////////////////////////
// Batch
///////////////////////////////////////////////////////

var batchUrl = baseUrl + "batches/"

// Returns a list of every batch
func GetBatchList() string {
	url :=  batchUrl + getAccessToken()
	res := makeRequest(url)
	return res
}

// Returns a particular batch
// Takes a batch ID
func GetBatch(id string) string {
	url := batchUrl + id + getAccessToken()
	res := makeRequest(url)
	return res
}

// Returns the details of every member of a batch
// Takes a batch ID
func GetBatchMembers(id string) string {
	url := batchUrl + id + "/people" + getAccessToken()
	res := makeRequest(url)
	return res
}
