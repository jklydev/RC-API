RC API Wrapper
=================

<a href='http://www.recurse.com' title='Made with love at the Recurse Center'><img src='https://cloud.githubusercontent.com/assets/2883345/11325206/336ea5f4-9150-11e5-9e90-d86ad31993d8.png' height='20px'/></a>

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


var rcConfig = rc_api.MakeConfig(
		"http://localhost:3000/recurse",
		"ID",
		"Secret")
```

### Step 3:
In your app, make sure you have a route that redirects the user to the url provided by you config object (`rcConfig.GetUrl()`). This sends them to a page on the RC website which asks them to log in (if they're not already), and then asks if they give your app permission to access their account.

### Step 4:
After Step 3 recurse.com will redirect the user back to the redirect url you provided, so make sure you have a handler for that as well. It will also pass you a code as a param in the url, your handler will need to take this code and use it to create an auth object by calling `rcConfig.MakeAuth(code)` with your config object. This package will then use that code to get an authorization token for your application and wrap it in an auth object.

It might look like this:
```Go
func handleRCRedirect(w http.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    rcAuth = rcConfig.MakeAuth(code)

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
config.GetUrl()
```

This gives you the url which you need to redirect users to in order to authenticate. 

#### Make Auth
```Go
var auth = config.MakeAuth(code)
```

Takes a code and hits RC's token endpoint in order to get a token object.

Auth objects wrap the [oauth2 token object](https://godoc.org/golang.org/x/oauth2#Token) and adds methods for hitting the RC api.

### Recursers

#### Get yourself
```Go
auth.GetMe()
```
Returns you! Or at least whoever is logged in. A json object containing all the details pertaining to you.

When I do it it looks like this:
```
{"id":1786,
"first_name":"John",
"middle_name":"",
"last_name":"Kiely",
"email":"jjkiely@gmail.com",
"twitter":null,
"github":"JKiely",
"batch_id":28,
"phone_number":"13476367122",
"has_photo":true,
"interests":null,
"before_rc":null,
"during_rc":null,
"is_faculty":false,
"is_hacker_schooler":true,
"job":null,
"image":"https://d29xw0ra2h4o4u.cloudfront.net/assets/people/john_kiely_150-fe5be2ceba10783d19ad16dd2511c2a118a66258926a367a2a6a40811b8b4729.jpg",
"batch":{"id":28,"name":"Summer 2, 2016","start_date":"2016-07-05","end_date":"2016-09-22"},
"pseudonym":"Pvc Drop",
"current_location":null,
"stints":[{"batch_id":28,"start_date":"2016-07-05","end_date":"2016-09-22"}],
"projects":[],"links":[],"skills":[],"bio":null}
```

#### Get any Recurser
```Go
auth.GetRecurser("id")
```
Gets a Recurser who you specify by their id (given as a string) or their email address.

(The `GetMe` function above is really just this one, but passing in "me" as the id.)

### Batches

#### Get a batch
```Go
auth.GetBatch("id")
```
Gets the given batch as a, as specified by the batch ID for example if called with "28" it will return:
```
{"id":28,
"name":"Summer 2, 2016",
"start_date":"2016-07-05",
"end_date":"2016-09-22"}
```

#### Get Recursers from batch
```Go
auth.GetBatchMembers("id")
```
Takes a batch id as a string and returns a list of every Recurser in a given batch.

#### Get every batch
```Go
auth.GetBatchList()
```
Returns a list of every batch.
