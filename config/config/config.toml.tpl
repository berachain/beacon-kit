# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

###############################################################################
###                           Base Configuration                            ###
###############################################################################

# default: the last 362880 states are kept, pruning at 10 block intervals
# nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
# everything: 2 latest states will be kept; pruning at 10 block intervals.
# custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'
pruning = "{{ .BaseConfig.Pruning }}"

# These are applied if and only if the pruning strategy is custom.
pruning-keep-recent = "{{ .BaseConfig.PruningKeepRecent }}"
pruning-interval = "{{ .BaseConfig.PruningInterval }}"

# HaltHeight contains a non-zero block height at which a node will gracefully
# halt and shutdown that can be used to assist upgrades and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-height = {{ .BaseConfig.HaltHeight }}

# HaltTime contains a non-zero minimum block time (in Unix seconds) at which
# a node will gracefully halt and shutdown that can be used to assist upgrades
# and testing.
#
# Note: Commitment of state will be attempted on the corresponding block.
halt-time = {{ .BaseConfig.HaltTime }}

# MinRetainBlocks defines the minimum block height offset from the current
# block being committed, such that all blocks past this offset are pruned
# from CometBFT. It is used as part of the process of determining the
# ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
# that no blocks should be pruned.
#
# This configuration value is only responsible for pruning CometBFT blocks.
# It has no bearing on application state pruning which is determined by the
# "pruning-*" configurations.
#
# Note: CometBFT block pruning is dependent on this parameter in conjunction
# with the unbonding (safety threshold) period, state pruning and state sync
# snapshot parameters to determine the correct minimum value of
# ResponseCommit.RetainHeight.
min-retain-blocks = {{ .BaseConfig.MinRetainBlocks }}

# InterBlockCache enables inter-block caching.
inter-block-cache = {{ .BaseConfig.InterBlockCache }}

# IavlCacheSize set the size of the iavl tree cache (in number of nodes).
iavl-cache-size = {{ .BaseConfig.IAVLCacheSize }}

# IAVLDisableFastNode enables or disables the fast node feature of IAVL.
# Default is false.
iavl-disable-fastnode = {{ .BaseConfig.IAVLDisableFastNode }}


###############################################################################
###                         Telemetry Configuration                         ###
###############################################################################

[telemetry]

# Enabled enables Prometheus metrics collection for all beacon-kit services.
enabled = {{ .Telemetry.Enabled }}

# ServiceName defines the namespace prefix for all Prometheus metrics.
service-name = "{{ .Telemetry.ServiceName }}"

# EnableHostnameLabel enables adding hostname as a constant label to all metrics.
enable-hostname-label = {{ .Telemetry.EnableHostnameLabel }}

# GlobalLabels defines a global set of name/value label tuples applied to all metrics.
{{ if gt (len .Telemetry.GlobalLabels) 0 }}global-labels = [{{ range $k, $v := .Telemetry.GlobalLabels }}
  ["{{index $v 0 }}", "{{ index $v 1}}"],{{ end }}
]{{ else }}global-labels = []{{ end }}
