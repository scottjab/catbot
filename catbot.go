package main

import (
	"github.com/nlopes/slack"
	"github.com/scottjab/catbot/types"
	"log"
	"os"
	"strings"
)

var commands = make(chan types.Command)

func main() {
	log.Println("Starting up catbot")
	// Load the config first
	if len(os.Args) > 1 {
		LoadConfig(os.Args[1])
	} else {
		LoadConfig("./config.json")
	}
	if CONFIG.SlackAPIKey == "" {
		log.Printf("Missing SLACK_API_KEY")
		os.Exit(1)
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
				log.Printf("Message: %v\n", ev)
				if ev.Text != "" && strings.Contains(ev.Text, " cat ") {
					var cmd types.Command
					cmd.Cmd = "cat"
					cmd.Target = ev.Channel
					cmd.SlackApi = api
					cmd.SlackRtm = rtm
					commands <- cmd
				}
				if ev.Text != "" && ev.Text[0] == '!' {
					args := strings.Split(ev.Text[1:], " ")
					var cmd types.Command
					cmd.Cmd = args[0]
					cmd.Args = args[1:]
					cmd.Target = ev.Channel
					cmd.SlackApi = api
					cmd.SlackRtm = rtm
					commands <- cmd
				}

			case *slack.LatencyReport:
				log.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
