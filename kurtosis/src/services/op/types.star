def get_l1(rpc_url, rpc_kind="any", ws_url="http://localhost:8546", auth_rpc_url="http://localhost:8551", block_time=6, chain_id=80087):
    return struct(
        rpc_url=rpc_url,
        rpc_kind=rpc_kind,
        ws_url=ws_url,
        auth_rpc_url=auth_rpc_url,
        block_time=block_time,
        chain_id=chain_id
    )

def get_l2(rpc_url):
    return struct(
        rpc_url=rpc_url
    )