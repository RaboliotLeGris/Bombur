package main

import (
	"context"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/raboliotlegris/bombur/db"
	"github.com/raboliotlegris/bombur/handlers"
)

func main() {
	log.SetLevel(log.DebugLevel)

	log.Info("Launching Bombur - Yet Another URL Shortener")

	// "Config"
	DB_URI := os.Getenv("BOMBUR_DB_URI")
	if DB_URI == "" {
		DB_URI = os.Getenv("DATABASE_URL") // Heroku default DB env var
	}
	if DB_URI == "" {
		log.Fatal("Missing BOMBUR_DB_URI env var")
	}

	var port uint64
	if tmp_port := os.Getenv("PORT"); tmp_port != "" {
		var err error
		if port, err = strconv.ParseUint(tmp_port, 10, 64); err != nil {
			log.Fatal("Unable to parse Port with error: ", err)
		}
	} else {
		port = 7777
	}

	log.Info("Initializing database")
	if err := db.InitDB(DB_URI); err != nil {
		log.Fatal("Failed to init DB with error : ", err)
	}

	log.Info("Creating DB pool")
	pool, err := pgxpool.Connect(context.Background(), DB_URI)
	if err != nil {
		log.Fatal("Failed to connect to db with error:", err)
	}

	log.Info("Creating routes")
	if err = handlers.Launch(handlers.Create_router(pool), port); err != nil {
		log.Fatal("Bombur crash with error: ", err)
	}
}
