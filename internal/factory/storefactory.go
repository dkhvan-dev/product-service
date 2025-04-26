package factory

import (
	service3 "github.com/dkhvan-dev/product-service/internal/fake/service"
	"github.com/dkhvan-dev/product-service/internal/lenses/service"
	"github.com/dkhvan-dev/product-service/internal/products"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
)

var ProductStoreFactory *StoreFactory

type StoreFactory struct {
	StoreMap       map[string]products.ProductService
	BestSellersMap map[string]products.ProductBestSellerService
}

func (s *StoreFactory) Register(productCategory string, store any) {
	s.StoreMap[productCategory] = store.(products.ProductService)
	//s.BestSellersMap[productCategory] = store.(products.ProductBestSellerService)
}

func (s *StoreFactory) Get(productCategory string) (any, *errors.CustomError) {
	store, exists := s.StoreMap[productCategory]
	if !exists {
		config.Logger.Error("Invalid product type " + productCategory)
		return nil, errors.BadRequestError("INVALID_INPUT_BODY", nil)
	}

	return store, nil
}

func InitStoreFactory() {
	factory := &StoreFactory{
		StoreMap: make(map[string]products.ProductService),
		//BestSellersMap: make(map[string]products.ProductBestSellerService),
	}

	factory.Register("LENSES", service.InitLensStore())
	factory.Register("FAKE", service3.InitFakeStore())
	ProductStoreFactory = factory
}
