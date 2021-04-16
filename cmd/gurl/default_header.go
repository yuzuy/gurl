package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type pattern string

func (p pattern) match(uri *url.URL) (bool, error) {
	if uri.Host == string(p) {
		return true, nil
	}

	re, err := p.regexp()
	if err != nil {
		return false, err
	}
	return re.MatchString(uri.Host + uri.Path), nil
}

func (p pattern) isValid() bool {
	_, err := p.regexp()
	return err == nil
}

func (p pattern) regexp() (*regexp.Regexp, error) {
	s := string(p)
	s = strings.ReplaceAll(s, "/", "\\/")
	s = strings.ReplaceAll(s, "*", ".*")
	return regexp.Compile("^" + s + "$")
}

type defaultHeader map[string]string

type defaultHeaders map[pattern]defaultHeader

func newDefaultHeader() defaultHeader {
	return make(defaultHeader)
}

func (dhs defaultHeaders) save() error {
	f, err := os.Create(defaultHeadersFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	dhsBytes, err := json.MarshalIndent(dhs, "", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(dhsBytes)
	return err
}

func (dhs defaultHeaders) set(p pattern, h string) error {
	if !p.isValid() {
		return errors.New("invalid pattern")
	}

	if _, ok := dhs[p]; !ok {
		dhs[p] = newDefaultHeader()
	}

	key, val := parseHeader(h)
	dhs[p][key] = val

	return nil
}

func (dhs defaultHeaders) remove(p pattern, k string) {
	delete(dhs[p], k)
}

var (
	configDirPath          = os.Getenv("HOME") + "/.config/gurl"
	defaultHeadersFilePath = configDirPath + "/default_headers.json"
)

func getDefaultHeaders() (defaultHeaders, error) {
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(defaultHeadersFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dhsBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(dhsBytes) == 0 {
		return make(defaultHeaders), nil
	}
	var dhs defaultHeaders
	err = json.Unmarshal(dhsBytes, &dhs)
	return dhs, err
}

func printDefaultHeaders() error {
	dhs, err := getDefaultHeaders()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for p, dh := range dhs {
		buf.WriteString(string(p) + "\n")
		for k, v := range dh {
			buf.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}
	log.Print(buf.String())
	return nil
}

func printDefaultHeader(p pattern) error {
	dh, err := getDefaultHeader(p)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(string(p + ":\n"))
	for k, v := range dh {
		buf.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}
	log.Print(buf.String())
	return nil
}

func getDefaultHeader(p pattern) (defaultHeader, error) {
	dhs, err := getDefaultHeaders()
	if err != nil {
		return nil, err
	}

	return dhs[p], nil
}

func setDefaultHeader(p pattern, h string) error {
	dhs, err := getDefaultHeaders()
	if err != nil {
		return err
	}

	if err := dhs.set(p, h); err != nil {
		return err
	}

	return dhs.save()
}

func deleteDefaultHeader(p pattern, key string) error {
	dhs, err := getDefaultHeaders()
	if err != nil {
		return err
	}

	dhs.remove(p, key)

	return dhs.save()
}
