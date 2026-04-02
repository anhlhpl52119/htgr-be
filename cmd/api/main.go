package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tiny-goclean/config"
	"tiny-goclean/internal/database"
	"tiny-goclean/internal/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	fmt.Println("App running on port:", cfg.Server.Port)

	helpers.InitLogger(cfg)
	log := helpers.GetLogger().Sugar()
	defer log.Sync()

	db, err := database.NewConnection(cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("db connection failed: %w", err)
	}

	fmt.Println("** Connect database succesed **")
	fmt.Println("Checking migration..")
	err = database.CheckMigration(db.Client.DB, "migrations/")
	if err != nil {
		log.Fatal(err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Timezone"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, _ := json.MarshalIndent(map[string]string{"status": "ok"}, " ", "")
		w.Write(b)
	})

	r.Mount("/", r)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Infow("server listening", "config", cfg.Server)
		if err := svr.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("server error", "error", err)
			os.Exit(1)
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		log.Errorw("forced shutdown", "error", err)

	}
	log.Infow("server stopped")
}
