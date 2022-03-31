package bot

import (
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

type sayer interface {
	Say(channel, text string)
}

type MessageHandler func(message twitch.PrivateMessage) (string, error)

var commands = make(map[string]MessageHandler)

func Register(cmd string, h MessageHandler) {
	commands[cmd] = h
}

func PrivateMessageHandler(c sayer) func(msg twitch.PrivateMessage) {
	return func(msg twitch.PrivateMessage) {
		log.Println("msg:", msg)

		if !strings.HasPrefix(msg.Message, "!") {
			return
		}

		parts := strings.Split(msg.Message, " ")
		if len(parts) > 0 {
			h, ok := commands[parts[0]]
			if !ok {
				log.Println("command not found")
				return
			}

			resp, err := h(msg)
			if err != nil {
				log.Println("command error", err)
			}

			if resp != "" {
				c.Say(msg.Channel, resp)
			}
		}
	}
}
