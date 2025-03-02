package database

import (
	"context"
	"database/sql"
	appConfig "github.com/dkhvan-dev/product-service/internal/config"
	"github.com/dkhvan-dev/web-commons/config"
	customErrors "github.com/dkhvan-dev/web-commons/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
	"net/http"
)

var DB *sqlx.DB

func InitDB(cfg appConfig.AppConfig) error {
	db, err := sqlx.Open("postgres", cfg.ComputeDBUrl())

	if err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	if err := db.Ping(); err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	DB = db
	migrationsPath := "internal/migrations"

	if err := goose.Up(DB.DB, migrationsPath); err != nil {
		config.Logger.Error("Failed to migrate sql", zap.Error(err))
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
