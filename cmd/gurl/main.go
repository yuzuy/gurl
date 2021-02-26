package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()
	logger := log.New(os.Stderr, "gurl: ", 0)

	if err := run(); err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	return nil
}
