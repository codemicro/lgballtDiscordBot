package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// TODO: Switch this file to being generated instead of manually updated

const configFileName = "botConfig.json"

type Info struct {
	MainGuildID                string           `json:"mainGuildID"`
	Token                      string           `json:"token"`
	Prefix                     string           `json:"prefix"`
	DbFileName                 string           `json:"dbFileName"`
	DebugMode                  bool             `json:"debug"`
	Statuses                   []string         `json:"statuses"`
	AdminRoles                 []string         `json:"adminRoles"`
	VerificationIDs            VerificationIds  `json:"verificationIds"`
	RedditFeeds                []RedditFeedInfo `json:"redditFeeds"`
	ChatChartChannelExclusions []string         `json:"ccExclusions"`
	OwnerIds                   []string         `json:"ownerIds"`
	PkApi                      PkApiInfo        `json:"pkApi"`
	Listeners                  ListenerInfo     `json:"listeners"`
	MuteMe                     MuteMeInfo       `json:"muteMe"`
	BioFields                  []string         `json:"bioFields"`
	ActionLogChannel           string           `json:"actionLogChannel"`
	PrometheusAddress          string           `json:"prometheusAddress"`
	AdminWebsiteAddress        string           `json:"adminWebsiteAddress"`
	AdminSite                  AdminInfo        `json:"adminSite"`
}

type AdminInfo struct {
	ServeAddress string `json:"serveAddress"`
	VisibleURL   string `json:"visibleURL"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
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
	MainGuildID                string
	Token                      string
	Prefix                     string
	DbFileName                 string
	Statuses                   []string
	DebugMode                  bool
	AdminRoles                 []string
	VerificationIDs            VerificationIds
	RedditFeeds                []RedditFeedInfo
	ChatChartChannelExclusions []string
	OwnerIds                   []string
	PkApi                      PkApiInfo
	Listeners                  ListenerInfo
	MuteMe                     MuteMeInfo
	BioFields                  []string
	ActionLogChannel           string
	PrometheusAddress          string
	AdminSite                  AdminInfo
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

	MainGuildID = cfg.MainGuildID
	Token = cfg.Token
	Prefix = cfg.Prefix
	DbFileName = cfg.DbFileName
	Statuses = cfg.Statuses
	DebugMode = cfg.DebugMode
	AdminRoles = cfg.AdminRoles
	VerificationIDs = cfg.VerificationIDs
	RedditFeeds = cfg.RedditFeeds
	ChatChartChannelExclusions = cfg.ChatChartChannelExclusions
	OwnerIds = cfg.OwnerIds
	PkApi = cfg.PkApi
	Listeners = cfg.Listeners
	MuteMe = cfg.MuteMe
	BioFields = cfg.BioFields
	ActionLogChannel = cfg.ActionLogChannel
	PrometheusAddress = cfg.PrometheusAddress
	AdminSite = cfg.AdminSite

	if PkApi.ApiUrl[len(PkApi.ApiUrl)-1] == '/' {
		PkApi.ApiUrl = PkApi.ApiUrl[:len(PkApi.ApiUrl)-1]
	}

	if DbFileName == "" {
		DbFileName = "lgballtBot.db"
	}
}
