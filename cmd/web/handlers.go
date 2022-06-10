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
		tokens := models.Token{}

		// create accesss token for user
		accessToken, err := app.TokenManager.NewJWT("id")
		if err != nil {
			app.ErrorLog.Fatalf("can't create token, error: %v", err.Error())
			app.ServerError(w, err)
			return
		}
		// create refresh token for user
		refreshToken, err := app.TokenManager.NewRefreshToken()
		if err != nil {
			app.ErrorLog.Fatalf("can't create refresh token, error: %v", err.Error())
			return
		}

		// encode Refresh Token with base64
		base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

		// fill token.Models for response it in body
		tokens.AccessToken, tokens.RefreshToken = accessToken, base64RefreshToken

		resp, err := json.Marshal(tokens)
		if err != nil {
			app.ErrorLog.Fatalf("can't Marshal %v, error: %v", tokens, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
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

		//jwt := app.JWTMaker

		// user, err := jwt.ValidateRefreshToken(token)
		// if err != nil {
		// 	fmt.Fprint(w, "invalid refresh token")
		// 	app.ClientError(w, http.StatusUnauthorized)
		// 	return
		// }

		// token, err = jwt.CreateToken(user)
		// if err != nil {
		// 	app.ClientError(w, http.StatusUnauthorized)
		// 	return
		// }

		// resp, err := json.Marshal(token)
		// if err != nil {
		// 	app.ServerError(w, err)
		// 	return
		// }

		// w.Header().Add("Content-type", "application/json")
		// w.Write(resp)
	}
}
