package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	*pgxpool.Pool
}

func Dial(connStr string) (*Connection, error) {
	if connStr == "" {
		log.Fatal("Database URL is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	log.Println("db connected")

	return &Connection{dbpool}, nil
}
