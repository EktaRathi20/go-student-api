package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"student-api/internal/config"
	"student-api/internal/handlers/student"
	"student-api/internal/storage/sql"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Hello, Students API!")

	// Load the configuration
	configuration := config.MustLoad()

	// Database connection
	_, err := sql.New(*configuration)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database connection established", slog.String("env", configuration.Env), slog.String(("version"), "1.0.0"))

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	server := http.Server{
		Addr:    configuration.Address,
		Handler: router,
	}

	slog.Info("Starting server", slog.String("address", configuration.Address))

	// gracefully handle shutdown
	doneChan := make(chan os.Signal, 1)

	signal.Notify(doneChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	<-doneChan

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", slog.String("error", err.Error()))
	}

	slog.Info("Server gracefully stopped")
}
