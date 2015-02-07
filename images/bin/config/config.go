package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	APP_CONFIG_NAME = "/../conf/image.conf"
)

type Config struct {
	Bind     string `json:"bind"`
	Protocol string `json:"protocol"`
}

var app_conf *Config

func GetConfig() *Config {

	if app_conf == nil {

		app_conf = new(Config)
		cwd, _ := os.Getwd()

		if _, err := os.Stat(cwd + APP_CONFIG_NAME); err == nil {

			app_conf.load(cwd + APP_CONFIG_NAME)

		} else {

			base_dir := os.Args[0]
			os.Chdir(filepath.Dir(base_dir))
			cwd, _ = os.Getwd()
			app_conf.load(cwd + APP_CONFIG_NAME)
		}
	}

	return app_conf
}

func (conf *Config) load(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
