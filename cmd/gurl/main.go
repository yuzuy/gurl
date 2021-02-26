package main

import (
	"flag"
	"fmt"
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
	default:
		var respBody string
		respBody, err = doHTTPRequest(flag.Arg(0))
		if err == nil {
			fmt.Print(respBody)
		}
	}
	return err
}
