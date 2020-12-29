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
	Token      string   `json:"token"`
	Prefix     string   `json:"prefix"`
	DbFileName string   `json:"dbFileName"`
	DebugMode  bool     `json:"debug"`
	Statuses   []string `json:"statuses"`
	AdminRole  string   `json:"adminRole"`
	VerificationIDs VerificationIds `json:"verificationIds"`
}

type VerificationIds struct {
	InputChannel string `json:"inputChannel"`
	OutputChannel string `json:"outputChannel"`
	RoleId string `json:"assignRoleId"`
}

// var Config Info

var Token string
var Prefix string
var DbFileName string
var Statuses []string
var DebugMode bool
var AdminRole string
var VerificationIDs VerificationIds

func init() {
	configFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Failed to read %s", configFileName))
		os.Exit(1)
	}

	var cfg Info

	err = json.Unmarshal(configFileBytes, &cfg)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Failed to parse %s", configFileName))
		os.Exit(1)
	}

	Token = cfg.Token
	Prefix = cfg.Prefix
	DbFileName = cfg.DbFileName
	Statuses = cfg.Statuses
	DebugMode = cfg.DebugMode
	AdminRole = cfg.AdminRole
	VerificationIDs = cfg.VerificationIDs

	if DbFileName == "" {
		DbFileName = "lgballtBot.db"
	}
}
