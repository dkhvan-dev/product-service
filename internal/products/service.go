package products

import (
	"encoding/json"
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/dkhvan-dev/web-commons/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	Upsert(input json.RawMessage) *errors.CustomError
	FindAll(pageable model.PageInput, selectedFields []string, search *model.ProductSearchInput) (*model.ProductPage, *errors.CustomError)
	Delete(id primitive.ObjectID) *errors.CustomError
}

type ProductBestSellerService interface {
	BestSellers() ([]*model.ProductBestSeller, *errors.CustomError)
}
