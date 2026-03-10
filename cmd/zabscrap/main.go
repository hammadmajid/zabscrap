package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zabscrap/internal/api"
	"zabscrap/internal/app"
	"zabscrap/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("env PORT is not defined")
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	handler := api.NewHandler(logger)

	zabscrap := app.NewApp(logger, handler)

	r := router.SetupRoutes(zabscrap)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: time.Minute,
	}

	zabscrap.Logger.Printf("Listening on http://localhost:%s", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zabscrap.Logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zabscrap.Logger.Println("Server exited")
}
