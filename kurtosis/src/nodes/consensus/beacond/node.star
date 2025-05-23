# Contains functionality for initializing and starting the nodes

def start(persistent_peers, is_seed, validator_index, config_settings, app_settings, kzg_impl):
    mv_genesis = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    set_config = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    set_config += '\nsed -i "s/^pprof_laddr = \\".*\\"/pprof_laddr = \\"0.0.0.0:6060\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/\\":26660/\\"0.0.0.0:26660/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^flush_throttle_timeout = \\".*\\"$/flush_throttle_timeout = \\"10ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose = \\".*\\"$/timeout_propose = \\"{}\\"/" {}/config/config.toml'.format(config_settings.timeout_propose, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_propose_delta = \\".*\\"$/timeout_propose_delta = \\"500ms\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_prevote = \\".*\\"$/timeout_prevote = \\"{}\\"/" {}/config/config.toml'.format(config_settings.timeout_prevote, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_precommit = \\".*\\"$/timeout_precommit = \\"{}\\"/" {}/config/config.toml'.format(config_settings.timeout_precommit, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^timeout_commit = \\".*\\"$/timeout_commit = \\"{}\\"/" {}/config/config.toml'.format(config_settings.timeout_commit, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^unsafe = false$/unsafe = true/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^type = \\".*\\"$/type = \\"nop\\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^discard_abci_responses = false$/discard_abci_responses = true/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^# other sinks such as Prometheus.\nenabled = false$/# other sinks such as Prometheus.\nenabled = true/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus-retention-time = 0$/prometheus-retention-time = 60/" {}/config/app.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^payload-timeout = \\".*\\"$/payload-timeout = \\"{}\\"/" {}/config/app.toml'.format(app_settings.payload_timeout, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^enable-optimistic-payload-builds = \\".*\\"$/enable-optimistic-payload-builds = \\"{}\\"/" {}/config/app.toml'.format(app_settings.enable_optimistic_payload_builds, "$BEACOND_HOME")
    set_config += '\nsed -i "s/^suggested-fee-recipient = \\"0x0000000000000000000000000000000000000000\\"/suggested-fee-recipient = \\"0x$(printf \"%040d\" {})\\"/" {}/config/app.toml'.format(validator_index, "$BEACOND_HOME")
    persistent_peers_option = ""
    seed_option = ""
    if persistent_peers != "":
        persistent_peers_option = "--p2p.seeds {}".format("$BEACOND_PERSISTENT_PEERS")

    if is_seed:
        set_config += '\nsed -i "s/^max_num_inbound_peers = 40$/max_num_inbound_peers = 200/" {}/config/config.toml'.format("$BEACOND_HOME")
        set_config += '\nsed -i "s/^max_num_outbound_peers = 10$/max_num_outbound_peers = 200/" {}/config/config.toml'.format("$BEACOND_HOME")
        seed_option = "--p2p.seed_mode true"
    else:
        set_config += '\nsed -i "s/^max_num_inbound_peers = 40$/max_num_inbound_peers = {}/" {}/config/config.toml'.format(config_settings.max_num_inbound_peers, "$BEACOND_HOME")
        set_config += '\nsed -i "s/^max_num_outbound_peers = 10$/max_num_outbound_peers = {}/" {}/config/config.toml'.format(config_settings.max_num_outbound_peers, "$BEACOND_HOME")

    start_node = "/usr/bin/beacond start --rpc.laddr tcp://0.0.0.0:26657 \
    --beacon-kit.chain-spec={} \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json \
    --beacon-kit.kzg.implementation={} \
    --beacon-kit.engine.rpc-dial-url {} \
    --beacon-kit.block-store-service.enabled \
    --beacon-kit.node-api.enabled --beacon-kit.node-api.logging --beacon-kit.node-api.address 0.0.0.0:3500 \
    --pruning=nothing \
    {} {}".format("$BEACOND_CHAIN_SPEC", kzg_impl, "$BEACOND_ENGINE_DIAL_URL", seed_option, persistent_peers_option)

    return "{} && {} && {}".format(mv_genesis, set_config, start_node)

def get_genesis_env_vars(cl_service_name, chain_id, chain_spec):
    return {
        "BEACOND_MONIKER": cl_service_name,
        "BEACOND_NET": "VALUE_2",
        "BEACOND_HOME": "/root/.beacond",
        "BEACOND_CHAIN_ID": "beacon-kurtosis-{}".format(chain_id),
        "BEACOND_DEBUG": "false",
        "BEACOND_KEYRING_BACKEND": "test",
        "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
        "BEACOND_ETH_CHAIN_ID": str(chain_id),
        "BEACOND_ENABLE_PROMETHEUS": "true",
        "ETH_GENESIS": "/root/eth_genesis/genesis.json",
        # For devnet/testing purposes, we use the same withdrawal address for all validators.
        # In production, each validator should use an address derived from their own withdrawal credentials.
        # This is fine for a local development network.
        "WITHDRAWAL_ADDRESS": "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4",
        "DEPOSIT_AMOUNT": "32000000000",
        "BEACOND_CHAIN_SPEC": chain_spec,
        "CHAIN_SPEC": chain_spec,
    }
