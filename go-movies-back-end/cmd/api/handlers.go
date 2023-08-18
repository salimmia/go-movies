package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	_ = app.WriteJSON(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request){
	// read json payload
	var requestPayload struct{
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil{
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)

	log.Println(user)

	if err != nil{
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return;
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid{
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil{
		app.ErrorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.WriteJSON(w, http.StatusAccepted, tokens)
}