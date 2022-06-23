package overlay

import (
	"context"
	"net/http"

	"github.com/gerifield/coderbot42/overlay/middleware"
	"github.com/go-chi/chi/v5"
)

func (l *Logic) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Auth(l.token))

	r.Handle("/overlay", http.StripPrefix("/overlay", http.FileServer(http.Dir(l.staticDir))))
	r.Get("/websocket", l.websocketHandler)

	return r
}

func (l *Logic) websocketHandler(w http.ResponseWriter, r *http.Request) {

}

func (l *Logic) Start() error {
	l.httpSrv.Handler = l.routes()

	return l.httpSrv.ListenAndServe()
}

func (l *Logic) Stop(ctx context.Context) error {
	return l.httpSrv.Shutdown(ctx)
}
