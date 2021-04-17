package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "0.2.0"

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
	case "dh":
		switch flag.Arg(1) {
		case "list":
			p := newPattern(flag.Arg(2))
			if p == "" {
				err = printDefaultHeaders()
			} else {
				err = printDefaultHeader(p)
			}
		case "add":
			p := newPattern(flag.Arg(2))
			header := flag.Arg(3)
			err = addDefaultHeader(p, header)
		case "rm":
			p := newPattern(flag.Arg(2))
			key := flag.Arg(3)
			err = removeDefaultHeader(p, key)
		}
	case "":
		flag.Usage()
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
