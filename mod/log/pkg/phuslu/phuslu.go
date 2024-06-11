package phuslu

import (
	"github.com/phuslu/log"
)

type Logger[KeyValT, ImplT any] struct {
	logger *log.Logger
	r      map[string]interface{}
}

func NewLogger[KeyValT, ImplT any](level string) *Logger[KeyValT, ImplT] {
	logger := &log.DefaultLogger
	logger.SetLevel(log.ParseLevel(level))
	return &Logger[KeyValT, ImplT]{
		logger: logger,
	}
}

func (l *Logger[KeyValT, ImplT]) Info(msg string, keyVals ...KeyValT) {
	l.logger.Info().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Warn(msg string, keyVals ...KeyValT) {
	l.logger.Warn().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Error(msg string, keyVals ...KeyValT) {
	l.logger.Error().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Debug(msg string, keyVals ...KeyValT) {
	l.logger.Debug().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Impl() any {
	return l.logger
}

func (l *Logger[KeyValT, ImplT]) With(keyVals ...KeyValT) ImplT {
	newLogger := *l
	r := make(map[string]interface{})
	for _, keyVal := range keyVals {
		r[keyVal.Key()] = keyVal.Value()
	}
	return any(&newLogger).(ImplT)
}
