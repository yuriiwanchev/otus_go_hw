package main

import (
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "success" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
