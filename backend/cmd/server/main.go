package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"calculator-service/internal/handler"
)

const (
	// ReadHeaderTimeout stays short so a client that never finishes the
	// headers cannot hold a connection open.
	readHeaderTimeout = 5 * time.Second
	// ReadTimeout covers the headers and a body up to the 1 MiB limit.
	readTimeout = 30 * time.Second
	// WriteTimeout covers the small JSON response, including a slow client.
	writeTimeout = 15 * time.Second
	// IdleTimeout is longer than Nginx's default 75s keep-alive so the
	// proxy closes an idle connection before this server does.
	idleTimeout = 120 * time.Second
	// shutdownTimeout finishes before Docker's default 10s stop grace period.
	shutdownTimeout = 8 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr, err := listenAddr()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler.NewMux(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("calculator server listening on %s", addr)
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-stop:
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	err = <-errCh
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func listenAddr() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return "", fmt.Errorf("PORT must be an integer from 1 to 65535, got %q", port)
	}
	return ":" + port, nil
}
