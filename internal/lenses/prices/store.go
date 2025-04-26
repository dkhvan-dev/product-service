package prices

import (
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

func Create(sc mongo.SessionContext, price primitive.Decimal128, lensId primitive.ObjectID) *errors.CustomError {
	newPrice := New()
	newPrice.Price = price
	newPrice.LensId = lensId

	if _, err := database.DB.Collection("lenses_price_history").InsertOne(sc, newPrice); err != nil {
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}

func DeleteByLensId(sc mongo.SessionContext, lensId primitive.ObjectID) *errors.CustomError {
	collection := database.DB.Collection("lenses_price_history")

	if _, err := collection.DeleteMany(sc, bson.M{"lensId": lensId}); err != nil {
		config.Logger.Error("Failed to delete lens price history", zap.Any("lensId", lensId), zap.Error(err))
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}
