package main

import (
	"context"
	"log"
	"time"

	"github.com/andygellermen/CEE4AI/internal/config"
	appdb "github.com/andygellermen/CEE4AI/internal/db"
	"github.com/andygellermen/CEE4AI/internal/importer"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := appdb.NewPool(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := importer.NewService(pool)
	if err := service.SeedDir(ctx, cfg.SeedsDir); err != nil {
		log.Fatal(err)
	}

	log.Println("seed import completed")
}
