package main

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request){
	var payload = struct{
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status: "active",
		Message: "Go movies up and running",
		Version: "1.0.0",
	}

	out, err := json.Marshal(payload)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request){
	var movies[] models.Movie

	rd, _:= time.Parse("2006-01-02", "1986-03-07")

	priotoma := models.Movie{
		ID: 1,
		Title: "Priotoma",
		ReleaseDate: rd,
		MPAARating: "R",
		RunTime: 116,
		Description: "A amazing movie",
		CreatedAT: time.Now(),
		UpdatedAt: time.Now(),
	}
	movies = append(movies, priotoma)

	rd, _ = time.Parse("2006-01-02", "1981-06-12")

	Avengers := models.Movie{
		ID: 2,
		Title: "Avengers",
		ReleaseDate: rd,
		MPAARating: "PG-13",
		RunTime: 115,
		Description: "Another amazing movie",
		CreatedAT: time.Now(),
		UpdatedAt: time.Now(),
	}

	movies = append(movies, Avengers)

	out, err := json.Marshal(movies)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}