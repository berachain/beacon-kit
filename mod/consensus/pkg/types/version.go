package types

type ConsensusVersion uint8

const (
	CometBFTConsensus ConsensusVersion = iota
	RollKitConsensus
)
