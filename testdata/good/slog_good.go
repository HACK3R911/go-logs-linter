package testdata

import (
	"log/slog"
)

var (
	goodCaseMsg1      = "starting server on port 8080"
	goodCaseMsg2      = "failed to connect to database"
	goodEngMsg1       = "starting server"
	goodEngMsg2       = "failed to connect to database"
	goodSpecialMsg1   = "server started"
	goodSpecialMsg2   = "connection failed"
	goodSpecialMsg3   = "something went wrong"
	goodSensitiveMsg1 = "user authentication successful"
	goodSensitiveMsg2 = "api request completed"
	goodSensitiveMsg3 = "session validated"
)

func GoodSlogMessages() {
	//
	slog.Info(goodCaseMsg1)
	slog.Error(goodCaseMsg2)

	slog.Info(goodEngMsg1)
	slog.Error(goodEngMsg2)

	slog.Info(goodSpecialMsg1)
	slog.Error(goodSpecialMsg2)
	slog.Warn(goodSpecialMsg3)

	slog.Info(goodSensitiveMsg1)
	slog.Debug(goodSensitiveMsg2)
	slog.Info(goodSensitiveMsg3)
}
