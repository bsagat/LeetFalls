package app

import (
	"errors"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/bsagat/envzilla"
)

func init() {
	if err := envzilla.Loader(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("Env file is not found ")
		} else {
			log.Fatal("Failed to load environment variables: ", err)
		}
	}

	flag.Parse()

	SetConfigs()
}
