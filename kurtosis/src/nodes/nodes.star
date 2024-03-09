def parse_validators_from_dict(vals):
    validator_list = []
    for val_configuration in vals:
        val_struct = parse_validator_from_dict(val_configuration)
        if 'replicas' not in val_configuration:
            validator_list.append(val_struct)
        else:
            for i in range(val_configuration['replicas']):
                validator_list.append(val_struct)

    return validator_list
        

def parse_validator_from_dict(val):
    return struct(
        el_type = val['el_type'],
        cl_type = val['cl_type'],
        cl_image = val['cl_image'],
    )