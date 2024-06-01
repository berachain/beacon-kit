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
    command = "/usr/bin/beacond genesis add-premined-deposit --home {}".format("$BEACOND_HOME")
    bash.exec_on_service(plan, cl_service_name, command)

def collect_validator(plan, cl_service_name):
    command = "/usr/bin/beacond genesis collect-premined-deposits --home {}".format("$BEACOND_HOME")
    bash.exec_on_service(plan, cl_service_name, command)

def start(persistent_peers, is_seed, validator_index):
    mv_genesis = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    set_config = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus_listen_addr = \\":26660\\"$/prometheus_listen_addr = \\"0.0.0.0:26660\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^flush_throttle_timeout = \\".*\\"$/flush_throttle_timeout = \\"10ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose = \\".*\\"$/timeout_propose = \\"3s\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose_delta = \\".*\\"$/timeout_propose_delta = \\"500ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_vote = \\".*\\"$/timeout_vote = \\"2s\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_commit = \\".*\\"$/timeout_commit = \\"1s\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^unsafe = false$/unsafe = true/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^type = \\".*\\"$/type = \\"nop\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^discard_abci_responses = false$/discard_abci_responses = true/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^# other sinks such as Prometheus.\nenabled = false$/# other sinks such as Prometheus.\nenabled = true/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus-retention-time = 0$/prometheus-retention-time = 60/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^payload-timeout = \\".*\\"$/payload-timeout = \\"1.5s\\"/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^suggested-fee-recipient = \\"0x0000000000000000000000000000000000000000\\"/suggested-fee-recipient = \\"0x$(printf \"%040d\" {})\\"/" {}/config/app.toml'.format(validator_index, "$BEACOND_HOME")
    persistent_peers_option = ""
    seed_option = ""
    if persistent_peers != "":
        persistent_peers_option = "--p2p.seeds {}".format("$BEACOND_PERSISTENT_PEERS")

    if is_seed:
        set_config += '\nsed -i "s/^max_num_inbound_peers = 40$/max_num_inbound_peers = 200/" {}/config/config.toml'.format("$BEACOND_HOME")
        set_config += '\nsed -i "s/^max_num_outbound_peers = 10$/max_num_outbound_peers = 200/" {}/config/config.toml'.format("$BEACOND_HOME")
        seed_option = "--p2p.seed_mode true"

    start_node = "/usr/bin/beacond start \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json \
    --beacon-kit.engine.rpc-dial-url {} \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --grpc.address 0.0.0.0:9090 --api.address tcp://0.0.0.0:1317 \
    --api.enable {} {}".format("$BEACOND_ENGINE_DIAL_URL", seed_option, persistent_peers_option)

    return "{} && {} && {}".format(mv_genesis, set_config, start_node)

def get_init_sh():
    genesis_file = "{}/config/genesis.json".format("$BEACOND_HOME")

    # Check if genesis file exists, if not then initialize the beacond
    command = "if [ ! -f {} ]; then /usr/bin/beacond init --chain-id {} {} --home {} --consensus-key-algo {}; fi".format(genesis_file, "$BEACOND_CHAIN_ID", "$BEACOND_MONIKER", "$BEACOND_HOME", "$BEACOND_CONSENSUS_KEY_ALGO")
    return command

def get_add_validator_sh():
    command = "/usr/bin/beacond genesis add-premined-deposit --home {}".format("$BEACOND_HOME")
    return command

def get_collect_validator_sh():
    command = "/usr/bin/beacond genesis collect-premined-deposits --home {}".format("$BEACOND_HOME")
    return command

def get_execution_payload_sh():
    command = "/usr/bin/beacond genesis execution-payload {} --home {}".format("$ETH_GENESIS", "$BEACOND_HOME")
    return command

def get_genesis_env_vars(cl_service_name):
    return {
        "BEACOND_MONIKER": cl_service_name,
        "BEACOND_NET": "VALUE_2",
        "BEACOND_HOME": "/root/.beacond",
        "BEACOND_CHAIN_ID": "beacon-kurtosis-80086",
        "BEACOND_DEBUG": "false",
        "BEACOND_KEYRING_BACKEND": "test",
        "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
        "BEACOND_ETH_CHAIN_ID": "80086",
        "BEACOND_ENABLE_PROMETHEUS": "true",
        "BEACOND_CONSENSUS_KEY_ALGO": "bls12_381",
        "ETH_GENESIS": "/root/eth_genesis/genesis.json",
    }
