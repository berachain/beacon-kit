package node

type Backend interface {
	GetVersionFromGenesis() (string, error)
}
