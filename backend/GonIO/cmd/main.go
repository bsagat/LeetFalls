package main

import (
	"GonIO/internal/app"
	"GonIO/internal/domain"
	"fmt"
	"log"
	"net/http"
)

func main() {
	h := app.SetHandler()

	url := fmt.Sprintf("%s:%s", domain.Host, domain.Port)
	log.Printf("Starting listening on %s", url)
	if err := http.ListenAndServe(url, h); err != nil {
		log.Fatal("Server listening error: ", err)
	}
}
