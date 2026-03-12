package testdata

import (
	"go.uber.org/zap"
)

func BadZapMessages() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	//
	sugar.Info(badCaseMsg1)
	sugar.Error(badCaseMsg2)

	sugar.Info(badEngMsg1)
	sugar.Error(badEngMsg2)

	sugar.Info(badSensitiveMsg1)
	sugar.Error(badSpecialMsg2)
	sugar.Warn(badSpecialMsg3)

	sugar.Info(badSensitiveMsg1)
	sugar.Debug(badSensitiveMsg2)
	sugar.Info(badSensitiveMsg3)
}
