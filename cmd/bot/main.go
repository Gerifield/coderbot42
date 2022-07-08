package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/gerifield/coderbot42/bot"
	"github.com/gerifield/coderbot42/command/autoraid"
	"github.com/gerifield/coderbot42/command/jatek"

	//_ "github.com/gerifield/coderbot42/command/kappa"
	"github.com/gerifield/coderbot42/token"
)

func main() {
	channelName := flag.String("channel", "bate81", "Twitch channel name")
	botName := flag.String("botName", "CoderBot42", "Bot name")
	clientID := flag.String("clientID", "", "Twitch App ClientID")
	clientSecret := flag.String("clientSecret", "", "Twitch App clientSecret")

	jatekosFile := flag.String("jatekosFile", "jatekosok.json", "Jatekosok listajanak tarolasi helye")
	channelsName := flag.String("channels", "moakaash,gibbonrike,streaminks,marinemammalrescue", "Twitch channels name to check")
	flag.Parse()

	tl := token.New(*clientID, *clientSecret)
	log.Println("Fetching token")
	accToken, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	client := twitch.NewClient(*botName, "oauth:"+accToken.AccessToken)

	sayFn := func(msg string) {
		client.Say(*channelName, msg)
	}

	l, err := jatek.NewLogic(sayFn, *jatekosFile)
	if err != nil {
		log.Println(err)
		return
	}
	bot.Register("!jatek", l.JatekHandler)
	bot.Register("!jatek-start", l.JatekStart)
	bot.Register("!jatek-stop", l.JatekStop)
	bot.Register("!jatek-sorsol", l.JatekSorsol)

	autoRaider, err := autoraid.New(sayFn, *clientID, accToken.AccessToken, strings.Split(*channelsName, ","))
	if err != nil {
		log.Println(err)
		return
	}
	bot.Register("!autoraid", autoRaider.Handler)

	client.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
		log.Println("[JOIN]", message)
		sayFn("Hello there!")
		//client.Say(*channelName, "Hello!")
	})

	commandHandler := bot.PrivateMessageHandler(client)
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		l.CheerHandler(m)
		commandHandler(m)
	})
	client.Join(*channelName)

	log.Println("Connect with client")
	err = client.Connect()
	if err != nil {
		log.Println(err)
	}
}
