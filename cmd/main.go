package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/caarlos0/env/v11"
	"github.com/dkhvan-dev/product-service/internal/api"
	appConfig "github.com/dkhvan-dev/product-service/internal/config"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/product-service/internal/factory"
	"github.com/dkhvan-dev/product-service/src/graph"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/dkhvan-dev/web-commons/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		config.Logger.Fatal("Failed to load location Asia/Almaty", zap.Error(err))
		errors.HandleInternalError(errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil))
	}

	time.Local = loc

	var cfg appConfig.AppConfig
	if err := env.Parse(&cfg); err != nil {
		config.Logger.Fatal("Failed parsing env variables to config struct: %v", zap.Error(err))
		errors.HandleInternalError(errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil))
	}

	// Init logger
	if err := config.InitLogger(os.Getenv("ENVIRONMENT_LOG")); err != nil {
		errors.HandleInternalError(errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil))
	}

	// Load envs
	if err := godotenv.Load(); err != nil {
		errors.HandleInternalError(errors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil))
	}

	errMsgsPaths := []string{
		os.Getenv("EN_ERRORS_PATH"),
		os.Getenv("RU_ERRORS_PATH"),
		os.Getenv("KK_ERRORS_PATH"),
	}

	if err := config.InitLocalization(errMsgsPaths); err != nil {
		errors.HandleInternalError(err)
	}

	if err := database.InitMongoDB(cfg); err != nil {
		errors.HandleInternalError(err)
	}

	// Init Store Factory
	factory.InitStoreFactory()
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(cfg)}))
	srv.SetErrorPresenter(middlewares.HandleGraphQLError)
	srv.AddTransport(transport.POST{})

	newApi := api.InitAPI()
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	if cfg.ReleaseMode {
		gin.DefaultWriter = io.Discard
	} else {
		srv.Use(extension.Introspection{})
	}

	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.ErrorHandler())
	newApi.AddRoutes(router, srv)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		config.Logger.Fatal("Failed to start server", zap.Error(err))
		errors.HandleInternalError(err)
	}
}
