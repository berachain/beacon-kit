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
    "config": {
        "timeout_propose": "3s",
        "timeout_vote": "2s",
        "timeout_commit": "1s",
        "max_num_inbound_peers": 40,
        "max_num_outbound_peers": 10,
    },
    "app": {
        "payload-timeout": "1.5s",
    },
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

    consensus_settings = parse_consensus_settings(settings)
    execution_settings = parse_execution_settings(settings)

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

def parse_consensus_settings(settings):
    consensus_settings = {}
    if "consensus_settings" in settings:
        consensus_settings = dict(settings["consensus_settings"])
    consensus_settings = parse_default_node_settings(consensus_settings, CONSENSUS_DEFAULT_SETTINGS)
    consensus_settings = parse_extra_consensus_settings(consensus_settings)

    return get_consensus_struct(consensus_settings)

def parse_extra_consensus_settings(settings):
    config_settings = parse_consensus_config_settings(settings["config"]) if "config" in settings else parse_consensus_config_settings(CONSENSUS_DEFAULT_SETTINGS["config"])
    app_settings = parse_consensus_app_settings(settings["app"]) if "app" in settings else parse_consensus_app_settings(CONSENSUS_DEFAULT_SETTINGS["app"])

    consensus_settings = dict(settings)
    consensus_settings["config"] = config_settings
    consensus_settings["app"] = app_settings

    return consensus_settings

def parse_consensus_config_settings(config_settings):
    config_settings = dict(config_settings)

    if "timeout_propose" not in config_settings:
        config_settings["timeout_propose"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_propose"]
    if "timeout_vote" not in config_settings:
        config_settings["timeout_vote"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_vote"]
    if "timeout_commit" not in config_settings:
        config_settings["timeout_commit"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_commit"]
    if "max_num_inbound_peers" not in config_settings:
        config_settings["max_num_inbound_peers"] = CONSENSUS_DEFAULT_SETTINGS["config"]["max_num_inbound_peers"]
    if "max_num_outbound_peers" not in config_settings:
        config_settings["max_num_outbound_peers"] = CONSENSUS_DEFAULT_SETTINGS["config"]["max_num_outbound_peers"]

    return struct(
        timeout_propose = config_settings["timeout_propose"],
        timeout_vote = config_settings["timeout_vote"],
        timeout_commit = config_settings["timeout_commit"],
        max_num_inbound_peers = config_settings["max_num_inbound_peers"],
        max_num_outbound_peers = config_settings["max_num_outbound_peers"],
    )

def parse_consensus_app_settings(app_settings):
    app_settings = dict(app_settings)

    if "payload-timeout" not in app_settings:
        app_settings["payload-timeout"] = CONSENSUS_DEFAULT_SETTINGS["app"]["payload-timeout"]

    return struct(
        payload_timeout = app_settings["payload-timeout"],
    )

def parse_execution_settings(settings):
    execution_settings = {}
    if "execution_settings" in settings:
        execution_settings = dict(settings["execution_settings"])
    execution_settings = parse_default_node_settings(execution_settings, EXECUTION_DEFAULT_SETTINGS)
    return get_execution_struct(execution_settings)

def parse_default_node_settings(settings, default_settings):
    node_settings = dict(settings)

    if "min_cpu" not in node_settings:
        node_settings["min_cpu"] = default_settings["min_cpu"]
    if "max_cpu" not in node_settings:
        node_settings["max_cpu"] = default_settings["max_cpu"]
    if "min_memory" not in node_settings:
        node_settings["min_memory"] = default_settings["min_memory"]
    if "max_memory" not in node_settings:
        node_settings["max_memory"] = default_settings["max_memory"]
    if "images" not in node_settings:
        node_settings["images"] = default_settings["images"]
    if "labels" not in node_settings:
        node_settings["labels"] = default_settings["labels"]
    if "node_selectors" not in node_settings:
        node_settings["node_selectors"] = default_settings["node_selectors"]

    return node_settings

def get_consensus_struct(consensus_settings):
    return struct(
        min_cpu = consensus_settings["min_cpu"],
        max_cpu = consensus_settings["max_cpu"],
        min_memory = consensus_settings["min_memory"],
        max_memory = consensus_settings["max_memory"],
        images = consensus_settings["images"],
        labels = consensus_settings["labels"],
        node_selectors = consensus_settings["node_selectors"],
        config = consensus_settings["config"],
        app = consensus_settings["app"],
    )

def get_execution_struct(execution_settings):
    return struct(
        min_cpu = execution_settings["min_cpu"],
        max_cpu = execution_settings["max_cpu"],
        min_memory = execution_settings["min_memory"],
        max_memory = execution_settings["max_memory"],
        images = execution_settings["images"],
        labels = execution_settings["labels"],
        node_selectors = execution_settings["node_selectors"],
    )

