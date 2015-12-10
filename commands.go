package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/patrickmn/go-cache"
	"github.com/scottjab/catbot/types"
	"time"
)

func Handler(commands <-chan types.Command) {
	catCache := cache.New(5*time.Minute, 30*time.Second)
	cmds := CONFIG.Reddit.Subreddits
	log.WithField("commands", cmds).Debug("commands and subreddits")
	for command := range commands {
		var response = ""
		if reddit, ok := cmds[command.Cmd]; ok {
			response = GetImage(reddit, catCache)
			if response != "" {
				log.WithFields(log.Fields{
					"response": response,
					"command":  command.Cmd,
					"target":   command.Target,
				}).Info("Command response")
				rtm := command.SlackRtm
				rtm.SendMessage(rtm.NewOutgoingMessage(response, command.Target))
			}
		}
	}
}
