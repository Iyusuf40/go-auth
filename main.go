package main

import (
	"fmt"
	"os"

	"github.com/Iyusuf40/go-auth/storage"
)

func main() {
	m, err := storage.MakeMongoWrapper("test", "test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Println(m)
}
