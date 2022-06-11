package main

import (
	"net/http"

	"github.com/dshurubtsov/cmd/config"
)

func Routes(app *config.Application) *http.ServeMux {

	// main routes for app
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home(app))
	mux.HandleFunc("/auth/sign-up", SignUp(app))

	// Main routes for test task
	mux.HandleFunc("/auth/get/tokens", GetTokensForUser(app))
	mux.HandleFunc("/auth/refresh", Refresh(app))

	return mux
}
