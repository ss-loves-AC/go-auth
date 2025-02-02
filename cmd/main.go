package main

import (
	"context"
	"go-auth/handler"
	"go-auth/router"
	"go-auth/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := storage.NewDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	cache, err := storage.NewCache()
	if err != nil {
		log.Fatalf("Cache initialization failed: %v", err)
	}

	authHandler := handler.NewAuthHandler(db, cache)
	tokenHandler := handler.NewTokenHandler(db, cache)

	r := router.NewRouter(authHandler, tokenHandler)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-signalChan
	log.Println("Received termination signal, shutting down...")

	log.Println("Server shutdown complete.")
}
