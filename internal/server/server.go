package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gekich/go-web/database"
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
)

type Server struct {
	cfg *config.Config

	db     *sql.DB
	router *chi.Mux

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
	s.NewDatabase()
	s.newRouter()
	s.InitUserRoutes()
	s.InitAuthRoutes()
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
