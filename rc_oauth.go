package RC_Oauth

import (
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	//"strings"
)

///////////////////////////////////////////////////////
// Config:
///////////////////////////////////////////////////////

var rcConfig = &oauth2.Config{
	Scopes: []string{"public"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://recurse.com/oauth/authorize",
		TokenURL: "https://recurse.com/oauth/token",
	},
}

var oauthStateString = getStateString(20)

// Generates a random 20 char string, as per the protocol
func getStateString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randString := make([]rune, n)
	for i := range randString {
		randString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randString)
}

// Takes the applications details and generates the Config object
func SetConfigVars(url, id, secret string) {
	rcConfig.RedirectURL = url
	rcConfig.ClientID = id
	rcConfig.ClientSecret = secret
}

// Generates the URL on which use user can give consent for the app to use their RC data
func GetUrl() string {
	url := rcConfig.AuthCodeURL(oauthStateString)
	return url
}

///////////////////////////////////////////////////////
// Redirect:
///////////////////////////////////////////////////////

//var postAuthRedirect = "/";
//var authObject = *RCAuth{};

// A default function to handle the auth redirect and set the token
//func HandleRedirect(w http.ResponseWriter, r *http.Request) {
//	state := r.FormValue("state")
//    if ! IsStateString(state) {
//        log.Fatal("invalid oauth state")
//        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
//        return
//	}
//	code := r.FormValue("code")
//	authObject.SetToken(code)
//	url := postAuthRedirect
//	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
//}

// Set the url you want the user redirected to after HandleRedirect
// Default is "/"
//func GenHandleRedirect (url string, auth *RCAuth) func(http.ResponseWriter, *http.Request) {
//	postAuthRedirect = url
//	authObject = auth
//	return HandleRedirect
//}

//func HandlerGen(url string) func(string, func(http.ResponseWriter, *http.Request)) {
//	postAuthRedirect = url
//	pattern := getPattern(rcConfig.RedirectURL)
//	return http.HandleFunc(pattern, HandleRedirect)
//}
//
//func getPattern(url string) string {
//	pattern := "/" + strings.SplitN(url, "/", 4)[3]
//	return pattern
//}

// Checks if the state string passed back matches the one the user sent
func IsStateString(state string) bool {
	return state == oauthStateString
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

func MakeAuth(code string) RCAuth {
	token, err := rcConfig.Exchange(oauth2.NoContext, code)
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

func (t *RCAuth) SetToken(code string) {
	token, err := rcConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(token)
	}
	t.AccessToken = token.AccessToken
	t.TokenType = token.TokenType
	t.RefreshToken = token.RefreshToken
	t.Expiry = token.Expiry
	t.TokenParam = genAccessParam(token.AccessToken)
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
