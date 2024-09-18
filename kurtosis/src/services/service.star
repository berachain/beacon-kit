def parse_service_from_dict(service):
    if "replicas" not in service:
        service["replicas"] = 1
    if "client" not in service:
        service["client"] = None
    return struct(
        name = service["name"],
        replicas = service["replicas"],
        client = service["client"],
    )
