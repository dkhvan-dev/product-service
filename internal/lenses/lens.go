package lenses

import (
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LensUpsert struct {
	Id                  *primitive.ObjectID  `bson:"_id,omitempty"`
	CreatedAt           time.Time            `bson:"createdAt,omitempty"`
	CreatedBy           int                  `bson:"createdBy,omitempty"`
	Name                string               `json:"name" bson:"name,omitempty"`
	Category            model.ProductType    `bson:"category,omitempty"`
	Brand               string               `json:"brand" bson:"brand"`
	Description         *string              `json:"description" bson:"description"`
	Price               primitive.Decimal128 `json:"price" bson:"actualPrice"`
	Diameter            primitive.Decimal128 `json:"diameter" bson:"diameter"`
	CurvatureRadius     primitive.Decimal128 `json:"curvatureRadius" bson:"curvatureRadius"`
	Color               string               `json:"color" bson:"color,omitempty" cms:"colors"`
	MinOpticalPower     primitive.Decimal128 `json:"minOpticalPower" bson:"minOpticalPower,omitempty"`
	MaxOpticalPower     primitive.Decimal128 `json:"maxOpticalPower" bson:"maxOpticalPower,omitempty"`
	OpticalPowerStep    primitive.Decimal128 `json:"opticalPowerStep" bson:"opticalPowerStep,omitempty"`
	HasZeroOpticalPower bool                 `json:"hasZeroOpticalPower" bson:"hasZeroOpticalPower,omitempty"`
	IsAvailable         bool                 `bson:"isAvailable"`
	IsDeleted           bool                 `bson:"isDeleted"`
	Quantity            int                  `json:"quantity" bson:"quantity"`
	SalesQuantity       int                  `bson:"salesQuantity"`
}

type LensEntity struct {
	Id                  *primitive.ObjectID  `bson:"_id,omitempty"`
	CreatedAt           time.Time            `bson:"createdAt,omitempty"`
	CreatedBy           int                  `bson:"createdBy,omitempty"`
	Name                string               `json:"name" bson:"name,omitempty"`
	Category            model.ProductType    `bson:"category,omitempty"`
	Brand               string               `json:"brand" bson:"brand"`
	Description         *string              `json:"description" bson:"description"`
	ActualPrice         primitive.Decimal128 `json:"actualPrice" bson:"actualPrice"`
	Diameter            primitive.Decimal128 `json:"diameter" bson:"diameter"`
	CurvatureRadius     primitive.Decimal128 `json:"curvatureRadius" bson:"curvatureRadius"`
	Color               string               `json:"color" bson:"color,omitempty"`
	MinOpticalPower     primitive.Decimal128 `json:"minOpticalPower" bson:"minOpticalPower,omitempty"`
	MaxOpticalPower     primitive.Decimal128 `json:"maxOpticalPower" bson:"maxOpticalPower,omitempty"`
	OpticalPowerStep    primitive.Decimal128 `json:"opticalPowerStep" bson:"opticalPowerStep,omitempty"`
	HasZeroOpticalPower bool                 `json:"hasZeroOpticalPower" bson:"hasZeroOpticalPower,omitempty"`
	IsAvailable         bool                 `bson:"isAvailable"`
	IsDeleted           bool                 `bson:"isDeleted"`
	Quantity            int                  `json:"quantity" bson:"quantity"`
	SalesQuantity       int                  `bson:"salesQuantity"`
}

func NewLens() LensUpsert {
	return LensUpsert{
		CreatedAt:   time.Now(),
		Category:    model.ProductTypeLenses,
		IsAvailable: true,
		IsDeleted:   false,
	}
}

func (l LensUpsert) TableName() string {
	return "lenses"
}
