package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type configFile struct {
	HostToConfig map[string]config
}

type config struct {
	Header map[string]string `json:"header"`
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
		cf := configFile{
			HostToConfig: make(map[string]config),
		}
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
