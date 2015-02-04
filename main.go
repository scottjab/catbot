package main

import (
	"./plugins"
	"./types"
	irc "github.com/fluffle/goirc/client"
	"github.com/kelseyhightower/envconfig"
	"log"
	"strings"
)

const (
	PREFIX = "!"
)

var commands = make(chan types.Command)

func commandDispatcher(conn *irc.Conn, line *irc.Line) {
	log.Println(line.Target(), ":", line.Nick, ": ", line.Text())
	if line.Text()[0:1] == PREFIX {
		args := strings.Split(line.Text()[1:], " ")
		var command types.Command
		command.Cmd = args[0]
		command.Args = args[1:]
		command.Conn = conn
		command.Line = line
		log.Println("Command: ", command.Cmd)
		commands <- command
	}
}

func handleConfig() types.ConfigSpec {
	var config types.ConfigSpec
	err := envconfig.Process("catbot", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	// is there a better way?

	if config.Server == "" {
		log.Fatal("Server is required.")
	}
	if config.Port == 0 {
		config.Port = 6697
	}
	if config.Username == "" {
		log.Fatal("Username is required.")
	}
	if config.Password == "" {
		log.Fatal("Password is required.")
	}
	if config.Prefix == "" {
		config.Prefix = "!"
	}
	return config
}

func main() {
	config := handleConfig()
	cfg := irc.NewConfig(config.Username)
	cfg.SSL = config.Ssl
	cfg.Server = config.Server
	cfg.Pass = config.Password
	c := irc.Client(cfg)
	go catgur.Handler(commands)
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			if config.Chans == "" {
				conn.Join("#cats")
			} else {
				for _, channel := range strings.Split(config.Chans, ",") {
					conn.Join(channel)
				}
			}
		})
	quit := make(chan bool)
	c.HandleFunc("privmsg", commandDispatcher)

	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.HandleFunc("part",
		func(conn *irc.Conn, line *irc.Line) {
			if line.Target() == "#cats" && line.Nick == "kelly" {
				log.Printf("%s: Inviting %s to #s", line.Time, line.Nick, line.Target())
				conn.Privmsg(line.Target(), "YOU CAN NEVER LEAVE!")
				conn.Invite(line.Nick, line.Target())
			}
		})

	if err := c.ConnectTo(config.Server); err != nil {
		log.Printf("%s", err)
	}

	<-quit
}
