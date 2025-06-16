package main

import (
	"GonIO/internal/app"
	"GonIO/internal/domain"
	"log"
	"net/http"
)

func main() {
	h := app.SetHandler()

	log.Printf("Starting listening on port %s", domain.Port)
	if err := http.ListenAndServe(domain.URLDomain+":"+domain.Port, h); err != nil {
		log.Fatalf("Server listening error: %s", err.Error())
	}
}
