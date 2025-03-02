package lenses

import (
	"github.com/dkhvan-dev/product-service/internal/lenses/prices"
	"github.com/dkhvan-dev/product-service/internal/lensesmodels"
	"time"
)

type LensUpsert struct {
	Id              *int    `db:"id"`
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
	Quantity        int     `json:"quantity" db:"quantity"`
	ActualPriceId   *int    `db:"actual_price_id"`
}

type LensView struct {
	Id              int                          `json:"id" db:"id"`
	CreatedBy       *int                         `json:"createdBy" db:"created_by"`
	CreatedAt       *time.Time                   `json:"createdAt" db:"created_at"`
	UpdatedAt       *time.Time                   `json:"updatedAt" db:"updated_at"`
	Name            *string                      `json:"name" db:"name"`
	Model           *lensesmodels.LensModelView  `json:"model" db:"model_id" isParent:"true"`
	Brand           *string                      `json:"brand" db:"brand"`
	Description     *string                      `json:"description" db:"description"`
	ActualPrice     *prices.LensPriceHistoryView `json:"actualPrice" db:"actual_price_id" isParent:"true"`
	Color           *string                      `json:"color" db:"color"`
	OpticalPower    *float32                     `json:"opticalPower" db:"optical_power"`
	Diameter        *float32                     `json:"diameter" db:"diameter"`
	CurvatureRadius *float32                     `json:"curvatureRadius" db:"curvature_radius"`
	IsAvailable     *bool                        `json:"isAvailable" db:"is_available"`
	Quantity        *int                         `json:"quantity" db:"quantity"`
}

func (l LensView) TableName() string {
	return "lenses"
}

func AvailableSortFields() []string {
	return []string{"id", "created_by", "created_at", "updated_at", "name", "model.name", "brand", "description",
		"actualPrice.price", "color", "optical_power", "diameter", "curvature_radius", "quantity"}
}
