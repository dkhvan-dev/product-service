package api

import "github.com/dkhvan-dev/product-service/internal/lensesmodels/service"

type API struct {
	store *service.LensModelStore
}

func InitAPI() *API {
	return &API{store: &service.LensModelStore{}}
}
