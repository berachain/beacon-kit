CONSENSUS_DEFAULT_SETTINGS = {
    "min_cpu": 0,
    "max_cpu": 2000,
    "min_memory": 0,
    "max_memory": 2048,
    "images": {
        "beaconkit": "beacond:kurtosis-local",
    },
    "labels": {},
    "node_selectors": {},
}

EXECUTION_DEFAULT_SETTINGS = {
    "min_cpu": 0,
    "max_cpu": 2000,
    "min_memory": 0,
    "max_memory": 2048,
    "images": {
        "besu": "hyperledger/besu:latest",
        "erigon": "thorax/erigon:latest",
        "ethereumjs": "ethpandaops/ethereumjs:stable",
        "geth": "ethereum/client-go:latest",
        "nethermind": "nethermind/nethermind:latest",
        "reth": "ghcr.io/paradigmxyz/reth:latest",
    },
    "labels": {},
    "node_selectors": {},
}

CL_TYPE = "beaconkit"

def parse_nodes_from_dict(vals, settings):
    node_type = vals["type"]
    node_list = []

    consensus_settings = parse_node_settings(settings["consensus_settings"])
    execution_settings = parse_node_settings(settings["execution_settings"])

    count = 0
    for val_configuration in vals["nodes"]:
        replicas = 1
        if "replicas" in val_configuration:
            replicas = val_configuration["replicas"]
        if replicas != 0:
            for i in range(replicas):
                val_struct = parse_node_from_dict(node_type, val_configuration, consensus_settings, execution_settings, count)
                node_list.append(val_struct)
                count += 1

    return node_list

def parse_node_from_dict(node_type, val, consensus_settings, execution_settings, index):
    return struct(
        node_type = node_type,
        el_type = val["el_type"],
        el_image = execution_settings.images[val["el_type"]],
        cl_type = CL_TYPE,
        cl_image = consensus_settings.images[CL_TYPE],
        index = index,
        cl_service_name = "cl-{}-{}-{}".format(node_type, CL_TYPE, index),
        el_service_name = "el-{}-{}-{}".format(node_type, val["el_type"], index),
        consensus_settings = consensus_settings,
        execution_settings = execution_settings,
    )

def parse_node_settings(settings):
    node_settings = dict(settings)

    if "min_cpu" not in node_settings:
        node_settings["min_cpu"] = CONSENSUS_DEFAULT_SETTINGS["min_cpu"]
    if "max_cpu" not in node_settings:
        node_settings["max_cpu"] = CONSENSUS_DEFAULT_SETTINGS["max_cpu"]
    if "min_memory" not in node_settings:
        node_settings["min_memory"] = CONSENSUS_DEFAULT_SETTINGS["min_memory"]
    if "max_memory" not in node_settings:
        node_settings["max_memory"] = CONSENSUS_DEFAULT_SETTINGS["max_memory"]
    if "images" not in node_settings:
        node_settings["images"] = CONSENSUS_DEFAULT_SETTINGS["images"]
    if "labels" not in node_settings:
        node_settings["labels"] = CONSENSUS_DEFAULT_SETTINGS["labels"]
    if "node_selectors" not in node_settings:
        node_settings["node_selectors"] = CONSENSUS_DEFAULT_SETTINGS["node_selectors"]

    return struct(
        min_cpu = node_settings["min_cpu"],
        max_cpu = node_settings["max_cpu"],
        min_memory = node_settings["min_memory"],
        max_memory = node_settings["max_memory"],
        images = node_settings["images"],
        labels = node_settings["labels"],
        node_selectors = node_settings["node_selectors"],
    )
