// Command api is the application entrypoint: it loads configuration, opens the
// database, wires the modules and runs the HTTP server with graceful shutdown.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/example/blog/internal/db"
	"github.com/example/blog/internal/modules"
	"github.com/example/blog/internal/platform/auth"
	"github.com/example/blog/internal/platform/config"
	"github.com/example/blog/internal/platform/database"
	"github.com/example/blog/internal/platform/server"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	validate := validator.New(validator.WithRequiredStructEnabled())
	mods := modules.New(modules.Deps{
		Pool:      pool,
		Queries:   db.New(pool),
		Tokens:    auth.NewTokenService(cfg.JWT),
		Validator: validate,
		Config:    cfg,
	})

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.HTTPPort),
		Handler:           server.New(cfg, mods),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
