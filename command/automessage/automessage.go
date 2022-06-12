package automessage

import (
	"math/rand"
	"time"
)

var messages = []string{}

type sayer func(string)

type Logic struct {
	say    sayer
	ticker *time.Ticker
}

func New(say sayer) *Logic {
	return &Logic{
		say: say,
	}
}

func (l *Logic) Start() {
	l.ticker = time.NewTicker(60 * time.Second)

	go func() {
		for range l.ticker.C {
			if len(messages) == 0 {
				continue
			}

			l.say(messages[rand.Intn(len(messages))])
		}
	}()
}

func (l *Logic) Stop() {
	l.ticker.Stop()
}
