package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/adriancrafter/todoapp/internal/am"
)

const (
	staticFSRoot = "assets/web"
	staticPrefix = "/static/"
	staticPath   = staticFSRoot + staticPrefix
)

func (srv *Server) AddFileServer(staticFS embed.FS) {
	srv.router.PathPrefix(staticPrefix).Handler(
		srv.staticHandler(staticPrefix, staticPath))
}

func (srv *Server) staticHandler(urlPrefix, fsPrefix string) http.Handler {
	staticFS, _ := srv.fs.Get(am.WebFSKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = fsPrefix + strings.TrimPrefix(r.URL.Path, urlPrefix)
		http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
	})
}
