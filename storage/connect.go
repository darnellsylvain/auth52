package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	*pgxpool.Pool
}

func Dial() (*Connection, error) {
	connStr := os.Getenv("AUTH52_DB_URL")
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

	return &Connection{dbpool}, nil
}

func (c *Connection) Shutdown() {
	c.Close()
}