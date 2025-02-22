package main

import (
	"github.com/dkhvan-dev/product-service/internal/api"
	"github.com/dkhvan-dev/product-service/internal/database"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/dkhvan-dev/web-commons/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
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

	if err := database.InitDB(); err != nil {
		errors.HandleInternalError(err)
	}

	newApi := api.InitAPI()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.ErrorHandler())
	newApi.AddRoutes(router)

	router.Run(":" + os.Getenv("SERVER_PORT"))
}
