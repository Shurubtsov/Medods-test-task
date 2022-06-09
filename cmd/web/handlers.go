package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dshurubtsov/cmd/config"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for correct input with right path
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		fmt.Fprint(w, "home page")
	}
}

func SignUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set POST method if req does not exist this
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

		id, err := app.UserModel.Insert(username, encodedPassword)
		if err != nil {
			app.ServerError(w, err)
		}
		fmt.Fprint(w, id)
	}
}

func SignIn(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

		res, err := app.UserModel.CreateUser(username, encodedPassword)
		if err != nil {
			app.ClientError(w, http.StatusUnauthorized)
		}

		fmt.Fprintf(w, "%s", res)
	}
}

func FindUser(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		res, err := app.UserModel.FindById(id)
		if err != nil {
			app.NotFound(w)
		}

		fmt.Fprintf(w, "%s", res)
	}
}

func Login(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		user, err := app.UserModel.FindById(id)
		if err != nil {
			app.NotFound(w)
		}

		resp, err := json.Marshal(user)
		if err != nil {
			app.ServerError(w, err)
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(resp)
	}
}
