package echo

import (
	"io"

	"github.com/berachain/beacon-kit/mod/log"
)

var _ Logger = (*logger)(nil)

type logger struct {
	log.Logger[any]
}

func (l *logger) Output() io.Writer {
	return nil
}

func (l *logger) SetOutput(w io.Writer) {}

func (l *logger) Prefix() string {
	return ""
}

func (l *logger) SetPrefix(p string) {}

func (l *logger) Level() log.Lvl {
	return 0
}

func (l *logger) SetLevel(v log.Lvl) {}

func (l *logger) SetHeader(h string) {}

func (l *logger) Print(i ...interface{}) {}

func (l *logger) Printf(format string, args ...interface{}) {}

func (l *logger) Printj(j log.JSON) {}

func (l *logger) Debug(i ...interface{}) {}

func (l *logger) Debugf(format string, args ...interface{}) {}

func (l *logger) Debugj(j log.JSON) {}

func (l *logger) Info(i ...interface{}) {}

func (l *logger) Infof(format string, args ...interface{}) {}

func (l *logger) Infoj(j log.JSON) {}

func (l *logger) Warn(i ...interface{}) {}

func (l *logger) Warnf(format string, args ...interface{}) {}

func (l *logger) Warnj(j log.JSON) {}

func (l *logger) Error(i ...interface{}) {}

func (l *logger) Errorf(format string, args ...interface{}) {}

func (l *logger) Errorj(j log.JSON) {}

func (l *logger) Fatal(i ...interface{}) {}

func (l *logger) Fatalj(j log.JSON) {}

func (l *logger) Fatalf(format string, args ...interface{}) {}

func (l *logger) Panic(i ...interface{}) {}

func (l *logger) Panicj(j log.JSON) {}

func (l *logger) Panicf(format string, args ...interface{}) {}
