package factory

import (
	"github.com/dkhvan-dev/product-service/internal/lenses/service"
	"github.com/dkhvan-dev/product-service/internal/products"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
)

var ProductStoreFactory *StoreFactory

type StoreFactory struct {
	storeMap map[string]products.ProductService
}

func (s *StoreFactory) Register(productType string, store products.ProductService) {
	s.storeMap[productType] = store
}

func (s *StoreFactory) Get(productType string) (products.ProductService, *errors.CustomError) {
	store, exists := s.storeMap[productType]
	if !exists {
		config.Logger.Error("Invalid product type " + productType)
		return nil, errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	return store, nil
}

func InitStoreFactory() {
	factory := &StoreFactory{storeMap: make(map[string]products.ProductService)}
	factory.Register("LENSES", &service.LensStore{})

	ProductStoreFactory = factory
}
