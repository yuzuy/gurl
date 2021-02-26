package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
)

var (
	hostFlag = flag.String("host", "default", "Using when you want to link info set by set command with the host")
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
		host := *hostFlag
		switch flag.Arg(1) {
		case "header":
			err = setDefaultHeader(flag.Arg(2), host)
		}
	}
	return err
}

func setDefaultHeader(header, host string) error {
	cf, err := getConfigFile()
	if err != nil {
		return err
	}

	if _, ok := cf.HostToConfig[host]; !ok {
		cf.HostToConfig[host] = config{
			Header: make(map[string]string),
		}
	}

	tmp := strings.Split(header, ":")
	if len(tmp) != 2 {
		return errors.New("invalid header format")
	}
	key := tmp[0]
	val := strings.TrimPrefix(tmp[1], " ")
	cf.HostToConfig[host].Header[key] = val

	return saveConfigFile(cf)
}
