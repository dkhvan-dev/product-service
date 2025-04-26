package products

import (
	"github.com/dkhvan-dev/product-service/src/graph/model"
	"github.com/dkhvan-dev/web-commons/errors"
	"slices"
)

type ProductStore struct {
}

func (s *ProductStore) BestSellers(storeMap map[string]ProductBestSellerService) ([]*model.ProductBestSeller, *errors.CustomError) {
	var bestSellers []*model.ProductBestSeller
	for _, v := range storeMap {
		result, err := v.BestSellers()
		if err != nil {
			return nil, err
		}

		bestSellers = append(bestSellers, result...)
	}

	slices.SortFunc(bestSellers, func(a, b *model.ProductBestSeller) int {
		return b.SalesQuantity - a.SalesQuantity
	})

	if len(bestSellers) > 6 {
		return bestSellers[:6], nil
	}

	return bestSellers, nil
}
