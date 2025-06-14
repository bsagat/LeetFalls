package domain

import "flag"

var Config = struct {
	Port     *string
	HelpFlag *bool
}{
	flag.String("port", "8080", "Default port number"),
	flag.Bool("help", false, "Shows help message"),
}
