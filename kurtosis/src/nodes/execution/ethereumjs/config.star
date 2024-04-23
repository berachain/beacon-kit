global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")

# NODE_TYPE = "ethereumjs"

IMAGE = "ethpandaops/ethereumjs:stable"
ENTRYPOINT = ["sh", "-c"]
GENESIS_FILENAME = "genesis.json"
FILES = {
    "/app/genesis": "genesis_file",
    "/jwt": "jwt_file",
}
# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER
CMD = [
    "node",
	"/usr/app/node_modules/.bin/ethereumjs",
    "--gethGenesis={0}".format("/app/genesis/{}".format(GENESIS_FILENAME)),
    "--rpc",
    "--rpcEngine",
    "--rpcEngineAddr",
    "0.0.0.0",
    "--jwtSecret",
    global_constants.JWT_MOUNT_PATH_ON_CONTAINER,
    "--dataDir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--logLevel",
    "debug",
]
BOOTNODE_CMD = "--bootnodes"
GLOBAL_FILES = []
USED_PORTS_TEMPLATE = defaults.USED_PORTS_TEMPLATE
METRICS_PATH = defaults.METRICS_PATH