def parse_service_from_dict(service):
    if "replicas" not in service:
        service["replicas"] = 1
    return struct(
        name = service["name"],
        replicas = service["replicas"],
    )
