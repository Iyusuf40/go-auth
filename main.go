package main

import (
	"fmt"
	"time"
)

func main() {
	t1 := time.Now().Unix()
	t2 := time.Now().Add(time.Duration(10) * time.Second).Unix()
	fmt.Println(t2 - t1)
}
