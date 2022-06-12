package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/gerifield/coderbot42/bot"
	"github.com/gerifield/coderbot42/command/automessage"
	"github.com/gerifield/coderbot42/command/autoraid"
	"github.com/gerifield/coderbot42/token"
)

func main() {
	channelName := flag.String("channel", "gerifield", "Twitch channel name")
	botName := flag.String("botName", "CoderBot42", "Bot name")
	clientID := flag.String("clientID", "", "Twitch App ClientID")
	clientSecret := flag.String("clientSecret", "", "Twitch App clientSecret")

	raidChannels := flag.String("raidChannels", "moakaash,gibbonrike,streaminks,marinemammalrescue", "Twitch channels name to check")
	flag.Parse()

	tl := token.New(*clientID, *clientSecret)
	log.Println("Fetching token")
	token, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	client := twitch.NewClient(*botName, "oauth:"+token.AccessToken)

	sayFn := func(msg string) {
		client.Say(*channelName, msg)
	}

	autoRaider, err := autoraid.New(sayFn, *clientID, token.AccessToken, strings.Split(*raidChannels, ","))
	if err != nil {
		log.Println(err)
		return
	}
	bot.Register("!autoraid", autoRaider.Handler)

	messager := automessage.New(sayFn)
	defer messager.Stop()

	client.OnConnect(func() {
		messager.Start()
	})

	client.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
		log.Println(message)
	})

	commandHandler := bot.PrivateMessageHandler(client)
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		commandHandler(m)
	})
	client.Join(*channelName)

	log.Println("Connect with client")
	err = client.Connect()
	if err != nil {
		log.Println(err)
	}
}
