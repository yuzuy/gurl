package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	hFlag = &headerFlag{}

	dFlag = flag.String("d", "", "Input the request body")
	xFlag = flag.String("X", "GET", "Input the http method")
)

func init() {
	flag.Var(hFlag, "H", "Input the request header")
}

type headerFlag []string

func (h *headerFlag) String() string {
	return "http header"
}

func (h *headerFlag) Set(v string) error {
	*h = append(*h, v)
	return nil
}

func doHTTPRequest(urlStr string) (respBody string, err error) {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("could not parse the url: %s", urlStr)
	}
	cf, err := getConfigFile()
	if err != nil {
		return "", err
	}
	conf, ok := cf.HostToConfig[uri.Host]
	if !ok {
		conf = config{}
	}

	req, err := makeHTTPRequest(urlStr, conf)
	if err != nil {
		return "", err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("could not read the response body")
	}
	return string(body), nil
}

func makeHTTPRequest(urlStr string, conf config) (*http.Request, error) {
	bodyStr := *dFlag
	var body io.Reader = nil
	if bodyStr != "" {
		body = bytes.NewReader([]byte(bodyStr))
	}
	method := *xFlag
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header = conf.header()
	header := *hFlag
	for _, v := range header {
		key, val, err := parseHeader(v)
		if err != nil {
			return nil, err
		}
		req.Header.Set(key, val)
	}

	return req, nil
}

func parseHeader(v string) (key, val string, err error) {
	tmp := strings.Split(v, ":")
	key = tmp[0]
	val = strings.Join(tmp[1:], ":")
	return
}
