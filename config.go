package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	SlackAPIKey string       `json:"SLACK_API_KEY"`
	Reddit      RedditConfig `json:"reddit"`
	Debug       bool         `json:"debug"`
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
	err = json.Unmarshal(file, &CONFIG)
	if err != nil {
		log.Println("Error Parsing config file: ", err)
		os.Exit(1)
	}
}
