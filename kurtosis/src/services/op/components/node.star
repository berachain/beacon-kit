def init(plan, image, env):
    plan.run_sh(
        image=image,
        run='mkdir datadir && genesis l2\
            --deploy-config contracts-bedrock/deploy-config/getting-started.json \
            --l1-deployments contracts-bedrock/deployments/getting-started/l1.json \
            --outfile.l2 genesis.json \
            --outfile.rollup rollup.json \
            --l1-rpc {}'.format(env["L1_RPC_URL"]),
        store=[
            StoreSpec(src="genesis.json", name="/l2/genesis.json"),
            StoreSpec(src="rollup.json", name="/rollup/rollup.json"),
            StoreSpec(src="datadir", name="/datadir"),
        ]
    )


def launch(plan, image, files, env, l1, l2, node_rpc_port):
    node_rpc_url = "http://{}:{}".format("0.0.0.0", node_rpc_port)
    return plan.add_service(
        name="op-node",
        config=ServiceConfig(
            image=image,
            cmd=[
                "op-node",
                "--l1={}".format(l1.rpc_url),
                "--l1.rpckind={}".format(l1.rpc_kind),
                "--l1.trustrpc=true",
                "--l2={}".format(l2.rpc_url),
                "--l2.jwt-secret=./config/jwt.txt",
                "--sequencer.enabled",
                "--sequencer.l1-confs=5",
                "--verifier.l1-confs=4",
                "--rollup.config=./rollup/rollup.json",
                "--rpc.addr=0.0.0.0",
                "--rpc.port={}".format(node_rpc_port),
                "--p2p.disable",
                "--rpc.enable-admin",
                "--p2p.sequencer.key={}".format(env["GS_SEQUENCER_PRIVATE_KEY"]),            
            ],
            ports={
                "rpc": PortSpec(number=node_rpc_port),
            },
            files={"/config": files.config, "/rollup": files.rollup},
            env_vars=env,
        ),
    )