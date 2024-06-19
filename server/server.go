package server

import (
	"fmt"
	"net/http"
)

var htmlDir = "."

type server struct {
	srv *http.Server
}

func New(dir string, port int) *server {
	htmlDir = dir
	mux := getHandler()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return &server{
		srv: srv,
	}
}

func (s *server) Listen() error {
	return s.srv.ListenAndServe()
}

func getHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", noCacheHeaders(index))
	mux.HandleFunc("/ws", ws)
	return mux
}

func noCacheHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")
		f(w, r)
	}
}
