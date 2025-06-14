package app

import (
	"flag"
	"log"

	"github.com/bsagat/envzilla"
)

func init() {
	if err := envzilla.Loader(".env"); err != nil {
		log.Fatal("Failed to load environment variables: ", err)
	}

	flag.Parse()
}
