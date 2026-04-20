package main

import (
	"context"
	"log"
	"time"

	"github.com/andygellermen/CEE4AI/internal/config"
	appdb "github.com/andygellermen/CEE4AI/internal/db"
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

	if err := appdb.RunMigrations(ctx, pool, cfg.MigrationsDir); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
