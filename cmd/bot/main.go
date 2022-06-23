package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/gerifield/coderbot42/bot"
	"github.com/gerifield/coderbot42/command/automessage"
	"github.com/gerifield/coderbot42/command/autoraid"
	"github.com/gerifield/coderbot42/config"
	"github.com/gerifield/coderbot42/overlay"
	"github.com/gerifield/coderbot42/token"
)

func main() {
	configFile := flag.String("config", "config.json", "Config file")
	flag.Parse()

	conf, err := config.Load(*configFile)
	if err != nil {
		log.Println(err)
		return
	}

	tl := token.New(conf.Secret.ClientID, conf.Secret.ClientSecret)
	log.Println("Fetching token")
	token, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	client := twitch.NewClient(conf.BotName, "oauth:"+token.AccessToken)

	sayFn := func(msg string) {
		client.Say(conf.Channel, msg)
	}

	autoRaider, err := autoraid.New(sayFn, conf.Secret.ClientID, token.AccessToken, conf.RaidChannels)
	if err != nil {
		log.Println(err)
		return
	}
	bot.Register("!autoraid", autoRaider.Handler)

	messager := automessage.New(sayFn, conf.AutoMessages)
	defer messager.Stop()

	overlayLogic, err := overlay.New(conf.Server)
	if err != nil {
		log.Println(err)
		return
	}
	go overlayLogic.Start()
	defer overlayLogic.Stop(context.Background())

	client.OnConnect(func() {
		messager.Start()
	})

	client.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
		log.Println(message)
	})

	commandHandler := bot.PrivateMessageHandler(client)
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		commandHandler(m)
		
		overlayLogic.Send(fmt.Sprintf("%s: %s", m.User.DisplayName, m.Message))
	})
	client.Join(conf.Channel)

	log.Println("Connect with client")
	err = client.Connect()
	if err != nil {
		log.Println(err)
	}
}
