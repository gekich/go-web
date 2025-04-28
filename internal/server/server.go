package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gekich/go-web/database"
	"github.com/gekich/go-web/internal/auth"
	intLogger "github.com/gekich/go-web/internal/logger"
	"github.com/gekich/go-web/internal/middleware"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5"

	"github.com/gekich/go-web/config"

	db "github.com/gekich/go-web/integrations/database"
	logger "github.com/gekich/go-web/integrations/logger"
)

type Server struct {
	cfg        *config.Config
	jwtManager *auth.JWTManager

	db     *sql.DB
	router *chi.Mux
	logger intLogger.Logger

	httpServer *http.Server
}

type Options func(opts *Server) error

func New(opts ...Options) *Server {
	s := defaultServer()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return s
}

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(),
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
	s.InitLogger()
	s.InitJwtManager()
	s.NewDatabase()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitServiceRoutes()
	s.InitUserRoutes()
	s.InitAuthRoutes()
	s.InitProfileRoutes()
}

func (s *Server) InitLogger() {
	s.logger = logger.New(s.cfg.Env)
}

func (s *Server) InitJwtManager() {
	s.jwtManager = auth.NewJWTManager(s.cfg.Secret, s.cfg.TokenExpiry)
}

func (s *Server) NewDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}

	s.db = db.New(s.cfg)
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.Use(chiMiddleware.RequestID)
	s.router.Use(middleware.Json)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(chiMiddleware.Recoverer)
}

func (s *Server) Migrate() {
	log.Println("migrating...")

	var databaseUrl string
	switch s.cfg.Database.Driver {
	case "postgres":
		databaseUrl = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
			s.cfg.Database.Driver,
			s.cfg.Database.User,
			s.cfg.Database.Password,
			s.cfg.Database.Host,
			s.cfg.Database.Port,
			s.cfg.Database.Name,
			s.cfg.Database.SslMode,
		)
	case "mysql":
		databaseUrl = fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
			s.cfg.Database.User,
			s.cfg.Database.Password,
			s.cfg.Database.Host,
			s.cfg.Database.Port,
			s.cfg.Database.Name,
		)
	}

	migrator := database.Migrator(s.db, database.WithDSN(databaseUrl))
	migrator.Up()

	log.Println("done migration.")
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:              s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Api.ReadHeaderTimeout,
	}

	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func start(s *Server) {
	log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
	log.Printf("http://%s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	s.closeResources()

	return nil
}

func (s *Server) closeResources() {
	_ = s.db.Close()
}
