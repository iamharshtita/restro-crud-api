package main

import (
	"fmt"
	"net/http"
	"os"
)

func login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		password := r.Header.Get("Password")

		if username != os.Getenv("adminUsername") || password != os.Getenv("adminPassword") {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Credentials do not match")
			return
		}

		token, err := GenerateToken(username)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Faled to generate JWT")
			return
		}

		//return the jwt
		fmt.Fprintln(w, token)
	}
}
