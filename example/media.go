package main

import (
	"fmt"
	"github.com/dvarelap/ig"
)

func main() {
	var (
		client       = ig.NewClient("<ACCESS_TOKEN>")
		profile, err = client.GetProfile()
	)

	checkIfErr(err)

	fmt.Println("Here's my profile:", profile)
}
