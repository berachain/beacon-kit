def init(plan, image, files):
    plan.run_sh(
        image=image,
        run="geth init --datadir=datadir /config/genesis.json",
        files={
            "/config/": files.config,
        },
        store={
            StoreSpec(src="datadir", dest="datadir")
        }
    )

def launch(plan, image, l1, files):
    rpc_port = l1.rpc_url.split(":")[-1]
    ws_port = l1.ws_url.split(":")[-1]
    auth_rpc_port = l1.auth_rpc_url.split(":")[-1]
    return plan.add_service(
        name="l1",
        config=ServiceConfig(
            image=image,
            cmd=[
                "geth",
                "--datadir", "./datadir",
                "--http",
                "--http.corsdomain=*",
                "--http.vhosts=*",
                "--http.addr=0.0.0.0",
                "--http.port={0}".format(rpc_port),
                "--http.api=web3,debug,eth,txpool,net,engine",
                "--ws",
                "--ws.addr=0.0.0.0",
                "--ws.port={0}".format(ws_port),
                "--ws.origins=*",
                "--ws.api=debug,eth,txpool,net,engine",
                "--syncmode=full",
                "--gcmode=archive",
                "--nodiscover",
                "--maxpeers=0",
                "--networkid={0}".format(l1.chain_id),
                "--authrpc.vhosts=*",
                "--authrpc.addr=0.0.0.0",
                "--authrpc.port={0}".format(auth_rpc_port),
                "--authrpc.jwtsecret=./config/jwt.txt",
                "--rollup.disabletxpoolgossip=true",
            ]
            ports={
                "rpc": PortSpec(l1.rpc_url),
                "ws": PortSpec(l1.ws_url),
            },
            files={
                "/config/": files.config,
            },
        ),
    )