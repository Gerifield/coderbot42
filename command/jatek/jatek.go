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

	// Subscriberek
	if isSubscirber(m) {
		// Mar van 2 reg, off
		if numOfRegs(l.users, usr) == 2 {
			return fmt.Sprintf("%s mar regisztralt.", usr), nil
		}

		// Beallitani 2 regre
		setRegs(&l.users, usr, 2)
	} else {

		// Nem sub
		if numOfRegs(l.users, usr) == 1 {
			return fmt.Sprintf("%s mar regisztralt.", usr), nil
		}

		// Beallitani 1-re
		setRegs(&l.users, usr, 1)
	}

	log.Println("User:", l.users, len(l.users))
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
	log.Println("Random gen:", rnd, winner, l.users)

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

func numOfRegs(users []string, usr string) int {
	usr = strings.ToLower(usr)
	regs := 0
	for _, u := range users {
		if strings.ToLower(u) == usr {
			regs = regs + 1
		}
	}

	return regs
}

func setRegs(users *[]string, usr string, desiredReg int) {
	for numOfRegs(*users, usr) < desiredReg {
		*users = append(*users, usr)
	}
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

func isSubscirber(m twitch.PrivateMessage) bool {
	for _, i := range strings.Split(m.Tags["badge-info"], ",") {
		parts := strings.Split(i, "/")

		if len(parts) < 1 {
			continue
		}

		if parts[0] == "subscriber" {
			return true
		}
	}

	return false
}
