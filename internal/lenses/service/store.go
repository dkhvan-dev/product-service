package service

import (
	"encoding/json"
	defaultErrors "errors"
	"fmt"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/product-service/internal/lenses"
	"github.com/dkhvan-dev/product-service/internal/lenses/prices"
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
	"slices"
)

type LensStore struct {
	collection *mongo.Collection
}

func InitLensStore() *LensStore {
	return &LensStore{
		collection: database.DB.Collection("lenses"),
	}
}

func (s *LensStore) Upsert(input json.RawMessage) *errors.CustomError {
	lensInput := lenses.NewLens()
	lensInput.CreatedBy = -10

	if err := json.Unmarshal(input, &lensInput); err != nil {
		return errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	if cmsErr := utils.ValidateCmsValue(lensInput); cmsErr != nil {
		return cmsErr
	}

	return database.StartMongoTransaction(func(sc mongo.SessionContext) *errors.CustomError {
		if err := exists(sc, s.collection, lensInput); err != nil {
			return err
		}

		var entity lenses.LensEntity
		if lensInput.Id != nil {
			filter := bson.M{
				"_id": lensInput.Id,
			}

			if err := s.collection.FindOne(sc, filter).Decode(&entity); err != nil {
				config.Logger.Error("Failed to find lens", zap.Any("lensId", lensInput.Id), zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}

			if _, err := s.collection.ReplaceOne(sc, bson.D{}, lensInput); err != nil {
				config.Logger.Error("Failed to update lens", zap.Any("body", lensInput), zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}
		} else {
			objId := primitive.NewObjectID()
			lensInput.Id = &objId

			if _, err := s.collection.InsertOne(sc, lensInput); err != nil {
				config.Logger.Error("Failed to insert lens", zap.Any("body", lensInput), zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}
		}

		if entity.ActualPrice != lensInput.Price {
			if err := prices.Create(sc, lensInput.Price, *lensInput.Id); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *LensStore) FindAll(pageable model.PageInput, selectedFields []string, search *model.ProductSearchInput) (*model.ProductPage, *errors.CustomError) {
	var lensArr []*model.Lens
	var totalElements int64

	err := database.StartMongoTransaction(func(sc mongo.SessionContext) *errors.CustomError {
		filter := buildFilter(search)
		projection := buildProjection(selectedFields)
		sort := buildSort(pageable.Sort)
		opts := options.Find().
			SetProjection(projection).
			SetSkip(int64(pageable.Page * pageable.Size)).
			SetLimit(int64(pageable.Size)).
			SetSort(sort)

		cursor, err := s.collection.Find(sc, filter, opts)
		if err != nil {
			config.Logger.Error("Failed to find lens models", zap.Error(err))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		defer cursor.Close(sc)

		for cursor.Next(sc) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				config.Logger.Error("Failed to decode lens model", zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}

			lensArr = append(lensArr, mapToResponse(doc))
		}

		if slices.Contains(selectedFields, "totalElements") {
			totalElements, err = s.collection.CountDocuments(sc, bson.D{})
			if err != nil {
				config.Logger.Error("Failed to accurate count lenses", zap.Error(err))
				return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	contentResponse := make([]model.ProductUnion, len(lensArr))
	for i := range lensArr {
		contentResponse[i] = lensArr[i]
	}

	return &model.ProductPage{Content: contentResponse, TotalElements: int(totalElements)}, nil
}

func buildSort(sortInputs []*model.SortInput) bson.D {
	sort := bson.D{}

	if len(sortInputs) == 0 {
		sort = append(sort, bson.E{Key: "id", Value: 1})
	}

	for _, s := range sortInputs {
		if s == nil {
			continue
		}

		direction := 1
		if s.Direction == model.SortDirectionDesc {
			direction = -1
		}

		sort = append(sort, bson.E{Key: s.Field, Value: direction})
	}

	return sort
}

func buildFilter(search *model.ProductSearchInput) bson.M {
	filter := bson.M{}
	filter["category"] = "LENSES"

	if search == nil {
		return filter
	}

	buildCommonFilters(filter, search)

	if search.LensFilters != nil {
		lensFilters := search.LensFilters
		if lensFilters.Color != nil && len(lensFilters.Color) != 0 {
			filter["color"] = bson.M{"$in": lensFilters.Color}
		}

		if lensFilters.Brand != nil && len(lensFilters.Brand) != 0 {
			filter["brand"] = bson.M{"$in": lensFilters.Brand}
		}

		if lensFilters.CurvatureRadius != nil {
			filter["curvatureRadius"], _ = primitive.ParseDecimal128(*lensFilters.CurvatureRadius)
		}

		if lensFilters.Diameter != nil {
			filter["diameter"], _ = primitive.ParseDecimal128(*lensFilters.Diameter)
		}

		if lensFilters.HasZeroOpticalPower != nil {
			filter["hasZeroOpticalPower"] = *lensFilters.HasZeroOpticalPower
		}

		if lensFilters.OpticalPower != nil {
			opticalPowerDecimal, err := primitive.ParseDecimal128(*lensFilters.OpticalPower)
			if err == nil {
				stepTolerance, _ := primitive.ParseDecimal128("0.000001")
				expr := bson.M{
					"$and": bson.A{
						bson.M{"$lte": bson.A{opticalPowerDecimal, "$maxOpticalPower"}},
						bson.M{"$gte": bson.A{opticalPowerDecimal, "$minOpticalPower"}},
						bson.M{
							"$lte": bson.A{
								bson.M{"$abs": bson.M{
									"$mod": bson.A{
										bson.M{"$subtract": bson.A{opticalPowerDecimal, "$minOpticalPower"}},
										"$opticalPowerStep",
									},
								}},
								stepTolerance,
							},
						},
					},
				}
				filter["$expr"] = expr
			}
		}
	}

	return filter
}

func buildCommonFilters(filter bson.M, search *model.ProductSearchInput) {
	if search.IsDeleted != nil {
		filter["isDeleted"] = *search.IsDeleted
	}

	if search.Name != nil {
		filter["name"] = bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", *search.Name),
			"$options": "i",
		}
	}

	if search.IsAvailable != nil {
		filter["isAvailable"] = *search.IsAvailable
	}

	if search.PriceFrom != nil || search.PriceTo != nil {
		priceFilter := bson.M{}

		if search.PriceFrom != nil {
			priceFilter["$gte"], _ = primitive.ParseDecimal128(*search.PriceFrom)
		}

		if search.PriceTo != nil {
			priceFilter["$lte"], _ = primitive.ParseDecimal128(*search.PriceTo)
		}

		filter["actualPrice"] = priceFilter
	}

	if search.Quantity != nil {
		filter["quantity"] = *search.Quantity
	}
}

func buildProjection(fields []string) bson.M {
	projection := bson.M{}

	for _, field := range fields {
		projection[field] = 1
	}

	return projection
}

func mapToResponse(doc bson.M) *model.Lens {
	return &model.Lens{
		ID:                  utils.ToString(doc["_id"]),
		Name:                *commonUtils.ToString(doc["name"]),
		CreatedAt:           commonUtils.ToTime(doc["createdAt"]),
		CreatedBy:           commonUtils.ToInt(doc["createdBy"]),
		Color:               *commonUtils.ToString(doc["color"]),
		Category:            model.ProductTypeLenses,
		MinOpticalPower:     utils.ToString(doc["minOpticalPower"]),
		MaxOpticalPower:     utils.ToString(doc["maxOpticalPower"]),
		OpticalPowerStep:    utils.ToString(doc["opticalPowerStep"]),
		Diameter:            utils.ToString(doc["diameter"]),
		CurvatureRadius:     utils.ToString(doc["curvatureRadius"]),
		HasZeroOpticalPower: commonUtils.ToBool(doc["hasZeroOpticalPower"]),
		ActualPrice:         utils.ToString(doc["actualPrice"]),
		Quantity:            commonUtils.ToInt(doc["quantity"]),
		SalesQuantity:       commonUtils.ToInt(doc["salesQuantity"]),
		IsAvailable:         commonUtils.ToBool(doc["isAvailable"]),
		IsDeleted:           commonUtils.ToBool(doc["isDeleted"]),
	}
}

func (s *LensStore) Delete(id primitive.ObjectID) *errors.CustomError {
	return database.StartMongoTransaction(func(sc mongo.SessionContext) *errors.CustomError {
		if err := existsById(sc, s.collection, id); err != nil {
			return err
		}

		if _, err := s.collection.DeleteOne(sc, bson.M{"_id": id}); err != nil {
			config.Logger.Error("Failed to delete lens", zap.Any("lensId", id), zap.Error(err))
			return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		if err := prices.DeleteByLensId(sc, id); err != nil {
			return err
		}

		return nil
	})
}

func exists(sc mongo.SessionContext, collection *mongo.Collection, lensInput lenses.LensUpsert) *errors.CustomError {
	filter := bson.M{
		"name":                lensInput.Name,
		"color":               lensInput.Color,
		"brand":               lensInput.Brand,
		"minOpticalPower":     lensInput.MinOpticalPower,
		"maxOpticalPower":     lensInput.MaxOpticalPower,
		"opticalPowerStep":    lensInput.OpticalPowerStep,
		"curvatureRadius":     lensInput.CurvatureRadius,
		"diameter":            lensInput.Diameter,
		"hasZeroOpticalPower": true,
	}

	if lensInput.Id != nil {
		filter["_id"] = bson.M{"$ne": *lensInput.Id}
	}

	err := collection.FindOne(sc, filter).Err()
	if err == nil {
		config.Logger.Error("Lens already exists", zap.Any("req_body", filter))
		return errors.BadRequestError("LENS_ALREADY_EXISTS", nil)
	} else if !defaultErrors.Is(err, mongo.ErrNoDocuments) {
		config.Logger.Error("Failed to find lens", zap.Any("body", filter), zap.Error(err))
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}

func existsById(sc mongo.SessionContext, collection *mongo.Collection, id primitive.ObjectID) *errors.CustomError {
	filter := bson.M{
		"_id": id,
	}

	err := collection.FindOne(sc, filter).Err()
	if err != nil && !defaultErrors.Is(err, mongo.ErrNoDocuments) {
		config.Logger.Error("Failed to find lens", zap.Any("body", filter), zap.Error(err))
		return errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	return nil
}
