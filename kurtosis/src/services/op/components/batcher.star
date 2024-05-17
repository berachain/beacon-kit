NAME = "op-batcher"

def launch(plan, image, env, l1, l2, batcher_rpc_port, node_rpc_url):
    service = plan.add_service(
        name = NAME,
        config = ServiceConfig(
            image = image,
            cmd = [
                NAME,
                "--l1-eth-rpc={}".format(l1.rpc_url),
                "--l2-eth-rpc={}".format(l2.rpc_url),
                "--rollup-rpc={}".format(node_rpc_url),
                "--poll-interval=1s",
                "--sub-safety-margin=6",
                "--num-confirmations=1",
                "--safe-abort-nonce-too-low-count=3",
                "--resubmission-timeout=30s",
                "--rpc.addr=0.0.0.0",
                "--rpc.port={}".format(batcher_rpc_port),
                "--rpc.enable-admin",
                "--max-channel-duration=1",
                "--private-key={}".format(env.batcher_pk),
            ],
            ports = {
                "rpc": PortSpec(number = int(batcher_rpc_port)),
            },
        ),
    )

    return service.ip_address
