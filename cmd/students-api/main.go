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

	"github.com/amitpoitrix/students-api/internal/config"
	"github.com/amitpoitrix/students-api/internal/http/handlers/student"
	"github.com/amitpoitrix/students-api/internal/storage/sqlite"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables from .env (optional)
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, continuing with environment variables")
	}
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentLists(storage))
	router.HandleFunc("PATCH /api/students/{id}", student.ModifyById(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Go Server started at", slog.String("address", cfg.Addr))

	/* Way to gracefully shutting down the server */
	/* Creating channel for synchronization as we're using go routines with size 1 */
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	/* Now below is non-blocking code */
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start the server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	/* creating context with 5 sec timeout */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*
		As server.shutdown(ctx) gracefully stops the server but it might fall in infinite loop due to
		waiting for request to complete so using context timeout to stop after specific time
	*/
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
