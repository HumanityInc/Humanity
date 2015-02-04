package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	APP_CONFIG_NAME = "/../conf/web.conf"
)

type Config struct {
	Process struct {
		Daemon bool   `json:"daemon"`
		ChRoot string `json:"chroot"`
	} `json:"process"`

	Network struct {
		Bind     string `json:"bind"`
		Protocol string `json:"protocol"`
	} `json:"network"`

	Storage struct {
		Postgresql string   `json:"postgresql"`
		Memcache   []string `json:"memcache"`
	} `json:"storage"`

	OAuth struct {
		Facebook struct {
			AppId       string `json:"app_id"`
			AppSecret   string `json:"app_secret"`
			RedirectUrl string `json:"redirect_url"`
		} `json:"facebook"`

		Twitter struct {
			ConsumerKey    string `json:"consumer_key"`
			ConsumerSecret string `json:"consumer_secret"`
			RedirectUrl    string `json:"redirect_url"`
		} `json:"twitter"`

		Google struct {
			ClientId     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			RedirectUrl  string `json:"redirect_url"`
		} `json:"google"`
	} `json:"oauth"`

	Sendmail struct {
		Mandrill struct {
			ApiKey string `json:"api_key"`
		} `json:"mandrill"`
	} `json:"sendmail"`
}

var app_conf *Config

func GetConfig() Config {

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

	return *app_conf
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
