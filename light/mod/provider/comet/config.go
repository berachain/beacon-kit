package comet

// import (
// 	"time"

// 	"github.com/cometbft/cometbft/libs/log"
// )

// type Config struct {
// 	Logger             log.Logger
// 	ChainID            string
// 	TrustingPeriod     time.Duration
// 	TrustedHeight      int64
// 	TrustedHash        []byte
// 	TrustLevel         string
// 	ListeningAddr      string
// 	Sequential         bool
// 	PrimaryAddr        string
// 	WitnessesAddrs     []string
// 	Directory          string
// 	MaxOpenConnections int
// 	ConfirmationFunc   func(string) bool
// }

// func NewConfig(
// 	logger log.Logger,
// 	chainID string,
// 	trustingPeriod time.Duration,
// 	trustedHeight int64,
// 	trustedHash []byte,
// 	trustLevel string,
// 	listeningAddr string,
// 	sequential bool,
// 	primaryAddr string,
// 	witnessesAddrs []string,
// 	directory string,
// 	maxOpenConnections int,
// 	confirmationFunc func(string) bool,
// ) *Config {
// 	return &Config{
// 		Logger:             logger,
// 		ChainID:            chainID,
// 		TrustingPeriod:     trustingPeriod,
// 		TrustedHeight:      trustedHeight,
// 		TrustedHash:        trustedHash,
// 		TrustLevel:         trustLevel,
// 		ListeningAddr:      listeningAddr,
// 		Sequential:         sequential,
// 		PrimaryAddr:        primaryAddr,
// 		WitnessesAddrs:     witnessesAddrs,
// 		Directory:          directory,
// 		MaxOpenConnections: maxOpenConnections,
// 		ConfirmationFunc:   confirmationFunc,
// 	}
// }
