execution = import_module("../../execution/execution.star")
shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
COMETBFT_RPC_PORT_NUM = 26657
COMETBFT_P2P_PORT_NUM = 26656
COMETBFT_REST_PORT_NUM = 1317
COMETBFT_PPROF_PORT_NUM = 6060
METRICS_PORT_NUM = 26660
ENGINE_RPC_PORT_NUM = 8551

COMETBFT_RPC_PORT_ID = "cometbft-rpc"
COMETBFT_P2P_PORT_ID = "cometbft-p2p"
COMETBFT_GRPC_PORT_ID = "cometbft-grpc"
COMETBFT_REST_PORT_ID = "cometbft-rest"
COMETBFT_PPROF_PORT_ID = "cometbft-pprof"
ENGINE_RPC_PORT_ID = "engine-rpc"
METRICS_PORT_ID = "metrics"
METRICS_PATH = "/metrics"

USED_PORTS = {
    COMETBFT_RPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_P2P_PORT_ID: shared_utils.new_port_spec(COMETBFT_P2P_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_REST_PORT_ID: shared_utils.new_port_spec(COMETBFT_REST_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_PPROF_PORT_ID: shared_utils.new_port_spec(COMETBFT_PPROF_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(ENGINE_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    METRICS_PORT_ID: shared_utils.new_port_spec(METRICS_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = None),
}

def create_node_config(plan, node_struct, node_settings, paired_el_client_name, jwt_file = None, kzg_trusted_setup_file = None):
    engine_dial_url = "http://{}:{}".format(paired_el_client_name, execution.ENGINE_RPC_PORT_NUM)

    plan.print("node_struct", str(node_struct))
    plan.print("node_settings", str(node_settings))
    config_settings = node_settings["consensus_settings"]["config"]
    app_settings = node_settings["consensus_settings"]["app"]
    kzg_impl = node_struct["kzg_impl"]

    cmd = "{}".format(start(plan, 0, config_settings, app_settings, kzg_impl))

    beacond_config = get_config(
        plan,
        node_struct,
        node_settings,
        engine_dial_url,
        entrypoint = ["bash", "-c"],
        cmd = [cmd],
        jwt_file = jwt_file,
        kzg_trusted_setup_file = kzg_trusted_setup_file,
    )

    plan.print(beacond_config)

    return beacond_config

def start(plan, validator_index, config_settings, app_settings, kzg_impl):
    plan.print("BEACOND_HOME", "$BEACOND_HOME")
    BEACOND_HOME = "/root/.beacond"
    plan.print("BEACOND_HOME", BEACOND_HOME)

    # mv_genesis = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    set_config = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    set_config += '\nsed -i "s/^pprof_laddr = \\".*\\"/pprof_laddr = \\"0.0.0.0:6060\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/:26660/0.0.0.0:26660/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^flush_throttle_timeout = \\".*\\"$/flush_throttle_timeout = \\"10ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose = \\".*\\"$/timeout_propose = \\"{}\\"/" {}/config/config.toml'.format(config_settings["timeout_propose"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose_delta = \\".*\\"$/timeout_propose_delta = \\"500ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_prevote = \\".*\\"$/timeout_prevote = \\"{}\\"/" {}/config/config.toml'.format(config_settings["timeout_prevote"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_precommit = \\".*\\"$/timeout_precommit = \\"{}\\"/" {}/config/config.toml'.format(config_settings["timeout_precommit"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_commit = \\".*\\"$/timeout_commit = \\"{}\\"/" {}/config/config.toml'.format(config_settings["timeout_commit"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^unsafe = false$/unsafe = true/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^type = \\".*\\"$/type = \\"nop\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^discard_abci_responses = false$/discard_abci_responses = true/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^# other sinks such as Prometheus.\nenabled = false$/# other sinks such as Prometheus.\nenabled = true/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus-retention-time = 0$/prometheus-retention-time = 60/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^payload-timeout = \\".*\\"$/payload-timeout = \\"{}\\"/" {}/config/app.toml'.format(app_settings["payload_timeout"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^enable-optimistic-payload-builds = \\".*\\"$/enable-optimistic-payload-builds = \\"{}\\"/" {}/config/app.toml'.format(app_settings["enable_optimistic_payload_builds"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^suggested-fee-recipient = \\"0x0000000000000000000000000000000000000000\\"/suggested-fee-recipient = \\"0x$(printf \"%040d\" {})\\"/" {}/config/app.toml'.format(validator_index, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^max_num_inbound_peers = 40$/max_num_inbound_peers = {}/" {}/config/config.toml'.format(config_settings["max_num_inbound_peers"], "$BEACOND_HOME")
    set_config += '\nsed -i "s/^max_num_outbound_peers = 10$/max_num_outbound_peers = {}/" {}/config/config.toml'.format(config_settings["max_num_outbound_peers"], "$BEACOND_HOME")

    start_node = "CHAIN_SPEC=testnet /usr/bin/beacond start --help && sleep infinity"
    # --home {} \
    # --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    # --beacon-kit.kzg.trusted-setup-path=/root/.beacond/config/kzg-trusted-setup.json \
    # --beacon-kit.kzg.implementation={} \
    # --beacon-kit.engine.rpc-dial-url {} \
    # --rpc.laddr tcp://0.0.0.0:26657 --api.address tcp://0.0.0.0:1317 \
    # --api.enable".format("$BEACOND_HOME",kzg_impl, "$BEACOND_ENGINE_DIAL_URL")
    return "{} && {}".format(set_config, start_node)

    # return "{} && {} && {}".format(mv_genesis, set_config, start_node)

def get_config(plan, node_struct, node_settings, engine_dial_url, entrypoint = [], cmd = [], persistent_peers = "", expose_ports = True, jwt_file = None, kzg_trusted_setup_file = None):
    exposed_ports = {}
    if expose_ports:
        exposed_ports = USED_PORTS

    files = {}
    if jwt_file:
        files["/root/jwt"] = jwt_file
    if kzg_trusted_setup_file:
        files["/root/kzg"] = kzg_trusted_setup_file

    genesis_file_execution = plan.upload_files(
        src = "../../../network/kurtosis-devnet/network-configs/genesis.json",
        name = "genesis_file_execution",
    )

    config_directory = plan.upload_files(
        src = "../../../network/80084/",
        name = "config_directory",
    )

    # Now map the entire directory to the container
    files["/root/.beacond/config"] = config_directory

    # app_toml = plan.upload_files(
    #     src = "../../../network/80084/app.toml",
    #     name = "app_toml",
    # )

    # config_toml = plan.upload_files(
    #     src = "../../../network/80084/config.toml",
    #     name = "config_toml",
    # )

    # files["/root/.beacond/config"] = genesis_file_execution
    # files["/root/.beacond/config"] = app_toml
    # files["/root/.beacond/config"] = config_toml

    settings = node_settings["consensus_settings"]

    cl_service_name = "cl-{}-{}".format("consensus", 0)

    config = ServiceConfig(
        image = settings["images"]["beaconkit"],
        files = files,
        entrypoint = entrypoint,
        cmd = cmd,
        min_cpu = settings["specs"]["min_cpu"],
        max_cpu = settings["specs"]["max_cpu"],
        min_memory = settings["specs"]["min_memory"],
        max_memory = settings["specs"]["max_memory"],
        env_vars = {
            "BEACOND_MONIKER": cl_service_name,
            "BEACOND_NET": "VALUE_2",
            "BEACOND_HOME": "/root/.beacond",
            "BEACOND_CHAIN_ID": "bartio-beacon-80084",
            "BEACOND_DEBUG": "true",
            "BEACOND_KEYRING_BACKEND": "test",
            "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
            "BEACOND_ENGINE_DIAL_URL": engine_dial_url,
            "BEACOND_ETH_CHAIN_ID": "80084",
            "BEACOND_PERSISTENT_PEERS": persistent_peers,
            "BEACOND_ENABLE_PROMETHEUS": "true",
            "BEACOND_CONSENSUS_KEY_ALGO": "bls12_381",
            "CHAIN_SPEC": "testnet",
        },
        ports = exposed_ports,
    )

    return config
