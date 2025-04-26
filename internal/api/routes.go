package api

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dkhvan-dev/web-commons/middlewares"
	"github.com/gin-gonic/gin"
)

func (a *API) AddRoutes(r *gin.Engine, srv *handler.Server) {
	r.POST("/graphql", middlewares.GraphQLMiddleware(srv))

	group := r.Group("/api/v1")
	addProductRoutes(group, a)
}

func addProductRoutes(r *gin.RouterGroup, a *API) {
	group := r.Group("/products")
	group.POST("", a.UpsertProductHandler)
	group.DELETE("/:id", a.DeleteProductHandler)
}
