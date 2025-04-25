package server

import (
	userdb "github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/handler"
	"github.com/gekich/go-web/internal/repository"
	"github.com/gekich/go-web/internal/service"
	validate "github.com/gekich/go-web/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) InitUserRoutes() {
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
		router.Put("/{id}", userHandler.UpdateUser)
		router.Delete("/{id}", userHandler.DeleteUser)
	})
}

func (s *Server) InitAuthRoutes() {
	queries := userdb.New(s.db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewAuthService(repo, s.cfg.Jwt.Secret, s.cfg.Jwt.TokenExpiry)
	validator := validate.New()
	authHandler := handler.NewAuthHandler(svc, validator)

	s.router.Post("/register", authHandler.Register)
	s.router.Post("/login", authHandler.Login)
}
