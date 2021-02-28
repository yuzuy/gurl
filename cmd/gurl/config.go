package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type configFile map[string]config

type config struct {
	Header map[string]string `json:"header"`
}

func newConfig() config {
	return config{
		Header: make(map[string]string),
	}
}

var (
	configDirPath  = os.Getenv("HOME") + "/.config/gurl"
	configFilePath = configDirPath + "/config.json"
)

func getConfigFile() (configFile, error) {
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return configFile{}, err
	}
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return configFile{}, err
	}
	defer f.Close()

	cfBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return configFile{}, err
	}
	if len(cfBytes) == 0 {
		cf := make(configFile)
		return cf, nil
	}
	var cf configFile
	err = json.Unmarshal(cfBytes, &cf)
	return cf, err
}

func saveConfigFile(cf configFile) error {
	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	cfBytes, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(cfBytes)
	return err
}

func printDefaultHeader(host string) error {
	header, err := getDefaultHeader(host)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(host + ":\n")
	for k, v := range header {
		buf.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}
	log.Print(buf.String())
	return nil
}

func getDefaultHeader(pattern string) (map[string]string, error) {
	cf, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	conf, ok := cf[pattern]
	if !ok {
		return map[string]string{}, nil
	}
	return conf.Header, nil
}

func setDefaultHeader(header, pattern string) error {
	cf, err := getConfigFile()
	if err != nil {
		return err
	}

	if _, ok := cf[pattern]; !ok {
		cf[pattern] = newConfig()
	}

	key, val, err := parseHeader(header)
	if err != nil {
		return err
	}
	cf[pattern].Header[key] = val

	return saveConfigFile(cf)
}

func deleteDefaultHeader(key, pattern string) error {
	cf, err := getConfigFile()
	if err != nil {
		return err
	}

	conf, ok := cf[pattern]
	if !ok {
		return nil
	}
	if _, ok := conf.Header[key]; !ok {
		return nil
	}
	delete(conf.Header, key)

	return saveConfigFile(cf)
}
