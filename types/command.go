package types

import (
	"github.com/nlopes/slack"
)

type Command struct {
	Cmd      string
	Args     []string
	Target   string
	User     string
	SlackApi *slack.Client
	SlackRtm *slack.RTM
}

type ConfigSpec struct {
	Server   string
	Username string
	Password string
	Prefix   string
	Ssl      bool
	Port     int
	Chans    string
}
