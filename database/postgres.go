package database

import (
	"context"
	"log"
	"prestasi_mhs/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PG *pgxpool.Pool

func ConnectPostgres() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, config.C.PostgresDSN)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	PG = pool
	log.Println("PostgreSQL connected")
	return nil
}
