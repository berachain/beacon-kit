execution = import_module("../../../nodes/execution/execution.star")

NAME = "op-geth"

NETWORK_ID = 42069
ARTIFACT_NAME = NAME
PATH = "/op-geth"

DATADIR = "{}/datadir".format(PATH)
JWT_PATH = "{}/jwt.txt".format(PATH)
GENESIS_PATH = "{}/genesis.json".format(PATH)

def init(plan, image, files):
    plan.run_sh(
        description = "Initializing op-geth",
        image = image,
        run = "cd {} && mkdir {} && geth init --datadir={} {}".format(
            PATH,
            DATADIR,
            DATADIR,
            GENESIS_PATH,
        ),
        files = {PATH: files.op_geth},
        store = [StoreSpec(src = PATH, name = files.op_geth)],
    )

def launch(plan, image, l1, files):
    service = plan.add_service(
        name = NAME,
        config = ServiceConfig(
            image = image,
            cmd = [
                "--datadir",
                DATADIR,
                "--http",
                "--http.corsdomain=*",
                "--http.vhosts=*",
                "--http.addr=0.0.0.0",
                "--http.port={}".format(execution.RPC_PORT_NUM),
                "--http.api=web3,debug,eth,txpool,net,engine",
                "--ws",
                "--ws.addr=0.0.0.0",
                "--ws.port={}".format(execution.WS_PORT_NUM),
                "--ws.origins=*",
                "--ws.api=debug,eth,txpool,net,engine",
                "--syncmode=full",
                "--gcmode=archive",
                "--nodiscover",
                "--maxpeers=0",
                "--networkid={}".format(NETWORK_ID),
                "--authrpc.vhosts=*",
                "--authrpc.addr=0.0.0.0",
                "--authrpc.port={}".format(execution.ENGINE_RPC_PORT_NUM),
                "--authrpc.jwtsecret={}".format(JWT_PATH),
                "--rollup.disabletxpoolgossip=true",
            ],
            ports = {
                "rpc": PortSpec(number = execution.RPC_PORT_NUM),
                "ws": PortSpec(number = execution.WS_PORT_NUM),
                "auth_rpc": PortSpec(number = execution.ENGINE_RPC_PORT_NUM),
            },
            files = {PATH: files.op_geth},
        ),
    )

    return service.ip_address

def generate_jwt_secret(plan):
    output = plan.run_sh(
        image = "alpine/openssl:latest",
        run = "mkdir {} && openssl rand -hex 32 > {}".format(PATH, JWT_PATH),
        store = [StoreSpec(src = JWT_PATH, name = ARTIFACT_NAME)],
    )

    return output.files_artifacts[0]
