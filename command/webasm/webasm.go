package webasm

import (
	"github.com/gempir/go-twitch-irc/v3"
)

// Logic .
type Logic struct {
}

// New .
func New() {

}

func handler(_ twitch.PrivateMessage) (string, error) {
	return "Kappa Kappa Kappa", nil
}
