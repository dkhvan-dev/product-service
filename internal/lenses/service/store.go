package service

import (
	"encoding/json"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/product-service/internal/lenses"
	"github.com/dkhvan-dev/product-service/internal/lenses/prices"
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/dkhvan-dev/product-service/src/utils"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"slices"
	"strings"
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
		config.QueryLogger(query)
		if !existsModel {
			config.Logger.Error("Lens models not found", zap.Int("model_id", lensInput.ModelId))
			return errors.BadRequestError("LENS_MODEL_NOT_FOUND", nil)
		}

		query = `
			insert into lenses (created_by, name, model_id, brand, description, 
			                    color, optical_power, diameter, curvature_radius, quantity)
			
			values (:created_by, :name, :model_id, :brand, :description, :color, 
			        :optical_power, :diameter, :curvature_radius, :quantity)
			returning id
		`

		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			config.Logger.Error("Failed prepare statement")
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		var lensId int
		config.QueryLogger(query)

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
	lensInput.Id = &id

	if err := json.Unmarshal(input, &lensInput); err != nil {
		return errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		if !exists(id) {
			config.Logger.Error("Lens not found", zap.Int("lens_id", id))
			return errors.NotFoundError("LENS_NOT_FOUND", nil)
		}

		var existsModel bool
		existsModelQuery := "select exists(select 1 from lenses_models where id = $1)"
		tx.QueryRow(existsModelQuery, lensInput.ModelId).Scan(&existsModel)
		config.QueryLogger(existsModelQuery)

		if !existsModel {
			config.Logger.Error("Lens models not found", zap.Int("model_id", lensInput.ModelId))
			return errors.BadRequestError("LENS_MODEL_NOT_FOUND", nil)
		}

		var queryBuilder strings.Builder
		queryBuilder.WriteString("update lenses set")
		queryBuilder.WriteString(`
			updated_at = now(), name = :name, model_id = :model_id, brand = :brand, description = :description, 
			color = :color, optical_power = :optical_power, diameter = :diameter, curvature_radius = :curvature_radius, 
			quantity = :quantity
		`)

		var existsPrice bool
		existingPriceQuery := "select exists(select 1 from lenses_price_history where lens_id = $1 and price = $2)"

		tx.Get(&existsPrice, existingPriceQuery, id, lensInput.Price)
		config.QueryLogger(existingPriceQuery)
		if !existsPrice {
			actualPriceId, createActualPriceErr := prices.CreateActualPrice(tx, lensInput.Price, id)
			if createActualPriceErr != nil {
				return createActualPriceErr
			}

			lensInput.ActualPriceId = actualPriceId
			queryBuilder.WriteString(", actual_price_id = :actual_price_id")
		}

		queryBuilder.WriteString(" where id = :id")
		config.QueryLogger(queryBuilder.String())
		if _, err := tx.NamedExec(queryBuilder.String(), &lensInput); err != nil {
			config.Logger.Error("Failed update lens", zap.String("db", err.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		return nil
	})
}

func assignActualPrice(tx *sqlx.Tx, actualPriceId, lensId int) *errors.CustomError {
	query := "update lenses set actual_price_id = $1 where id = $2"
	_, err := tx.Exec(query, actualPriceId, lensId)
	config.QueryLogger(query)

	if err != nil {
		config.Logger.Error("Failed update actual lens price", zap.String("db", err.Error()))
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}

func exists(id int) bool {
	var entityExists bool
	query := "select exists(select 1 from lenses where id = $1)"
	database.DB.QueryRow(query, id).Scan(&entityExists)
	config.QueryLogger(query)

	return entityExists
}

func (s *LensStore) FindAll(pageable model.PageInput, selectedFields []string) (*model.ProductPage, *errors.CustomError) {
	var lensesContent []lenses.LensView
	var totalElements int

	err := database.StartReadTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		query := buildQuery(selectedFields, pageable)
		selectErr := tx.Select(&lensesContent, query, pageable.Size, pageable.Page)
		config.QueryLogger(query)

		if selectErr != nil {
			config.Logger.Error("Failed search lenses", zap.String("db", selectErr.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		if slices.Contains(selectedFields, "totalElements") {
			totalElementsQuery := "select count(id) from lenses where is_deleted is false"
			config.QueryLogger(totalElementsQuery)

			if err := tx.Get(&totalElements, totalElementsQuery); err != nil {
				config.Logger.Error("Failed calculate count lenses", zap.String("db", err.Error()))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	contentResponse := make([]model.ProductInterface, len(lensesContent))
	for i, lensView := range lensesContent {
		var lens model.Lens
		copier.Copy(&lens, &lensView)
		contentResponse[i] = &lens
	}

	return &model.ProductPage{Content: contentResponse, TotalElements: totalElements}, nil
}

func buildQuery(selectedFields []string, pageable model.PageInput) string {
	var query strings.Builder
	query.WriteString(utils.GenerateSQL[lenses.LensView](selectedFields))
	query.WriteString(" where lenses.is_deleted is false order by ")
	query.WriteString(utils.ToSqlSort[lenses.LensView](pageable.Sort))
	query.WriteString(" limit $1 offset $2")

	return query.String()
}

func (s *LensStore) Delete(id int) *errors.CustomError {
	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		query := "update lenses set updated_at = now(), is_deleted = true, deleted_at = now() where id = $1"
		config.QueryLogger(query)

		if _, err := tx.Exec(query, id); err != nil {
			config.Logger.Error("Failed delete lens", zap.String("db", err.Error()))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		return nil
	})
}
