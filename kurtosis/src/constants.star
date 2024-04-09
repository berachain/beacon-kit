KURTOSIS_IP_ADDRESS_PLACEHOLDER = "KURTOSIS_IP_ADDR_PLACEHOLDER"

GLOBAL_LOG_LEVEL = struct(
    info = "info",
    error = "error",
    warn = "warn",
    debug = "debug",
    trace = "trace",
)

JWT_MOUNT_PATH_ON_CONTAINER = "/jwt/jwt-secret.hex"
JWT_FILEPATH = "/kurtosis/src/nodes/jwt-secret.hex"
KZG_TRUSTED_SETUP_MOUNT_PATH_ON_CONTAINER = "/kzg/kzg-trusted-setup.json"
KZG_TRUSTED_SETUP_FILEPATH = "/kurtosis/src/nodes/kzg-trusted-setup.json"

def new_prefunded_account(address, private_key):
    return struct(address = address, private_key = private_key)

PRE_FUNDED_ACCOUNTS = [
    new_prefunded_account(
        "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4",
        "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
    ),
]
