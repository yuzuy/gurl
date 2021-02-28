package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	hFlag = &headerFlag{}

	dFlag = flag.String("d", "", "Input the request body")
	vFlag = flag.Bool("v", false, "Output the verbose log")
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
	req, err := makeHTTPRequest(uri, cf)
	if err != nil {
		return "", err
	}
	if *vFlag {
		logRequest(req)
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
	if *vFlag {
		log.Println()
		logResponseExceptForBody(resp)
	}
	return string(body), nil
}

func makeHTTPRequest(uri *url.URL, cf configFile) (*http.Request, error) {
	bodyStr := *dFlag
	var body io.Reader = nil
	if bodyStr != "" {
		body = bytes.NewReader([]byte(bodyStr))
	}
	method := *xFlag
	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	defaultHeader, err := makeDefaultHeader(uri, cf)
	if err != nil {
		return nil, err
	}
	req.Header = defaultHeader

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

func makeDefaultHeader(uri *url.URL, cf configFile) (http.Header, error) {
	header := make(http.Header)
	for pattern, conf := range cf {
		matchHost, err := path.Match(pattern, uri.Host)
		if err != nil {
			return nil, err
		}
		matchPattern, err := path.Match(pattern, uri.Host+uri.Path)
		if err != nil {
			return nil, err
		}
		if !matchHost && !matchPattern {
			continue
		}

		for k, v := range conf.Header {
			header.Set(k, v)
		}
	}

	return header, nil
}

func parseHeader(v string) (key, val string, err error) {
	tmp := strings.Split(v, ":")
	key = tmp[0]
	val = strings.Join(tmp[1:], ":")
	return
}

func logRequest(req *http.Request) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s %s %s\n", req.Method, req.URL.Path, req.Proto))
	for k, v := range req.Header {
		vs := strings.Join(v, ",")
		buf.WriteString(fmt.Sprintf("%s: %s\n", k, vs))
	}

	body := *dFlag
	if body != "" {
		buf.WriteString("\n" + body + "\n")
	}

	log.Println(buf.String())
}

func logResponseExceptForBody(resp *http.Response) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status))
	for k, v := range resp.Header {
		vs := strings.Join(v, ",")
		buf.WriteString(fmt.Sprintf("%s: %s\n", k, vs))
	}

	log.Println(buf.String())
}
