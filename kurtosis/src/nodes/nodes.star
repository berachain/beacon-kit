def parse_nodes_from_dict(vals):
    node_list = []
    for val_configuration in vals:
        val_struct = parse_node_from_dict(val_configuration)
        if "replicas" not in val_configuration:
            node_list.append(val_struct)
        else:
            for i in range(val_configuration["replicas"]):
                node_list.append(val_struct)

    return node_list

def parse_node_from_dict(val):
    return struct(
        el_type = val["el_type"],
        cl_type = val["cl_type"],
        cl_image = val["cl_image"],
    )
