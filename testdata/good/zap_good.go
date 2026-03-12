package testdata

import (
	"go.uber.org/zap"
)

func GoodZapMessages() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	//
	sugar.Info("starting server on port 8080")
	sugar.Error("failed to connect to database")

	sugar.Info("starting server")
	sugar.Error("failed to connect to database")

	sugar.Info("server started")
	sugar.Error("connection failed")
	sugar.Warn("something went wrong")

	sugar.Info("user login successful")
	sugar.Debug("api request completed")
	sugar.Info("session validated")
}
