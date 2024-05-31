package cmd

import "errors"

var (
	errInvalidEthGenesisFormat = errors.New("invalid eth genesis format")
	errPredeployFlagsLength    = errors.New("predeploy flags must have the same length")
	errAccountFlagsLength      = errors.New("account flags must have the same length")

	errInvalidAccountBalance = errors.New("invalid account balance")
)
