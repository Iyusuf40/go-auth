package main

import (
	"github.com/Iyusuf40/go-auth/api"
	"github.com/Iyusuf40/go-auth/auth"
)

func main() {

	wait := make(chan int)
	go func() {
		api.ServeAPI()
	}()

	go func() {
		auth.ServeAUTH()
	}()
	<-wait
}
