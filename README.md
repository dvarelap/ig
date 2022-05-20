[![Build Status](https://travis-ci.org/dvarelap/ig.svg?branch=master)](https://travis-ci.org/dvarelap/ig)

# Instagram Basic Display API in Golang

This is a Go package that fully supports the [Instagram Basic Display API.](https://developers.facebook.com/docs/instagram-basic-display-api/)

Feel free to create an issue or send me a pull request if you have any "how-to" question or bug or suggestion when using this package. I'll try my best to reply it.

## Install

Install this package with `go get github.com/dvarelap/ig.git`

## Usage

### Get profile from a ig user.

```go
package main

import (
	"fmt"
	"github.com/dvarelap/ig.git"
)

func main() {
	var (
		client       = ig.NewClient("<ACCESS_TOKEN>")
		profile, err = client.GetProfile()
	)
        
    	checkIfErr(err)

	fmt.Println("Here's my profile:", profile)
}
``` 

The type of **profile** is `ig.Profile` struct, which has the following structure

```go
type Profile struct {
	ID          string
	Username    string
	AccountType string
	MediaCount  string
}
```

and represents fields on [User's fields](https://developers.facebook.com/docs/instagram-basic-display-api/reference/user#fields)

### Get user's media.
```go
package main

import (
	"fmt"
	"github.com/dvarelap/ig.git"
)

func main() {
	var (
		client       = ig.NewClient("<ACCESS_TOKEN>")
		media, err = client.GetMedia()
	)

	checkIfErr(err)

	fmt.Println("User's media:")
	for _, entry := range media {
		fmt.Println(entry)
	}
}
```

**Media** is a slice of type `ig.Entry` struct,  which has the following structure
```go
type Entry struct {
	ID           string
	Username     string
	Caption      string
	MediaType    string
	MediaURL     string
	Permalink    string
	ThumbnailURL string
	Timestamp    string
}
```

and represents fields on [Media](https://developers.facebook.com/docs/instagram-basic-display-api/reference/media#fields)


### Exchange access token for a long lived token.

```go
package main

import (
	"fmt"
	"github.com/dvarelap/ig.git"
)

func main() {
	var (
		client     = ig.NewClient("<ACCESS_TOKEN>")
		token, err = client.LongLivedToken("<CLIENT_SECRET>")
	)

	checkIfErr(err)

	fmt.Println("This is my token:", token)
}
```

The type of **token** is `ig.LongLiveToken` struct, which has the following structure

```go
type LongLiveToken struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64 
}
```

it represents Long Lived Token from described [here](https://developers.facebook.com/docs/instagram-basic-display-api/guides/long-lived-access-tokens).

## Licence
This package is licensed under the MIT license. See LICENSE for details.
