package logger

import (
	"github.com/berachain/beacon-kit/mod/log"
)

type Logger[KeyValT any, ImplT any] struct {
	log.AdvancedLogger[KeyValT, ImplT]
}
