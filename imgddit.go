package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	token         AuthToken
	lastTokenTime time.Time
	apiCache      *cache.Cache
	userAgent     = "Catbot/1 by cattebot"
)

type Image struct {
	url    string
	title  string
	domain string
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func getToken() {
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("username", CONFIG.Reddit.Username)
	data.Add("password", CONFIG.Reddit.Password)
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		log.WithError(err).Error("Failed to build Authentication POST")
	}
	req.Header.Add("User-Agent", userAgent)
	req.SetBasicAuth(CONFIG.Reddit.AppId, CONFIG.Reddit.AppSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to login to Reddit")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(contents, &token)
	lastTokenTime = time.Now()
	log.WithFields(log.Fields{
		"AccessToken":   token.AccessToken,
		"ExpiresIn":     token.ExpiresIn,
		"Scope":         token.Scope,
		"tokenType":     token.TokenType,
		"lastTokenTime": lastTokenTime,
	}).Debug("Got token!")
}

func randInt(min int, max int) int {
	if max <= min {
		return max
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func getReddit(sub string) RedditResponse {
	log.WithField("subreddit", sub).Debug("Getting Reddit")
	log.WithField("token", token).Debug("token value")
	if token.ExpiresIn == 0 || time.Since(lastTokenTime).Seconds() >= float64(token.ExpiresIn) {
		getToken()
	}
	client := &http.Client{}
	reqUrl := fmt.Sprintf("https://oauth.reddit.com/r/%s.json", sub)
	log.WithField("url", reqUrl).Debug("Request URL")
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.WithError(err).Error("Failed request to reddit")
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Getting subreddit failed")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	log.WithField("rawjson", string(contents)).Debug("Raw response from reddit")
	var subReddit RedditResponse
	json.Unmarshal(contents, &subReddit)
	return subReddit
}

func checkForImage(url string) bool {
	whitelist := [...]string{"imgur.com", "imgur", "giphy", "flickr", "photobucket", "youtube", "youtu.be", "gif", "gifv", "png", "jpg", "tiff", "webem", "bmp", "flv", "mpg", "mpeg", "avi"}
	for _, thing := range whitelist {
		if strings.Contains(url, thing) {
			log.WithField("url", url).Debug("Found Image")
			return true
		}
	}
	log.WithField("url", url).Debug("Didn't Find Image")
	return false
}

func cleanURL(url string) string {
	if strings.Contains(url, "imgur") {
		log.WithField("url", url).Debug("Found imgur url")
		if url[len(url)-3:] == "gif" {
			url = url + "v"
			log.WithField("url", url).Debug("Converting to gifv")
		}
	}
	return url
}
func GetImage(sub string, cat_cache *cache.Cache) string {
	var submissions RedditResponse
	if subs, found := cat_cache.Get(sub); !found {
		log.WithFields(log.Fields{
			"cache":     false,
			"subreddit": sub,
		}).Info("Subreddit not found in cache.")
		submissions = getReddit(sub)
		cat_cache.Set(sub, submissions, cache.DefaultExpiration)
		log.WithField("subreddit", submissions.Data.Children[0].Data.URL).Debug("Subreddit value")
	} else {
		log.WithFields(log.Fields{
			"cache":     true,
			"subreddit": sub,
		}).Info("Subreddit Found in Cache.")

		submissions = subs.(RedditResponse)

	}
	size := len(submissions.Data.Children)
	count := 0
	noImage := true
	for noImage {
		count += 1
		random := randInt(0, size-1)
		s := submissions.Data.Children[random].Data
		if !s.Over18 {
			if checkForImage(s.URL) {
				noImage = false
				return s.URL
			}
		} else {
			log.WithField("nsfw", "true").Info("NSFW Link Found.")
		}
		if count >= size {
			log.WithFields(log.Fields{
				"subreddit": sub,
				"size":      size,
			}).Info("I ran out of links")
			return ""
		}
	}
	return ""
}
