package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/sudo-init-do/okies_core/internal/auth"
	"github.com/sudo-init-do/okies_core/internal/wallet"
	"github.com/sudo-init-do/okies_core/middleware"
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

	// ---------- Global Middleware ----------
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// ---------- Auth Routes ----------
	authHandler := auth.NewHandler()
	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)

	// ---------- Wallet Routes (protected by JWT) ----------
	walletHandler := wallet.NewHandler()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/wallet/balance", walletHandler.GetBalance)
		r.Post("/wallet/fund", walletHandler.Fund)
		r.Post("/wallet/withdraw", walletHandler.Withdraw)
	})

	// ---------- Health Check ----------
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.Write(w, http.StatusOK, "okies_core is healthy", nil)
	})

	fmt.Println("okies_core running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
