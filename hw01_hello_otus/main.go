package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	baseString := "Hello, OTUS!"
	reverseString := stringutil.Reverse(baseString)
	fmt.Print(reverseString)
}
