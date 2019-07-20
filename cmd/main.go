package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/mvochoa/graphqldoc"
)

func main() {
	if len(os.Args) > 1 {
		u, err := url.ParseRequestURI(os.Args[1])
		if err == nil {
			fmt.Print("Generating documentation...\n")
			graphqldoc.HTTP(u.String())
			return
		}

		fmt.Println("ERROR: The endpoint is incorrect.")
	}
	fmt.Print("\n	Usage: graphdoc <ENDPOINT GRAPHQL>\n\n")
}
