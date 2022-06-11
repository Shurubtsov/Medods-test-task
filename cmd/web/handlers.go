package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dshurubtsov/cmd/config"
	"github.com/dshurubtsov/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// Test endpoint for check server errors
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

// Endpoint for registration users to database
func SignUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set POST method if req does not exist this
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Add("Content-type", "application/json")
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		// read body request for find Username and Pass
		user := models.User{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}
		json.Unmarshal(b, &user)

		encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))
		//fmt.Println("[CREATE-USER] encoded password: ", encodedPassword)

		id, err := app.UserModel.CreateUser(user.Username, encodedPassword)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		fmt.Fprint(w, fmt.Sprint("id: ", strings.Trim(id, `ObjectID(")`)))
	}
}

func GetTokensForUser(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-type", "application/json")
			app.ClientError(w, http.StatusMethodNotAllowed)
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		user, err := app.UserModel.FindById(id)
		if err != nil {
			app.ErrorLog.Fatalf("%v, can't find user", err.Error())
		}

		tokens, err := createTokens(*app, user)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// bcrypt token for storage it in database
		bcryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), 14)
		if err != nil {
			app.ErrorLog.Fatal("can't bcrypt token")
			return
		}

		// add refresh token in db
		app.UserModel.UpdateUserToken(id, string(bcryptedRefreshToken))

		// encode Refresh Token with base64 for response
		tokens.RefreshToken = base64.StdEncoding.EncodeToString([]byte(tokens.RefreshToken))

		// create body of response
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
		tokens.RefreshToken = string(decodedBase64RefreshToken)

		// find user by refresh token
		user, err := app.UserModel.FindByRefreshToken(string(decodedBase64RefreshToken))
		if err != nil {
			app.NotFound(w)
			return
		}

		// check for valid token
		ok, err := app.TokenManager.ValidateRefreshToken(tokens)
		if !ok || err != nil {
			app.ClientError(w, http.StatusBadRequest)
			fmt.Fprint(w, "invalid token")
			return
		}

		tokens, err = createTokens(*app, user)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// bcrypt token for storage it in database
		bcryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), 14)
		if err != nil {
			app.ErrorLog.Fatal("can't bcrypt token")
			return
		}

		// encode Refresh Token with base64 for response
		tokens.RefreshToken = base64.StdEncoding.EncodeToString([]byte(tokens.RefreshToken))

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

func createTokens(app config.Application, user models.User) (models.Token, error) {

	tokens := models.Token{}

	accessToken, err := app.TokenManager.NewJWT(user.Username)
	if err != nil {
		app.ErrorLog.Fatal("can't create new tokens")
		return tokens, err
	}

	refreshToken, err := app.TokenManager.NewRefreshToken(accessToken)
	if err != nil {
		app.ErrorLog.Fatalf("can't create refresh token, error: %v", err.Error())
		return tokens, err
	}

	// fill token.Models for response it in body
	tokens.AccessToken, tokens.RefreshToken = accessToken, refreshToken

	return tokens, nil
}
