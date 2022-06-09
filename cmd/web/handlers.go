package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dshurubtsov/cmd/config"
	"github.com/dshurubtsov/pkg/models"
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

		id, err := app.UserModel.CreateUser(username, encodedPassword)
		if err != nil {
			app.ServerError(w, err)
		}

		fmt.Fprint(w, id)
	}
}

func Login(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		user, err := app.UserModel.FindById(id)
		if err != nil {
			app.NotFound(w)
			return
		}

		token, err := app.JWTMaker.CreateToken(user)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		resp, err := json.Marshal(token)
		if err != nil {
			app.ServerError(w, err)
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(resp)
	}
}

func Refresh(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set POST method if req does not exist this
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-type", "application/json")
			app.ClientError(w, http.StatusMethodNotAllowed)
		}
		token := models.Token{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.ServerError(w, err)
		}
		json.Unmarshal(body, &token)

		jwt := app.JWTMaker

		user, err := jwt.ValidateRefreshToken(token)
		if err != nil {
			fmt.Fprint(w, "invalid refresh token")
			app.ClientError(w, http.StatusUnauthorized)
			return
		}

		token, err = jwt.CreateToken(user)
		if err != nil {
			app.ClientError(w, http.StatusUnauthorized)
			return
		}

		resp, err := json.Marshal(token)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.Write(resp)
	}
}
