JWT_FILEPATH = "/testing/sync-network/network/jwt-secret.hex"
GENESIS_FILEPATH = "/testing/sync-network/network/80084/genesis.json"

def upload_global_files(plan, node_modules):
    genesis_file = plan.upload_files(
        src = "../../network/kurtosis-devnet/network-configs/genesis.json",
        name = "genesis_file",
    )

    jwt_file = plan.upload_files(
        src = JWT_FILEPATH,
        name = "jwt_file",
    )
    for node_module in node_modules.values():
        for global_file in node_module.GLOBAL_FILES:
            plan.upload_files(
                src = global_file[0],
                name = global_file[1],
            )

    return jwt_file
