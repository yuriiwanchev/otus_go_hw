package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	input := "Hello, OTUS!"
	reverseInput := reverse.String(input)
	fmt.Println(reverseInput)
}
