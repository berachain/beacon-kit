def init(plan, image, files):
    plan.run_sh(
        image = image,
        run = "mkdir /config/datadir && geth init --datadir=/config/datadir /config/genesis.json",
        files = {
            "/config": files.config,
        },
        store = [
            StoreSpec(src = "/config", name = files.config),
        ],
    )

def launch(plan, image, l1, files):
    rpc_base, rpc_port = get_url_parts(l1.rpc_url)
    ws_base, ws_port = get_url_parts(l1.ws_url)
    auth_rpc_base, auth_rpc_port = get_url_parts(l1.auth_rpc_url)
    service = plan.add_service(
        name = "op-geth",
        config = ServiceConfig(
            image = image,
            cmd = [
                "geth",
                "--datadir",
                "/config/datadir",
                "--http",
                "--http.corsdomain=*",
                "--http.vhosts=*",
                "--http.addr={}".format(rpc_base),
                "--http.port={}".format(rpc_port),
                "--http.api=web3,debug,eth,txpool,net,engine",
                "--ws",
                "--ws.addr={}".format(ws_base),
                "--ws.port={}".format(ws_port),
                "--ws.origins=*",
                "--ws.api=debug,eth,txpool,net,engine",
                "--syncmode=full",
                "--gcmode=archive",
                "--nodiscover",
                "--maxpeers=0",
                "--networkid={}".format(l1.chain_id),
                "--authrpc.vhosts=*",
                "--authrpc.addr={}".format(auth_rpc_base),
                "--authrpc.port={}".format(auth_rpc_port),
                "--authrpc.jwtsecret=./config/jwt.txt",
                "--rollup.disabletxpoolgossip=true",
            ],
            ports = {
                "rpc": PortSpec(
                    number = int(rpc_port),
                    url = l1.rpc_url,
                ),
                "ws": PortSpec(
                    number = int(ws_port),
                    url = l1.ws_url,
                ),
                "auth_rpc": PortSpec(
                    number = int(auth_rpc_port),
                    url = l1.auth_rpc_url,
                ),
            },
            files = {
                "/config": files.config,
            },
        ),
    )

    return service.ports["rpc"].url

def get_url_parts(url):
    parts = url.split(":")
    if parts[1].startswith("//"):
        parts[1] = parts[1][2:]

    return parts[1], parts[-1]
