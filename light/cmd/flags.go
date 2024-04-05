package cmd

// import "time"

// // Flag Names.
// const (
// 	listenAddr         = "listening-address"
// 	engineURL          = "engine-url"
// 	primaryAddr        = "primary-addr"
// 	primaryEthAddr     = "primary-eth-addr"
// 	witnessAddrsJoined = "witness-addr"
// 	dir                = "dir"
// 	maxOpenConnections = "max-open-connections"

// 	sequential     = "sequential-verification"
// 	trustingPeriod = "trust-period"
// 	trustedHeight  = "height"
// 	trustedHash    = "hash"
// 	trustLevel     = "trust-level"

// 	logLevel = "log-level"
// )

// // Default Flag Values.
// const (
// 	defaultListeningAddress = "tcp://localhost:26658"
// 	defaultEngineURL        = "http://localhost:8551"
// 	defaultPrimaryAddress   = "tcp://localhost:26657"
// 	defaultPrimaryEthAddr   = "http://localhost:8545"
// 	defaultWitnessAddresses = "http://localhost:26657"
// 	defaultDir              = ".tmp/.beacon-light"
// 	defaultMaxOpenConn      = 900
// 	defaultTrustPeriod      = 168 * time.Hour
// 	defaultTrustedHeight    = 1
// 	defaultLogLevel         = "info"
// 	defaultTrustLevel       = "1/3"
// 	defaultSequential       = false
// )

// // Flag Descriptions.
// const (
// 	listenAddrDesc         = "serve the proxy on the given address"
// 	engineURLDesc          = "URL of the engine client to connect to"
// 	primaryAddrDesc        = "connect to a Tendermint node at this address"
// 	primaryEthAddrDesc     = "connect to an Ethereum execution client at this address"
// 	witnessAddrsJoinedDesc = "tendermint nodes to cross-check the primary node, comma-separated"
// 	dirDesc                = "specify the directory"
// 	maxOpenConnectionsDesc = "maximum number of simultaneous connections (including WebSocket)"
// 	trustingPeriodDesc     = `trusting period that headers can be verified within.
// 	Should be significantly less than the unbonding period`
// 	trustedHeightDesc = "Trusted header's height"
// 	trustedHashDesc   = "Trusted header's hash"
// 	logLevelDesc      = "Log level, info or debug (Default: info) "
// 	trustLevelDesc    = "trust level. Must be between 1/3 and 3/3"
// 	sequentialDesc    = `sequential verification.
// 	Verify all headers sequentially as opposed to using skipping verification`
// )
