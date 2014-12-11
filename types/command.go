package types

import (
	irc "github.com/fluffle/goirc/client"
)

type Command struct {
	Cmd  string
	Args []string
	Conn *irc.Conn
	Line *irc.Line
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
