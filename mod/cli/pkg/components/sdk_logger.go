package components

import (
	sdklog "cosmossdk.io/log"
)

type SDKLogger struct {
	*Logger
}

func (l *SDKLogger) With(keyVals ...any) sdklog.Logger {
	logger := l.Logger.With(keyVals...)
	return &SDKLogger{logger}
}

type SDKLoggerInput struct {
	Logger *Logger
}

func ProvideSDKLogger(in SDKLoggerInput) *SDKLogger {
	return &SDKLogger{in.Logger}
}
