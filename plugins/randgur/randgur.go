package randgur

import (
	"bitbucket.org/liamstask/go-imgur/imgur"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"time"
)

type ImgurSpec struct {
	ClientID     string
	ClientSecret string
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomImageFromSubReddit(subreddit string) string {
	var config ImgurSpec
	envconfig.Process("imgur", &config)
	rand.Seed(time.Now().UTC().UnixNano())
	client := imgur.NewClient(nil, config.ClientID, config.ClientSecret)
	results, err := client.Gallery.Subreddit(subreddit, "time", "day", 0)
	if err != nil {
		log.Fatal(err)
	}
	return results[randInt(0, len(results)-1)].Link
}
