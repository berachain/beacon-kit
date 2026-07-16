#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
#
# Tests a coordinated network upgrade end to end: a local devnet of NUM_VALS validators halts itself at a
# configured halt point, every beacond is swapped for a new binary, and the chain must resume from the same
# data directories. Only beacond is swapped, each validator's bera-reth execution client runs throughout.
#
# The test asserts that no validator ever commits past the halt point, that a halted node restarted with the
# halt flag still set refuses to start, and that transactions keep landing both before the halt and after the
# swap. See main() at the bottom for the phase-by-phase flow.
#
# Usage:
#   halt-swap-resume-test.sh [--num-vals 4] [--halt-height 8] [--halt-time-offset 0] [--keep] [--no-load]
#
#   --halt-time-offset N halts on --halt-time (now + N seconds) instead of --halt-height.
#
# Required environment:
#   OLD_BIN   beacond the chain starts and halts on
#   NEW_BIN   beacond the chain resumes on, may equal OLD_BIN to smoke-test the halt/restart mechanics alone
#   RETH_BIN  bera-reth used as every validator's execution client

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# ---------------------------------------------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------------------------------------------

NUM_VALS=7           # validators in the devnet, each paired with its own execution client
HALT_HEIGHT=10       # halt point for --halt-height mode (ignored when HALT_TIME_OFFSET > 0)
HALT_TIME_OFFSET=0   # when > 0, halt on --halt-time now+offset instead; the halt height is derived afterwards
KEEP=0               # keep RUN_DIR after a passing run instead of deleting it
LOAD=1               # run the background tx load and assert txs landed on both sides of the swap
RESUME_BLOCKS=5      # blocks the network must produce past the halt height after the swap

while [[ $# -gt 0 ]]; do
    case "$1" in
        --num-vals) NUM_VALS="$2"; shift 2 ;;
        --halt-height) HALT_HEIGHT="$2"; shift 2 ;;
        --halt-time-offset) HALT_TIME_OFFSET="$2"; shift 2 ;;
        --keep) KEEP=1; shift ;;
        --no-load) LOAD=0; shift ;;
        *) echo "unknown argument: $1" >&2; exit 2 ;;
    esac
done

: "${OLD_BIN:?set OLD_BIN to the current-fork beacond binary}"
: "${NEW_BIN:?set NEW_BIN to the upgraded beacond binary}"
: "${RETH_BIN:?set RETH_BIN to a bera-reth binary}"

CHAIN_SPEC_ARG="--beacon-kit.chain-spec devnet"             # chain spec passed to every beacond invocation
CHAIN_ID="beacond-2061"                                     # CometBFT chain id of the devnet
DEPOSIT_AMOUNT=32000000000                                  # premined deposit per validator, in gwei
WITHDRAWAL_ADDRESS=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4  # withdrawal credentials for those deposits

# Prefunded devnet account (same keys as kurtosis/src/constants.star).
LOAD_ADDR=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4                        # tx-load sender
LOAD_KEY=fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306   # its private key
LOAD_RECIPIENT=0x000000000000000000000000000000000000dEaD                   # transfer recipient
LOAD_INTERVAL=0.4                                                           # seconds between transfers

ETH_GENESIS="$REPO_ROOT/testing/files/eth-genesis.json"  # execution genesis every EL is initialized from
JWT="$REPO_ROOT/testing/files/jwt.hex"                   # engine-API JWT secret shared by all CL/EL pairs
KZG="$REPO_ROOT/testing/files/kzg-trusted-setup.json"    # KZG trusted setup for beacond
RUN_DIR="${RUN_DIR:-$REPO_ROOT/.tmp/upgrade-test}"       # node homes and datadirs, overridable via env
LOG_DIR="$RUN_DIR/logs"                                  # per-node and per-phase log files
RUN_DIR_MARKER=".halt-swap-resume-test"                  # marks RUN_DIR as safe for this harness to delete

# Budget for the network to reach the halt point: fixed startup slack plus consensus time.
if ((HALT_TIME_OFFSET > 0)); then
    HALT_TIMEOUT_SECS=$((120 + HALT_TIME_OFFSET))
else
    HALT_TIMEOUT_SECS=$((120 + HALT_HEIGHT * 5))
fi
RESUME_TIMEOUT_SECS=240  # budget for all nodes to reach TARGET_HEIGHT after the swap

# Per-node port layout, node i:
cl_rpc_port()  { echo $((26657 + $1 * 10)); }
cl_p2p_port()  { echo $((26656 + $1 * 10)); }
el_http_port() { echo $((8545 + $1 * 10)); }
el_auth_port() { echo $((8551 + $1 * 10)); }
el_p2p_port()  { echo $((30303 + $1)); }

# ---------------------------------------------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------------------------------------------

log() { printf '\n==> %s\n' "$*"; }
fail() { echo "$*" >&2; exit 1; }

# wait_for <deadline-secs> <description> <command...>
wait_for() {
    local deadline=$1 desc=$2; shift 2
    local start elapsed
    start=$(date +%s)
    until "$@"; do
        elapsed=$(( $(date +%s) - start ))
        if (( elapsed > deadline )); then
            echo "timeout after ${deadline}s waiting for: $desc" >&2
            return 1
        fi
        sleep 1
    done
}

port_free() { ! (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null; }
proc_gone() { ! kill -0 "$1" 2>/dev/null; }

# el_rpc <index> <method> [params-json]
el_rpc() {
    curl -sf "http://127.0.0.1:$(el_http_port "$1")" -H 'Content-Type: application/json' \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"$2\",\"params\":${3:-[]},\"id\":1}"
}
el_ready() { el_rpc "$1" eth_chainId >/dev/null 2>&1; }

cl_status() { curl -sf "http://127.0.0.1:$(cl_rpc_port "$1")/status"; }
node_height() { cl_status "$1" | jq -r '.result.sync_info.latest_block_height' 2>/dev/null || echo 0; }
node_p2p_version() { cl_status "$1" | jq -r '.result.node_info.protocol_version.p2p' 2>/dev/null || echo "?"; }
node_resumed() {
    local h
    h="$(node_height "$1")"
    [[ "$h" =~ ^[0-9]+$ ]] && (( h >= TARGET_HEIGHT ))
}

# start_cl <index> <binary> <logfile> [extra args...] — starts a validator, prints its pid.
start_cl() {
    local i=$1 bin=$2 logfile=$3; shift 3
    "$bin" start $CHAIN_SPEC_ARG --home "$RUN_DIR/cl$i" \
        --beacon-kit.engine.rpc-dial-url "http://127.0.0.1:$(el_auth_port "$i")" \
        --beacon-kit.engine.jwt-secret-path "$JWT" --beacon-kit.kzg.trusted-setup-path "$KZG" \
        --beacon-kit.logger.log-level info "$@" >"$LOG_DIR/$logfile" 2>&1 &
    echo $!
}

# Validators that exited after logging the halt line, respectively ones still running.
halted_count() {
    local n=0 k
    for ((k = 0; k < NUM_VALS; k++)); do
        if ! kill -0 "${CL_PIDS[k]}" 2>/dev/null &&
            grep -q "halting node per configuration" "$LOG_DIR/cl$k.old.log"; then
            n=$((n + 1))
        fi
    done
    echo "$n"
}
any_halted() { (($(halted_count) >= 1)); }
alive_count() {
    local n=0 k
    for ((k = 0; k < NUM_VALS; k++)); do
        if kill -0 "${CL_PIDS[k]}" 2>/dev/null; then n=$((n + 1)); fi
    done
    echo "$n"
}

last_committed_height() {
    # `|| true` so a log without a match yields "" under pipefail; the caller fails loudly on an empty result.
    grep -o "Committed state.*height=[0-9]*" "$LOG_DIR/cl$1.old.log" | grep -o "[0-9]*$" | tail -1 || true
}

# el_tx_count_range <first> <last> — load txs in el0 blocks first..last, filtered by sender since every block
# also carries a system transaction.
el_tx_count_range() {
    local total=0 h n
    for ((h = $1; h <= $2; h++)); do
        n="$(el_rpc 0 eth_getBlockByNumber "[\"$(printf '0x%x' "$h")\",true]" | jq -r --arg from "$LOAD_ADDR" \
            '[.result.transactions[] | select(.from == $from)] | length' 2>/dev/null)" || n=0
        total=$((total + ${n:-0}))
    done
    echo "$total"
}

PIDS=()
cleanup() {
    local code=$?
    for pid in "${PIDS[@]:-}"; do
        kill "$pid" >/dev/null 2>&1 || true
    done
    # Bounded shutdown, escalating to SIGKILL so one wedged child cannot hang the run (or its CI job).
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
    local bin tool
    for bin in "$OLD_BIN" "$NEW_BIN" "$RETH_BIN"; do
        [[ -x "$bin" ]] || fail "not executable: $bin"
    done
    for tool in jq curl; do
        command -v "$tool" >/dev/null || fail "missing tool: $tool"
    done
    if ((LOAD)) && ! command -v cast >/dev/null; then
        fail "missing tool: cast (foundry) is required for tx load, or pass --no-load"
    fi

    # Node 0 uses the standard devnet ports, so a leftover devnet (or --keep run) would answer our readiness
    # probes and the test would silently run against the wrong chain. Fail fast if any port is already taken.
    local busy="" i port
    for ((i = 0; i < NUM_VALS; i++)); do
        for port in "$(cl_rpc_port "$i")" "$(cl_p2p_port "$i")" \
            "$(el_http_port "$i")" "$(el_auth_port "$i")" "$(el_p2p_port "$i")"; do
            port_free "$port" || busy+=" $port"
        done
    done
    [[ -z "$busy" ]] || fail "ports already in use:$busy (stop the running devnet or previous --keep run)"
}

# RUN_DIR is env-overridable and gets wiped, so only delete a directory this harness marked as its own,
# the default location under the repo's .tmp, or an empty one. Anything else is not ours to destroy.
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

# Every validator adds a premined deposit, cl0 aggregates them into the final genesis.json (copied back to
# every home) and writes the eth-genesis with the deposit storage set.
run_genesis_ceremony() {
    log "Genesis ceremony for $NUM_VALS validators"
    local i
    for ((i = 0; i < NUM_VALS; i++)); do
        "$OLD_BIN" init "val$i" --chain-id "$CHAIN_ID" --home "$RUN_DIR/cl$i" \
            $CHAIN_SPEC_ARG >"$LOG_DIR/init-cl$i.log" 2>&1
        "$OLD_BIN" genesis add-premined-deposit "$DEPOSIT_AMOUNT" "$WITHDRAWAL_ADDRESS" \
            --home "$RUN_DIR/cl$i" $CHAIN_SPEC_ARG >>"$LOG_DIR/init-cl$i.log" 2>&1
    done
    {
        for ((i = 1; i < NUM_VALS; i++)); do
            cp "$RUN_DIR/cl$i/config/premined-deposits/premined-deposit"* "$RUN_DIR/cl0/config/premined-deposits/"
        done
        "$OLD_BIN" genesis collect-premined-deposits --home "$RUN_DIR/cl0" $CHAIN_SPEC_ARG
        "$OLD_BIN" genesis set-deposit-storage "$ETH_GENESIS" --home "$RUN_DIR/cl0" $CHAIN_SPEC_ARG
        "$OLD_BIN" genesis execution-payload "$RUN_DIR/cl0/eth-genesis.json" --home "$RUN_DIR/cl0" $CHAIN_SPEC_ARG
        for ((i = 1; i < NUM_VALS; i++)); do
            cp "$RUN_DIR/cl0/config/genesis.json" "$RUN_DIR/cl$i/config/genesis.json"
        done
    } >"$LOG_DIR/genesis-ceremony.log" 2>&1
}

wire_cometbft_configs() {
    log "Wiring CometBFT configs (ports and persistent peers)"
    local i j peers cfg node_ids=()
    for ((i = 0; i < NUM_VALS; i++)); do
        node_ids[i]="$("$OLD_BIN" comet show-node-id --home "$RUN_DIR/cl$i")"
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        peers=""
        for ((j = 0; j < NUM_VALS; j++)); do
            [[ $i -eq $j ]] && continue
            peers+="${peers:+,}${node_ids[j]}@127.0.0.1:$(cl_p2p_port "$j")"
        done
        cfg="$RUN_DIR/cl$i/config/config.toml"
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
    for ((i = 0; i < NUM_VALS; i++)); do
        "$RETH_BIN" init --datadir "$RUN_DIR/el$i" --chain "$RUN_DIR/cl0/eth-genesis.json" >"$LOG_DIR/el$i.log" 2>&1
        "$RETH_BIN" node --datadir "$RUN_DIR/el$i" --chain "$RUN_DIR/cl0/eth-genesis.json" \
            --http --http.addr 127.0.0.1 --http.port "$(el_http_port "$i")" --http.api admin,eth,net,web3 \
            --authrpc.addr 127.0.0.1 --authrpc.port "$(el_auth_port "$i")" --authrpc.jwtsecret "$JWT" \
            --port "$(el_p2p_port "$i")" --discovery.port "$(el_p2p_port "$i")" \
            --txpool.max-account-slots 1000 --ipcdisable >>"$LOG_DIR/el$i.log" 2>&1 &
        PIDS+=($!)
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        wait_for 60 "el$i JSON-RPC endpoint" el_ready "$i"
    done

    log "Peering execution clients"
    local enode0
    enode0="$(el_rpc 0 admin_nodeInfo | jq -r '.result.enode' | sed -E 's/@[^:@?]+:/@127.0.0.1:/')"
    for ((i = 1; i < NUM_VALS; i++)); do
        el_rpc "$i" admin_addPeer "[\"$enode0\"]" >/dev/null
    done
}

# Submits a steady stream of transfers to el0 for the rest of the run. The pending nonce is re-fetched per send
# so the stream self-heals after a rejected or dropped transaction. The execution clients stay up across the
# halt, so the load also exercises EL mempool carryover across the consensus binary swap.
start_tx_load() {
    log "Starting tx load (transfers to el0 every ${LOAD_INTERVAL}s)"
    local rpc="http://127.0.0.1:$(el_http_port 0)"
    (
        while :; do
            nonce="$(cast nonce --block pending "$LOAD_ADDR" --rpc-url "$rpc" 2>/dev/null)" || nonce=""
            if [[ -n "$nonce" ]]; then
                cast send "$LOAD_RECIPIENT" --value 1 --private-key "$LOAD_KEY" --nonce "$nonce" \
                    --gas-limit 21000 --gas-price 10gwei --async --rpc-url "$rpc" \
                    >/dev/null 2>>"$LOG_DIR/load.log" || true
            fi
            sleep "$LOAD_INTERVAL"
        done
    ) &
    PIDS+=($!)
    # Drop the job from bash's job table so cleanup's SIGTERM does not print a "Terminated: 15" job notice.
    disown
}

# ---------------------------------------------------------------------------------------------------------------
# Phases
# ---------------------------------------------------------------------------------------------------------------

phase1_halt_on_old() {
    # Halt time is computed here rather than at flag parsing so genesis and EL setup time does not eat the offset.
    if ((HALT_TIME_OFFSET > 0)); then
        HALT_FLAGS=(--halt-time $(($(date +%s) + HALT_TIME_OFFSET)))
    else
        HALT_FLAGS=(--halt-height "$HALT_HEIGHT")
    fi

    log "Phase 1: starting $NUM_VALS validators on OLD binary with ${HALT_FLAGS[*]}"
    local i
    CL_PIDS=()
    for ((i = 0; i < NUM_VALS; i++)); do
        CL_PIDS[i]="$(start_cl "$i" "$OLD_BIN" "cl$i.old.log" "${HALT_FLAGS[@]}")"
        PIDS+=("${CL_PIDS[i]}")
    done
    if ((LOAD)); then start_tx_load; fi

    log "Waiting for the first validator to self-halt (timeout ${HALT_TIMEOUT_SECS}s)"
    wait_for "$HALT_TIMEOUT_SECS" "a validator to exit at the halt point" any_halted
    stop_stragglers
    derive_and_check_halt_height
}

# Committing the halt block requires >2/3 precommits, but a validator exits ~immediately after its own commit,
# racing the gossip of the final precommits. So any subset of validators can wedge one block short of the halt
# height with their peers gone. That mirrors a real coordinated halt: the chain stops at the halt height, some
# nodes self-halt, and operators stop the wedged rest before swapping binaries. Success means at least one node
# self-halted and NO node committed past the halt height (asserted in derive_and_check_halt_height).
stop_stragglers() {
    # Validators only ever exit here, so instead of a fixed sleep, wait until the set of live processes has
    # been stable for a while: fast runs move on as soon as everyone halted, slow runners get more room.
    log "Waiting for the remaining validators to self-halt"
    local i prev_alive alive stable_since
    prev_alive="$(alive_count)"
    stable_since=$SECONDS
    while ((prev_alive > 0 && SECONDS - stable_since < 15)); do
        sleep 1
        alive="$(alive_count)"
        if ((alive != prev_alive)); then
            prev_alive=$alive
            stable_since=$SECONDS
        fi
    done

    STRAGGLERS=()
    for ((i = 0; i < NUM_VALS; i++)); do
        if kill -0 "${CL_PIDS[i]}" 2>/dev/null; then
            STRAGGLERS+=("cl$i")
            kill "${CL_PIDS[i]}" 2>/dev/null || true
        fi
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        wait_for 30 "cl$i process to exit" proc_gone "${CL_PIDS[i]}"
    done
}

derive_and_check_halt_height() {
    local i h

    # In halt-time mode the halt block is wherever the chain was when the halt time passed. Every self-halted
    # validator must have stopped at the same height, which becomes the halt height for the remaining phases.
    if ((HALT_TIME_OFFSET > 0)); then
        HALT_HEIGHT=""
        for ((i = 0; i < NUM_VALS; i++)); do
            grep -q "halting node per configuration" "$LOG_DIR/cl$i.old.log" || continue
            h="$(last_committed_height "$i")"
            [[ -n "$h" ]] || fail "could not read a committed height for self-halted cl$i"
            if [[ -n "$HALT_HEIGHT" && "$h" != "$HALT_HEIGHT" ]]; then
                fail "self-halted nodes disagree on the halt height ($h vs $HALT_HEIGHT)"
            fi
            HALT_HEIGHT="$h"
        done
        log "Halt time reached at height $HALT_HEIGHT"
    fi

    for ((i = 0; i < NUM_VALS; i++)); do
        h="$(last_committed_height "$i")"
        [[ -n "$h" ]] || fail "could not read a committed height for cl$i from its log," \
            "did the 'Committed state' log line change format?"
        ((h <= HALT_HEIGHT)) || fail "cl$i committed height $h past the halt height $HALT_HEIGHT"
    done

    if ((${#STRAGGLERS[@]} > 0)); then
        log "$(halted_count)/$NUM_VALS validators self-halted at height $HALT_HEIGHT," \
            "stopped stragglers: ${STRAGGLERS[*]} (they catch up after the swap)"
    else
        log "All $NUM_VALS validators halted after committing height $HALT_HEIGHT"
    fi
}

# A node that halted and is restarted with the halt flag still set must refuse to enter consensus instead of
# advancing or crash-looping.
phase2_restart_with_flag_refused() {
    local i node="" pid
    for ((i = 0; i < NUM_VALS; i++)); do
        if grep -q "halting node per configuration" "$LOG_DIR/cl$i.old.log"; then
            node=$i
            break
        fi
    done
    log "Phase 2: restarting cl$node with ${HALT_FLAGS[*]} still set (must refuse to start)"
    pid="$(start_cl "$node" "$OLD_BIN" "cl$node.haltcheck.log" "${HALT_FLAGS[@]}")"
    PIDS+=("$pid")
    wait_for 90 "cl$node to refuse startup at the halt point" proc_gone "$pid"
    grep -q "reached the configured halt point" "$LOG_DIR/cl$node.haltcheck.log" ||
        fail "cl$node restarted past the halt point instead of refusing, see $LOG_DIR/cl$node.haltcheck.log"
}

phase3_resume_on_new() {
    log "Phase 3: restarting all validators on NEW binary (same data dirs, no halt flags)"
    local i
    CL_PIDS=()
    for ((i = 0; i < NUM_VALS; i++)); do
        CL_PIDS[i]="$(start_cl "$i" "$NEW_BIN" "cl$i.new.log")"
        PIDS+=("${CL_PIDS[i]}")
    done

    TARGET_HEIGHT=$((HALT_HEIGHT + RESUME_BLOCKS))
    log "Waiting for all nodes to advance past height $TARGET_HEIGHT (timeout ${RESUME_TIMEOUT_SECS}s)"
    for ((i = 0; i < NUM_VALS; i++)); do
        wait_for "$RESUME_TIMEOUT_SECS" "cl$i to reach height $TARGET_HEIGHT" node_resumed "$i"
    done
    for ((i = 0; i < NUM_VALS; i++)); do
        if grep -q "CONSENSUS FAILURE" "$LOG_DIR/cl$i.new.log"; then
            fail "cl$i hit a consensus failure after the swap, see $LOG_DIR/cl$i.new.log"
        fi
    done

    EL0_HEIGHT="$(($(el_rpc 0 eth_blockNumber | jq -r '.result')))"
    if ((LOAD)); then
        PRE_SWAP_TXS="$(el_tx_count_range 1 "$HALT_HEIGHT")"
        POST_SWAP_TXS="$(el_tx_count_range $((HALT_HEIGHT + 1)) "$EL0_HEIGHT")"
        ((PRE_SWAP_TXS > 0)) || fail "no transactions were included before the halt, see $LOG_DIR/load.log"
        ((POST_SWAP_TXS > 0)) || fail "no transactions were included after the swap, see $LOG_DIR/load.log"
    fi
}

print_result() {
    log "RESULT"
    local i mined
    for ((i = 0; i < NUM_VALS; i++)); do
        echo "  cl$i: height=$(node_height "$i") p2p_version=$(node_p2p_version "$i")"
    done
    echo "  el0: block=$EL0_HEIGHT"
    if ((LOAD)); then
        mined="$(($(el_rpc 0 eth_getTransactionCount "[\"$LOAD_ADDR\",\"latest\"]" | jq -r '.result')))"
        echo "  load: $mined txs mined in total: $PRE_SWAP_TXS in blocks 1..$HALT_HEIGHT (old binary)," \
            "$POST_SWAP_TXS in blocks $((HALT_HEIGHT + 1))..$EL0_HEIGHT (new binary)"
    fi
    echo
    echo "PASS: the network halted at height $HALT_HEIGHT on the old binary (${#STRAGGLERS[@]} straggler(s))," \
        "resumed from the same data dirs on the new binary, and all $NUM_VALS validators advanced past" \
        "height $TARGET_HEIGHT."
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

    phase1_halt_on_old
    phase2_restart_with_flag_refused
    phase3_resume_on_new

    print_result
}

main
