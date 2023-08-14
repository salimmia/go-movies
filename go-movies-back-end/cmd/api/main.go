package main

import (
	"fmt"
	"log"
	"net/http"
)

var port = 8080

type application struct{
	Domain string
}

func main(){
	// set application config
	var app application

	app.Domain = "example.com"

	// read from command line

	// connect to database

	// start a web server

	log.Println("Starting application on port", port)

	http.HandleFunc("/", app.Home)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil{
		log.Fatal(err)
		return
	}
}