package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sudo-init-do/okies_core/pkg/db"
	"github.com/sudo-init-do/okies_core/pkg/logger"
	"github.com/sudo-init-do/okies_core/pkg/response"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		logger.Error.Println("No .env file found, using system env")
	}

	// Load port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect DB
	if err := db.Connect(); err != nil {
		logger.Error.Fatalf("failed to connect database: %v", err)
	}

	// Router
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.Write(w, http.StatusOK, "okies_core is healthy ", nil)
	})

	fmt.Println("okies_core running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
