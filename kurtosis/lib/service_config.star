"""
A library for manipulating Service Configs as dictionaries, prior to instantiating the actual ServiceConfig object

This is necessary because Starlark structs are largely immutable. To pass a config around and make manipulations to it
we represent the Service Config as a dictionary, and then convert it to a ServiceConfig object when we're ready to use it.
We additionally support validations to ensure that the dictionary is well-formed.

References:
    - https://docs.kurtosis.com/api-reference/starlark-reference/service-config/
    - https://github.com/kurtosis-tech/kurtosis/blob/473d0ee07f2b16c39cf9a453c3c28afdb1e2493d/core/server/api_container/server/startosis_engine/kurtosis_types/service_config/service_config.go
"""

builtins = import_module('./builtins.star')

def get_service_config_template(
    image,
    ports=None,
    files=None,
    entrypoint=None,
    cmd=None,
    env_vars=None,
    private_ip_address_placeholder=None,
    max_cpu=None,
    min_cpu=None,
    max_memory=None,
    min_memory=None,
    ready_conditions=None,
    labels=None,
    user=None,
    tolerations=None,
    node_selectors=None,
):
    service_config = {
        "image": image,
        "ports": ports,
        "files": files,
        "entrypoint": entrypoint,
        "cmd": cmd,
        "env_vars": env_vars,
        "private_ip_address_placeholder": private_ip_address_placeholder,
        "max_cpu": max_cpu,
        "min_cpu": min_cpu,
        "max_memory": max_memory,
        "min_memory": min_memory,
        "ready_conditions": ready_conditions,
        "labels": labels,
        "user": user,
        "tolerations": tolerations,
        "node_selectors": node_selectors,
    }

    validate_service_config_types(service_config)
    return service_config

def get_port_spec_template(
    number,
    transport=None,
    application_protocol=None,
    wait=None,
):
    port_spec = {
        "number": number,
        "transport": transport,
        "application_protocol": application_protocol,
        "wait": wait,
    }

    validate_port_spec(port_spec)
    return port_spec

def validate_service_config_types(service_config):
    if type(service_config) != builtins.types.dict:
        fail("Service config must be a dict, not {0}".format(type(service_config)))

    if type(service_config["image"]) != builtins.types.string:
        fail("Service config image must be a string, not {0}".format(type(service_config['image'])))

    if service_config["ports"] != None:
        if type(service_config["ports"]) != builtins.types.dict:
            fail("Service config ports must be a dict, not {0}".format(type(service_config['ports'])))
        for port_key, port_spec in service_config["ports"].items():
            if type(port_key) != builtins.types.string:
                fail("Service config port key must be an int, not {0}".format(type(port_key)))
            validate_port_spec(port_spec)

    if service_config["files"] != None:
        if type(service_config["files"]) != builtins.types.dict:
            fail("Service config files must be a dict, not {0}".format(type(service_config['files'])))
        for path, content in service_config["files"].items():
            if type(path) != builtins.types.string:
                fail("Service config file path must be a string, not {0}".format(type(path)))
            if type(content) not in [builtins.types.string, builtins.types.directory]:
                fail("Service config file content must be a string or a Directory object, not {0}".format(type(content)))

    if service_config["entrypoint"] != None:
        if type(service_config["entrypoint"]) != builtins.types.list:
            fail("Service config entrypoint must be a list, not {0}".format(type(service_config['entrypoint'])))
        for entrypoint in service_config["entrypoint"]:
            if type(entrypoint) != builtins.types.string:
                fail("Service config entrypoint must be a string, not {0}".format(type(entrypoint)))

    if service_config["cmd"] != None:
        if type(service_config["cmd"]) != builtins.types.list:
            fail("Service config cmd must be a list, not {0}".format(type(service_config['cmd'])))
        for cmd in service_config["cmd"]:
            if type(cmd) != builtins.types.string:
                fail("Service config cmd must be a string, not {0}".format(type(cmd)))

    if service_config["env_vars"] != None:
        if type(service_config["env_vars"]) != builtins.types.dict:
            fail("Service config env_vars must be a dict, not {0}".format(type(service_config["env_vars"])))
        for env_var_key, env_var_value in service_config["env_vars"].items():
            if type(env_var_key) != builtins.types.string:
                fail("Service config env_var key must be a string, not {0}".format(type(env_var_key)))
            if type(env_var_value) != builtins.types.string:
                fail("Service config env_var value must be a string, not {0}".format(type(env_var_value)))

    if service_config["private_ip_address_placeholder"] != None:
        if type(service_config["private_ip_address_placeholder"]) != builtins.types.string:
            fail("Service config private_ip_address_placeholder must be a string, not {0}".format(type(service_config['private_ip_address_placeholder'])))

    if service_config["max_cpu"] != None:
        if type(service_config["max_cpu"]) != builtins.types.int:
            fail("Service config max_cpu must be a int, not {0}".format(type(service_config['max_cpu'])))
    if service_config["min_cpu"] != None:
        if type(service_config["min_cpu"]) != builtins.types.int:
            fail("Service config min_cpu must be a int, not {0}".format(type(service_config['min_cpu'])))
    if service_config["max_memory"] != None:
        if type(service_config["max_memory"]) != builtins.types.int:
            fail("Service config max_memory must be a int, not {0}".format(type(service_config['max_memory'])))
    if service_config["min_memory"] != None:
        if type(service_config["min_memory"]) != builtins.types.int:
            fail("Service config min_memory must be a int, not {0}".format(type(service_config['min_memory'])))

    # TODO(validation): Implement validation for ready_conditions
    # TODO(validation): Implement validation for labels
    # TODO(validation): Implement validation for user
    # TODO(validation): Implement validation for tolerations
    # TODO(validation): Implement validation for node_selectors

def validate_port_spec(port_spec):
    if type(port_spec) != builtins.types.dict:
        fail("Port spec must be a dict, not {0}".format(type(port_spec)))

    if type(port_spec["number"]) != builtins.types.int:
        fail("Port spec number must be an int, not {0}".format(type(port_spec['number'])))

    if port_spec["transport"] != None:
        if type(port_spec["transport"]) != builtins.types.string:
            fail("Port spec transport must be a string, not {0}".format(type(port_spec['transport'])))

    if port_spec["application_protocol"] != None:
        if type(port_spec["application_protocol"]) != builtins.types.string:
            fail("Port spec application_protocol must be a string, not {0}".format(type(port_spec['application_protocol'])))

    if port_spec["wait"] != None:
        if type(port_spec["wait"]) != builtins.types.bool:
            fail("Port spec wait must be a bool, not {0}".format(type(port_spec['wait'])))