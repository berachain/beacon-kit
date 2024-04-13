# Contains functionality for initializing and starting the nodes

bash = import_module("./bash.star")  # Import helper module

def init_beacond(plan, chain_id, moniker, home, is_first_validator, cl_service_name):
    genesis_file = "{}/config/genesis.json".format(home)

    # Check if genesis file exists, if not then initialize the beacond
    init_node = "if [ ! -f {} ]; then /usr/bin/beacond init --chain-id {} {} --home {} --beacon-kit.accept-tos; fi".format(genesis_file, chain_id, moniker, home)
    bash.exec_on_service(plan, cl_service_name, init_node)
    if is_first_validator == True:
        create_beacond_config_directory(plan, home, cl_service_name)
    add_validator(plan, home, cl_service_name)
    collect_validator(plan, home, cl_service_name)

def create_beacond_config_directory(plan, home, cl_service_name):
    GENESIS = "{}/config/genesis.json".format(home)
    TMP_GENESIS = "{}/config/tmp_genesis.json".format(home)
    command = 'jq \'.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]\' {} > {}'.format(GENESIS, TMP_GENESIS)
    bash.exec_on_service(plan, cl_service_name, command)
    mv = "mv {} {}".format(TMP_GENESIS, GENESIS)
    bash.exec_on_service(plan, cl_service_name, mv)

def add_validator(plan, home, cl_service_name):
    command = "/usr/bin/beacond genesis add-validator --home {} --beacon-kit.accept-tos".format(home)
    bash.exec_on_service(plan, cl_service_name, command)

def collect_validator(plan, home, cl_service_name):
    command = "/usr/bin/beacond genesis collect-validators --home {}".format(home)
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
