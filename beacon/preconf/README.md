# Preconfirmations (preconf)

## Overview

The preconf package lets a designated **sequencer** node build execution
payloads ahead of time and hand them to the validator that is about to propose,
instead of each validator building its own payload under the tight
block-production deadline.

In BeaconKit a proposer has only about a second inside `PrepareProposal` to
assemble a block, which includes a round trip to the execution client to build
the execution payload. The sequencer takes that round trip off the critical
path. It starts building optimistically as soon as it knows which whitelisted
validator proposes next, and serves the finished payload over HTTP. When that
proposer reaches `PrepareProposal`, it fetches the pre-built payload with a
single HTTP call rather than driving the build itself.

"Preconfirmation" here refers to this sequencer-builds-ahead mechanism. It is
not a user-facing inclusion or ordering promise.

### Flow

1. The sequencer sees that a whitelisted validator is the upcoming proposer and
   optimistically builds its payload (sends a forkchoice update with attributes
   to the execution client).
2. During `PrepareProposal`, that validator calls the sequencer's `/payload`
   endpoint for its `(slot, parent_block_root)`.
3. On success it proposes the fetched payload. On any failure (sequencer down,
   timeout, payload not built yet) it falls back to building locally.
4. While the sequencer is unreachable, the validator probes the health endpoint
   every `health-check-interval` and resumes fetching once it recovers.

## Roles

A node operates in one of two preconf roles, controlled by config:

- **Sequencer** (`sequencer-mode = true`): runs the preconf HTTP API server,
  builds payloads optimistically for whitelisted proposers, and serves them.
- **Validator / fetcher** (`sequencer-url` set): runs the preconf client, which
  fetches payloads from the sequencer during `PrepareProposal`. Falls back to
  local building if the fetch fails.

A single node is one or the other, never both.

## Configuration

All settings live under `[beacon-kit.preconf]`. Preconf is **off by default**.
`enabled = true` is required before anything else takes effect. A node is a
sequencer when `enabled && sequencer-mode`, and a fetcher when `enabled` is set
together with `sequencer-url`.

| Field | Side | Default | Purpose |
|-------|------|---------|---------|
| `enabled` | both | `false` | Master switch. Nothing runs unless this is true. |
| `sequencer-mode` | sequencer | `false` | Run as the sequencer (build and serve payloads). |
| `whitelist-path` | sequencer | none | JSON list of validator pubkeys allowed to fetch. Required in sequencer mode. |
| `validator-jwts-path` | sequencer | none | JSON map of validator pubkey to JWT secret. Required in sequencer mode. |
| `api-port` | sequencer | `9090` | Port the preconf HTTP API listens on. |
| `tls-cert-path` / `tls-key-path` | sequencer | none | TLS cert and key. Both or neither (see TLS). |
| `sequencer-url` | validator | none | URL of the sequencer's preconf API. Setting it makes this node a fetcher. |
| `sequencer-jwt-path` | validator | none | This validator's JWT secret for authenticating to the sequencer. Required when `sequencer-url` is set. |
| `sequencer-ca-cert-path` | validator | none | Optional CA cert that pins the sequencer (see TLS). |
| `fetch-timeout` | validator | `500ms` | Per-request timeout for fetching a payload. |
| `health-check-interval` | validator | `10s` | How often to probe the sequencer while it is unavailable. |

## Trust model

The sequencer is a shared block builder, so the `/payload` API is guarded on
three axes, all enforced server-side on every request:

- **Whitelist (who).** Only validator pubkeys in `whitelist-path` may fetch. An
  unknown or removed key is rejected with `403`. The whitelist hot-reloads on
  SIGHUP.
- **Proposer binding (what).** A validator may only fetch the payload for a slot
  it is the elected proposer for. This stops one validator from pulling another
  validator's payload and extracting its transaction ordering. A mismatch is
  rejected with `403`.
- **Authentication (proof).** Every request carries a per-validator JWT, see
  [Authentication](#authentication) below.

**Liveness is never at the sequencer's mercy.** If the sequencer is down, slow,
malicious, or returns nothing, the validator falls back to building the block
locally. The sequencer can speed up block production but cannot halt the chain
or prevent a validator from proposing. This is the key safety property of the
design.

Because JWTs are validated against an issued-at (`iat`) time window (currently
five minutes either side of now), the sequencer and validators should run NTP.
Large clock skew between them causes authentication failures.

## API

The server exposes two endpoints:

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `GET`  | `/eth/v1/preconf/health`  | none | Liveness check |
| `POST` | `/eth/v1/preconf/payload` | JWT  | Fetch the payload for a `(slot, parent_block_root)` |

### Authentication

Validators authenticate to the `/payload` endpoint with a JWT bearer token
(HMAC-SHA256, `iat`-based validity window). Each validator has its own shared
secret. The sequencer loads a `pubkey -> secret` map from `validator-jwts-path`.
After verifying the token, the server also checks that the caller is (a) on the
whitelist and (b) the expected proposer for the requested slot.

## TLS

By default the preconf API runs over **plaintext HTTP**. This is fine for
local/devnet, but **not** for any internet-facing deployment, since JWT tokens
and MEV-sensitive payloads would transit in the clear. TLS is opt-in and
configured entirely through file paths.

### How it works (high level)

Three files, all referenced by path in config:

- **`tls-key-path`** is the server private key. Secret, sequencer-only.
- **`tls-cert-path`** is the public server cert. It binds the key to the
  sequencer's hostname/IP via its SAN and carries a CA signature.
- **`sequencer-ca-cert-path`** is the validator-side trust anchor used to verify
  that signature.

On each connection the validator checks that the cert chains to a CA it trusts
and that the SAN matches the host it dialed, then the sequencer proves it holds
the matching key. JWTs and payloads flow only once that succeeds.

**Certificate pinning.** When `sequencer-ca-cert-path` is set, the validator
trusts **only** that CA, not the ~150 in the system store. This blocks CA
mis-issuance and BGP-hijack attacks, where an attacker gets a technically-valid
cert for the sequencer's hostname from another CA. Left empty, the validator
uses the system CA store, which is fine for publicly-trusted commercial certs.

### Config

Sequencer side (`[beacon-kit.preconf]`):

```toml
tls-cert-path = "/etc/beacond/preconf/server-cert.pem"
tls-key-path  = "/etc/beacond/preconf/server-key.pem"
```

Validator side (`[beacon-kit.preconf]`):

```toml
sequencer-url          = "https://sequencer.example.com:9090"
sequencer-ca-cert-path = "/etc/beacond/preconf/ca-cert.pem"  # optional, pins to this CA
```

Config rules enforced at startup (`Config.Validate`):

- `tls-cert-path` and `tls-key-path` must **both** be set or both be empty.
  No half-configured TLS.
- When TLS is configured, the server uses **only** HTTPS. It never falls back to
  plaintext on the same port (no dual-mode listener).
- `sequencer-ca-cert-path` requires `sequencer-url` to use the `https://` scheme.
- Cert and key files are checked for accessibility at startup. The node fails
  fast with a clear error if they're missing.

### Generating certificates

**Dev / devnet, self-signed** (the cert acts as its own CA, so pin it directly):

```bash
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
  -keyout server-key.pem -out server-cert.pem \
  -days 365 -nodes -subj "/CN=sequencer" \
  -addext "subjectAltName=DNS:sequencer,IP:1.2.3.4"
```

The SAN is required. Go's TLS client rejects certs whose SAN doesn't match the
dialed host. Set `tls-cert-path`/`tls-key-path` on the sequencer and point each
validator's `sequencer-ca-cert-path` at `server-cert.pem`.

**Production, internal CA** (stable anchor that lets you rotate the server cert
without re-touching validators):

```bash
# 1. Create the CA once (keep ca-key.pem offline, NOT on the sequencer)
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
  -keyout ca-key.pem -out ca-cert.pem -days 3650 -nodes \
  -subj "/CN=preconf-ca"

# 2. Create the server key + CSR
openssl req -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
  -keyout server-key.pem -out server.csr -nodes \
  -subj "/CN=sequencer"

# 3. Define the extensions the CA will set on the leaf. The CA controls these
#    explicitly rather than copying whatever the CSR requested.
cat > server-ext.cnf <<'EOF'
subjectAltName = DNS:sequencer.example.com,IP:1.2.3.4
basicConstraints = critical, CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
EOF

# 4. Sign the server cert with the CA, applying the extension file.
openssl x509 -req -in server.csr -CA ca-cert.pem -CAkey ca-key.pem \
  -CAcreateserial -out server-cert.pem -days 365 -extfile server-ext.cnf
```

Distribute `ca-cert.pem` to validators (as `sequencer-ca-cert-path`). It's the
anchor they pin once. Reissue `server-cert.pem` as needed, and validators keep
trusting it as long as the same CA signs it.

### Rotation

- **CA cert:** generate once, long-lived (e.g. 10 years). Regenerate only on CA
  key compromise, approaching expiry, or an algorithm change. Re-pinning every
  validator is the expensive operation, so keep this rare.
- **Server cert:** rotate routinely (every ~90 days, or annually if manual), and
  immediately on suspected key compromise or a host/IP/DNS change. The server
  cert's validity must not outlast the CA cert.
- **No restart needed:** replace the cert and key files in place, then send
  `SIGHUP` (see below). The new cert is served on the next handshake. Existing
  connections are unaffected. When validators pin the CA (`sequencer-ca-cert-path`),
  a server cert signed by the same CA needs no change on the validator side.

## Observability

The sequencer emits two Prometheus counters:

- `beacon_kit.preconf.server.payload_request_total`, labeled by `result`
  (`ok`, `unauthorized`, `not_whitelisted`, `wrong_proposer`,
  `payload_not_found`, `internal_error`, `bad_request`, and so on). A spike in
  any rejection label is the first signal of a misconfiguration or an attack.
- `beacon_kit.preconf.proposer_tracker.check_total`, labeled by `matched`
  (`true`/`false`), counting expected-proposer matches versus mismatches.

## Operational hot-reload (SIGHUP)

On `SIGHUP` (`kill -HUP <pid>`) the sequencer reloads, without a restart:

- the **validator whitelist** from `whitelist-path`, and
- the **TLS server certificate** from `tls-cert-path`/`tls-key-path` (when TLS is
  configured).

The cert swap uses the running listener's `GetCertificate` hook, so the new cert
takes effect on the next TLS handshake while in-flight connections continue on
the old one. If a reload fails (for example a malformed cert file), the previous
value is kept and the failure is logged, so a bad file never breaks the listener.

Validator JWT secrets are **not** hot-reloaded and still require a restart to
pick up changes.
