package lensesmodels

import "github.com/dkhvan-dev/product-service/internal/common"

type LensModelEntity struct {
	common.AuditEntity
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"isAvailable" db:"is_available"`
}
