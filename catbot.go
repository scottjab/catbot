package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kaneshin/pigeon"
	"github.com/mvdan/xurls"
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"github.com/scottjab/catbot/types"
	"google.golang.org/api/vision/v1"

	"os"
	"strings"
	"time"
)

var (
	commands = make(chan types.Command)
	findCats = make(chan *DetectionEvent)

	channelCache = cache.New(12*time.Hour, 1*time.Minute)
	userCache    = cache.New(12*time.Hour, 1*time.Minute)
)

type DetectionEvent struct {
	ItemRef  *slack.ItemRef
	Message  string
	Username string
}

func init() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

}

func addReact(emoji string, msg *slack.ItemRef, api *slack.Client) {
	err := api.AddReaction(emoji, *msg)
	if err != nil {
		log.WithError(err).Error("reaction failed")
	}

}

func react(msg *slack.ItemRef, username string, labels *vision.AnnotateImageResponse, api *slack.Client, rtm *slack.RTM) {
	log.Debug("reacting")
	//cats := ["cat", "cat2", "kittycat"]
	//dogs := ["dog", "dog2"]
	found := false
	for _, label := range labels.LabelAnnotations {
		log.WithFields(log.Fields{"desc": label.Description,
			"score": label.Score}).Debug("Found label")
		if label.Score > 0.5 {
			switch label.Description {
			case "cat":
				found = true
				addReact("cat", msg, api)
			case "dog":
				found = true
				addReact("dog", msg, api)
			case "shiba inu":
				found = true
				addReact("doge", msg, api)
			default:
			}
		}
	}
	if found {
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("!m %s", username), msg.Channel))
	}

}

// gorutine for finding cats.
func findCatGifs(lines <-chan *DetectionEvent, api *slack.Client, rtm *slack.RTM) {
	log.Debug("Starting cat finder")
	whitelist := [...]string{".jpg", "imgur", "imgur", "photobucket"}
	found := false
	client, err := pigeon.New(nil)
	feature := pigeon.NewFeature(pigeon.LabelDetection)
	feature.MaxResults = 15
	if err != nil {
		log.Fatalf("Unable to retrieve vision service: %v\n", err)
	}
	log.Debug("Starting to consume from channel")
	for event := range lines {
		log.WithField("message", event.Message).Debug("FOUND ONE")
		for _, thing := range whitelist {
			if strings.Contains(event.Message, thing) {
				found = true
				break
			}
		}
		if found {
			log.Debug("found link")
			url := xurls.Relaxed.FindString(event.Message)
			log.WithField("url", url).Debug("Found url")
			if url != "" {
				req, err := client.NewBatchAnnotateImageRequest([]string{url}, feature)
				if err != nil {
					log.WithError(err).Error("Shits fucked yo")
				}

				res, err := client.ImagesService().Annotate(req).Do()
				if err != nil {
					log.WithError(err).Warn("google hates me")
				}
				react(event.ItemRef, event.Username, res.Responses[0], api, rtm)
			}
		}
	}
}

func getChannelName(channelId string, api *slack.Client) string {
	var channelName string
	if channel, found := channelCache.Get(channelId); found {
		channelName = channel.(*slack.Channel).Name
		log.WithFields(log.Fields{
			"channelName": channelName,
			"channelId":   channelId,
		}).Debug("Found channel in cache")
	} else {
		log.WithField("channel", channelId).Debug("Channel missing in cache")
		channelInfo, err := api.GetChannelInfo(channelId)
		if err != nil {
			log.WithField("error", err).Warn("Slack channel ID lookup error")
			return ""
		}
		channelCache.Set(channelId, channelInfo, cache.DefaultExpiration)
		channelName = channelInfo.Name
	}
	return channelName
}

func getUserInfo(userId string, api *slack.Client) string {
	var userName string
	if user, found := userCache.Get(userId); found {
		if user == nil {
			return ""
		}
		userName = user.(*slack.User).Name
		log.WithFields(log.Fields{
			"userName": userName,
			"userId":   userId,
		}).Debug("Found user in cache")
	} else {
		log.WithField("user", userId).Debug("user missing in cache")
		userInfo, err := api.GetUserInfo(userId)
		if err != nil {
			log.WithField("error", err).Warn("Slack user ID lookup error")
			userCache.Set(userId, nil, cache.DefaultExpiration)
			return ""
		}
		userCache.Set(userId, userInfo, cache.DefaultExpiration)
		log.Info(userInfo)
		userName = userInfo.Name
	}
	return userName
}
func main() {
	log.Info("Starting up catbot")
	// Load the config first
	if len(os.Args) > 1 {
		log.WithField("configFile", os.Args[1]).Debug("Loading config from arguement.")
		LoadConfig(os.Args[1])
	} else {
		log.WithField("configFile", "./config.json").Debug("Loading default config.")
		LoadConfig("./config.json")
	}
	if CONFIG.SlackAPIKey == "" {
		log.Fatal("No Slack API Key")
		os.Exit(1)
	}
	prefix := "!"
	if CONFIG.Prefix != "" {
		log.WithField("prefix", CONFIG.Prefix).Debug("Prefix Found in Config.")
		prefix = CONFIG.Prefix
	} else {
		log.WithField("prefix", prefix).Debug("Using default prefix.")
	}
	api := slack.New(CONFIG.SlackAPIKey)
	api.SetDebug(CONFIG.Debug)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	go Handler(commands)
	go findCatGifs(findCats, api, rtm)

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				// Send the event to the findcat thread
				log.Debug("Sending event to finder")
				findCats <- &DetectionEvent{
					ItemRef: &slack.ItemRef{Channel: ev.Channel,
						Timestamp: ev.Timestamp,
					},
					Username: getUserInfo(ev.User, api),
					Message:  ev.Text}
				log.Debug("Sent log to finder")
				log.WithFields(log.Fields{
					"user":    getUserInfo(ev.User, api),
					"channel": getChannelName(ev.Channel, api),
					"message": ev.Text,
				}).Info(fmt.Sprintf("#%s <%s>: %s", getChannelName(ev.Channel, api), getUserInfo(ev.User, api), ev.Text))
				if ev.Text != "" && strings.Contains(ev.Text, " cat ") {
					log.Warn("CAT DETECTED!")
					var cmd types.Command
					cmd.Cmd = "cat"
					cmd.Target = ev.Channel
					cmd.User = ev.User
					cmd.SlackApi = api
					cmd.SlackRtm = rtm
					commands <- cmd
				}

				if ev.Text != "" && string(ev.Text[0]) == prefix {
					args := strings.Split(ev.Text[1:], " ")
					log.WithFields(log.Fields{
						"cmd":    args[0],
						"args":   args[1:],
						"target": ev.Channel,
					}).Debug("Command found!")

					var cmd types.Command
					cmd.Cmd = args[0]
					cmd.Args = args[1:]
					cmd.Target = ev.Channel
					cmd.User = ev.User
					cmd.SlackApi = api
					cmd.SlackRtm = rtm
					commands <- cmd
				}

			case *slack.LatencyReport:
				log.WithField("latency", ev.Value).Info("Latency Report")

			case *slack.RTMError:
				log.WithField("error", ev.Error()).Warn("RTM Error!")

			case *slack.InvalidAuthEvent:
				log.Fatal("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
