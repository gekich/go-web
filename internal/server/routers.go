package server

import (
	userdb "github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/handler"
	"github.com/gekich/go-web/internal/middleware"
	"github.com/gekich/go-web/internal/repository"
	"github.com/gekich/go-web/internal/service"
	validate "github.com/gekich/go-web/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s Server) InitServiceRoutes() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not found"}`))
	})

	s.router.Route("/ping", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})
}

func (s *Server) InitUserRoutes() {
	queries := userdb.New(s.db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(repo)
	validator := validate.New()
	userHandler := handler.NewUserHandler(svc, validator, s.logger)
	jwtMiddleware := middleware.JWTAuthMiddleware(s.jwtManager)

	s.router.Route("/users", func(router chi.Router) {
		router.Use(jwtMiddleware)
		router.Get("/{id}", userHandler.GetUserByID)
		router.Get("/", userHandler.ListUsers)
		router.Put("/{id}", userHandler.UpdateUser)
		router.Delete("/{id}", userHandler.DeleteUser)
	})
}

func (s *Server) InitAuthRoutes() {
	queries := userdb.New(s.db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewAuthService(repo, s.jwtManager)
	validator := validate.New()
	authHandler := handler.NewAuthHandler(svc, validator)

	s.router.Post("/register", authHandler.Register)
	s.router.Post("/login", authHandler.Login)
}

func (s *Server) InitProfileRoutes() {
	queries := userdb.New(s.db)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(repo)
	profileHandler := handler.NewProfileHandler(svc)
	jwtMiddleware := middleware.JWTAuthMiddleware(s.jwtManager)

	s.router.With(jwtMiddleware).Get("/profile", profileHandler.Profile)
}
