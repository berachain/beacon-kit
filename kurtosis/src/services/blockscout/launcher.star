shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
postgres = import_module("github.com/kurtosis-tech/postgres-package/main.star")

IMAGE_NAME_BLOCKSCOUT = "blockscout/blockscout:7.0.2.commit.900ef697"
IMAGE_NAME_BLOCKSCOUT_VERIF = "ghcr.io/blockscout/smart-contract-verifier:v1.9.2-arm"
IMAGE_NAME_BLOCKSCOUT_FRONTEND = "ghcr.io/blockscout/frontend:v2.3.5"

SERVICE_NAME_BLOCKSCOUT = "blockscout"
SERVICE_NAME_BLOCKSCOUT_FRONTEND = "blockscout-frontend"

HTTP_PORT_ID = "http"
HTTP_PORT_NUMBER = 4000
HTTP_PORT_NUMBER_VERIF = 8050
HTTP_PORT_NUMBER_FRONTEND = 3000

BLOCKSCOUT_MIN_CPU = 100
BLOCKSCOUT_MAX_CPU = 1000
BLOCKSCOUT_MIN_MEMORY = 1024
BLOCKSCOUT_MAX_MEMORY = 4096

BLOCKSCOUT_VERIF_MIN_CPU = 10
BLOCKSCOUT_VERIF_MAX_CPU = 1000
BLOCKSCOUT_VERIF_MIN_MEMORY = 10
BLOCKSCOUT_VERIF_MAX_MEMORY = 1024

BLOCKSCOUT_FRONTEND_MIN_CPU = 10
BLOCKSCOUT_FRONTEND_MAX_CPU = 500
BLOCKSCOUT_FRONTEND_MIN_MEMORY = 256
BLOCKSCOUT_FRONTEND_MAX_MEMORY = 1024

USED_PORTS = {
    HTTP_PORT_ID: shared_utils.new_port_spec(
        HTTP_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
}

VERIF_USED_PORTS = {
    HTTP_PORT_ID: shared_utils.new_port_spec(
        HTTP_PORT_NUMBER_VERIF,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
}

FRONTEND_USED_PORTS = {
    HTTP_PORT_ID: shared_utils.new_port_spec(
        HTTP_PORT_NUMBER_FRONTEND,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
}

def launch_blockscout(
        plan,
        full_node_el_clients,
        client_from_user,
        persistent):
    postgres_output = postgres.run(
        plan,
        service_name = "{}-postgres".format(SERVICE_NAME_BLOCKSCOUT),
        database = "blockscout",
        extra_configs = ["max_connections=1000"],
        persistent = persistent,
    )
    el_client_info = {}

    # Get the full_node_el_clients that match the client_from_user
    for full_node_el_client_name, full_node_el_client_service in full_node_el_clients.items():
        if full_node_el_client_name in client_from_user:
            name = full_node_el_client_name

            el_client_info = get_el_client_info(
                full_node_el_client_name,
                8545,
                8546,
                name,
            )
            break

    config_verif = get_config_verif()
    verif_service_name = "{}-verif".format(SERVICE_NAME_BLOCKSCOUT)
    verif_service = plan.add_service(verif_service_name, config_verif)
    verif_url = "http://{}:{}/api".format(
        verif_service.hostname,
        verif_service.ports["http"].number,
    )

    config_backend = get_config_backend(
        postgres_output,
        el_client_info.get("RPC_Url"),
        el_client_info.get("WS_Url"),
        verif_url,
        el_client_info.get("Eth_Type"),
    )
    blockscout_service = plan.add_service(SERVICE_NAME_BLOCKSCOUT, config_backend)

    # NEXT_PUBLIC_API_HOST must be "hostname:port" with NO scheme prefix.
    # The frontend prepends NEXT_PUBLIC_API_PROTOCOL to build the full URL and
    # CSP connect-src. Passing a full URL (e.g. "http://...") creates a
    # double-scheme like "http://http://..." which invalidates the entire CSP
    # directive and causes the browser to block all connections.
    backend_api_host = "{}:{}".format(
        blockscout_service.hostname,
        blockscout_service.ports["http"].number,
    )

    config_frontend = get_config_frontend(backend_api_host)
    frontend_service = plan.add_service(SERVICE_NAME_BLOCKSCOUT_FRONTEND, config_frontend)

    # public_ports pins container:3000 → host:3000 deterministically, so this is always correct.
    frontend_url = "http://localhost:{}".format(HTTP_PORT_NUMBER_FRONTEND)
    plan.print("Blockscout frontend available at: {}".format(frontend_url))

    return frontend_url

def get_config_verif():
    return ServiceConfig(
        image = IMAGE_NAME_BLOCKSCOUT_VERIF,
        ports = VERIF_USED_PORTS,
        env_vars = {
            "SMART_CONTRACT_VERIFIER__SERVER__HTTP__ADDR": "0.0.0.0:{}".format(
                HTTP_PORT_NUMBER_VERIF,
            ),
        },
        min_cpu = BLOCKSCOUT_VERIF_MIN_CPU,
        max_cpu = BLOCKSCOUT_VERIF_MAX_CPU,
        min_memory = BLOCKSCOUT_VERIF_MIN_MEMORY,
        max_memory = BLOCKSCOUT_VERIF_MAX_MEMORY,
    )

def get_config_backend(
        postgres_output,
        el_client_rpc_url,
        el_client_ws_url,
        verif_url,
        el_client_name):
    database_url = "{protocol}://{user}:{password}@{hostname}:{port}/{database}".format(
        protocol = "postgresql",
        user = postgres_output.user,
        password = postgres_output.password,
        hostname = postgres_output.service.hostname,
        port = postgres_output.port.number,
        database = postgres_output.database,
    )

    return ServiceConfig(
        image = IMAGE_NAME_BLOCKSCOUT,
        ports = USED_PORTS,
        cmd = [
            "/bin/sh",
            "-c",
            'bin/blockscout eval "Elixir.Explorer.ReleaseTasks.create_and_migrate()" && bin/blockscout start',
        ],
        env_vars = {
            "ETHEREUM_JSONRPC_VARIANT": el_client_name,
            "ETHEREUM_JSONRPC_HTTP_URL": el_client_rpc_url,
            "ETHEREUM_JSONRPC_FALLBACK_HTTP_URL": el_client_rpc_url,
            "ETHEREUM_JSONRPC_WS_URL": el_client_ws_url,
            "ETHEREUM_JSONRPC_TRACE_URL": el_client_rpc_url,
            "DATABASE_URL": database_url,
            "COIN": "ETH",
            "MICROSERVICE_SC_VERIFIER_ENABLED": "true",
            "MICROSERVICE_SC_VERIFIER_URL": verif_url,
            "MICROSERVICE_SC_VERIFIER_TYPE": "sc_verifier",
            "INDEXER_DISABLE_PENDING_TRANSACTIONS_FETCHER": "true",
            "INDEXER_DISABLE_INTERNAL_TRANSACTIONS_FETCHER": "true",
            "INDEXER_DISABLE_COIN_BALANCES_FETCHER": "true",
            "INDEXER_COIN_BALANCES_BATCH_SIZE": "0",
            "COIN_BALANCE_HISTORY_DAYS": "0",
            "EXCHANGE_RATES_ENABLED": "false",
            "EXCHANGE_RATES_COINGECKO_API_KEY": "",
            "FIRST_BLOCK": "1",
            "ECTO_USE_SSL": "false",
            "NETWORK": "Kurtosis",
            "SUBNETWORK": "Kurtosis",
            "API_V2_ENABLED": "true",
            "PORT": "{}".format(HTTP_PORT_NUMBER),
            "SECRET_KEY_BASE": "56NtB48ear7+wMSf0IQuWDAAazhpb31qyc7GiyspBP2vh7t5zlCsF5QDv76chXeN",
        },
        min_cpu = BLOCKSCOUT_MIN_CPU,
        max_cpu = BLOCKSCOUT_MAX_CPU,
        min_memory = BLOCKSCOUT_MIN_MEMORY,
        max_memory = BLOCKSCOUT_MAX_MEMORY,
    )

def get_config_frontend(backend_api_host):
    return ServiceConfig(
        image = IMAGE_NAME_BLOCKSCOUT_FRONTEND,
        ports = FRONTEND_USED_PORTS,
        # Pin container port 3000 → host port 3000 deterministically.
        # Without this, Kurtosis assigns a random ephemeral host port. The
        # frontend constructs absolute proxy URLs as
        # http://NEXT_PUBLIC_APP_HOST:NEXT_PUBLIC_APP_PORT/node-api/proxy/...
        # which must match the browser's origin for CSP 'self' to allow it.
        # If the host port doesn't match, the browser sees a cross-origin
        # request and CSP blocks it. Requires port 3000 to be free on the host.
        public_ports = FRONTEND_USED_PORTS,
        env_vars = {
            # "hostname:port" only — no "http://" prefix. The frontend
            # prepends NEXT_PUBLIC_API_PROTOCOL to build both the API URL and
            # the CSP connect-src directive. A full URL here produces a
            # double-scheme ("http://http://...") which invalidates the CSP
            # and causes the browser to block all API calls.
            "NEXT_PUBLIC_API_HOST": backend_api_host,
            "NEXT_PUBLIC_API_PROTOCOL": "http",
            "NEXT_PUBLIC_API_WEBSOCKET_PROTOCOL": "ws",
            # NEXT_PUBLIC_APP_HOST:NEXT_PUBLIC_APP_PORT is the browser-accessible
            # address of this frontend. The proxy URL is constructed as
            # http://localhost:3000/node-api/proxy/... which must be same-origin
            # for CSP 'self' to allow it. public_ports above guarantees host:3000.
            "NEXT_PUBLIC_APP_HOST": "localhost",
            "NEXT_PUBLIC_APP_PORT": "{}".format(HTTP_PORT_NUMBER_FRONTEND),
            "NEXT_PUBLIC_APP_PROTOCOL": "http",
            # Network info
            "NEXT_PUBLIC_NETWORK_NAME": "Kurtosis",
            "NEXT_PUBLIC_NETWORK_SHORT_NAME": "Kurtosis",
            "NEXT_PUBLIC_NETWORK_ID": "80087",
            "NEXT_PUBLIC_NETWORK_CURRENCY_NAME": "Ether",
            "NEXT_PUBLIC_NETWORK_CURRENCY_SYMBOL": "ETH",
            "NEXT_PUBLIC_NETWORK_CURRENCY_DECIMALS": "18",
            "NEXT_PUBLIC_IS_TESTNET": "true",
        },
        min_cpu = BLOCKSCOUT_FRONTEND_MIN_CPU,
        max_cpu = BLOCKSCOUT_FRONTEND_MAX_CPU,
        min_memory = BLOCKSCOUT_FRONTEND_MIN_MEMORY,
        max_memory = BLOCKSCOUT_FRONTEND_MAX_MEMORY,
    )

def get_el_client_info(service_name, rpc_port_num, ws_port_num, full_name):
    el_client_rpc_url = "http://{}:{}".format(
        service_name,
        rpc_port_num,
    )
    el_client_ws_url = "ws://{}:{}".format(
        service_name,
        ws_port_num,
    )
    el_client_type = full_name.split("-")[2]

    # Blockscout has no reth-specific config; reth is geth-compatible for JSON-RPC
    variant_map = {"reth": "geth", "erigon": "geth"}
    blockscout_variant = variant_map.get(el_client_type, el_client_type)

    return {
        "RPC_Url": el_client_rpc_url,
        "WS_Url": el_client_ws_url,
        "Eth_Type": blockscout_variant,
    }
