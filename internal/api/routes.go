package api

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dkhvan-dev/web-commons/middlewares"
	"github.com/gin-gonic/gin"
)

func (a *API) AddRoutes(r *gin.Engine, srv *handler.Server) {
	r.POST("/query", middlewares.GraphQLMiddleware(srv))

	group := r.Group("/api/v1")
	addLensModelRoutes(group, a)
	addProductRoutes(group, a)
}

func addLensModelRoutes(r *gin.RouterGroup, a *API) {
	group := r.Group("/lenses-models")
	group.POST("", a.lensModelApi.SaveLensModelHandler)
}

func addProductRoutes(r *gin.RouterGroup, a *API) {
	group := r.Group("/products")
	group.POST("", a.CreateProductHandler)
	group.PUT("/:id", a.UpdateProductHandler)
	//group.GET("/", a.FindAllProductHandler)
	group.DELETE("/:id", a.DeleteProductHandler)
}
