package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/patrickmn/go-cache"
	"github.com/scottjab/catbot/types"
	"sort"
	"strings"
	"time"
)

func Handler(commands <-chan types.Command) {
	catCache := cache.New(5*time.Minute, 30*time.Second)
	cmds := CONFIG.Reddit.Subreddits
	log.WithField("commands", cmds).Debug("commands and subreddits")
	for command := range commands {
		var response = ""
		log.WithField("command", command.Cmd).Debug("Possible command.")
		if command.Cmd == "help" {

			response = "🐱 Commands: "
			var commands []string
			commands = make([]string, 1)
			for key, _ := range cmds {
				commands = append(commands, key)
			}
			sort.Strings(commands)
			response = response + strings.Join(commands, ", ")[2:]

			log.WithFields(log.Fields{
				"command":  "help",
				"target":   command.User,
				"response": response,
			}).Info("Help Message.")
			_, _, target, err := command.SlackApi.OpenIMChannel(command.User)
			if err != nil {
				log.WithError(err).Warn("Could not create IM channel")
			}
			rtm := command.SlackRtm
			rtm.SendMessage(rtm.NewOutgoingMessage(response, target))
		}
		if reddit, ok := cmds[command.Cmd]; ok {
			response = GetImage(reddit, catCache)
			if response != "" {
				log.WithFields(log.Fields{
					"response": response,
					"command":  command.Cmd,
					"user":     command.User,
					"target":   command.Target,
				}).Info("Command response")
				rtm := command.SlackRtm
				rtm.SendMessage(rtm.NewOutgoingMessage(response, command.Target))
			}
		}
	}
}
