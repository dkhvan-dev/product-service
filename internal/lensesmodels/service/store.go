package service

import (
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
)

type LensModelStore struct {
}

func (s *LensModelStore) Save(name string) *errors.CustomError {
	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var existsModel bool
		query := "select exists(select 1 from lenses_models where name = $1)"

		if err := tx.QueryRow(query, name).Scan(&existsModel); err != nil {
			config.Logger.Error("Failed scan existing lens models", zap.String("db", err.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		if !existsModel {
			if _, err := tx.Exec("insert into lenses_models(created_by, name) values ($1, $2)", -10, name); err != nil {
				config.Logger.Error("Failed create lens models", zap.String("db", err.Error()))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}
		}

		return nil
	})
}
