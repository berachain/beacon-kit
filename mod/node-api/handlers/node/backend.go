package node

type Backend interface {
	GetNodeVersion() (string, error)
}
