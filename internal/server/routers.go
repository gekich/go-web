package server

import (
	userdb "github.com/gekich/go-web/internal/db/user"
	handler "github.com/gekich/go-web/internal/handler/user"
	repository "github.com/gekich/go-web/internal/repository/user"
	service "github.com/gekich/go-web/internal/service/user"
	validate "github.com/gekich/go-web/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) InitRoutes() {
	queries := userdb.New(s.db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(repo)
	validator := validate.New()
	userHandler := handler.NewUserHandler(svc, validator)

	s.router.Route("/ping", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})
	s.router.Route("/users", func(router chi.Router) {
		router.Get("/{id}", userHandler.GetUserByID)
		router.Get("/", userHandler.ListUsers)
		router.Post("/", userHandler.CreateUser)
		router.Put("/{id}", userHandler.UpdateUser)
		router.Delete("/{id}", userHandler.DeleteUser)
	})
}
