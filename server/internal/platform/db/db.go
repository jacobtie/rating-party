package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func New(un, pw, host, dbName string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", un, pw, host, dbName)
	mysql, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("[db.New] could not open connection to DB: %w", err)
	}
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := mysql.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[db.New] could not complete initial DB ping: %w", err)
	}
	return &DB{mysql}, nil
}
