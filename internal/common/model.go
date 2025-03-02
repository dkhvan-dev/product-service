package common

import (
	"time"
)

type AuditEntity struct {
	Id        int        `json:"id" db:"id"`
	CreatedBy int        `db:"created_by"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	IsDeleted bool       `json:"isDeleted" db:"is_deleted"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type TableNameInterface interface {
	TableName() string
}
