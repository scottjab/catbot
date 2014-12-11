package catgur

import (
	"../types"
	"./randgur"
	"log"
)

func Handler(commands <-chan types.Command) {
	cmds := map[string]string{
		"cat":      "catgifs",
		"standup":  "catsstandingup",
		"startled": "StartledCats",
		"kitten":   "kittengifs",
		"catte":    "catts",
		"lolcat":   "lolcats",
		"space":    "spacecats",
		"chubby":   "delightfullychubby",
		"sink":     "catsinsinks",
		"catt":     "CatPics",
		"sitting":  "sittinglikehumans",
		"tuckedin": "tuckedinkitties",
		"bigcat": "bigcats",
		"beard": "beardsandcats",
		"bacon": "catswithbacon",
		"catbelly": "catbellies",
		"tuxedo": "tuxedocats",
	}

	for command := range commands {
		var response = ""
		if reddit, ok := cmds[command.Cmd]; ok {
			response = randgur.RandomImageFromSubReddit(reddit)
			log.Println(command.Line.Target(), ":", response)
			command.Conn.Privmsg(command.Line.Target(), response)
		}
	}
}
