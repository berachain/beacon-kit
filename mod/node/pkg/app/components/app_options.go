package components

// This is a struct containing all options determined
// solely via cli flags (not available in our config).
//
// This is useful for options that are not meant to be
// persisted to the config file while not coupling the
// user to the format of the input of the info.

// To use this struct, we must supply it in the depinject process.

type AppOptions struct {
	HomeDir string
}
