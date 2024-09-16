package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DarRo9/Tenders/internal/config"
	httphandler "github.com/DarRo9/Tenders/internal/handlers/http"
	"github.com/DarRo9/Tenders/internal/repository/postgres"
	"github.com/DarRo9/Tenders/internal/server"
	service "github.com/DarRo9/Tenders/internal/services"
)

func main() {
	log := server.SetupLogger()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	

	if err := server.Migrate(&cfg.PG, log); err != nil {
		log.Fatal(err)
	}

	repo, err := postgres.New(&cfg.PG)
	if err != nil {
		log.Fatalf("pool connection error: %v", err)
	}

	srv := service.New(repo, log)
	handler := httphandler.New(srv, log)

	app := server.New(handler.CreateRoutes(), &cfg.Server)
	go func() {
		log.Infof("start server on %v", cfg.Server.Address)

		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("running http server error: %s", err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	<-done

	if err := app.Shutdown(context.Background()); err != nil {
		log.Errorf("server shutting down error: %s", err)
	}

	log.Info("the server stopped")
}
