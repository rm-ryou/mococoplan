package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rm-ryou/mococoplan/internal/api/router"
	"github.com/rm-ryou/mococoplan/internal/config"
	"github.com/rm-ryou/mococoplan/pkg/mysql"
)

func main() {
	cfg := config.NewConfig()

	dsn := mysql.CreateDSN(cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.Port)
	db, err := mysql.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer db.Close()

	router := router.Setup(db)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	errCh := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// TODO: logging
	select {
	case err := <-errCh:
		log.Printf("HTTP server ListenAndServe: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
}
