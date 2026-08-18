package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"example.com/calibration-vault/internal/httpapi"
	"example.com/calibration-vault/internal/service"
	"example.com/calibration-vault/internal/store"
)

func main() {
	if len(os.Args) > 1 {
		os.Exit(runCLI(os.Args[1:]))
	}
	address := os.Getenv("CALIBRATION_ADDR")
	if address == "" {
		address = "127.0.0.1:18080"
	}
	repository := store.NewMemory()
	application := service.New(repository, service.SystemClock{})
	server := &http.Server{
		Addr:              address,
		Handler:           httpapi.New(application),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("calibration-vault listening on %s", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
