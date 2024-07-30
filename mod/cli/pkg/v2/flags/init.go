package flags

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"
	// FlagSeed defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"
)

const (
	DefaultOverwrite = false
	DefaultRecover   = false
)

const (
	OverwriteDescription = "Overwrite an existing genesis JSON file"
	RecoverDescription   = "Recover a private validator key from a seed"
)
