package main

import (
	"fmt"
	"github.com/dvarelap/ig"
)

func main() {
	var (
		client     = ig.NewClient("<ACCESS_TOKEN>")
		token, err = client.LongLivedToken("<CLIENT_SECRET>")
	)

	checkIfErr(err)

	fmt.Println("This is my token:", token)
}
