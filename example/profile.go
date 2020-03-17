package main

import (
	"fmt"
	"github.com/dvarelap/ig"
)

func main() {
	var (
		client     = ig.NewClient("<ACCESS_TOKEN>")
		media, err = client.GetMedia()
	)

	checkIfErr(err)

	fmt.Println("User's media:")
	for _, entry := range media {
		fmt.Println(entry)
	}
}
