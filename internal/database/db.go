package database

import (
	"context"
	"errors"
	appConfig "github.com/dkhvan-dev/product-service/internal/config"
	"github.com/dkhvan-dev/web-commons/config"
	customErrors "github.com/dkhvan-dev/web-commons/errors"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var DB *mongo.Database

func InitMongoDB(cfg appConfig.AppConfig) error {
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			config.Logger.Info(e.Command.String())
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			config.Logger.Info(e.Reply.String())
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			config.Logger.Error(e.Failure)
		},
	}

	clientOpts := options.Client().ApplyURI(cfg.MONGO_URI).SetMonitor(monitor)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		config.Logger.Error(err.Error())
		return err
	}

	DB = client.Database(cfg.MONGO_DATABASE)
	return nil
}

func StartMongoTransaction(txFunc func(sc mongo.SessionContext) *customErrors.CustomError) *customErrors.CustomError {
	client := DB.Client()
	session, err := client.StartSession()

	if err != nil {
		config.Logger.Error("Failed to start Mongo session", zap.Error(err))
		return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
	}

	defer session.EndSession(context.Background())

	transactionErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			config.Logger.Error("Failed to start Mongo transaction", zap.Error(err))
			return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		txErr := txFunc(sc)
		if txErr != nil {
			session.AbortTransaction(sc)
			return txErr
		}

		if err := session.CommitTransaction(sc); err != nil {
			config.Logger.Error("Failed to commit Mongo transaction", zap.Error(err))
			return customErrors.NewCustomError("INTERNAL", http.StatusInternalServerError, nil)
		}

		return nil
	})

	if transactionErr != nil {
		var customErr *customErrors.CustomError
		errors.As(transactionErr, &customErr)
		statusCode := customErr.Code
		key := customErr.Key

		return customErrors.NewCustomError(key, statusCode, nil)
	}

	return nil
}
