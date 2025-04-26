package service

import (
	"encoding/json"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/product-service/internal/fake"
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/dkhvan-dev/product-service/src/utils"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	commonUtils "github.com/dkhvan-dev/web-commons/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
)

type FakeStore struct {
	collection *mongo.Collection
}

func InitFakeStore() *FakeStore {
	return &FakeStore{
		collection: database.DB.Collection("fake"),
	}
}

func (s *FakeStore) Upsert(input json.RawMessage) *errors.CustomError {
	var lensInput fake.FakeCreate
	if err := json.Unmarshal(input, &lensInput); err != nil {
		return errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	lensInput.Id = primitive.NewObjectID()
	lensInput.Category = model.ProductTypeFake

	return database.StartMongoTransaction(func(sc mongo.SessionContext) *errors.CustomError {
		_, err := s.collection.InsertOne(sc, lensInput)
		if err != nil {
			config.Logger.Error("Failed to insert lenses model", zap.Any("body", lensInput), zap.Error(err))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		return nil
	})
}

func (s *FakeStore) FindAll(pageable model.PageInput, selectedFields []string, search *model.ProductSearchInput) (*model.ProductPage, *errors.CustomError) {
	var models []*model.Fake

	err := database.StartMongoTransaction(func(sc mongo.SessionContext) *errors.CustomError {
		filter := bson.M{}
		filter["category"] = model.ProductTypeFake

		projection := buildProjection(selectedFields)
		opts := options.Find().SetProjection(projection)

		cursor, err := s.collection.Find(sc, filter, opts)
		if err != nil {
			config.Logger.Error("Failed to find fakes", zap.Error(err))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		defer cursor.Close(sc)

		for cursor.Next(sc) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				config.Logger.Error("Failed to decode fakes", zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}

			models = append(models, mapToResponse(doc))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	contentResponse := make([]model.ProductUnion, len(models))
	for i := range models {
		contentResponse[i] = models[i]
	}

	return &model.ProductPage{Content: contentResponse, TotalElements: 100}, nil
}

func mapToResponse(doc bson.M) *model.Fake {
	response := &model.Fake{
		ID:        utils.ToString(doc["_id"]),
		Name:      *commonUtils.ToString(doc["name"]),
		CreatedAt: commonUtils.ToTime(doc["createdAt"]),
		Category:  model.ProductTypeFake,
		CreatedBy: commonUtils.ToInt(doc["createdBy"]),
	}

	if doc["description"] != nil {
		response.Description = commonUtils.ToString(doc["description"])
	}

	return response
}

func buildProjection(fields []string) bson.M {
	projection := bson.M{}

	for _, field := range fields {
		projection[field] = 1
	}

	return projection
}

func (s *FakeStore) Delete(id primitive.ObjectID) *errors.CustomError {
	return nil
}
