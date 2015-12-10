package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jzelinskie/geddit"
	"github.com/patrickmn/go-cache"
	"math/rand"
	"strings"
	"time"
)

type Image struct {
	url    string
	title  string
	domain string
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
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
	var submissions []*geddit.Submission
	if subs, found := cat_cache.Get(sub); !found {
		log.WithFields(log.Fields{
			"cache":     false,
			"subreddit": sub,
		}).Info("Subreddit not found in cache.")
		session, err := geddit.NewLoginSession(
			CONFIG.Reddit.Username,
			CONFIG.Reddit.Password,
			"linux:com.catbot:1 (by /u/cattebot)",
		)

		if err != nil {
			log.WithField("error", err).Warn("Reddit connection error.")
			return ""
		}
		subOpts := geddit.ListingOptions{
			Limit: 25,
		}
		submissions, _ = session.SubredditSubmissions(sub, geddit.DefaultPopularity, subOpts)
		cat_cache.Set(sub, submissions, cache.DefaultExpiration)
	} else {
		log.WithFields(log.Fields{
			"cache":     true,
			"subreddit": sub,
		}).Info("Subreddit Found in Cache.")

		submissions = subs.([]*geddit.Submission)

	}
	size := len(submissions)
	count := 0
	noImage := true
	for noImage {
		count += 1
		random := randInt(0, size-1)
		s := submissions[random]
		if !s.IsNSFW {
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
