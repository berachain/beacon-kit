builtins = import_module('./builtins.star')

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