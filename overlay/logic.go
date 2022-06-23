package overlay

import (
	"net/http"

	"github.com/gerifield/coderbot42/config"
)

type Logic struct {
	httpSrv   *http.Server
	token     string
	staticDir string
}

func New(conf config.Server) (*Logic, error) {
	return &Logic{
		httpSrv:   &http.Server{Addr: conf.Addr},
		token:     conf.Token,
		staticDir: conf.StaticDir,
	}, nil
}
