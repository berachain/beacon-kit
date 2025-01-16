el_cl_genesis_data = import_module(
    "github.com/ethpandaops/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_data.star",
)

NETWORKS_DIR_PATH = "/kurtosis/src/networks/"

NETWORKS = struct(
    kurtosis_devnet = "kurtosis-devnet/",
)

def get_genesis_data(plan, chain_spec = None):
    """
    Generates genesis data for the network.
    
    Args:
        plan: The execution plan
        chain_spec: Optional chain specification to use
    Returns:
        Dictionary containing genesis data
    """
    if chain_spec:
        return chain_spec
        
    # Default genesis configuration if none provided
    return {
        "config": {
            "chainId": 2061,
            "homesteadBlock": 0,
            "eip150Block": 0,
            "eip155Block": 0,
            "eip158Block": 0,
            "byzantiumBlock": 0,
            "constantinopleBlock": 0,
            "petersburgBlock": 0,
            "istanbulBlock": 0,
            "muirGlacierBlock": 0,
            "berlinBlock": 0,
            "londonBlock": 0,
            "arrowGlacierBlock": 0,
            "grayGlacierBlock": 0,
            "mergeNetsplitBlock": 0,
            "shanghaiTime": 0,
            "cancunTime": 0,
            "pragueTime": 0,
            "verkleTime": 0,
        },
        "nonce": "0x0",
        "timestamp": "0x0",
        "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "gasLimit": "0x1c9c380",
        "difficulty": "0x1",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "coinbase": "0x0000000000000000000000000000000000000000",
        "alloc": {},
    }

def get_genesis_path(network = "kurtosis-devnet"):
    return NETWORKS_DIR_PATH + network
