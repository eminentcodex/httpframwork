package main

import (
	"log"

	"httpframwork/app"
)

// main application entry point
func main() {
	var application app.Application
	var err error
	application, err = app.New()
	if err != nil {
		log.Fatal("failed to start the server: " + err.Error())
	}

	application.Run()
}
