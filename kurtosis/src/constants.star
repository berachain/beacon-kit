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
