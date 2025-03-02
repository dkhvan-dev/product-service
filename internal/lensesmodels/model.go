package lensesmodels

import (
	"time"
)

type LensModelView struct {
	Id          int       `json:"id" db:"id"`
	CreatedBy   int       `json:"createdBy" db:"created_by"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Name        string    `json:"name" db:"name"`
	IsAvailable bool      `json:"isAvailable" db:"is_available"`
}

func (l *LensModelView) TableName() string {
	return "lenses_models"
}
