package overlay

import (
	"log"
	"net/http"
	"sync"

	"github.com/gerifield/coderbot42/config"
)

type Logic struct {
	httpSrv   *http.Server
	token     string
	staticDir string

	event chan string

	connectionsLock *sync.Mutex
	connections     map[int64]chan string
}

func New(conf config.Server) (*Logic, error) {
	l := &Logic{
		httpSrv:         &http.Server{Addr: conf.Addr},
		token:           conf.Token,
		staticDir:       conf.StaticDir,
		event:           make(chan string, 20),
		connectionsLock: &sync.Mutex{},
		connections:     make(map[int64]chan string),
	}

	// go func() {
	// 	ticker := time.NewTicker(5 * time.Second)
	// 	i := 0
	// 	for range ticker.C {
	// 		l.event <- fmt.Sprintf("hello-%d", i)
	// 		i++
	// 	}
	// }()

	go l.eventLoop()

	return l, nil
}

func (l *Logic) Send(msg string) {
	l.event <- msg
}

func (l *Logic) eventLoop() {
	for evt := range l.event {
		l.connectionsLock.Lock()
		for id, ch := range l.connections {
			select {
			case ch <- evt:
			default:
				log.Printf("message skipped, cid: %d\n", id)
			}

		}
		l.connectionsLock.Unlock()
	}
}
