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

func newPattern(v string) pattern {
	return pattern(v)
}

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

type defaultHeaderList map[pattern]defaultHeader

func newDefaultHeader() defaultHeader {
	return make(defaultHeader)
}

func (dhs defaultHeaders) save() error {
	f, err := os.Create(defaultHeadersFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	dhlBytes, err := json.MarshalIndent(dhl, "", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(dhlBytes)
	return err
}

func (dhl defaultHeaderList) add(p pattern, h string) error {
	if !p.isValid() {
		return errors.New("invalid pattern")
	}

	if _, ok := dhl[p]; !ok {
		dhl[p] = newDefaultHeader()
	}

	key, val, err := parseHeader(h)
	if err != nil {
		return err
	}
	dhl[p][key] = val

	return nil
}

func (dhl defaultHeaderList) remove(p pattern, k string) {
	delete(dhl[p], k)
}

var (
	configDirPath          = os.Getenv("HOME") + "/.config/gurl"
	defaultHeadersFilePath = configDirPath + "/default_headers.json"
)

func getDefaultHeaderList() (defaultHeaderList, error) {
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(defaultHeadersFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dhlBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(dhlBytes) == 0 {
		return make(defaultHeaderList), nil
	}
	var dhl defaultHeaderList
	err = json.Unmarshal(dhlBytes, &dhl)
	return dhl, err
}

func printDefaultHeaders() error {
	dhl, err := getDefaultHeaderList()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for p, dh := range dhl {
		buf.WriteString(string(p) + ":\n")
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
	dhl, err := getDefaultHeaderList()
	if err != nil {
		return nil, err
	}

	return dhl[p], nil
}

func addDefaultHeader(p pattern, h string) error {
	dhl, err := getDefaultHeaderList()
	if err != nil {
		return err
	}

	if err := dhl.add(p, h); err != nil {
		return err
	}

	return dhl.save()
}

func removeDefaultHeader(p pattern, key string) error {
	dhl, err := getDefaultHeaderList()
	if err != nil {
		return err
	}

	dhl.remove(p, key)

	return dhl.save()
}
