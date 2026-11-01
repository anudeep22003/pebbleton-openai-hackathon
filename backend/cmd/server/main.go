package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"pebbleton/backend/internal/game"
	"pebbleton/backend/internal/httpapi"
)

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func main() {
	seed, err := strconv.ParseInt(env("SEED", "42"), 10, 64)
	if err != nil {
		log.Fatal("SEED must be an integer")
	}
	app, err := httpapi.New(game.Mock{}, env("DATA_DIR", "./data"), seed, "mock")
	if err != nil {
		log.Fatal(err)
	}
	address := "127.0.0.1:" + env("PORT", "8080")
	server := &http.Server{Addr: address, Handler: app.Handler("frontend/dist"), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Pebbleton · mock mode · http://%s", address)
	log.Fatal(server.ListenAndServe())
}
