package autoraid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type sayer func(string)

// Logic .
type Logic struct {
	preparedReq *http.Request
	channels    []string
	say         sayer
}

type streamInfo struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIds       []string  `json:"tag_ids"`
	IsMature     bool      `json:"is_mature"`
}

const streamsURL = "https://api.twitch.tv/helix/streams"

// New .
func New(say sayer, clientID string, accessToken string, channels []string) (*Logic, error) {
	if len(channels) == 0 {
		return nil, nil
	}

	req, err := http.NewRequest(http.MethodGet, streamsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Client-Id", clientID)

	params := req.URL.Query()
	for _, v := range channels {
		params.Add("user_login", v)
	}
	req.URL.RawQuery = params.Encode()

	return &Logic{
		preparedReq: req,
		channels:    channels,
		say:         say,
	}, nil
}

func (l *Logic) Handler(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	if l == nil {
		return "No channels set for raid :(", nil
	}

	streamInfos, err := l.fetchChannels()
	if err != nil {
		return "", err
	}

	activeChannels := make([]streamInfo, 0, len(l.channels))
	for _, c := range l.channels {
		info, ok := matchInfo(c, streamInfos)
		if !ok {
			continue
		}

		activeChannels = append(activeChannels, info)
	}

	if len(activeChannels) == 0 {
		return "No live channel to raid :(", nil
	}

	raidTarget := activeChannels[0]

	l.say("/raid " + raidTarget.UserLogin)

	return fmt.Sprintf("Rading %s soon!", raidTarget.UserLogin), nil
}

func matchInfo(channelName string, streamInfos []streamInfo) (streamInfo, bool) {
	for _, s := range streamInfos {
		if strings.EqualFold(s.UserLogin, channelName) {
			return s, true
		}
	}

	return streamInfo{}, false
}

func (l *Logic) fetchChannels() ([]streamInfo, error) {

	resp, err := http.DefaultClient.Do(l.preparedReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		StreamInfos []streamInfo `json:"data"`
		// Pagination
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.StreamInfos, nil
}

func isAdmin(nick string) bool {
	admins := []string{"gerifield"}

	for _, a := range admins {
		if strings.EqualFold(a, nick) {
			return true
		}
	}

	return false
}
