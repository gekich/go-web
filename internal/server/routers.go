package server

import (
	userdb "github.com/gekich/go-web/internal/db/user"
	handler "github.com/gekich/go-web/internal/handler/user"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) InitRoutes() {
	userQueries := userdb.New(s.db)
	userHandler := handler.NewUserHandler(userQueries)

	s.router.Route("/ping", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})
	s.router.Route("/users", func(router chi.Router) {
		router.Get("/{id}", userHandler.GetUserByID)
	})
}
