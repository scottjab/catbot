package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/scottjab/catbot/types"
	"github.com/scottjab/catlady"
	"sort"
)

func Handler(commands <-chan types.Command) {
	catLady := catlady.NewCatLady(CONFIG.Reddit.Username, CONFIG.Reddit.Password, CONFIG.Reddit.AppId, CONFIG.Reddit.AppSecret, CONFIG.Reddit.Subreddits, log.GetLevel())

	cmds := CONFIG.Reddit.Subreddits
	log.WithField("commands", cmds).Debug("commands and subreddits")
	for command := range commands {
		var response = ""
		log.WithField("command", command.Cmd).Debug("Possible command.")
		if command.Cmd == "cathelp" {

			var commands []string
			commands = make([]string, 1)
			for key, _ := range cmds {
				commands = append(commands, key)
			}
			sort.Strings(commands)
			_, _, target, err := command.SlackApi.OpenIMChannel(command.User)
			if err != nil {
				log.WithError(err).Warn("Could not create IM channel.")
			}
			rtm := command.SlackRtm
			rtm.SendMessage(rtm.NewOutgoingMessage("🐱 Commands: ", target))

			for i, cmd := range commands {
				if i%25 == 1 {
					rtm.SendMessage(rtm.NewOutgoingMessage(response[:len(response)-2], target))
					response = ""
				}
				response = response + cmd + ", "
			}
			if response != "" {
				rtm.SendMessage(rtm.NewOutgoingMessage(response[:len(response)-2], target))
			}
		}
		if reddit, ok := cmds[command.Cmd]; ok {
			response = catLady.GetImage(reddit)
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
