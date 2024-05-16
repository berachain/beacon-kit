NAME="op-geth"

ARTIFACT_NAME = NAME
PATH = "/op-geth"

DATADIR = "{}/datadir".format(PATH)
JWT_PATH = "{}/jwt.txt".format(PATH)
GENESIS_PATH = "{}/genesis.json".format(PATH)

def init(plan, image, files):
    plan.run_sh(
        image = image,
        run = "cd {} && mkdir {} && geth init --datadir={} {}".format(
            PATH,
            DATADIR,
            DATADIR,
            GENESIS_PATH,
        ),
        files = {PATH: files.op_geth},
        store = [StoreSpec(src=PATH, name=files.op_geth)],
    )

def launch(plan, image, l1, files):
    rpc_base, rpc_port = get_url_parts(l1.rpc_url)
    ws_base, ws_port = get_url_parts(l1.ws_url)
    auth_rpc_base, auth_rpc_port = get_url_parts(l1.auth_rpc_url)
    service = plan.add_service(
        name = NAME,
        config = ServiceConfig(
            image = image,
            cmd = [
                "geth",
                "--datadir",
                DATADIR,
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
                "--authrpc.jwtsecret={}".format(JWT_PATH),
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
            files = {PATH: files.op_geth},
        ),
    )

    return service.ports["rpc"].url

def get_url_parts(url):
    parts = url.split(":")
    if parts[1].startswith("//"):
        parts[1] = parts[1][2:]

    return parts[1], parts[-1]

def generate_jwt_secret(plan):
    output = plan.run_sh(
        image = "alpine/openssl:latest",
        run = "mkdir {} && openssl rand -hex 32 > {}".format(PATH, JWT_PATH),
        store = [StoreSpec(src=JWT_PATH, name=ARTIFACT_NAME)],
    )

    return output.files_artifacts[0]
