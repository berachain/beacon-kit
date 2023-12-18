package config

const (
	PolarisConfigTemplate = `
###############################################################################
###                                 Polaris                                 ###
###############################################################################
# General Polaris settings
[polaris]

[polaris.execution-client]
# HTTP url of the execution client JSON-RPC endpoint.
rpc-dial-url = "{{ .Polaris.ExecutionClient.RPCDialURL }}"

# RPC timeout for execution client requests.
rpc-timeout = "{{ .Polaris.ExecutionClient.RPCTimeout }}"

# Number of retries before shutting down consensus client.
rpc-retries = "{{.Polaris.ExecutionClient.RPCRetries}}"

# Path to the execution client JWT-secret
jwt-secret-path = "{{.Polaris.ExecutionClient.JWTSecretPath}}"

[polaris.beacon-config]
# Altair fork epoch
altair-fork-epoch = {{.Polaris.BeaconConfig.AltairForkEpoch}}

# Bellatrix fork epoch
bellatrix-fork-epoch = {{.Polaris.BeaconConfig.BellatrixForkEpoch}}

# Capella fork epoch
capella-fork-epoch = {{.Polaris.BeaconConfig.CapellaForkEpoch}}

# Deneb fork epoch
deneb-fork-epoch = {{.Polaris.BeaconConfig.DenebForkEpoch}}
`
)
