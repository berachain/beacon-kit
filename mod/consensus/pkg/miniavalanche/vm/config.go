package vm

import (
	"github.com/ava-labs/avalanchego/snow/validators"
)

type Config struct {
	Validators validators.Manager
}
