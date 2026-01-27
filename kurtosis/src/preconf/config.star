# Preconf configuration module for generating whitelist and JWT files.
# This module handles the setup of preconfirmation infrastructure in the devnet.

# Default port for the preconf API server
DEFAULT_PRECONF_API_PORT = 9090

# JWT secret length (32 bytes in hex = 64 chars + 0x prefix)
JWT_SECRET_LENGTH = 66

def generate_preconf_config(plan, validators, sequencer_index = 0):
    """
    Generate preconfirmation configuration files for the devnet.

    This creates:
    1. A whitelist.json file containing all validator pubkeys
    2. A validator-jwts.json file mapping pubkeys to JWT secrets (for sequencer)
    3. Individual JWT secret files for each validator

    Args:
        plan: The Kurtosis plan object
        validators: List of validator structs from the node parsing
        sequencer_index: Index of the validator that will run as sequencer (default: 0)

    Returns:
        A struct containing:
        - whitelist_file: Artifact name for the whitelist JSON
        - validator_jwts_file: Artifact name for the validator JWTs mapping
        - validator_jwt_secrets: Dict mapping validator index to JWT secret artifact
        - sequencer_index: The index of the sequencer validator
        - sequencer_url: The URL validators should use to connect to sequencer
    """
    num_validators = len(validators)

    # Generate JWT secrets for each validator (deterministic for reproducibility)
    # In production, these should be randomly generated and securely distributed
    # Pre-defined secrets for up to 10 validators (32 bytes hex each)
    predefined_secrets = [
        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
        "0x234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1",
        "0x34567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12",
        "0x4567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123",
        "0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234",
        "0x67890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345",
        "0x7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456",
        "0x890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567",
        "0x90abcdef1234567890abcdef1234567890abcdef1234567890abcdef12345678",
        "0x0abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456789",
    ]
    jwt_secrets = []
    for i in range(num_validators):
        if i < len(predefined_secrets):
            jwt_secrets.append(predefined_secrets[i])
        else:
            # Fallback for more than 10 validators - use base secret with index suffix
            jwt_secrets.append("0xabcdef00000000000000000000000000000000000000000000000000000000" + str(i))

    # Create whitelist.json content
    # We'll use placeholder pubkeys that will be replaced after genesis ceremony
    # For now, create a template that marks all validators as whitelisted
    whitelist_content = "[]"  # Will be populated after we get actual pubkeys

    # Create validator-jwts.json content (for sequencer)
    # This maps validator pubkeys to their JWT secrets
    validator_jwts_content = "{}"  # Will be populated after we get actual pubkeys

    # Store the JWT secrets as individual files for each validator
    validator_jwt_artifacts = {}
    for i in range(num_validators):
        jwt_content = jwt_secrets[i] + "\n"
        artifact_name = "validator-jwt-secret-{}".format(i)

        artifact = plan.render_templates(
            config = {
                "jwt-secret.hex": struct(
                    template = jwt_content,
                    data = {},
                ),
            },
            name = artifact_name,
            description = "Creating JWT secret for validator {}".format(i),
        )
        validator_jwt_artifacts[i] = artifact

    # Get the sequencer service name for URL construction
    sequencer_service_name = validators[sequencer_index].cl_service_name

    return struct(
        jwt_secrets = jwt_secrets,
        validator_jwt_artifacts = validator_jwt_artifacts,
        sequencer_index = sequencer_index,
        sequencer_service_name = sequencer_service_name,
        num_validators = num_validators,
    )

def create_preconf_files_from_configs(plan, num_validators, jwt_secrets):
    """
    Create whitelist.json and validator-jwts.json by extracting pubkeys from config artifacts.

    This function runs in the EXECUTION phase, extracting actual pubkeys from the
    premined deposit files and generating the preconf config files.

    Args:
        plan: The Kurtosis plan object
        num_validators: Number of validators
        jwt_secrets: List of JWT secret strings (0x-prefixed hex)

    Returns:
        Struct with whitelist_artifact and validator_jwts_artifact
    """

    # Build the files mount dictionary - mount all validator config artifacts
    files_mount = {}
    for i in range(num_validators):
        files_mount["/config{}".format(i)] = "node-beacond-config-{}".format(i)

    # Build JWT secrets as a shell variable
    jwt_secrets_str = " ".join(jwt_secrets)

    # Shell script that extracts pubkeys and generates both JSON files.
    # Uses jq for all JSON generation to avoid shell escaping issues.
    # Note: Using single-line commands to avoid shell line continuation issues in Kurtosis.
    script = "mkdir -p /out && touch /tmp/pubkeys.txt && "

    # Build the loop to extract pubkeys
    for i in range(num_validators):
        script += "jq -r '.pubkey' /config{idx}/.beacond/config/premined-deposits/premined-deposit-*.json >> /tmp/pubkeys.txt 2>/dev/null || true && ".format(idx = i)

    # Generate whitelist.json and validator-jwts.json
    script += "jq -R -s 'split(\"\\n\") | map(select(length > 0))' /tmp/pubkeys.txt > /out/whitelist.json && "
    script += "echo '{jwt_secrets}' | tr ' ' '\\n' > /tmp/jwts.txt && ".format(jwt_secrets = jwt_secrets_str)
    script += "jq -n --rawfile pks /tmp/pubkeys.txt --rawfile jwts /tmp/jwts.txt '($pks | split(\"\\n\") | map(select(length > 0))) as $pubkeys | ($jwts | split(\"\\n\") | map(select(length > 0))) as $secrets | [range($pubkeys | length)] | map({($pubkeys[.]): $secrets[.]}) | add' > /out/validator-jwts.json && "
    script += "cat /out/whitelist.json && cat /out/validator-jwts.json"

    # Run the script and store the output files
    result = plan.run_sh(
        run = script,
        image = "badouralix/curl-jq",
        files = files_mount,
        store = [
            StoreSpec(src = "/out/whitelist.json", name = "preconf-whitelist"),
            StoreSpec(src = "/out/validator-jwts.json", name = "preconf-validator-jwts"),
        ],
        description = "Extracting pubkeys and generating preconf config files",
    )

    return struct(
        whitelist_artifact = result.files_artifacts[0],
        validator_jwts_artifact = result.files_artifacts[1],
    )

def get_preconf_start_flags(is_sequencer, sequencer_url = "", preconf_enabled = True):
    """
    Get the preconf-related flags for the beacond start command.

    Args:
        is_sequencer: Whether this node runs as the sequencer
        sequencer_url: URL of the sequencer's preconf API (for non-sequencer validators)
        preconf_enabled: Whether preconf is enabled at all

    Returns:
        String of CLI flags to append to beacond start command
    """
    if not preconf_enabled:
        return ""

    flags = "--beacon-kit.preconf.enabled=true"

    if is_sequencer:
        # Sequencer mode flags
        flags += " --beacon-kit.preconf.sequencer-mode=true"
        flags += " --beacon-kit.preconf.whitelist-path=/root/preconf/whitelist.json"
        flags += " --beacon-kit.preconf.validator-jwts-path=/root/preconf/validator-jwts.json"
        flags += " --beacon-kit.preconf.api-port={}".format(DEFAULT_PRECONF_API_PORT)
    else:
        # Validator mode flags (fetch from sequencer)
        if sequencer_url != "":
            flags += " --beacon-kit.preconf.sequencer-url={}".format(sequencer_url)
            flags += " --beacon-kit.preconf.sequencer-jwt-path=/root/preconf/jwt-secret.hex"
            flags += " --beacon-kit.preconf.fetch-timeout=2s"

    return flags
