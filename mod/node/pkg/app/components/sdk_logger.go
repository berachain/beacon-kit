package components

import (
	sdklog "cosmossdk.io/log"
)

type SDKLogger struct {
	*Logger
}

func (l *SDKLogger) With(keyvals ...interface{}) sdklog.Logger {
	return &SDKLogger{l.Logger.With(keyvals...)}
}

func ProvideSDKLogger(logger *Logger) *SDKLogger {
	return &SDKLogger{logger}
}
