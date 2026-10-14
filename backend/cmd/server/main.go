package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"calculator-service/internal/handler"
)

func main() {
	addr, err := listenAddr()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler.NewMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("calculator server listening on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
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
