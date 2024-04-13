def start():
    mv_command = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    sed = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    sed += '\nsed -i "s/^prometheus_listen_addr = \":26660\"$/prometheus_listen_addr = \"0.0.0.0:26660\"/" {}/config/config.toml'.format("$BEACOND_HOME")
    sed += '\nsed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")


    command = '/usr/bin/beacond start \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json \
    --beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url {} \
    --beacon-kit.engine.required-chain-id {} \
    --p2p.persistent_peers {} \
    --rpc.laddr tcp://127.0.0.1:26657 \
    --grpc.address localhost:9090 --api.address tcp://localhost:1317 \
    --api.enable '.format("$BEACOND_ENGINE_DIAL_URL", "$BEACOND_ETH_CHAIN_ID", "$BEACOND_PERSISTENT_PEERS")
    # final_command = "{} && {} && {} && {} && {}".format(mv_command, sed_1,sed_2,sed_3,command)
    final_command = "{} && {} && {}".format(mv_command, sed,command)

    # commands = [
    #     "mv /root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json",
    #     'sed -i "s/^prometheus = false$/prometheus = $BEACOND_ENABLE_PROMETHEUS/" $BEACOND_HOME/config/config.toml',
    #     'sed -i "s/^prometheus_listen_addr = \\":26660\\"$/prometheus_listen_addr = \\"0.0.0.0:26660\\"/" $BEACOND_HOME/config/config.toml',
    #     'sed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "$BEACOND_HOME/config/config.toml"',
    #     '/usr/bin/beacond start --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex --beacon-kit.kzg.trusted-setup-path=/root/kzg/kzg-trusted-setup.json --beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url $BEACOND_ENGINE_DIAL_URL --beacon-kit.engine.required-chain-id $BEACOND_ETH_CHAIN_ID --p2p.persistent_peers "$BEACOND_PERSISTENT_PEERS" --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090 --api.address tcp://0.0.0.0:1317 --api.enable'
    # ]

    return final_command

def change_config_toml():
    mv_command = "mv root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json"
    sed_1 = 'sed -i "s/^prometheus = false$/prometheus = {}/" {}/config/config.toml'.format("$BEACOND_ENABLE_PROMETHEUS", "$BEACOND_HOME")
    sed_2 = 'sed -i "s/^prometheus_listen_addr = ":26660"$/prometheus_listen_addr = "0.0.0.0:26660"/" {}/config/config.toml'.format("$BEACOND_HOME")
    sed_3 = 'sed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "{}/config/config.toml"'.format("$BEACOND_HOME")
    return "{} && {} && {} && {}".format(mv_command, sed_1, sed_2, sed_3)