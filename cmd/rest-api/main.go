package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atharv3221/rest-api/internal/config"
	"github.com/atharv3221/rest-api/internal/http/handlers/student"
	"github.com/atharv3221/rest-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	storage, err := sqlite.New(cfg)
	if err != nil {

		log.Fatal("storage initialization failed", err.Error())
	}

	// setup router
	slog.Info("Setting up routes")

	router := http.NewServeMux()

	router.Handle("GET /api/students/health", student.HealthCheck())

	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	router.HandleFunc("GET /api/students", student.GetList(storage))

	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))

	router.HandleFunc("PUT /api/students", student.Update(storage))

	slog.Info("Routes setup completed")
	// setup server

	slog.Info("Starting server")

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("server startup failed")
		}
	}()

	slog.Info("Sever started successfully at ", slog.String("address", cfg.Addr))

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
