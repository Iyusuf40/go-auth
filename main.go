package main

import (
	"fmt"
)

func main() {
	var x = map[string]string{"a": "abc"}

	y := c(x)

	y["a"] = "xyz"

	fmt.Printf("%p %p", &x, &y)
	fmt.Println(x, y)
}

func c(x map[string]string) map[string]string {
	return x
}
