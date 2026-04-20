package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andygellermen/CEE4AI/internal/config"
	appdb "github.com/andygellermen/CEE4AI/internal/db"
	apphttp "github.com/andygellermen/CEE4AI/internal/http"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
	Router http.Handler
}

func New(cfg config.Config) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := appdb.NewPool(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	router := apphttp.NewRouter(pool)

	return &App{
		Config: cfg,
		DB:     pool,
		Router: router,
	}, nil
}

func (a *App) Close() {
	if a == nil || a.DB == nil {
		return
	}

	a.DB.Close()
}
