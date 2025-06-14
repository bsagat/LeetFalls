package app

import (
	"fmt"
	"leetFalls/internal/domain"
	"log"
	"os"
	"strconv"
)

func CheckFlags() {
	conf := domain.Config

	if *conf.HelpFlag {
		PrintHelp()
	}

	port, err := strconv.Atoi(*conf.Port)
	if err != nil {
		log.Fatal("Failed to parse port number: ", domain.ErrInvalidPort)
	}

	if port < 1024 || port > 65000 {
		log.Fatal("Failed to check flags: ", domain.ErrPortRange)
	}
}

// Prints Help message and Exit program with 0 code
func PrintHelp() {
	t := ` ./LeetFalls --help
LeetFalls

Usage:
  LeetFalls [--port <N>]  
  LeetFalls --help

Options:
  --help       Show this screen.
  --port N     Port number.
	`
	fmt.Println(t)
	os.Exit(0)
}
