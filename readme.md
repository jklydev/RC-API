RC API Wrapper
=================

<a href='http://www.recurse.com' title='Made with love at the Recurse Center'><img src='https://cloud.githubusercontent.com/assets/2883345/11325206/336ea5f4-9150-11e5-9e90-d86ad31993d8.png' height='20px'/></a>

[![Go Report Card](https://goreportcard.com/badge/github.com/JKiely/RC-API)](https://goreportcard.com/report/github.com/JKiely/RC-API)

A wrapper for the [Recurse Center API](https://github.com/recursecenter/wiki/wiki/Recurse-Center-API) written in Go. It can handle authenticating with OAuth2 and making requests.

[Here's](https://gist.github.com/JKiely/267083e9fa776eb7d35c38fb8447e57c) an example app using it.


Getting Started
---------------

### Step 0:
[Apply](https://www.recurse.com/apply/retreat) and get accepted to the [Recurse Center](https://www.recurse.com/). Done? Excellent.

### Step 1:
Go to your [RC Settings page](https://www.recurse.com/settings/oauth), look in OAuth applications, and create a new application. You'll need to give it a name and redirect url (it can just be http://localhost:3000/recurse (with whatever port/extension you want) while you're testing.

### Step 2:
This will give you back an ID and a Secret, call `rc_api.MakeConfig` with these, and your redirect url in order to get a config object:

```Go
import (
    //...
	rc_api "github.com/JKiely/RC-API"
    //...
)


var config = rc_api.MakeConfig(
		"http://localhost:3000/recurse",
		"ID",
		"Secret")
```

### Step 3:
In your app, make sure you have a route that redirects the user to the url provided by you config object (`config.GetUrl()`). This sends them to a page on the RC website which asks them to log in (if they're not already), and then asks if they give your app permission to access their account.

### Step 4:
After Step 3 recurse.com will redirect the user back to the redirect url you provided, so make sure you have a handler for that as well. It will also pass you a code as a param in the url, your handler will need to take this code and use it to create an auth object by calling `config.MakeAuth(code)` with your config object. This package will then use that code to get an authorization token for your application and wrap it in an auth object.

It might look like this:
```Go
func handleRCRedirect(w http.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    auth = config.MakeAuth(code)

    //...

    // Whatever you want to do next
}
```


Using the API
-------------

### Config

The first thing you need to do is get a config object:

```Go
var config = rc_api.MakeConfig(
		"http://localhost:3000/recurse",
		"ID",
		"Secret")
```

This wraps the [oauth2 config object](https://godoc.org/golang.org/x/oauth2#Config).

#### Get URL
```Go
config.GetUrl("state")
```

This gives you the url which you need to redirect users to in order to authenticate. You need to pass it a state string that it will give back to you after the redirect. This is to project against [Cross-Site Request Forgery](https://tools.ietf.org/html/rfc6749#section-10.12).

#### Make Auth
```Go
var auth = config.MakeAuth(code)
```

Takes a code and hits RC's token endpoint in order to get a token object.

Auth objects wrap the [oauth2 token object](https://godoc.org/golang.org/x/oauth2#Token) and adds methods for hitting the RC api.

### Recursers

#### Get yourself
```Go
auth.Me()
```
Returns you! Or at least whoever is logged in. A struct containing all the details pertaining to you.

It contains the following fields:
```Go
type Recurser struct {
	Id                 int
	First_name         string
	Middle_name        string
	Last_name          string
	Email              string
	Twitter            string
	Github             string
	Batch_id           int
	Phone_number       string
	Has_photo          bool
	Interests          string
	Before_rc          string
	During_rc          string
	Is_faculty         bool
	Is_hacker_schooler bool
	Job                string
	Image              string
	Batch              Batch
	Pseudonym          string
	Current_location   Location
	Stints             []Stint
	Projects           []string
	Links              []string
	Skills             []string
	Bio                string
}

type Stint struct {
	Id         int
	Start_date string
	End_date   string
	Type       string
}

type Location struct {
	Geoname_id int
	Name       string
	Short_name string
	Ascii_name string

```
Though the final four are [depreciated](https://github.com/recursecenter/wiki/wiki/Recurse-Center-API#people) and will be `nil` for all more recent Recursers.

#### Get any Recurser
```Go
auth.Recurser("id")
```
Gets a Recurser who you specify by their id (given as a string) or their email address.

(The `Me` function above is really just this one, but passing in "me" as the id.)

### Batches

#### Get a batch
```Go
auth.Batch("id")
```
Gets the given batch, as specified by the batch ID. Contains the following fields:
```Go
type Batch struct {
	Id         int
	Name       string
	Start_date string
	End_date   string
}
```

#### Get Recursers from batch
```Go
auth.BatchMembers("id")
```
Takes a batch id as a string and returns a list of every Recurser in a given batch.

#### Get every batch
```Go
auth.BatchList()
```
Returns a list of every batch.
