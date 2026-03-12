package testdata

import (
	"log/slog"
)

func GoodSlogMessages() {
	// Default tests
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")

	slog.Info("starting server")
	slog.Error("failed to connect to database")

	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")

	slog.Info("user authentication successful")
	slog.Debug("api request completed")
	slog.Info("session validated")
}
