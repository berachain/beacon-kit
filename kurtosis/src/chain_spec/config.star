def get_chain_spec(chain_spec_config = {}):
    """
    Returns chain specification based on provided configuration.
    
    Args:
        chain_spec_config: Dictionary containing chain specification parameters
    Returns:
        Dictionary containing the chain specification
    """
    default_config = {
        "chain_id": 2061,
        "homestead_block": 0,
        "eip150_block": 0,
        "eip155_block": 0,
        "eip158_block": 0,
        "byzantium_block": 0,
        "constantinople_block": 0,
        "petersburg_block": 0,
        "istanbul_block": 0,
        "muir_glacier_block": 0,
        "berlin_block": 0,
        "london_block": 0,
        "arrow_glacier_block": 0,
        "gray_glacier_block": 0,
        "merge_netsplit_block": 0,
        "shanghai_time": 0,
        "cancun_time": 0,
        "prague_time": 0,
        "verkle_time": 0,
    }
    
    # Merge user provided config with defaults
    config = default_config.copy()
    config.update(chain_spec_config)
    
    return {
        "config": config,
        "nonce": "0x0",
        "timestamp": "0x0",
        "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "gasLimit": "0x1c9c380",
        "difficulty": "0x1",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "coinbase": "0x0000000000000000000000000000000000000000",
        "alloc": {},
    } 
