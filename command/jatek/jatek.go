package jatek

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/gerifield/coderbot42/bot"

	"github.com/gempir/go-twitch-irc/v3"
)

func init() {
	l := newLogic()
	bot.Register("!jatek", l.jatekHandler)
	bot.Register("!jatek-start", l.jatekStart)
	bot.Register("!jatek-stop", l.jatekStop)
	bot.Register("!jatek-sorsol", l.jatekSorsol)
}

type logic struct {
	active bool
	users  []string
}

func newLogic() *logic {
	return &logic{
		users: make([]string, 0),
	}
}

func (l *logic) jatekHandler(m twitch.PrivateMessage) (string, error) {
	if !l.active {
		return "", nil
	}

	usr := m.User.DisplayName

	if inUsers(l.users, usr) {
		return fmt.Sprintf("%s mar regisztralt.", usr), nil
	}

	l.users = append(l.users, usr)
	log.Println("User:", l.users)
	return fmt.Sprintf("%s regisztralt a jatekra!", usr), nil
}

func (l *logic) jatekStart(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	l.active = true
	l.users = make([]string, 0)
	return "Elindult a jatek!", nil
}

func (l *logic) jatekStop(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	l.active = false
	return "Vege a jateknak!", nil
}

func (l *logic) jatekSorsol(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	rnd, err := genRandNum(int64(len(l.users)))
	if err != nil {
		return "Random generalasi hiba", err
	}

	winner := l.users[int(rnd)]

	return fmt.Sprintf("Nyertes: %s", winner), nil
}

func genRandNum(max int64) (int64, error) {
	bg := big.NewInt(max - 0)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}

	return n.Int64(), nil
}

func inUsers(users []string, usr string) bool {
	for _, u := range users {
		if u == usr {
			return true
		}
	}

	return false
}

func isAdmin(nick string) bool {
	admins := []string{"Bate81", "gerifield"}

	for _, a := range admins {
		if strings.ToLower(a) == strings.ToLower(nick) {
			return true
		}
	}

	return false
}
