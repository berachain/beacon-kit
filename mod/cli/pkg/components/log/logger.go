package log

import (
	sdklog "cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log"
)

type SDKLogger[LoggerT log.AdvancedLogger[any, LoggerT]] struct {
	Logger LoggerT
}

func WrapSDKLogger[LoggerT log.AdvancedLogger[any, LoggerT]](
	logger LoggerT,
) *SDKLogger[LoggerT] {
	return &SDKLogger[LoggerT]{
		Logger: logger,
	}
}

func (l *SDKLogger[LoggerT]) Info(msg string, keyVals ...any) {
	l.Logger.Info(msg, keyVals...)
}

func (l *SDKLogger[LoggerT]) Warn(msg string, keyVals ...any) {
	l.Logger.Warn(msg, keyVals...)
}

func (l *SDKLogger[LoggerT]) Error(msg string, keyVals ...any) {
	l.Logger.Error(msg, keyVals...)
}

func (l *SDKLogger[LoggerT]) Debug(msg string, keyVals ...any) {
	l.Logger.Debug(msg, keyVals...)
}

func (l *SDKLogger[LoggerT]) With(keyVals ...any) sdklog.Logger {
	return &SDKLogger[LoggerT]{Logger: l.Logger.With(keyVals...)}
}

func (l *SDKLogger[LoggerT]) Impl() any {
	return l.Logger
}
