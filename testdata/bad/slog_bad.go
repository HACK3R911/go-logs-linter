package testdata

import (
	"log/slog"
)

var (
	password string
	apiKey   string
	token    string
)

var (
	badCaseMsg1 = "Starting server on port 8080"
	badCaseMsg2 = "Failed to connect to database"

	badEngMsg1       = "запуск сервера"
	badEngMsg2       = "ошибка подключения к базе данных"
	badSpecialMsg1   = "server started!🚀"
	badSpecialMsg2   = "connection failed!!!"
	badSpecialMsg3   = "warning: something went wrong..."
	badSensitiveMsg1 = "user password: " + password
	badSensitiveMsg2 = "api_key=" + apiKey
	badSensitiveMsg3 = "token: " + token
)

func BadSlogMessages() {
	// Default tests
	slog.Info(badCaseMsg1)
	slog.Error(badCaseMsg2)

	slog.Info(badEngMsg1)
	slog.Error(badEngMsg2)

	slog.Info(badSpecialMsg1)
	slog.Error(badSpecialMsg2)
	slog.Warn(badSpecialMsg3)

	slog.Info(badSensitiveMsg1)
	slog.Debug(badSensitiveMsg2)
	slog.Info(badSensitiveMsg3)
}
