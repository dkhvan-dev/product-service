package api

import (
	"github.com/dkhvan-dev/product-service/internal/factory"
	lensModelApi "github.com/dkhvan-dev/product-service/internal/lensesmodels/api"
)

type API struct {
	storeFactory *factory.StoreFactory
	lensModelApi *lensModelApi.API
}

func InitAPI() *API {
	return &API{
		storeFactory: factory.Init(),
		lensModelApi: lensModelApi.InitAPI(),
	}
}
