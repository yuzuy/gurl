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
	var err error
	switch flag.Arg(0) {
	case "set":
		switch flag.Arg(1) {
		case "header":
			err = setDefaultHeader(flag.Arg(2))
		}
	}
	return err
}
