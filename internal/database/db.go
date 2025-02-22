package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dkhvan-dev/web-commons/config"
	customErrors "github.com/dkhvan-dev/web-commons/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var DB *sqlx.DB

func InitDB() error {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	if err := db.Ping(); err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	DB = db
	if err := migrateSql(); err != nil {
		return err
	}

	return nil
}

func migrateSql() error {
	migrationsPath := "file://internal/migrations"
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})

	if err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	dbName := os.Getenv("DATABASE_NAME")
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, dbName, driver)

	if err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	ver, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		config.Logger.Error(err.Error())
		return err
	}

	if dirty {
		if err := m.Force(int(ver) - 1); err != nil {
			config.Logger.Error(err.Error())
			return err
		}

		if err := m.Down(); err != nil {
			config.Logger.Error("Failed to rollback dirty migration: " + err.Error())
			return err
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		config.Logger.Error(err.Error())
		return err
	}

	return nil
}

func StartTransaction(txFunc func(tx *sqlx.Tx) *customErrors.CustomError) *customErrors.CustomError {
	tx, err := DB.Beginx()
	if err != nil {
		config.Logger.Error(err.Error())
		return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	var txErr *customErrors.CustomError
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if txErr != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				config.Logger.Error("Failed to commit transaction", zap.Error(commitErr))
			}
		}
	}()

	txErr = txFunc(tx)
	return txErr
}

func StartReadTransaction(txFunc func(tx *sqlx.Tx) *customErrors.CustomError) *customErrors.CustomError {
	ctx := context.Background()
	tx, err := DB.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})

	if err != nil {
		config.Logger.Error(err.Error())
		return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return txFunc(tx)
}
