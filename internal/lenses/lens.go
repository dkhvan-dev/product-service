package lenses

import (
	"github.com/dkhvan-dev/product-service/internal/common"
	"github.com/dkhvan-dev/product-service/internal/lenses/prices"
	"github.com/dkhvan-dev/product-service/internal/lensesmodels"
)

type LensEntity struct {
	common.AuditEntity
	Name            string                        `json:"name" db:"name"`
	Model           lensesmodels.LensModelEntity  `json:"models" db:"model_id"`
	Brand           string                        `json:"brand" db:"brand"`
	Description     *string                       `json:"description" db:"description"`
	Price           *prices.ActualLensPriceEntity `json:"actualPrice" db:"actual_price"`
	Color           string                        `json:"color" db:"color"`
	OpticalPower    float32                       `json:"opticalPower" db:"optical_power"`
	Diameter        float32                       `json:"diameter" db:"diameter"`
	CurvatureRadius float32                       `json:"curvatureRadius" db:"curvature_radius"`
	IsAvailable     bool                          `json:"isAvailable" db:"is_available"`
}

// DTOs

type LensUpsert struct {
	CreatedBy       int     `db:"created_by"`
	Name            string  `json:"name" db:"name"`
	ModelId         int     `json:"modelId" db:"model_id"`
	Brand           string  `json:"brand" db:"brand"`
	Description     *string `json:"description" db:"description"`
	Price           float64 `json:"price"`
	Color           string  `json:"color" db:"color"`
	OpticalPower    float32 `json:"opticalPower" db:"optical_power"`
	Diameter        float32 `json:"diameter" db:"diameter"`
	CurvatureRadius float32 `json:"curvatureRadius" db:"curvature_radius"`
}
