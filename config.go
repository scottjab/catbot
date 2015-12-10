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
	Subreddits map[string]string `json:"subreddits"`
}

var CONFIG Config

func LoadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Unable to open config file at: ", path)
		os.Exit(1)
	}
	log.WithField("configPath", path).Debug("Loading Config")
	err = json.Unmarshal(file, &CONFIG)
	log.WithField("configPath", path).Debug("Finished Loading Config")
	if err != nil {
		log.WithField("error", err).Fatal("Error Parsing config file: ")
		os.Exit(1)
	}
}
