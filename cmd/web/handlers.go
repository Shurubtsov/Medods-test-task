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
		}

		// decode body request for find Username and Pass
		user := models.User{}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}
		json.Unmarshal(data, &user)

		encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))

		id, err := app.UserModel.CreateUser(user.Username, encodedPassword)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, fmt.Sprintln("id: ", id))
	}
}

// Endpoint for create couple tokens for user with ID from request Query
func GetTokensForUser(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-type", "application/json")
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		// find user in our database with id from request
		user, err := app.UserModel.FindById(id)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		tokens, err := createTokens(app, user)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		//fmt.Println("[TEST TOKEN]access_token is: ", test)

		// bcrypt token for storage it in database
		bcryptedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), 14)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// add refresh token in db
		err = app.UserModel.UpdateUserToken(id, string(bcryptedRefreshToken))
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// encode Refresh Token with base64 for response
		tokens.RefreshToken = base64.StdEncoding.EncodeToString([]byte(tokens.RefreshToken))

		// create body of response
		resp, err := json.Marshal(tokens)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.Write(resp)
	}
}

// Endpoint for Refresh access token for user. Need Json struct in body kind {access_token: "", refresh_token: ""}
func Refresh(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set POST method if req does not exist this
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-type", "application/json")
		}
		tokens := models.Token{}

		// decode body json request
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}
		json.Unmarshal(data, &tokens)

		decodedBase64RefreshToken, err := base64.StdEncoding.DecodeString(tokens.RefreshToken)
		if err != nil {
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

		tokens, err = createTokens(app, user)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		//fmt.Println("[HANDLER REFRESH]Tokens -> ", tokens)

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

		// create body of response
		resp, err := json.Marshal(tokens)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.Write(resp)
	}
}

// Create couple tokens with binding each other
func createTokens(app *config.Application, user models.User) (models.Token, error) {

	tokens := models.Token{}

	//fmt.Println("[TOKENS]user: ", user)
	// encode JWT with payload of sub from ID user for his identification in service
	accessToken, err := app.TokenManager.NewJWT(user.ID.Hex())
	if err != nil {
		app.ErrorLog.Fatal("can't create new tokens")
		return tokens, err
	}
	//fmt.Println("[TOKENS]access: ", accessToken)

	// encoded with salt of access token for his binding to refresh token
	refreshToken, err := app.TokenManager.NewRefreshToken(accessToken)
	if err != nil {
		app.ErrorLog.Fatalf("can't create refresh token, error: %v", err.Error())
		return tokens, err
	}

	// fill token.Models for response it in body
	tokens.AccessToken, tokens.RefreshToken = accessToken, refreshToken

	return tokens, nil
}
