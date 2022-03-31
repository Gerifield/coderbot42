package kappa

import (
	"github.com/gerifield/coderbot42/bot"

	"github.com/gempir/go-twitch-irc/v3"
)

func init() {
	bot.Register("!kappa", handler)
}

func handler(_ twitch.PrivateMessage) (string, error) {
	return "Kappa Kappa Kappa", nil
}
