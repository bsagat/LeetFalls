package domain

import "flag"

var Config = struct {
	Port          *string
	HelpFlag      *bool
	DatabaseDsn   string
	TemplatesPath string
}{
	Port:     flag.String("port", "8080", "Default port number"),
	HelpFlag: flag.Bool("help", false, "Shows help message"),
}

type Code int
