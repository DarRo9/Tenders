package postgres

import (
	"context"

	"github.com/DarRo9/Tenders/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(cfg *config.PGConfig) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), cfg.Conn)
	if err != nil {
		return nil, err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Postgres{DB: pool}, nil
}
