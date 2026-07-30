#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
#
# Chaos-style verification of blob distribution after the blob-consensus transition, on a fully local devnet:
# NUM_VALS beacond validators (each paired with its own bera-reth) run from local binaries, no kurtosis and no
# docker. A background load submits blob transactions, so every few blocks carries blob sidecars that, past the
# devnet's blob-consensus enable height (2), no longer ride inside CometBFT blocks.
#
# After liveness and layout checks (post-transition blocks must carry exactly ONE consensus tx while blob
# sidecars are still served), the harness loops for DURATION seconds: pick a random validator (never cl0, the
# reference), stop its beacond while blob blocks keep flowing, restart it from the same data directory, and
# require that it catches up, drains its blob fetch queue (is_syncing back to false), and serves the exact same
# sidecar counts as the reference for every blob slot produced during the outage. Those sidecars can only have
# arrived via the blob reactor / EL fetch. Execution clients run throughout.
#
# Usage:
#   blob-sync-local-test.sh [--num-vals 4] [--duration 300] [--keep]
#
# Environment:
#   BEACOND_BIN  beacond binary (default: <repo>/build/bin/beacond, `make build` produces it)
#   RETH_BIN     bera-reth binary (default: bera-reth on PATH)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# ---------------------------------------------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------------------------------------------

NUM_VALS=4      # validators in the devnet; one may be down at a time, so 4 keeps >2/3 voting power live
DURATION=300    # seconds to keep running random restart cycles
KEEP=0          # keep RUN_DIR after a passing run instead of deleting it
MIN_OUTAGE=15   # minimum seconds a victim stays down
MAX_OUTAGE=35   # maximum seconds a victim stays down

while [[ $# -gt 0 ]]; do
    case "$1" in
        --num-vals) NUM_VALS="$2"; shift 2 ;;
        --duration) DURATION="$2"; shift 2 ;;
        --keep) KEEP=1; shift ;;
        *) echo "unknown argument: $1" >&2; exit 2 ;;
    esac
done

BEACOND_BIN="${BEACOND_BIN:-$REPO_ROOT/build/bin/beacond}"
RETH_BIN="${RETH_BIN:-$(command -v bera-reth || true)}"

CHAIN_SPEC_ARG="--beacon-kit.chain-spec devnet"             # devnet spec: blob-consensus enable-height = 2
BLOB_ENABLE_HEIGHT=2                                        # must match the devnet chain spec
CHAIN_ID="beacond-2061"                                     # CometBFT chain id of the devnet
DEPOSIT_AMOUNT=32000000000                                  # premined deposit per validator, in gwei
WITHDRAWAL_ADDRESS=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4  # withdrawal credentials for those deposits

# Prefunded devnet account (same keys as kurtosis/src/constants.star).
LOAD_ADDR=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
LOAD_KEY=fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
LOAD_RECIPIENT=0x000000000000000000000000000000000000dEaD
BLOB_INTERVAL=1   # seconds between blob transactions

ETH_GENESIS="$REPO_ROOT/testing/files/eth-genesis.json"
JWT="$REPO_ROOT/testing/files/jwt.hex"
KZG="$REPO_ROOT/testing/files/kzg-trusted-setup.json"
RUN_DIR="${RUN_DIR:-$REPO_ROOT/.tmp/blob-sync-test}"
LOG_DIR="$RUN_DIR/logs"
RUN_DIR_MARKER=".blob-sync-local-test"

# Per-node port layout, node i:
cl_rpc_port()  { echo $((26657 + $1 * 10)); }
cl_p2p_port()  { echo $((26656 + $1 * 10)); }
cl_api_port()  { echo $((3500 + $1 * 10)); }
el_http_port() { echo $((8545 + $1 * 10)); }
el_auth_port() { echo $((8551 + $1 * 10)); }
el_p2p_port()  { echo $((30303 + $1)); }

# ---------------------------------------------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------------------------------------------

log() { printf '\n==> %s\n' "$*"; }
fail() { echo "FAIL: $*" >&2; exit 1; }
now() { date +%s; }

# Command tracing: every state-changing command is echoed to the terminal via fd 3, so traces surface even from
# blocks whose output is redirected to log files (and from command substitutions). Read-only polling probes
# (curl/jq) are deliberately not traced to keep the output readable.
exec 3>&1
trace() { printf '    $ %s\n' "$*" >&3; }
run() { trace "$@"; "$@"; }

# wait_for <deadline-secs> <description> <command...>
wait_for() {
    local deadline=$1 desc=$2; shift 2
    local start
    start=$(now)
    until "$@"; do
        if (( $(now) - start > deadline )); then
            echo "timeout after ${deadline}s waiting for: $desc" >&2
            return 1
        fi
        sleep 1
    done
}

port_free() { ! (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null; }
proc_gone() { ! kill -0 "$1" 2>/dev/null; }

el_rpc() {
    curl -sf "http://127.0.0.1:$(el_http_port "$1")" -H 'Content-Type: application/json' \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"$2\",\"params\":${3:-[]},\"id\":1}"
}
el_ready() { el_rpc "$1" eth_chainId >/dev/null 2>&1; }

api()        { echo "http://127.0.0.1:$(cl_api_port "$1")"; }
head_slot()  { curl -sf --max-time 5 "$(api "$1")/eth/v1/beacon/headers/head" | jq -r '.data.header.message.slot' 2>/dev/null || echo 0; }
blob_count() { curl -sf --max-time 5 "$(api "$1")/eth/v1/beacon/blob_sidecars/$2" 2>/dev/null | jq '.data | length' 2>/dev/null || echo 0; }
is_syncing() { curl -sf --max-time 5 "$(api "$1")/eth/v1/node/syncing" 2>/dev/null | jq -r '.data.is_syncing' 2>/dev/null || echo unknown; }
comet_txs()  { curl -sf --max-time 5 "http://127.0.0.1:$(cl_rpc_port "$1")/block?height=$2" | jq '.result.block.data.txs | length'; }

node_live() { [ "$(head_slot "$1")" -ge "$2" ] 2>/dev/null; }
node_synced_past() {
    local i=$1 target=$2
    [ "$(head_slot "$i")" -ge "$target" ] 2>/dev/null && [ "$(is_syncing "$i")" = "false" ]
}

# start_cl <index> <logfile> — starts a validator, prints its pid.
start_cl() {
    local i=$1 logfile=$2
    local cl_args=(start $CHAIN_SPEC_ARG --home "$RUN_DIR/cl$i"
        --beacon-kit.engine.rpc-dial-url "http://127.0.0.1:$(el_auth_port "$i")"
        --beacon-kit.engine.jwt-secret-path "$JWT" --beacon-kit.kzg.trusted-setup-path "$KZG"
        --beacon-kit.node-api.enabled --beacon-kit.node-api.address "127.0.0.1:$(cl_api_port "$i")"
        --beacon-kit.logger.log-level info)
    trace "$BEACOND_BIN" "${cl_args[@]}"
    "$BEACOND_BIN" "${cl_args[@]}" >>"$LOG_DIR/$logfile" 2>&1 &
    echo $!
}

PIDS=()
cleanup() {
    local code=$?
    for pid in "${PIDS[@]:-}"; do
        kill "$pid" >/dev/null 2>&1 || true
    done
    local waited=0 alive=1
    while ((alive && waited < 15)); do
        alive=0
        for pid in "${PIDS[@]:-}"; do
            if kill -0 "$pid" 2>/dev/null; then alive=1; fi
        done
        if ((alive)); then
            sleep 1
            waited=$((waited + 1))
        fi
    done
    for pid in "${PIDS[@]:-}"; do
        kill -9 "$pid" >/dev/null 2>&1 || true
    done
    wait >/dev/null 2>&1 || true
    if [[ $code -ne 0 ]]; then
        echo "FAILED (exit $code). Logs are in $LOG_DIR" >&2
    elif [[ $KEEP -eq 0 ]]; then
        remove_run_dir
    fi
}

# ---------------------------------------------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------------------------------------------

check_preconditions() {
    [[ -x "$BEACOND_BIN" ]] || fail "beacond not executable: $BEACOND_BIN (run 'make build' or set BEACOND_BIN)"
    [[ -n "$RETH_BIN" && -x "$RETH_BIN" ]] || fail "bera-reth not found; set RETH_BIN"
    local tool
    for tool in jq curl cast; do
        command -v "$tool" >/dev/null || fail "missing tool: $tool"
    done
    ((NUM_VALS >= 4)) || fail "need at least 4 validators so one can be down without losing >2/3 voting power"

    local busy="" i port
    for ((i = 0; i < NUM_VALS; i++)); do
        for port in "$(cl_rpc_port "$i")" "$(cl_p2p_port "$i")" "$(cl_api_port "$i")" \
            "$(el_http_port "$i")" "$(el_auth_port "$i")" "$(el_p2p_port "$i")"; do
            port_free "$port" || busy+=" $port"
        done
    done
    [[ -z "$busy" ]] || fail "ports already in use:$busy (a devnet or previous --keep run is still up)"
}

remove_run_dir() {
    [[ -e "$RUN_DIR" ]] || return 0
    if [[ ! -f "$RUN_DIR/$RUN_DIR_MARKER" && "$RUN_DIR" != "$REPO_ROOT/.tmp/"* ]] &&
        [[ -n "$(ls -A "$RUN_DIR" 2>/dev/null)" ]]; then
        fail "refusing to delete non-empty $RUN_DIR: no $RUN_DIR_MARKER marker, was it created by this harness?"
    fi
    rm -rf "$RUN_DIR"
}

prepare_run_dir() {
    log "Preparing run directory $RUN_DIR"
    remove_run_dir
    mkdir -p "$RUN_DIR" "$LOG_DIR"
    touch "$RUN_DIR/$RUN_DIR_MARKER"
}

# Premined-deposit genesis ceremony: every home is initialized, cl0 aggregates the deposits of all homes and
# ends up with the final genesis.json (copied back to every home) and an eth-genesis.json with the deposit
# storage set. Same flow as kurtosis/src/nodes/consensus/beacond/scripts/multiple-premined-deposits-cl.sh.
run_genesis_ceremony() {
    log "Genesis ceremony for $NUM_VALS validators"
    local i lead="$RUN_DIR/cl0" home
    {
        for ((i = 0; i < NUM_VALS; i++)); do
            run "$BEACOND_BIN" init "val$i" --chain-id "$CHAIN_ID" --home "$RUN_DIR/cl$i" $CHAIN_SPEC_ARG
            run "$BEACOND_BIN" genesis add-premined-deposit "$DEPOSIT_AMOUNT" "$WITHDRAWAL_ADDRESS" \
                --home "$RUN_DIR/cl$i" $CHAIN_SPEC_ARG
        done
        for ((i = 1; i < NUM_VALS; i++)); do
            run cp "$RUN_DIR/cl$i"/config/premined-deposits/premined-deposit* "$lead/config/premined-deposits/"
        done
        run "$BEACOND_BIN" genesis collect-premined-deposits --home "$lead" $CHAIN_SPEC_ARG
        run "$BEACOND_BIN" genesis set-deposit-storage "$ETH_GENESIS" --home "$lead" $CHAIN_SPEC_ARG
        run "$BEACOND_BIN" genesis execution-payload "$lead/eth-genesis.json" --home "$lead" $CHAIN_SPEC_ARG
        for ((i = 1; i < NUM_VALS; i++)); do
            run cp "$lead/config/genesis.json" "$RUN_DIR/cl$i/config/genesis.json"
        done
    } >"$LOG_DIR/genesis-ceremony.log" 2>&1
}

wire_cometbft_configs() {
    log "Wiring CometBFT configs (ports and persistent peers)"
    local i j peers cfg node_ids=()
    for ((i = 0; i < NUM_VALS; i++)); do
        node_ids[i]="$(run "$BEACOND_BIN" comet show-node-id --home "$RUN_DIR/cl$i")"
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        peers=""
        for ((j = 0; j < NUM_VALS; j++)); do
            [[ $i -eq $j ]] && continue
            peers+="${peers:+,}${node_ids[j]}@127.0.0.1:$(cl_p2p_port "$j")"
        done
        cfg="$RUN_DIR/cl$i/config/config.toml"
        trace "sed -i (rpc/p2p ports, persistent_peers) $cfg"
        sed -i.bak \
            -e "s|^laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://127.0.0.1:$(cl_rpc_port "$i")\"|" \
            -e "s|^laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://127.0.0.1:$(cl_p2p_port "$i")\"|" \
            -e "s|^persistent_peers = \"\"|persistent_peers = \"$peers\"|" \
            -e 's|^addr_book_strict = true|addr_book_strict = false|' \
            -e 's|^allow_duplicate_ip = false|allow_duplicate_ip = true|' \
            -e 's|^pprof_laddr = .*|pprof_laddr = ""|' \
            "$cfg"
        rm -f "$cfg.bak"
    done
}

start_execution_clients() {
    log "Starting $NUM_VALS execution clients (bera-reth)"
    local i
    local reth_args
    for ((i = 0; i < NUM_VALS; i++)); do
        run "$RETH_BIN" init --datadir "$RUN_DIR/el$i" --chain "$RUN_DIR/cl0/eth-genesis.json" >"$LOG_DIR/el$i.log" 2>&1
        reth_args=(node --datadir "$RUN_DIR/el$i" --chain "$RUN_DIR/cl0/eth-genesis.json"
            --http --http.addr 127.0.0.1 --http.port "$(el_http_port "$i")" --http.api admin,eth,net,web3
            --authrpc.addr 127.0.0.1 --authrpc.port "$(el_auth_port "$i")" --authrpc.jwtsecret "$JWT"
            --port "$(el_p2p_port "$i")" --discovery.port "$(el_p2p_port "$i")"
            --txpool.max-account-slots 1000 --ipcdisable)
        trace "$RETH_BIN" "${reth_args[@]}"
        "$RETH_BIN" "${reth_args[@]}" >>"$LOG_DIR/el$i.log" 2>&1 &
        PIDS+=($!)
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        wait_for 60 "el$i JSON-RPC endpoint" el_ready "$i"
    done

    log "Peering execution clients"
    local enode0
    enode0="$(el_rpc 0 admin_nodeInfo | jq -r '.result.enode' | sed -E 's/@[^:@?]+:/@127.0.0.1:/')"
    for ((i = 1; i < NUM_VALS; i++)); do
        run el_rpc "$i" admin_addPeer "[\"$enode0\"]" >/dev/null
    done
}

start_validators() {
    log "Starting $NUM_VALS validators"
    local i
    CL_PIDS=()
    for ((i = 0; i < NUM_VALS; i++)); do
        CL_PIDS[i]="$(start_cl "$i" "cl$i.log")"
        PIDS+=("${CL_PIDS[i]}")
    done
}

# Submits a steady stream of blob transactions to el0 for the rest of the run. The pending nonce is re-fetched
# per send so the stream self-heals after a rejected or dropped transaction.
start_blob_load() {
    log "Starting blob tx load (one blob tx to el0 every ${BLOB_INTERVAL}s)"
    local rpc="http://127.0.0.1:$(el_http_port 0)"
    run head -c 4096 /dev/urandom >"$RUN_DIR/blob.bin"
    trace "cast send $LOAD_RECIPIENT --blob --path $RUN_DIR/blob.bin --value 1 --nonce <pending> --async --rpc-url $rpc  # every ${BLOB_INTERVAL}s"
    (
        while :; do
            nonce="$(cast nonce --block pending "$LOAD_ADDR" --rpc-url "$rpc" 2>/dev/null)" || nonce=""
            if [[ -n "$nonce" ]]; then
                cast send "$LOAD_RECIPIENT" --blob --path "$RUN_DIR/blob.bin" --value 1 \
                    --private-key "$LOAD_KEY" --nonce "$nonce" --async --rpc-url "$rpc" \
                    >/dev/null 2>>"$LOG_DIR/blob-load.log" || true
            fi
            sleep "$BLOB_INTERVAL"
        done
    ) &
    PIDS+=($!)
    disown
}

# ---------------------------------------------------------------------------------------------------------------
# Checks
# ---------------------------------------------------------------------------------------------------------------

# scan_blob_slots <node> <from> <to> — echoes "slot:count" for every blob slot in range.
scan_blob_slots() {
    local s c
    for ((s = $2; s <= $3; s++)); do
        c="$(blob_count "$1" "$s")"
        [ "$c" -gt 0 ] && echo "$s:$c"
    done
    true
}

check_layout_and_blobs() {
    log "Waiting for chain liveness"
    local i
    for ((i = 0; i < NUM_VALS; i++)); do
        wait_for 120 "cl$i to reach slot 3" node_live "$i" 3
    done

    log "Waiting for blob production"
    blob_seen=""
    local t s
    for ((t = 0; t < 90; t++)); do
        s="$(head_slot 0)"
        [[ -n "$(scan_blob_slots 0 $((s > 8 ? s - 8 : 2)) "$s")" ]] && { blob_seen=1; break; }
        sleep 2
    done
    [[ -n "$blob_seen" ]] || fail "no blob blocks observed, see $LOG_DIR/blob-load.log"

    log "Checking tx layout (post-transition blocks must carry exactly one consensus tx)"
    local h ntx
    for h in "$BLOB_ENABLE_HEIGHT" "$(head_slot 0)"; do
        ntx="$(comet_txs 0 "$h")"
        [[ "$ntx" == "1" ]] || fail "block $h carries $ntx txs (want 1)"
        echo "    block $h: 1 tx ok"
    done
    echo "    blob sidecars are being produced and served while blocks carry a single tx"
}

chaos_loop() {
    local deadline=$(( $(now) + DURATION )) iteration=0
    local victim outage from to v r pair s want got outage_blobs

    while (( $(now) < deadline )); do
        iteration=$((iteration + 1))
        victim=$((1 + RANDOM % (NUM_VALS - 1)))   # never cl0, the reference
        outage=$(( MIN_OUTAGE + RANDOM % (MAX_OUTAGE - MIN_OUTAGE + 1) ))
        log "[$iteration] stopping cl$victim for ~${outage}s"

        from=$(( $(head_slot 0) + 1 ))
        run kill "${CL_PIDS[victim]}" 2>/dev/null || true
        wait_for 30 "cl$victim process to exit" proc_gone "${CL_PIDS[victim]}"
        sleep "$outage"
        to="$(head_slot 0)"

        CL_PIDS[victim]="$(start_cl "$victim" "cl$victim.log")"
        PIDS+=("${CL_PIDS[victim]}")
        echo "    restarted; outage covered slots $from..$to"

        outage_blobs=""
        (( to >= from )) && outage_blobs="$(scan_blob_slots 0 "$from" "$to")"
        [[ -n "$outage_blobs" ]] || echo "    (no blob slots during this outage; still checking catch-up)"

        # Catch up to the reference head, then require fully-synced (which includes an empty blob fetch queue).
        wait_for 180 "cl$victim to catch up and finish blob fetching" node_synced_past "$victim" "$(head_slot 0)"
        echo "    caught up and reports synced (slot $(head_slot "$victim"))"

        for pair in $outage_blobs; do
            s=${pair%%:*}; want=${pair##*:}
            got="$(blob_count "$victim" "$s")"
            [[ "$got" == "$want" ]] || fail "[$iteration] slot $s: cl$victim has $got sidecars, cl0 has $want"
            echo "    ok: slot $s has $got sidecars on both nodes"
        done
    done

    ITERATIONS=$iteration
}

print_result() {
    log "RESULT"
    local i
    for ((i = 0; i < NUM_VALS; i++)); do
        echo "  cl$i: slot=$(head_slot "$i") syncing=$(is_syncing "$i")"
    done
    echo
    echo "  blob fetch activity on restarted nodes (sample):"
    grep -hE 'Fetched blob sidecars by range|Fetched and stored blob sidecars' "$LOG_DIR"/cl*.log | tail -5 || true
    echo
    echo "PASS: $ITERATIONS random restart cycle(s) survived; blob distribution and sync verified"
}

# ---------------------------------------------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------------------------------------------

main() {
    check_preconditions
    trap cleanup EXIT

    prepare_run_dir
    run_genesis_ceremony
    wire_cometbft_configs
    start_execution_clients
    start_validators
    start_blob_load

    check_layout_and_blobs
    chaos_loop
    print_result
}

main
