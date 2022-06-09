package main

import (
	"net/http"

	"github.com/dshurubtsov/cmd/config"
)

func Routes(app *config.Application) *http.ServeMux {

	// main routes for app
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home(app))
	mux.HandleFunc("/auth/signup", SignUp(app))
	mux.HandleFunc("/auth/signin", SignIn(app))

	return mux
}
