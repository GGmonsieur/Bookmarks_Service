package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionData struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	DBName   string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg *ConnectionData) (*DB, error) {
	if cfg == nil {
		return nil, errors.New("connection data is empty")
	}

	connString := fmt.Sprintf("user=%s password='%s' host=%s port=%d dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db == nil || db.Pool == nil {
		return
	}
	db.Pool.Close()
}

func (db *DB) GetPool() *pgxpool.Pool {
	if db == nil {
		return nil
	}
	return db.Pool
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.Pool.Query(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, query, args...)
}
