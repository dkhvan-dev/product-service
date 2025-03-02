package products

import (
	"encoding/json"
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/dkhvan-dev/web-commons/errors"
)

type ProductService interface {
	Create(input json.RawMessage) *errors.CustomError
	Update(id int, input json.RawMessage) *errors.CustomError
	FindAll(pageable model.PageInput, selectedFields []string) (*model.ProductPage, *errors.CustomError)
	Delete(id int) *errors.CustomError
}
