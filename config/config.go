package config

import (
	"encoding/json"
	"os"
)

type Server struct {
	Addr      string `json:"addr"`
	Token     string `json:"token"`
	StaticDir string `json:"staticDir"`
}

type Config struct {
	Channel string `json:"channel"`
	BotName string `json:"botName"`

	Secret struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
	} `json:"secret"`

	Server Server `json:"server"`

	RaidChannels []string `json:"raidChannels"`
	AutoMessages []string `json:"autoMessages"`
}

func Load(name string) (Config, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = json.Unmarshal(b, &c)

	return c, err
}
