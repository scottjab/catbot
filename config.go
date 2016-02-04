package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Config struct {
	SlackAPIKey string       `json:"SLACK_API_KEY"`
	Reddit      RedditConfig `json:"reddit"`
	Debug       bool         `json:"debug"`
	Prefix      string       `json:"prefix"`
}

type RedditConfig struct {
	Username   string            `json:"username"`
	Password   string            `json:"password"`
	AppId      string            `json:"appid"`
	AppSecret  string            `json:"appsecret"`
	Subreddits map[string]string `json:"subreddits"`
}

var CONFIG Config

func CheckForEnvVars() {
	if CONFIG.SlackAPIKey == "" || CONFIG.SlackAPIKey == "SLACKAPIKEY" {
		CONFIG.SlackAPIKey = os.Getenv("CATBOT_SLACK_APIKEY")
		log.WithField("CATBOT_SLACK_APIKEY", CONFIG.SlackAPIKey).Debug("Found Slack api key")
	}
	if CONFIG.Reddit.Username == "" || CONFIG.Reddit.Username == "REDDITUSERNAME" {
		CONFIG.Reddit.Username = os.Getenv("CATBOT_REDDIT_USERNAME")
		log.WithField("CATBOT_REDDIT_USERNAME", CONFIG.Reddit.Username).Debug("Found Reddit username")
	}
	if CONFIG.Reddit.Password == "" || CONFIG.Reddit.Password == "REDDITPASSWORD" {
		CONFIG.Reddit.Password = os.Getenv("CATBOT_REDDIT_PASSWORD")
		log.WithField("CATBOT_REDDIT_PASSWORD", CONFIG.Reddit.Password).Debug("Found Reddit password")
	}
	if CONFIG.Reddit.AppId == "" || CONFIG.Reddit.AppId == "REDDITAPPID" {
		CONFIG.Reddit.AppId = os.Getenv("CATBOT_REDDIT_APPID")
		log.WithField("CATBOT_REDDIT_APPID", CONFIG.Reddit.AppId).Debug("Found appid")
	}
	if CONFIG.Reddit.AppSecret == "" || CONFIG.Reddit.AppSecret == "REDDITAPPSECRET" {
		CONFIG.Reddit.AppSecret = os.Getenv("CATBOT_REDDIT_APPSECRET")
		log.WithField("CATBOT_REDDIT_APPSECRET", CONFIG.Reddit.AppSecret).Debug("Found appsecret")
	}

}

func LoadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Unable to open config file at: ", path)
		os.Exit(1)
	}
	log.WithField("configPath", path).Debug("Loading Config")
	err = json.Unmarshal(file, &CONFIG)
	CheckForEnvVars()
	log.WithField("configPath", path).Debug("Finished Loading Config")
	if err != nil {
		log.WithField("error", err).Fatal("Error Parsing config file: ")
		os.Exit(1)
	}
}
