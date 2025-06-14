package main

import (
	"leetFalls/internal/app"
)

func main() {
	srv, cleanup := app.Setup()
	defer cleanup()

	app.StartServer(srv)

	app.WaitForShutDown(srv)
}
