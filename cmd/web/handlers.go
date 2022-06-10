package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dshurubtsov/cmd/config"
	"github.com/dshurubtsov/pkg/models"
	"golang.org/x/crypto/bcrypt"
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
			app.ErrorLog.Fatalf("something went wrong when server attempt to create user, error: %v", err.Error())
			app.ServerError(w, err)
			return
		}

		fmt.Fprint(w, id)
	}
}

func Login(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-type", "application/json")
			app.ClientError(w, http.StatusMethodNotAllowed)
		}

		id := r.URL.Query().Get("id")
		tokens := models.Token{}
		user, err := app.UserModel.FindById(id)
		if err != nil {
			app.ErrorLog.Fatalf("%v, can't find user", err.Error())
		}

		// create accesss token for user
		accessToken, err := app.TokenManager.NewJWT(user.Username)
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

		// bcrypt token for storage it in database
		bcryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 14)
		if err != nil {
			app.ErrorLog.Fatal("can't bcrypt token")
			return
		}

		app.UserModel.UpdateUserToken(id, string(bcryptedRefreshToken))

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
		tokens := models.Token{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.ServerError(w, err)
		}
		json.Unmarshal(body, &tokens)

		decodedBase64RefreshToken, err := base64.StdEncoding.DecodeString(tokens.RefreshToken)
		if err != nil {
			app.ErrorLog.Fatal("can't decode refresh token")
			app.ServerError(w, err)
			return
		}

		// find user by refresh token
		user, err := app.UserModel.FindByRefreshToken(string(decodedBase64RefreshToken))
		if err != nil {
			app.NotFound(w)
			return
		}
		//fmt.Println("[HANDLER:137]User after find ref token is ", user)

		// update tokenы if user finded
		accessToken, err := app.TokenManager.NewJWT(user.Username)
		if err != nil {
			app.ErrorLog.Fatal("can't create new tokens")
			return
		}

		refreshToken, err := app.TokenManager.NewRefreshToken()
		if err != nil {
			app.ErrorLog.Fatalf("can't create refresh token, error: %v", err.Error())
			return
		}

		// encode Refresh Token with base64
		base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))

		// fill token.Models for response it in body
		tokens.AccessToken, tokens.RefreshToken = accessToken, base64RefreshToken

		// bcrypt token for storage it in database
		bcryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 14)
		if err != nil {
			app.ErrorLog.Fatal("can't bcrypt token")
			return
		}

		//fmt.Println("[HANDLER:165]user ID: ", user.ID)
		app.UserModel.UpdateUserToken(user.ID.Hex(), string(bcryptedRefreshToken))

		resp, err := json.Marshal(tokens)
		if err != nil {
			app.ErrorLog.Fatalf("can't Marshal %v, error: %v", tokens, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resp)
	}
}
