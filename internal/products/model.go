package products

import (
	"github.com/dkhvan-dev/product-service/internal/prices"
	"time"
)

type AuditEntity struct {
	Id        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	IsDeleted bool       `json:"isDeleted" db:"is_deleted"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type ProductEntity struct {
	AuditEntity
	Name        string                    `json:"name" db:"name"`
	Description *string                   `json:"description" db:"description"`
	Price       *prices.ActualPriceEntity `json:"actualPrice" db:"actual_price"`
	Quantity    int                       `json:"quantity" db:"quantity"`
}

// DTOs
type ProductCreate struct {
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity" db:"quantity"`
}
