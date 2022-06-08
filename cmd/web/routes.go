package main

import (
	"net/http"

	"github.com/dshurubtsov/cmd/config"
)

func Routes(app *config.Application) *http.ServeMux {

	// main routes for app
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home(app))

	return mux
}
