package genesis

import "errors"

var (
	errAccountAlreadyExists   = errors.New("account already exists")
	errPredeployAlreadyExists = errors.New("predeploy already exists")
)
