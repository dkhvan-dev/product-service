package products

import (
	"encoding/json"
	"github.com/dkhvan-dev/web-commons/errors"
)

type ProductService interface {
	Create(input json.RawMessage) *errors.CustomError
	Update(id int, input json.RawMessage) *errors.CustomError
}
