package main

import (
	"flag"
	"log"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/gerifield/coderbot42/bot"
	"github.com/gerifield/coderbot42/command/jatek"
	_ "github.com/gerifield/coderbot42/command/kappa"
	"github.com/gerifield/coderbot42/token"
)

func main() {
	channelName := flag.String("channel", "gerifield", "Twitch channel name")
	botName := flag.String("botName", "CoderBot42", "Bot name")
	clientID := flag.String("clientID", "", "Twitch App ClientID")
	clientSecret := flag.String("clientSecret", "", "Twitch App clientSecret")
	flag.Parse()

	tl := token.New(*clientID, *clientSecret)
	log.Println("Fetching token")
	token, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	client := twitch.NewClient(*botName, "oauth:"+token.AccessToken)

	l := jatek.NewLogic(client, *channelName)
	bot.Register("!jatek", l.JatekHandler)
	bot.Register("!jatek-start", l.JatekStart)
	bot.Register("!jatek-stop", l.JatekStop)
	bot.Register("!jatek-sorsol", l.JatekSorsol)

	client.OnPrivateMessage(bot.PrivateMessageHandler(client))
	client.Join(*channelName)

	log.Println("Connect with client")
	err = client.Connect()
	if err != nil {
		log.Println(err)
	}
}
