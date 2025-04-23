package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/gekich/go-web/config"
)

func New(cfg *config.Config) *sql.DB {
	var dsn string

	switch cfg.Database.Driver {
	case "postgres":
		dsn = fmt.Sprintf("%s://%s:%d/%s?sslmode=%s&user=%s&password=%s",
			cfg.Database.Driver,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SslMode,
			cfg.Database.User,
			cfg.Database.Password)
	default:
		log.Fatal("Must choose a database driver")
	}

	db, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxConnectionPool)

	return db
}
