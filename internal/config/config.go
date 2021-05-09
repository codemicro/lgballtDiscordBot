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
	Token                      string           `json:"token"`
	Prefix                     string           `json:"prefix"`
	DbFileName                 string           `json:"dbFileName"`
	DebugMode                  bool             `json:"debug"`
	Statuses                   []string         `json:"statuses"`
	AdminRole                  string           `json:"adminRole"`
	VerificationIDs            VerificationIds  `json:"verificationIds"`
	RedditFeeds                []RedditFeedInfo `json:"redditFeeds"`
	ChatChartChannelExclusions []string         `json:"ccExclusions"`
	OwnerId                    string           `json:"ownerId"`
	PkApi                      PkApiInfo        `json:"pkApi"`
	Listeners                  ListenerInfo     `json:"listeners"`
	MuteMe                     MuteMeInfo       `json:"muteMe"`
	BioFields                  []string         `json:"bioFields"`
	ActionLogChannel           string           `json:"actionLogChannel"`
}

type ListenerInfo struct {
	RoleId          string   `json:"roleId"`
	AllowedChannels []string `json:"allowedChannels"`
}

type VerificationIds struct {
	InputChannel         string   `json:"inputChannel"`
	OutputChannel        string   `json:"outputChannel"`
	RoleId               string   `json:"assignRoleId"`
	ModlogChannel        string   `json:"modlogChannel"`
	ExcludedPronounRoles []string `json:"excludedPronounRoles"`
	ExtraValidRoles      []string `json:"extraValidRoles"`
}

type RedditFeedInfo struct {
	Webhook  string `json:"webhook"`
	RssUrl   string `json:"rssUrl"`
	Interval int    `json:"interval"`
	IconUrl  string `json:"iconUrl"`
}

type PkApiInfo struct {
	ContactEmail    string `json:"contactEmail"`
	ApiUrl          string `json:"apiUrl"`
	MinRequestDelay int    `json:"minReqDelay"`
	NumWorkers      int    `json:"numWorkers"`
}

type MuteMeInfo struct {
	TimeoutRole   string   `json:"timeoutRole"`
	RolesToRemove []string `json:"rolesToRemove"`
}

var (
	Token                      string
	Prefix                     string
	DbFileName                 string
	Statuses                   []string
	DebugMode                  bool
	AdminRole                  string
	VerificationIDs            VerificationIds
	RedditFeeds                []RedditFeedInfo
	ChatChartChannelExclusions []string
	OwnerId                    string
	PkApi                      PkApiInfo
	Listeners                  ListenerInfo
	MuteMe                     MuteMeInfo
	BioFields                  []string
	ActionLogChannel           string
)

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
	RedditFeeds = cfg.RedditFeeds
	ChatChartChannelExclusions = cfg.ChatChartChannelExclusions
	OwnerId = cfg.OwnerId
	PkApi = cfg.PkApi
	Listeners = cfg.Listeners
	MuteMe = cfg.MuteMe
	BioFields = cfg.BioFields
	ActionLogChannel = cfg.ActionLogChannel

	if PkApi.ApiUrl[len(PkApi.ApiUrl)-1] == '/' {
		PkApi.ApiUrl = PkApi.ApiUrl[:len(PkApi.ApiUrl)-1]
	}

	if DbFileName == "" {
		DbFileName = "lgballtBot.db"
	}
}
