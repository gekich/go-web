package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) InitRoutes() {
	s.router.Route("/ping", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})
}
