package main

import "github.com/ex0rcist/metflix/pkg/staticlint"

func main() {
	lint := staticlint.New()
	lint.Run()
}
