package testdata

import (
	"go.uber.org/zap"
)

func BadZapMessages() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	sugar.Info("Starting server on port 8080")
	sugar.Error("Failed to connect to database")

	sugar.Info("запуск сервера")
	sugar.Error("ошибка подключения к базе данных")

	sugar.Info("server started!🚀")
	sugar.Error("connection failed!!!")
	sugar.Warn("warning: something went wrong...")

	sugar.Info("user password: secret123")
	sugar.Debug("api_key=abc123")
	sugar.Info("token: xyz")

	password := "secret123"
	apiKey := "abc123"
	token := "xyz"

	sugar.Infow("user login", "password", password)
	sugar.Debugw("api request", "api_key", apiKey)
	sugar.Infow("token info", "token", token)
}
