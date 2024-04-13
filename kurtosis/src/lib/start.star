def start(plan,persistent_peers):
    mv_genesis = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    set_config = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    set_config += '\nsed -i "s/^prometheus_listen_addr = \":26660\"$/prometheus_listen_addr = \"0.0.0.0:26660\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    set_config += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    persistent_peers_option = ''
    if persistent_peers != "":
        persistent_peers_option = '--p2p.persistent_peers {}'.format("$BEACOND_PERSISTENT_PEERS")

    start_node = '/usr/bin/beacond start \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json \
    --beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url {} \
    --beacon-kit.engine.required-chain-id {} \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --grpc.address 0.0.0.0:9090 --api.address tcp://0.0.0.0:1317 \
    --api.enable {}'.format("$BEACOND_ENGINE_DIAL_URL", "$BEACOND_ETH_CHAIN_ID", persistent_peers_option)

    return "{} && {} && {}".format(mv_genesis, set_config,start_node)