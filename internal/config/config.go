package config

import (
	"encoding/json"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"io/ioutil"
	"os"
)

const configFileName = "botConfig.json"

type Info struct {
	Token    string   `json:"token"`
	Prefix   string   `json:"prefix"`
	DbFileName string `json:"dbFileName"`
	Statuses []string `json:"statuses"`
}

var Config Info

func init() {
	configFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Failed to read %s", configFileName))
		os.Exit(1)
	}

	err = json.Unmarshal(configFileBytes, &Config)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Failed to parse %s", configFileName))
		os.Exit(1)
	}

	if Config.DbFileName == "" {
		Config.DbFileName = "lgballtBot.db"
	}
}
