package main

import (
	"github.com/gekich/go-web/config"
	"github.com/gekich/go-web/database"
	db "github.com/gekich/go-web/integrations/database"
	"log"
)

func main() {
	log.Println("Migration started")

	cfg := config.New()
	store := db.New(cfg)
	migrator := database.Migrator(store)

	// todo: accept cli flag for other operations
	migrator.Down()
}
