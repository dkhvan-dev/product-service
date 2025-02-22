package service

import (
	"encoding/json"
	"fmt"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/product-service/internal/lenses"
	"github.com/dkhvan-dev/product-service/internal/lenses/prices"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
)

type LensStore struct {
}

func (s *LensStore) Create(input json.RawMessage) *errors.CustomError {
	var lensInput lenses.LensUpsert
	if err := json.Unmarshal(input, &lensInput); err != nil {
		return errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	// TODO: refactor CreatedBy
	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		lensInput.CreatedBy = -10

		query := "select exists(select 1 from lenses_models where id = $1)"
		var existsModel bool

		tx.QueryRow(query, lensInput.ModelId).Scan(&existsModel)
		if !existsModel {
			config.Logger.Error("Lens models not found", zap.Int("model_id", lensInput.ModelId))
			return errors.BadRequestError("LENS_MODEL_NOT_FOUND", nil)
		}

		query = `
			insert into lenses (created_by, name, model_id, brand, description, 
			                    color, optical_power, diameter, curvature_radius)
			
			values (:created_by, :name, :model_id, :brand, :description, :color, 
			        :optical_power, :diameter, :curvature_radius)
			returning id
		`

		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			config.Logger.Error("Failed prepare statement")
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		var lensId int
		if err := stmt.Get(&lensId, lensInput); err != nil {
			config.Logger.Error("Failed create lens", zap.String("db", err.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		actualPriceId, createActualPriceErr := prices.CreateActualPrice(tx, lensInput.Price, lensId)
		if createActualPriceErr != nil {
			return createActualPriceErr
		}

		if updateErr := assignActualPrice(tx, *actualPriceId, lensId); updateErr != nil {
			return updateErr
		}

		return nil
	})
}

func (s *LensStore) Update(id int, input json.RawMessage) *errors.CustomError {
	var lensInput lenses.LensUpsert
	if err := json.Unmarshal(input, &lensInput); err != nil {
		return errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		if !exists(id) {
			config.Logger.Error("Lens not found", zap.Int("lens_id", id))
			return errors.NotFoundError("LENS_NOT_FOUND", nil)
		}

		var existsModel bool
		tx.QueryRow("select exists(select 1 from lenses_models where id = $1)", lensInput.ModelId).Scan(&existsModel)
		if !existsModel {
			config.Logger.Error("Lens models not found", zap.Int("model_id", lensInput.ModelId))
			return errors.BadRequestError("LENS_MODEL_NOT_FOUND", nil)
		}

		query := "update lenses set"

		var existsPrice bool
		tx.Get(&existsPrice, "select exists(select 1 from lenses_price_history where lens_id = $1 and price = $2)", id, lensInput.Price)
		if !existsPrice {
			actualPriceId, createActualPriceErr := prices.CreateActualPrice(tx, lensInput.Price, id)
			if createActualPriceErr != nil {
				return createActualPriceErr
			}

			query += fmt.Sprintf(" actual_price_id = %d,", *actualPriceId)
		}

		query += `
			updated_at = now(), name = :name, model_id = :model_id, brand = :brand, description = :description, 
			color = :color, optical_power = :optical_power, diameter = :diameter, curvature_radius = :curvature_radius
			where id = %d
		`

		if _, err := tx.NamedExec(fmt.Sprintf(query, id), &lensInput); err != nil {
			config.Logger.Error("Failed update lens", zap.String("db", err.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		return nil
	})
}

func assignActualPrice(tx *sqlx.Tx, actualPriceId, lensId int) *errors.CustomError {
	_, err := tx.Exec("update lenses set actual_price_id = $1 where id = $2", actualPriceId, lensId)
	if err != nil {
		config.Logger.Error("Failed update actual lens price", zap.String("db", err.Error()))
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}

func exists(id int) bool {
	var entityExists bool
	database.DB.QueryRow("select exists(select 1 from lenses where id = $1)", id).Scan(&entityExists)
	return entityExists
}
