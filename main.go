// Package main is the entry point of the application for the mortgage calculation service.
//
// It initializes the application by calling the Run function from the app package,
// which starts the server, loads configurations, and handles incoming HTTP requests.
//
// The main package is responsible for launching the application and is the entry point when the program is executed.
package main

import (
	"sber/internal/app"
)

func main() {
	// Calls the Run function from the app package to start the application, which includes server setup and execution.
	app.Run()
}
