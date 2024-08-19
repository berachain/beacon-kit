builtins = import_module("./builtins.star")

DEFAULT_TRANSPORT_PROTOCOL = "tcp"
DEFAULT_APPLICATION_PROTOCOL = ""
DEFAULT_WAIT = "15s"

def get_port_spec_template(
        number,
        transport_protocol = DEFAULT_TRANSPORT_PROTOCOL,
        application_protocol = DEFAULT_APPLICATION_PROTOCOL,
        wait = DEFAULT_WAIT):
    port_spec = {
        "number": number,
        "transport_protocol": transport_protocol,
        "application_protocol": application_protocol,
        "wait": wait,
    }

    validate_port_spec(port_spec)
    return port_spec

def validate_port_spec(port_spec):
    if type(port_spec) != builtins.types.dict:
        fail("Port spec must be a dict, not {0}".format(type(port_spec)))

    if type(port_spec["number"]) != builtins.types.int:
        fail("Port spec number must be an int, not {0}".format(type(port_spec["number"])))

    if port_spec["transport_protocol"] != None:
        if type(port_spec["transport_protocol"]) != builtins.types.string:
            fail("Port spec transport_protocol must be a string, not {0}".format(type(port_spec["transport_protocol"])))

    if port_spec["application_protocol"] != None:
        if type(port_spec["application_protocol"]) != builtins.types.string:
            fail("Port spec application_protocol must be a string, not {0}".format(type(port_spec["application_protocol"])))

    if port_spec["wait"] != None:
        if type(port_spec["wait"]) != builtins.types.string:
            fail("Port spec wait must be a bool, not {0}".format(type(port_spec["wait"])))

def create_port_specs_from_config(config):
    ports = {}
    for port_key, port_spec in config["ports"].items():
        ports[port_key] = create_port_spec(port_spec)

    return ports

def create_port_spec(port_spec_dict):
    return PortSpec(
        number = port_spec_dict["number"],
        transport_protocol = port_spec_dict["transport_protocol"],
        application_protocol = port_spec_dict["application_protocol"],
        wait = port_spec_dict["wait"],
    )
