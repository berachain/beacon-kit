# Contains functionality for initializing and starting the nodes

bash = import_module("../../../lib/bash.star")  # Import helper module

def init_beacond(plan, is_first_validator, cl_service_name):
    init_sh = get_init_sh
    bash.exec_on_service(plan, cl_service_name, init_sh)
    if is_first_validator == True:
        create_beacond_config_directory(plan, cl_service_name)
    add_validator(plan, cl_service_name)
    collect_validator(plan, cl_service_name)

def create_beacond_config_directory(plan, cl_service_name):
    GENESIS = "{}/config/genesis.json".format("$BEACOND_HOME")
    TMP_GENESIS = "{}/config/tmp_genesis.json".format("$BEACOND_HOME")

def add_validator(plan, cl_service_name):
    command = "/usr/bin/beacond genesis add-validator --home {} --beacon-kit.accept-tos".format("$BEACOND_HOME")
    bash.exec_on_service(plan, cl_service_name, command)

def collect_validator(plan, cl_service_name):
    command = "/usr/bin/beacond genesis collect-validators --home {}".format("$BEACOND_HOME")
    bash.exec_on_service(plan, cl_service_name, command)

def start(persistent_peers):
    mv_genesis = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    set_config = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus_listen_addr = ":26660"$/prometheus_listen_addr = "0.0.0.0:26660"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    persistent_peers_option = ""
    if persistent_peers != "":
        persistent_peers_option = "--p2p.persistent_peers {}".format("$BEACOND_PERSISTENT_PEERS")

    start_node = "/usr/bin/beacond start \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json \
    --beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url {} \
    --beacon-kit.engine.required-chain-id {} \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --grpc.address 0.0.0.0:9090 --api.address tcp://0.0.0.0:1317 \
    --api.enable {}".format("$BEACOND_ENGINE_DIAL_URL", "$BEACOND_ETH_CHAIN_ID", persistent_peers_option)

    return "{} && {} && {}".format(mv_genesis, set_config, start_node)

def get_init_sh():
    genesis_file = "{}/config/genesis.json".format("$BEACOND_HOME")

    # Check if genesis file exists, if not then initialize the beacond
    command = "if [ ! -f {} ]; then /usr/bin/beacond init --chain-id {} {} --home {} --consensus-key-algo {} --beacon-kit.accept-tos; fi".format(genesis_file, "$BEACOND_CHAIN_ID", "$BEACOND_MONIKER", "$BEACOND_HOME", "$BEACOND_CONSENSUS_KEY_ALGO")
    return command

def get_add_validator_sh():
    command = "/usr/bin/beacond genesis add-validator --home {} --beacon-kit.accept-tos".format("$BEACOND_HOME")
    return command

def get_collect_validator_sh():
    command = "/usr/bin/beacond genesis collect-validators --home {}".format("$BEACOND_HOME")
    return command

def get_execution_payload_sh():
    command = "/usr/bin/beacond genesis execution-payload {} --home {}".format("$ETH_GENESIS", "$BEACOND_HOME")
    return command

def get_genesis_env_vars(cl_service_name):
    return {
        "BEACOND_MONIKER": cl_service_name,
        "BEACOND_NET": "VALUE_2",
        "BEACOND_HOME": "/root/.beacond",
        "BEACOND_CHAIN_ID": "beacon-kurtosis-80087",
        "BEACOND_DEBUG": "false",
        "BEACOND_KEYRING_BACKEND": "test",
        "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
        "BEACOND_ETH_CHAIN_ID": "80087",
        "BEACOND_ENABLE_PROMETHEUS": "true",
        "BEACOND_CONSENSUS_KEY_ALGO": "bls12_381",
        "ETH_GENESIS": "/root/eth_genesis/genesis.json",
    }
