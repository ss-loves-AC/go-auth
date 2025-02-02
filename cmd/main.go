package main

import (
	"context"
	"go-auth/handler"
	"go-auth/router"
	"go-auth/storage"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
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
	r := router.NewRouter(authHandler)
	r.Run()

}
