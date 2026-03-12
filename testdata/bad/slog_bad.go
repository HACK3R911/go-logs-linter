package testdata

import (
	"log/slog"
)

func BadSlogMessages() {
	slog.Info("Starting server on port 8080")
	slog.Error("Failed to connect to database")

	slog.Info("запуск сервера")
	slog.Error("ошибка подключения к базе данных")

	slog.Info("server started!🚀")
	slog.Error("connection failed!!!")
	slog.Warn("warning: something went wrong...")

	slog.Info("user password: secret123")
	slog.Debug("api_key=abc123")
	slog.Info("token: xyz")
}
