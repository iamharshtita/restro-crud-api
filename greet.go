package main

import (
	"fmt"
	"net/http"
)

func greet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the PPBET Cafe!")
	}
}
