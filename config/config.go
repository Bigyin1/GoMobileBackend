package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

var ConfPath string

type Config struct {
	Port           int    `json:"port"`
	StoragePath    string `json:"storagePath"`
	FileURIPrefix  string `json:"fileUriPrefix"`
	Debug          bool   `json:"debugMode"`
	GmailTokenPath string `json:"gmailTokenPath"`
	GmailCredsPath string `json:"gmailCredsPath"`
	GmailAddr      string `json:"gmailAddr"`
	UploadSubject  string `json:"uploadSubject"`
	MailTmplPath   string `json:"mappingTmplPath"`
	PollingPeriod  int    `json:"pollingPeriodSeconds"`
}

// AsString represents config as string
func (conf *Config) AsString() string {
	data, _ := json.Marshal(conf)
	return string(data)
}

func LoadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func LoadEnvironment() {
	cloudMobile, exists := os.LookupEnv("CLOUD_GO_MOBILE")
	cloud, _ := strconv.ParseBool(cloudMobile)

	if exists && cloud { // TODO add cloud config

	} else {
		ConfPath = "local-config.json" // local config path
	}
}
