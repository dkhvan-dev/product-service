package prices

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LensPriceHistory struct {
	Id        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	Price     primitive.Decimal128 `json:"price" bson:"price"`
	LensId    primitive.ObjectID   `json:"lensId" bson:"lensId"`
}

func New() LensPriceHistory {
	return LensPriceHistory{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
	}
}

func (l *LensPriceHistory) TableName() string {
	return "lenses_price_history"
}
