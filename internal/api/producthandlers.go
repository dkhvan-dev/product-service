package api

import (
	"github.com/dkhvan-dev/product-service/internal/factory"
	"github.com/dkhvan-dev/product-service/internal/products"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (a *API) UpsertProductHandler(ctx *gin.Context) {
	var request products.ProductUpsert

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Set("error", errors.NewCustomError("INTERNAL", http.StatusInternalServerError, ctx))
		return
	}

	if request.Category == "" {
		ctx.Set("error", errors.BadRequestError("PRODUCT_MISSING_TYPE", ctx))
		return
	}

	productStore, factoryErr := factory.ProductStoreFactory.Get(request.Category)
	if factoryErr != nil {
		ctx.Set("error", factoryErr)
		return
	}

	if err := productStore.(products.ProductService).Upsert(request.Input); err != nil {
		ctx.Set("error", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (a *API) DeleteProductHandler(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.Set("error", err)
		return
	}

	productCategory := ctx.Query("productCategory")

	if productCategory == "" {
		ctx.Set("error", errors.BadRequestError("PRODUCT_MISSING_TYPE", ctx))
		return
	}

	productStore, factoryErr := factory.ProductStoreFactory.Get(productCategory)
	if factoryErr != nil {
		ctx.Set("error", factoryErr)
		return
	}

	if deleteErr := productStore.(products.ProductService).Delete(id); deleteErr != nil {
		ctx.Set("error", deleteErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
