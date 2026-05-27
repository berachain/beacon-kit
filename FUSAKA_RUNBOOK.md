# Fusaka Runbook

Live-network sanity checks for the Fusaka hard fork (BRIP-10).

Fusaka couples a CL fork (beacon-kit, internally named **Fulu**, fork version `0x06000000`)
with an EL fork (bera-reth, internally named **Osaka** to match upstream revm).
Both must activate at the same `block.timestamp` for the chain to stay live.

| Network | Chain ID | Fusaka activation (unix) | Fusaka activation (UTC) |
| --- | --- | --- | --- |
| Bepolia (testnet) | `80069` | `1779897600` | 2026-05-27 16:00:00 |
| Mainnet | `80094` | not yet scheduled (`9999999999999999`) | TBD |

Releases under test:

- CL: [beacon-kit v1.4.0-rc3](https://github.com/berachain/beacon-kit/releases/tag/v1.4.0-rc3)
- EL: [bera-reth v1.4.0-rc2](https://github.com/berachain/bera-reth/releases/tag/v1.4.0-rc2)
- Spec: [BRIP-10](https://github.com/berachain/BRIPs/blob/main/meta/BRIP-0010.md)

Endpoints assumed throughout (substitute the node's real host/port as needed):

| Service | Default URL |
| --- | --- |
| EL JSON-RPC (bera-reth) | `http://localhost:8545` |
| CL Beacon-API (beacon-kit `node-api`) | `http://localhost:3500` |
| CL CometBFT RPC | `http://localhost:26657` |

> The CometBFT RPC port is **`26657`** by default (CometBFT's standard). If a
> particular deployment maps it to `26557`, substitute accordingly. All comet
> commands below use `26657`.

A quick environment helper to paste once per shell:

```bash
export EL=http://localhost:8545
export CL=http://localhost:3500
export CMT=http://localhost:26657
export DEPOSIT_CONTRACT=0x4242424242424242424242424242424242424242
export FORK_TIME=1779897600   # Bepolia
```

---

## 0. Pre-flight: confirm the fork has actually activated

The whole runbook assumes Fusaka has activated, i.e. the head block's timestamp
is `>= FORK_TIME`.

```bash
# EL: head block timestamp
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["latest",false]}' \
  | jq -r '.result.timestamp' | xargs -I{} printf "%d\n" {}
```

```bash
# CL: head fork version (must be 0x06000000 after Fusaka)
curl -s $CL/eth/v1/beacon/states/head/fork | jq .
```

Expected once active:

```json
{
  "data": {
    "previous_version": "0x05010000",
    "current_version": "0x06000000",
    "epoch": "0"
  }
}
```

```bash
# CL: current spec name (string "fulu" once active)
curl -s $CL/eth/v1/beacon/headers/head | jq '.version'
```

If `current_version` is still `0x05010000` (electra1) after `FORK_TIME`, the
node has not transitioned — investigate before continuing.

---

## 1. EL sanity checks (`:8545`)

### 1.1 EIP-7910 `eth_config` returns the Fusaka/Osaka fork

`bera-reth` ships an `eth_config` RPC (EIP-7910) that exposes per-fork
precompiles, system contracts, blob schedule, chain id and fork id.

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_config","params":[]}' | jq .
```

Sanity checks on the response:

- `result.current.chainId == "0x138d7"` (`80087`, devnet) /
  `"0x138c5"` (`80069`, Bepolia) / `"0x138de"` (`80094`, mainnet) — confirms
  the node is on the right network.
- `result.current.activationTime` equals the configured Osaka time for the
  network (Bepolia: `1779897600`). After fork, the "current" bucket has rolled
  forward from Prague* to Osaka.
- `result.current.precompiles` contains an entry **`"P256VERIFY":
  "0x0000000000000000000000000000000000000100"`**. This is the smoking gun
  for EIP-7951.
- `result.current.systemContracts.DEPOSIT_CONTRACT_ADDRESS ==
  "0x4242424242424242424242424242424242424242"` (the fix called out in
  v1.4.0-rc2 release notes).
- `result.current.forkId` is a non-zero 4-byte hex string and matches what
  peers advertise in devp2p (see `admin_nodeInfo` if exposed).
- `result.next` is `null` once Osaka is the latest configured fork (nothing
  scheduled after it). `result.last` is also `null` in that case.

### 1.2 EIP-7951 — P-256 (secp256r1) precompile at `0x0…0100`

A valid-signature test vector (from revm's secp256r1 suite):

| Field | Bytes |
| --- | --- |
| `msg_hash` | `4cee90eb86eaa050036147a12d49004b6b9c72bd725d39d4785011fe190f0b4d` |
| `r` | `a73bd4903f0ce3b639bbbf6e8e80d16931ff4bcf5993d58468e8fb19086e8cac` |
| `s` | `36dbcd03009df8c59286b162af3bd7fcc0450c9aa81be5d10d312af6c66b1d60` |
| `pubkey_x` | `4aebd3099c618202fcfe16ae7770b0c49ab5eadf74b754204a3bb6060e44eff3` |
| `pubkey_y` | `7618b065f9832de4ca6ca971a7a1adc826d0f7c00181a5fb2ddf79ae00b4e10e` |

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[{
    "to":"0x0000000000000000000000000000000000000100",
    "data":"0x4cee90eb86eaa050036147a12d49004b6b9c72bd725d39d4785011fe190f0b4da73bd4903f0ce3b639bbbf6e8e80d16931ff4bcf5993d58468e8fb19086e8cac36dbcd03009df8c59286b162af3bd7fcc0450c9aa81be5d10d312af6c66b1d604aebd3099c618202fcfe16ae7770b0c49ab5eadf74b754204a3bb6060e44eff37618b065f9832de4ca6ca971a7a1adc826d0f7c00181a5fb2ddf79ae00b4e10e"
  },"latest"]}' | jq -r .result
```

Expected: `0x0000000000000000000000000000000000000000000000000000000000000001` (32-byte word, last byte `0x01`).

Invalid-signature sanity check — flip one bit of `msg_hash`:

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[{
    "to":"0x0000000000000000000000000000000000000100",
    "data":"0xb3ee90eb86eaa050036147a12d49004b6b9c72bd725d39d4785011fe190f0b4da73bd4903f0ce3b639bbbf6e8e80d16931ff4bcf5993d58468e8fb19086e8cac36dbcd03009df8c59286b162af3bd7fcc0450c9aa81be5d10d312af6c66b1d604aebd3099c618202fcfe16ae7770b0c49ab5eadf74b754204a3bb6060e44eff37618b065f9832de4ca6ca971a7a1adc826d0f7c00181a5fb2ddf79ae00b4e10e"
  },"latest"]}' | jq -r .result
```

Expected: `0x` (empty output — signature does not verify).

Gas check — the precompile must charge `6900` gas (Osaka variant), not the
pre-Osaka `3450`:

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_estimateGas","params":[{
    "to":"0x0000000000000000000000000000000000000100",
    "data":"0x4cee90eb86eaa050036147a12d49004b6b9c72bd725d39d4785011fe190f0b4da73bd4903f0ce3b639bbbf6e8e80d16931ff4bcf5993d58468e8fb19086e8cac36dbcd03009df8c59286b162af3bd7fcc0450c9aa81be5d10d312af6c66b1d604aebd3099c618202fcfe16ae7770b0c49ab5eadf74b754204a3bb6060e44eff37618b065f9832de4ca6ca971a7a1adc826d0f7c00181a5fb2ddf79ae00b4e10e"
  }]}' | jq -r .result
```

Expected gas: roughly `0x76dc` (~30460) = 21000 base + ~2560 calldata + 6900
precompile.

Pre-fork (negative) verification: the same call against a block whose
timestamp is strictly **before** `FORK_TIME` must return `0x` and use ~3450
gas. Use `eth_call` with a `blockNumber` corresponding to the last pre-fork
block.

### 1.3 EIP-7939 — `CLZ` (0x1e) opcode

Test bytecode: load 32 bytes of calldata, run `CLZ`, return the result.

```
PUSH1 0x00 CALLDATALOAD CLZ PUSH1 0x00 MSTORE PUSH1 0x20 PUSH1 0x00 RETURN
0x60 00     35           1e  60 00      52     60 20      60 00       f3
```

Hex: `0x6000351e60005260206000f3`

Use an `eth_call` with a state override (`stateDiff`/`stateOverride` —
bera-reth supports the geth-style override map) to plant this code at a dead
address and call it. Input `0x00…01` should return `255` leading zeros.

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[
    {
      "to":"0x000000000000000000000000000000000000C12E",
      "data":"0x0000000000000000000000000000000000000000000000000000000000000001"
    },
    "latest",
    {
      "0x000000000000000000000000000000000000C12E":{
        "code":"0x6000351e60005260206000f3"
      }
    }
  ]}' | jq -r .result
```

Expected: `0x00000000000000000000000000000000000000000000000000000000000000ff` (`= 255`).

Additional inputs to verify:

| Calldata (32-byte BE) | Expected CLZ |
| --- | --- |
| `0x00…00` | `256` (`0x100`) |
| `0x00…01` | `255` (`0xff`) |
| `0x00…ff` | `248` (`0xf8`) |
| `0xff…ff` | `0` |
| `1 << 128` (`0x00…00 80 00…00`) | `127` (`0x7f`) |

Pre-fork (negative) check: the same call against a pre-fork block must
**revert** — `0x1e` is an invalid opcode pre-Osaka.

### 1.4 EIP-7883 — MODEXP minimum gas raised to 500

Call MODEXP at `0x05` with `base_len=1, exp_len=1, mod_len=1, base=0x02, exp=0x03, mod=0x05` (computes `2^3 mod 5 = 3`):

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[{
    "to":"0x0000000000000000000000000000000000000005",
    "data":"0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001020305"
  },"latest"]}' | jq -r .result
```

Expected result: `0x03`.

Gas sanity check:

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_estimateGas","params":[{
    "to":"0x0000000000000000000000000000000000000005",
    "data":"0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001020305"
  }]}' | jq -r .result
```

Expected gas: `>= 21500` (21000 tx base + 500 MODEXP floor). Pre-Osaka the same
call costs `21200` (200 MODEXP floor). The 500-floor is the headline signal.

### 1.5 EIP-7823 — MODEXP input bounds (1024 bytes per field)

Construct an input declaring `base_len = 1025` (one byte past the cap). The
call must consume all gas and fail.

```bash
# 96-byte header: base_len=1025, exp_len=0, mod_len=0, then 1025 zero bytes.
BASE_LEN_1025=$(printf '%064x' 1025)
ZERO=$(printf '%064x' 0)
BASE_PAYLOAD=$(python3 -c 'print("00"*1025)')
DATA="0x${BASE_LEN_1025}${ZERO}${ZERO}${BASE_PAYLOAD}"

curl -s -X POST $EL -H 'content-type: application/json' --data "{
  \"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_call\",\"params\":[{
    \"to\":\"0x0000000000000000000000000000000000000005\",
    \"data\":\"$DATA\",
    \"gas\":\"0xf4240\"
  },\"latest\"]}" | jq .
```

Expected: an `error` (or empty result with `"error": ...`) and that
`eth_estimateGas` against the same input fails because the call OOGs. The
exact error string varies; the requirement is simply that the call does **not**
succeed under Osaka.

### 1.6 BRIP-10 — `MAX_CODE_SIZE = 32 KB`, `MAX_INITCODE_SIZE = 64 KB`

The simplest live check is to deploy a contract whose runtime code is between
`24576` and `32768` bytes. Pre-Osaka this fails; post-Osaka it succeeds.

```bash
# 28 KB of STOPs as runtime code, wrapped in a constructor that returns it.
# Init prefix is exactly 15 bytes / 30 hex chars (see osaka_eip_tests
# create_init_code_returning_stops in bera-reth src/evm/mod.rs):
#   PUSH2 len | PUSH2 0x000f | PUSH1 0 | CODECOPY | PUSH2 len | PUSH1 0 | RETURN
#   61 <hi><lo>  61 00 0f      60 00     39         61 <hi><lo>  60 00     f3
LEN_HEX=$(printf '%04x' 28000)                       # e.g. "6d60"
INIT_PREFIX="61${LEN_HEX}61000f60003961${LEN_HEX}6000f3"
PAYLOAD=$(python3 -c 'print("00"*28000, end="")')    # 56000 hex chars, no newline
INITCODE="0x${INIT_PREFIX}${PAYLOAD}"

# Sanity-check the prefix length before sending: must be 30 hex chars.
[ ${#INIT_PREFIX} -eq 30 ] || { echo "bad init prefix: ${#INIT_PREFIX} chars"; exit 1; }

# Send via a funded test account:
cast send --rpc-url $EL --private-key $PRIV_KEY --create $INITCODE
```

Expected: tx mined, receipt has a non-zero `contractAddress`, and
`eth_getCode` against it returns 56,000 hex chars (28 KB).

A pre-fork repeat of the same tx (use a fork archive node or trace at an
earlier block via `eth_call` with `from`) must fail with "max code size
exceeded".

### 1.7 EIP-7934 — block RLP cap 10 MiB (8 MiB enforced)

BRIP-10 cites the spec value of 10 MiB. In practice bera-reth enforces
`MAX_RLP_BLOCK_SIZE = 8 MiB = 8_388_608 bytes` (upstream reth's safety margin
under the spec cap) — see `src/engine/builder.rs` for the payload builder
checks. This one is hard to provoke without a payload-builder under load; the
practical sanity check is "head blocks are still produced".

```bash
# RLP size of a recent block in bytes
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["latest",false]}' \
  | jq -r '.result.size' | xargs -I{} printf "%d\n" {}
```

Should always print a value `<= 8388608`. Exceeding the spec 10 MiB cap would
cause peer-level rejection — confirm peers stay connected (`net_peerCount`
stable, `admin_peers` healthy).

### 1.8 EIP-7825 — transaction gas cap 2^24 = 16,777,216

Try to submit (or `eth_estimateGas`) a tx with `gas > 16_777_216`. It must be
rejected.

```bash
# Should fail post-Osaka:
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_estimateGas","params":[{
    "to":"0x0000000000000000000000000000000000000000",
    "gas":"0x1000001"
  }]}' | jq .
```

Expected: an `error` mentioning gas limit cap. A request at exactly
`0x1000000 = 16777216` must be accepted.

The same field is also visible in `eth_config` under the current fork
(`max_tx_gas_limit` / equivalent — depending on alloy-eip7910 schema version).

### 1.9 EIP-6110 — deposit event on the deposit contract

The Berachain deposit contract uses a **custom 5-field event** (not the
mainnet 1-field packed encoding):

```solidity
event Deposit(
    bytes pubkey,        // 48 bytes BLS pubkey
    bytes credentials,   // 32 bytes withdrawal creds
    uint64 amount,       // gwei
    bytes signature,     // 96 bytes BLS sig
    uint64 index
);
```

- Contract: `0x4242424242424242424242424242424242424242`
- `topic0 = keccak256("Deposit(bytes,bytes,uint64,bytes,uint64)")`
  = **`0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46`**

Watch for events:

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getLogs","params":[{
    "fromBlock":"latest","toBlock":"latest",
    "address":"0x4242424242424242424242424242424242424242",
    "topics":["0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46"]
  }]}' | jq .
```

Read the on-chain deposit counter (slot 0 of the deposit contract; lower 8
bytes are the `uint64` count):

```bash
curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getStorageAt","params":[
    "0x4242424242424242424242424242424242424242",
    "0x0","latest"]}' | jq -r .result
```

Run this twice with a deposit submitted in between — the trailing `uint64`
must increment by exactly the number of `Deposit` events emitted.

### 1.10 EVM inflation withdrawal (Fulu)

Each block under Fulu credits a fixed amount of native BERA to the
`EVMInflationAddressFulu` via a system withdrawal.

| Network | Address | Per block (gwei) |
| --- | --- | --- |
| Mainnet | `0x1AE7dD7AE06F6C58B4524d9c1f816094B1bcCD8e` | `1_705_000_000` (`1.705 × 10⁹`) |
| Bepolia | (see `testing/networks/80069/spec.toml`) | (see same file) |

```bash
# Spot-check: the balance must increase by exactly INFLATION_GWEI * N_BLOCKS
# over a window of N blocks.
ADDR=0x1AE7dD7AE06F6C58B4524d9c1f816094B1bcCD8e
b0=$(curl -s -X POST $EL -H 'content-type: application/json' --data "{
  \"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_getBalance\",\"params\":[\"$ADDR\",\"latest\"]}" | jq -r .result)
sleep 12
b1=$(curl -s -X POST $EL -H 'content-type: application/json' --data "{
  \"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_getBalance\",\"params\":[\"$ADDR\",\"latest\"]}" | jq -r .result)
python3 -c "print((int('$b1',16) - int('$b0',16))/1e9, 'gwei across the window')"
```

Expected: roughly `1.705 * N_BLOCKS` gwei (mainnet value; substitute Bepolia
value as configured). If this is zero, the inflation withdrawal is not being
posted — block production may not be on Fulu yet.

---

## 2. CL Beacon-API sanity checks (`:3500`)

### 2.1 Genesis + spec + fork

```bash
curl -s $CL/eth/v1/beacon/genesis | jq .
curl -s $CL/eth/v1/config/spec | jq .
curl -s $CL/eth/v1/beacon/states/head/fork | jq .
```

Sanity:

- `genesis.data.genesis_fork_version` is the chain's genesis version
  (`0x04000000` for Bepolia/mainnet).
- `states/head/fork.data.current_version == "0x06000000"` after Fusaka.
- `config/spec` returns the live spec values:
  - `HYSTERESIS_QUOTIENT = 100`
  - `HYSTERESIS_UPWARD_MULTIPLIER = 10`
  - `HYSTERESIS_DOWNWARD_MULTIPLIER = 1`
  - `DEPOSIT_CONTRACT_ADDRESS = 0x4242…4242`
  - `FULU_FORK_VERSION = 0x06000000`
  - `FULU_FORK_EPOCH` / `FULU_FORK_TIME` matches the schedule

If `HYSTERESIS_QUOTIENT` is still `4`, BRIP-8 is not active — confirm the fork
has actually transitioned.

### 2.2 Head + sync status

```bash
curl -s $CL/eth/v1/node/syncing | jq .
curl -s $CL/eth/v1/node/version | jq .
curl -s $CL/eth/v1/beacon/headers/head | jq .
```

Expected: `is_syncing: false`, `version` reports `v1.4.0-rc3` or later, and
`headers/head` returns a block with `slot` consistent with comet height.

### 2.3 Validator set + balances

```bash
# Spot-check a known validator (or query the whole active set)
curl -s "$CL/eth/v1/beacon/states/head/validators?status=active_ongoing" | jq '.data | length'

# A single validator by index or pubkey:
curl -s "$CL/eth/v1/beacon/states/head/validators/0" | jq .
```

After Fusaka, validator effective balances should drift only when actual
balance moves outside the new **[−100 BERA, +1000 BERA]** hysteresis window
(quotient=100, downward=1, upward=10 with `MAX_EFFECTIVE_BALANCE = 10000 BERA`).
Observable behaviour: small balance jitter no longer churns `effective_balance`.

### 2.4 EIP-6110 deposit processing path

`beacon-kit` v1.4.0-rc3 changes how validator deposits flow into the CL. Pre
Fulu, the CL polls the EL deposit contract at `eth1_follow_distance`. Starting
on the first Fulu block, the CL:

1. Drains the pre-Fulu deposit queue once (`CatchupFuluDeposits` in
   `beacon/deposits/deposits.go`).
2. From the second Fulu block onward, deposits are sourced directly from
   `body.execution_requests.deposits` (EIP-7685 request-type `0x00`), produced
   by bera-reth's receipt scanner (`bera-reth/src/deposits.rs`).

Quickest live signal: submit a deposit transaction to
`0x4242…4242`, observe the `Deposit(bytes,bytes,uint64,bytes,uint64)` log, and
confirm the validator count increases on the next CL block — **without** the
node logging "depositFetcher" / "Found deposits on execution layer" lines (those
are the pre-Fulu path).

End-to-end deposit test:

```bash
# 1. Read on-chain deposit counter before:
BEFORE_HEX=$(curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getStorageAt","params":[
    "0x4242424242424242424242424242424242424242","0x0","latest"]}' | jq -r .result)

# 2. Read CL validator count before:
CL_BEFORE=$(curl -s "$CL/eth/v1/beacon/states/head/validators?status=active_ongoing" | jq '.data | length')

# 3. Submit a deposit (use the chain's standard `deposit create-validator`
#    tooling — see `beacond deposit create-validator --help`):
beacond deposit create-validator \
  --beacon-kit.chain-spec testnet \
  ...

# 4. Wait for inclusion, then re-read:
AFTER_HEX=$(curl -s -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_getStorageAt","params":[
    "0x4242424242424242424242424242424242424242","0x0","latest"]}' | jq -r .result)

CL_AFTER=$(curl -s "$CL/eth/v1/beacon/states/head/validators?status=active_ongoing" | jq '.data | length')

echo "EL counter: $BEFORE_HEX -> $AFTER_HEX"
echo "CL active validators: $CL_BEFORE -> $CL_AFTER"
```

Expected: `EL counter` increments by exactly 1, and the corresponding validator
appears in the CL state within a few slots. Tail the CL logs for
`"Processed deposit to set Eth 1 deposit index"` to confirm the in-protocol
processing path was used (not the pre-Fulu fetcher).

### 2.5 Block headers + randao + withdrawals

```bash
# Latest block header
curl -s $CL/eth/v1/beacon/headers/head | jq .

# Randao for current state
curl -s $CL/eth/v1/beacon/states/head/randao | jq .

# Pending partial withdrawals (EIP-7002 path, active since Electra)
curl -s $CL/eth/v1/beacon/states/head/pending_partial_withdrawals | jq .
```

Note: `/eth/v2/beacon/blocks/{block_id}` is **not implemented** in v1.4.0-rc3.
To inspect a block body (and the new `execution_requests.deposits` array), use
the CometBFT routes below.

---

## 3. CometBFT sanity checks (`:26657`)

The CometBFT JSON-RPC is unauthenticated. Use it to cross-check height,
app_hash and chain id between CL and EL.

```bash
curl -s $CMT/status | jq '.result.sync_info'
curl -s $CMT/abci_info | jq .
curl -s $CMT/health | jq .
curl -s $CMT/net_info | jq '.result.n_peers'
curl -s "$CMT/block?height=$(curl -s $CMT/status | jq -r '.result.sync_info.latest_block_height')" | jq '.result.block.header'
```

Sanity:

- `status.node_info.network` reports the comet chain-id string (e.g.
  `bepolia-beacon-80069`).
- `status.sync_info.catching_up == false` once synced.
- `status.sync_info.latest_block_height` advances every ~2 s
  (`BERACHAIN_BLOCK_TIME_SECONDS = 2`).
- `abci_info.response.last_block_app_hash == sync_info.latest_app_hash`.
- `status.sync_info.latest_block_time` is monotonically increasing and crosses
  the `FORK_TIME` boundary cleanly (no consecutive heights with identical app
  hashes around the fork).

beacon-kit also exposes a pass-through under the node-API itself:

```bash
curl -s $CL/cometbft/v1/block/HEIGHT | jq .
curl -s $CL/cometbft/v1/signed_header/HEIGHT | jq .
```

These return the SSZ-decoded beacon block at a given comet height — useful
for inspecting `execution_requests.deposits` and the execution payload's
withdrawal list (look for the `EVMInflationAddressFulu` entry).

---

## 4. Cross-layer consistency checks

These tie the CL and EL together. Run after sections 1–3 each look healthy.

1. **Heights match.** comet `latest_block_height` == beacon-API `slot` ==
   `eth_blockNumber` on the EL. They are 1:1 in beacon-kit.

   ```bash
   echo "comet: $(curl -s $CMT/status | jq -r '.result.sync_info.latest_block_height')"
   echo "beacon: $(curl -s $CL/eth/v1/beacon/headers/head | jq -r '.data.header.message.slot')"
   echo "eth: $(curl -s -X POST $EL -H 'content-type: application/json' --data '{
     "jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}' | jq -r .result | xargs -I{} printf '%d\n' {})"
   ```

2. **Fork id agreement.** EL `eth_config.current.fork_id.hash` should match
   the post-Fulu `ForkID` peers advertise (devp2p Hello). If they disagree,
   peers will drop the connection.

3. **App hash == beacon state root.** The hash returned by `/abci_info`'s
   `last_block_app_hash` is the SSZ hash-tree-root of the beacon state at that
   height; it must equal the `state_root` field of the corresponding
   `/eth/v1/beacon/headers/HEIGHT`.

4. **No fork-boundary stall.** Inspect the block timestamps either side of
   `FORK_TIME` — `block(N).timestamp < FORK_TIME <= block(N+1).timestamp` —
   and confirm `N` and `N+1` both finalised within the usual 2 s window. A
   gap larger than ~10 s here indicates a transition bug.

---

## 5. Negative / regression checks

These should all **fail** if the node is correctly running Fusaka:

| Action | Expected response |
| --- | --- |
| `eth_call` to `0x…0100` with no/empty input | empty `0x` (or revert) |
| `eth_call` to a contract using `0x1e` (CLZ) pinned to a pre-fork block | reverts (invalid opcode) |
| `eth_sendRawTransaction` with `gas > 0x1000000` | rejected |
| `eth_call` to MODEXP with `base_len > 1024` | fails (consumes all gas) |
| Deploy of a 33 KB contract (> `MAX_CODE_SIZE_OSAKA`) | reverts ("max code size exceeded") |
| Submit a deposit with `amount < 1 BERA` (1e9 gwei) | deposit contract reverts with selector `0x0e1eddda` |
| EL responds to `engine_getPayloadV5` for a pre-Osaka payload | returns `UnsupportedFork` error, **not** a panic (v1.4.0-rc2 fix) |

For the last item, the only way to trigger it is by replaying an old engine
call; in practice, watch the EL logs for `engine_getPayloadV5` calls and
confirm no `panic` lines around the fork boundary.

---

## 6. Where to look when something is wrong

- **EL won't start / fork id mismatch:** the Bepolia/Mainnet genesis JSON
  must set `"osakaTime": 1779897600` (Bepolia) — bera-reth reads this via the
  upstream Ethereum `osakaTime` genesis field, not a Bera-specific override
  (see `src/chainspec/mod.rs` and `tests/fixtures/bepolia-genesis.json`).
- **CL stays on `0x05010000` past FORK_TIME:** the spec.toml on the node is
  out of date — confirm `fulu-fork-time = 1779897600` in
  `testing/networks/80069/spec.toml` (or whatever was distributed with the
  release).
- **Deposits stop being processed at the fork:** look for
  `"Found deposits to catchup for Fulu"` once in the CL logs at the first
  Fulu block. If absent, the `CatchupFuluDeposits` path didn't run; deposits
  posted in the last pre-Fulu block may be stuck.
- **`engine_*` failures:** v1.4.0-rc2 fixed two regressions
  (`engine_getPayloadV5` panic, `V4P11` response envelope). If you see panics
  or "envelope" decode errors in the EL log, you're on the older RC.
- **`eth_config` deposit contract is wrong:** v1.4.0-rc2 also fixed the
  deposit contract address it reports — confirm you're on rc2+.

---

## 7. Quick "is Fusaka live and well?" smoke test (copy/paste)

```bash
set -e
: "${EL:=http://localhost:8545}"
: "${CL:=http://localhost:3500}"
: "${CMT:=http://localhost:26657}"

echo "== Fork version (CL) =="
curl -fsS $CL/eth/v1/beacon/states/head/fork | jq '.data.current_version'

echo "== Hysteresis (CL spec) =="
curl -fsS $CL/eth/v1/config/spec | jq '{HYSTERESIS_QUOTIENT,HYSTERESIS_UPWARD_MULTIPLIER,HYSTERESIS_DOWNWARD_MULTIPLIER}'

echo "== P256VERIFY at 0x100 (EL) =="
curl -fsS -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[{
    "to":"0x0000000000000000000000000000000000000100",
    "data":"0x4cee90eb86eaa050036147a12d49004b6b9c72bd725d39d4785011fe190f0b4da73bd4903f0ce3b639bbbf6e8e80d16931ff4bcf5993d58468e8fb19086e8cac36dbcd03009df8c59286b162af3bd7fcc0450c9aa81be5d10d312af6c66b1d604aebd3099c618202fcfe16ae7770b0c49ab5eadf74b754204a3bb6060e44eff37618b065f9832de4ca6ca971a7a1adc826d0f7c00181a5fb2ddf79ae00b4e10e"
  },"latest"]}' | jq -r .result

echo "== CLZ via state override (EL) =="
curl -fsS -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_call","params":[
    {"to":"0x000000000000000000000000000000000000C12E",
     "data":"0x0000000000000000000000000000000000000000000000000000000000000001"},
    "latest",
    {"0x000000000000000000000000000000000000C12E":{"code":"0x6000351e60005260206000f3"}}
  ]}' | jq -r .result

echo "== eth_config current fork (EL) =="
curl -fsS -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_config","params":[]}' \
  | jq '.result.current | {chainId, activationTime, forkId, p256_at_0x100: (.precompiles.P256VERIFY == "0x0000000000000000000000000000000000000100")}'

echo "== Heights agree =="
echo "comet : $(curl -fsS $CMT/status | jq -r '.result.sync_info.latest_block_height')"
echo "beacon: $(curl -fsS $CL/eth/v1/beacon/headers/head | jq -r '.data.header.message.slot')"
echo "eth   : $(curl -fsS -X POST $EL -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}' | jq -r .result | xargs -I{} printf '%d\n' {})"
```

Pass criteria: `current_version == "0x06000000"`, P256VERIFY returns
`0x…01`, CLZ returns `0x…ff`, `precompiles_has_p256 == true`,
all three height counters within one of each other.
