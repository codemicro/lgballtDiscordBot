package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	PrometheusAddress          string           `json:"prometheusAddress"`
}

type ListenerInfo struct {
	RoleId          string   `json:"roleId"`
	AllowedChannels []string `json:"allowedChannels"`
}

type VerificationIds struct {
	InputChannel         string   `json:"inputChannel"`
	OutputChannel        string   `json:"outputChannel"`
	ArchiveChannel       string   `json:"archiveChannel"`
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
	PrometheusAddress          string
)

func init() {
	configFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s: %v", configFileName, err))
	}

	var cfg Info

	err = json.Unmarshal(configFileBytes, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse %s: %v", configFileName, err))
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
	PrometheusAddress = cfg.PrometheusAddress

	if PkApi.ApiUrl[len(PkApi.ApiUrl)-1] == '/' {
		PkApi.ApiUrl = PkApi.ApiUrl[:len(PkApi.ApiUrl)-1]
	}

	if DbFileName == "" {
		DbFileName = "lgballtBot.db"
	}
}
