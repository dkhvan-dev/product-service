package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var DB *pgxpool.Pool

func InitDB() {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		panic(err) // todo: заменить на handler errors
	}

	DB = conn
	migrateSql()
}

func migrateSql() {
	migrationsPath := "file://internal/migrations"

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err) // todo: заменить на handler errors
	}

	dbName := os.Getenv("DATABASE_NAME")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, dbName, driver)
	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err) // todo: заменить на handler errors
	}
}

func StartTransaction(txFunc func(*pgx.Tx) error) error {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	return txFunc(&tx)
}

func StartReadTransaction(txFunc func(*pgx.Tx) error) error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		panic(err) // todo: заменить на handler errors
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	return txFunc(&tx)
}
