package jatek

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

var activeMessages = []string{
	"Regisztralj a jatekra, a kovetkezo paranccsal: !jatek",
	"Ne feledd hasznalni a kovetkezo parancsot: !jatek",
	"Gyere jatszani a !jatek paranccsal!",
}

type logic struct {
	active bool
	users  []string
	ticker *time.Ticker
}

func NewLogic(c *twitch.Client, channel string) *logic {
	l := &logic{
		users:  make([]string, 0),
		ticker: time.NewTicker(15 * time.Minute),
	}

	go func() {
		for range l.ticker.C {
			if l.active {
				rnd, _ := genRandNum(int64(len(activeMessages)))
				c.Say(channel, activeMessages[rnd])
			}
		}
	}()

	return l
}

func (l *logic) JatekHandler(m twitch.PrivateMessage) (string, error) {
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

func (l *logic) JatekStart(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	l.active = true
	l.users = make([]string, 0)
	return "Elindult a jatek!", nil
}

func (l *logic) JatekStop(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	l.active = false
	return "Vege a jateknak!", nil
}

func (l *logic) JatekSorsol(m twitch.PrivateMessage) (string, error) {
	if !isAdmin(m.User.Name) {
		return "", nil
	}

	if len(l.users) == 0 {
		return "Nincs regisztralt jatekos", nil
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
