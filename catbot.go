package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"github.com/scottjab/catbot/types"

	"fmt"
	"os"
	"strings"
	"time"
)

var (
	commands     = make(chan types.Command)
	channelCache = cache.New(12*time.Hour, 1*time.Minute)
	userCache    = cache.New(12*time.Hour, 1*time.Minute)
)

func init() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

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
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
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
