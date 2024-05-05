package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/egreb/boilerplate/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func dsn(c config.DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
	)
}

func Connect(ctx context.Context, c config.DBConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dsn(c))
	if err != nil {
		return nil, fmt.Errorf("failed connecting to the db: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed pinging the database: %w", err)
	}

	return db, nil
}

func GetDB(c config.DBConfig) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn(c))

	if err != nil {
		return nil, fmt.Errorf("failed opening db: %w", err)
	}

	return conn, nil
}
