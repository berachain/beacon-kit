geth = import_module("./kurtosis/src/nodes/execution/geth/launcher.star")
shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")

def run(plan, args = {}):
    my_ports = {"rpc": 8545, "ws": 8546, "p2p": 30303, "metrics": 6060}
    plan.print(my_ports)
    my_ports['rpc'] = 9010
    plan.print(my_ports)
    my_ports['ws'] = shared_utils.new_port_spec(9999, shared_utils.TCP_PROTOCOL)
    plan.print(my_ports)
    my_ports['new'] = 10000
    plan.print(my_ports)
    my_ports['ws'] = shared_utils.new_port_spec(4200, shared_utils.TCP_PROTOCOL)
    plan.print(my_ports)
    my_ports['ws'] = 1703
    plan.print(my_ports)

    geth_sc = geth.get_default_service_config()
    plan.print(geth_sc)
    
