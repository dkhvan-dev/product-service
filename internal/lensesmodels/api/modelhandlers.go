package api

import (
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *API) SaveLensModelHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.Set("error", errors.BadRequestError("INVALID_INPUT_BODY", ctx))
		return
	}

	if err := a.store.Save(name); err != nil {
		ctx.Set("error", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
