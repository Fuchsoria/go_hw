package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	source := "Hello, OTUS!"
	reverse := stringutil.Reverse(source)

	fmt.Println(reverse)
}
