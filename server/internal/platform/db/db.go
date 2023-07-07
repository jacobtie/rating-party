package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func New(un, pw, host, dbName string) (*DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s", un, pw, host, dbName)
	postgres, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("[db.New] could not open connection to DB: %w", err)
	}
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := postgres.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[db.New] could not complete initial DB ping: %w", err)
	}
	return &DB{postgres}, nil
}

func (db *DB) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("[db.WithTransaction] failed to begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("[db.WithTransaction] failed to rollback transaction: %w", err)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[db.WithTransaction] failed to commit transaction: %w", err)
	}
	return nil
}
