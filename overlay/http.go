package overlay

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gerifield/coderbot42/overlay/middleware"
	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

func (l *Logic) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Auth(l.token))

	r.Handle("/overlay/*", http.StripPrefix("/overlay", http.FileServer(http.Dir(l.staticDir))))
	r.Get("/websocket", l.websocketHandler)

	return r
}

func (l *Logic) websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("[ERROR] websocket accept failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	connectionID := rand.Int63n(1000000000)
	log.Printf("new connection ID: %d\n", connectionID)

	event := make(chan string, 20)
	l.registerConnection(connectionID, event)
	l.logStats()

	defer func() {
		l.unregisterConnection(connectionID)
		close(event)
		l.logStats()
	}()

	go func() {
		for evt := range event {
			log.Printf("send event to %d, event: %s\n", connectionID, evt)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = c.Write(ctx, websocket.MessageText, []byte(evt))
			if err != nil {
				log.Println("[ERROR] websocket write failed", err)
				cancel()

				return
			}
			cancel()
		}
	}()

	for {
		_, _, err := c.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusGoingAway {
				log.Println("[ERROR] websocket read failed", err)
			} else {
				log.Printf("%d client disconnected\n", connectionID)
			}

			return
		}
	}
}

func (l *Logic) Start() error {
	l.httpSrv.Handler = l.routes()

	return l.httpSrv.ListenAndServe()
}

func (l *Logic) Stop(ctx context.Context) error {
	close(l.event)

	return l.httpSrv.Shutdown(ctx)
}

func (l *Logic) logStats() {
	l.connectionsLock.Lock()
	conns := len(l.connections)
	l.connectionsLock.Unlock()

	log.Printf("Connections: %d\n", conns)
}

func (l *Logic) registerConnection(cid int64, ch chan string) {
	l.connectionsLock.Lock()
	l.connections[cid] = ch
	l.connectionsLock.Unlock()
}

func (l *Logic) unregisterConnection(cid int64) {
	l.connectionsLock.Lock()
	delete(l.connections, cid)
	l.connectionsLock.Unlock()
}
