package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(connString string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Postgres")

	return pool, nil
}
