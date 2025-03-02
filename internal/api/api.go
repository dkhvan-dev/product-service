package api

import (
	lensModelApi "github.com/dkhvan-dev/product-service/internal/lensesmodels/api"
)

type API struct {
	lensModelApi *lensModelApi.API
}

func InitAPI() *API {
	return &API{
		lensModelApi: lensModelApi.InitAPI(),
	}
}
