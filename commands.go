package main

import (
	"github.com/patrickmn/go-cache"
	"github.com/scottjab/catbot/types"
	"time"
)

func Handler(commands <-chan types.Command) {
	catCache := cache.New(5*time.Minute, 30*time.Second)
	cmds := CONFIG.Reddit.Subreddits
	for command := range commands {
		var response = ""
		if reddit, ok := cmds[command.Cmd]; ok {
			response = GetImage(reddit, catCache)
			if response != "" {
				rtm := command.SlackRtm
				rtm.SendMessage(rtm.NewOutgoingMessage(response, command.Target))
			}
		}
	}
}
