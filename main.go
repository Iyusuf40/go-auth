package main

import "fmt"

func main() {
	var x any = map[string]any{"a": 1}

	if obj, ok := x.(map[string]any); ok {
		fmt.Println(obj["a"], "again")
	}
}
