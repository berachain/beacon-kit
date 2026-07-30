#!/usr/bin/env bash
# SPDX-License-Identifier: BUSL-1.1
#
# Copyright (C) 2025, Berachain Foundation. All rights reserved.
# Use of this software is governed by the Business Source License included
# in the LICENSE file of this repository and at www.mariadb.com/bsl11.
#
# ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
# TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
# VERSIONS OF THE LICENSED WORK.
#
# THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
# LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
# LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
#
# TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
# AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
# EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
# TITLE.

set -euo pipefail

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }

# Suppressed advisories, with the reason above each entry.
SUPPRESSED=(
  # x/crypto/openpgp is unmaintained with no fix, only reached via cosmos-sdk keyring armor.
  'GO-2026-5932'
)

# Scan with the toolchain declared in go.mod.
export GOTOOLCHAIN="go$(sed -n 's/^go //p' go.mod)"

report="$(mktemp)"
trap 'rm -f "$report"' EXIT

# Scan all non-test packages. JSON mode always exits 0; pass/fail is decided below after dropping suppressed advisories.
go run golang.org/x/vuln/cmd/govulncheck@latest -format json $(go list ./... | grep -v '/testing/') > "$report"

# Advisories our code reaches at symbol level. A finding is symbol level when its first trace frame has a function.
found="$(jq -r 'select(.finding.trace[0].function != null) | .finding.osv' "$report" | sort -u | grep -vxF -f <(printf '%s\n' "${SUPPRESSED[@]}") || true)"

if [ -n "$found" ]; then
  echo "Reachable vulnerabilities:"
  jq -r --arg found "$found" 'select(.osv.id | IN($found | split("\n")[]))
    | "  \(.osv.id): \(.osv.summary)\n    https://pkg.go.dev/vuln/\(.osv.id)"' "$report"
  echo "For full traces run: go run golang.org/x/vuln/cmd/govulncheck@latest ./..."
  exit 3
fi

echo "OK (suppressed advisories: ${SUPPRESSED[*]})"
