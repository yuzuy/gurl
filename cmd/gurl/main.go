package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Println("gurl: " + err.Error())
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
		case "header":
			switch flag.Arg(3) {
			case "get":
				err = printDefaultHeader(host)
			case "set":
				header := flag.Arg(4)
				err = setDefaultHeader(header, host)
			case "delete":
				key := flag.Arg(4)
				err = deleteDefaultHeader(key, host)
			}
		}
	default:
		var respBody string
		respBody, err = doHTTPRequest(flag.Arg(0))
		if err == nil {
			fmt.Println(respBody)
		}
	}
	return err
}

func init() {
	log.SetFlags(0)
}
