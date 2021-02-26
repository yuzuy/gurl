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
	case "config":
		host := flag.Arg(1)
		switch flag.Arg(2) {
		case "set":
			switch flag.Arg(3) {
			case "header":
				header := flag.Arg(4)
				err = setDefaultHeader(header, host)
			}
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
