package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"pebbleton/backend/internal/ai"
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
	mock, err := strconv.ParseBool(env("MOCK_AI", "true"))
	if err != nil {
		log.Fatal("MOCK_AI must be true or false")
	}
	var provider game.Provider = game.Mock{}
	mode := "mock"
	if !mock {
		if os.Getenv("OPENAI_API_KEY") == "" {
			log.Fatal("OPENAI_API_KEY is required when MOCK_AI=false")
		}
		provider = ai.New(os.Getenv("OPENAI_API_KEY"), env("OPENAI_TEXT_MODEL", "gpt-5.4-mini"), env("OPENAI_DECISIONS_MODEL", "gpt-6-luna"), env("OPENAI_IMAGE_MODEL", "gpt-image-2.5-flare"))
		mode = "live"
	}
	app, err := httpapi.New(provider, env("DATA_DIR", "./data"), seed, mode)
	if err != nil {
		log.Fatal(err)
	}
	generateImages, err := strconv.ParseBool(env("GENERATE_IMAGES", "false"))
	if err != nil {
		log.Fatal("GENERATE_IMAGES must be true or false")
	}
	if generateImages && !mock {
		app.EnableIllustrations(provider.(*ai.OpenAI))
	}
	address := "127.0.0.1:" + env("PORT", "8080")
	server := &http.Server{Addr: address, Handler: app.Handler("frontend/dist"), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Pebbleton · %s mode · http://%s", mode, address)
	log.Fatal(server.ListenAndServe())
}
