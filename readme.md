RC API Wrapper
=================

A wrapper for the [Recurse Center API](https://github.com/recursecenter/wiki/wiki/Recurse-Center-API) written in Go. It can handle authenticating with OAuth2 and making requests.

[Here's](https://gist.github.com/JKiely/267083e9fa776eb7d35c38fb8447e57c) an example app using it.


Getting Started
---------------

### Step 0:
[Apply](https://www.recurse.com/apply/retreat) and get accepted to the [Recurse Center](https://www.recurse.com/). Done? Excellent!

### Step 1:
Go to your [RC Settings page](https://www.recurse.com/settings/oauth), look in OAuth applications, and create a new application. You'll need to give it a name and redirect url (it can just be http://localhost:3000/recurse (with whatever port/extension you want) while you're testing.

### Step 2:
This will give you back an ID and a Secret, take these (and you redirect url) and call `rc_api.MakeConfig` on them, like this:

```Go
MakeConfig(
		"http://localhost:3000/recurse",
		"ID",
		"Secret",
	)
```
You really don't want to expose these so it's probably best to save them as environment variables and use `os.Getenv("varName")` to retrieve them.

### Step 3:
In your app, make sure you have a route that redirects the user to the url provided by `rc_api.GetUrl()`. This sends them to a page on the RC website which asks them to log in (if they're not already), and then asks if they give your app permission to access their account.

### Step 4:
After Step 3 recurse.com will redirect the user back to the redirect url you provided, so make sure you have a handler for that as well. It will also pass you a code as a param in the url, your handler will need to take this code and pass it to `rc_api.SetToken(code)`. This package will then use that code to authorize your application and get an authorization token that it will use when you make API requests.

It might look like this:
```Go
func handleRCRedirect(w http.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    rc_api.SetToken(code)

    ...

    // Whatever you want to do next
}
```

However, if you don't want to do anything in particular with the token, you really just want to set it inside the app and make api calls with it, then the package provides a default redirect handler that you can use, that sets the token and redirects the user again back the page of your choice. To use this you need to call `rc_api.SetPostAuthRedirect("url")` (the default is `"/"`), and pass `rc_api.HandleRedirect` to your handler function:

```Go
func main() {
    ...
	m.SetPostAuthRedirect("/me")
    http.HandleFunc("/RCCallback", m.HandleRedirect)
    ...
}
```

Using the API
-------------

### Auth


### Recursers

#### Get yourself
```Go
rc_api.GetMe()
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
rc_api.GetRecurser("id")
```
Gets a Recurser who you specify by their id (given as a string) or their email address.

(The `GetMe` function above is really just this one, but passing in "me" as the id.)

### Batches

#### Get a batch
```Go
rc_api.GetBatch("id")
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
rc_api.GetBatchMembers("id")
```
Takes a batch id as a string and returns a list of every Recurser in a given batch.

#### Get every batch
```Go
rc_api.GetBatchList()
```
Returns a list of every batch.
