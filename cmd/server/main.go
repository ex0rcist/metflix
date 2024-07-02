package main

import (
	"github.com/ex0rcist/metflix/internal/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		panic(err)
	}

	err2 := srv.Run()
	if err2 != nil {
		panic(err2)
	}
}
