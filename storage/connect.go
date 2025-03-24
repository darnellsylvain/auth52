package storage

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	*pgxpool.Pool
}

func Dial() (*Connection, error) {
	os.Setenv("DATABASE_URL", "postgres://auth52:password@localhost/auth52")
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		// Handle this better
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		log.Printf("db ping was not successful: %v", err)
	}

	return &Connection{dbpool}, nil
}

func (c *Connection) Shutdown() {
	c.Close()
}