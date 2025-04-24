package main

import (
	"database/sql"
	"github.com/gekich/go-web/config"
	"github.com/gekich/go-web/database"
	db "github.com/gekich/go-web/integrations/database"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

func main() {
	log.Println("Migration started")

	cfg := config.New()
	store := db.New(cfg)
	migrator := database.Migrator(store)

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		log.Println("Running UP migrations")
		migrator.Up()

	case "down":
		log.Println("Running DOWN migration")
		migrator.Down()

	case "status":
		log.Println("Migration STATUS")
		showStatus(store)

	default:
		log.Fatalf("unknown migration command: %q", cmd)
	}
}

func showStatus(db *sql.DB) {
	if err := goose.Status(db, "migrations"); err != nil {
		log.Fatalf("goose status failed: %v", err)
	}
}
