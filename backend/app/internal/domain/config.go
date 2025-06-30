package domain

import "flag"

var Config = struct {
	Port             *string
	HelpFlag         *bool
	StorageHost      string // s3 storage host
	StoragePort      string // s3 storage port
	GravityFallsHost string // external API host
	GravityFallsPort string // external API port
	DatabaseDsn      string // Database data sourse name
	TemplatesPath    string // frontend templates path
}{
	Port:     flag.String("port", "8080", "Default port number"),
	HelpFlag: flag.Bool("help", false, "Shows help message"),
}
