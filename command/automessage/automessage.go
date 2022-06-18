package automessage

import (
	"math/rand"
	"time"
)

type sayer func(string)

type Logic struct {
	say      sayer
	ticker   *time.Ticker
	messages []string
}

func New(say sayer, messages []string) *Logic {
	return &Logic{
		say:      say,
		messages: messages,
	}
}

func (l *Logic) Start() {
	l.ticker = time.NewTicker(60 * time.Second)

	go func() {
		for range l.ticker.C {
			if len(l.messages) == 0 {
				continue
			}

			l.say(l.messages[rand.Intn(len(l.messages))])
		}
	}()
}

func (l *Logic) Stop() {
	l.ticker.Stop()
}
