package api

import (
	"github.com/dkhvan-dev/product-service/internal/factory"
	"github.com/dkhvan-dev/product-service/internal/products"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (a *API) CreateProductHandler(ctx *gin.Context) {
	var request products.ProductUpsert

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Set("error", errors.NewCustomError("INTERNAL", http.StatusInternalServerError, ctx))
		return
	}

	if request.Type == "" {
		ctx.Set("error", errors.BadRequestError("PRODUCT_MISSING_TYPE", ctx))
		return
	}

	productStore, factoryErr := factory.ProductStoreFactory.Get(request.Type)
	if factoryErr != nil {
		ctx.Set("error", factoryErr)
		return
	}

	if err := productStore.Create(request.Input); err != nil {
		ctx.Set("error", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (a *API) UpdateProductHandler(ctx *gin.Context) {
	var request products.ProductUpsert

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Set("error", errors.NewCustomError("INTERNAL", http.StatusInternalServerError, ctx))
		return
	}

	if request.Type == "" {
		ctx.Set("error", errors.BadRequestError("PRODUCT_MISSING_TYPE", ctx))
		return
	}

	productStore, factoryErr := factory.ProductStoreFactory.Get(request.Type)
	if factoryErr != nil {
		ctx.Set("error", factoryErr)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Set("error", errors.BadRequestError("INVALID_PARAM_TYPE", ctx))
		return
	}

	if err := productStore.Update(id, request.Input); err != nil {
		ctx.Set("error", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (a *API) DeleteProductHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Set("error", err)
		return
	}

	productType := ctx.Query("productType")

	if productType == "" {
		ctx.Set("error", errors.BadRequestError("PRODUCT_MISSING_TYPE", ctx))
		return
	}

	productStore, factoryErr := factory.ProductStoreFactory.Get(productType)
	if factoryErr != nil {
		ctx.Set("error", factoryErr)
		return
	}

	if deleteErr := productStore.Delete(id); deleteErr != nil {
		ctx.Set("error", deleteErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
