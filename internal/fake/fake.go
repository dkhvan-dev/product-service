package fake

import (
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FakeView struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedBy     *int               `json:"createdBy" db:"created_by"`
	CreatedAt     *time.Time         `json:"createdAt" db:"created_at"`
	Name          *string            `json:"name" db:"name"`
	Description   *string            `json:"description" db:"description"`
	ActualPrice   *FakeHistoryView   `json:"actualPrice" db:"actual_price_id" isParent:"true"`
	SalesQuantity *int               `json:"salesQuantity" db:"sales_quantity"`
}

type FakeCreate struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description *string            `json:"description" bson:"description"`
	Category    model.ProductType  `bson:"category,omitempty"`
}

func (l FakeView) TableName() string {
	return "fake"
}

type FakeHistoryView struct {
	Id        int             `json:"id" db:"id"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	Price     decimal.Decimal `json:"price" db:"price"`
}

func (l *FakeHistoryView) TableName() string {
	return "fake_price_history"
}
