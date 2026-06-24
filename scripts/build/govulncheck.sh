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

# govulncheck v1.4.0 bundles golang.org/x/tools v0.46.0, whose SSA callgraph
# builder panics ("ForEachElement called on type containing *types.TypeParam")
# on our generics under Go 1.26. Fixed upstream in golang/go#80055 (x/tools CL
# 786280), not yet in a tagged release. Build govulncheck against the patched
# commit (MVS selects it over v0.46.0). Drop this once a govulncheck release
# ships with the fix and replace with `go run .../govulncheck@latest`.
XTOOLS_FIX="golang.org/x/tools@d711ac7849d4f5456228745090323144c4c2d190"

# Pin the analysis toolchain to the version declared in go.mod.
export GOTOOLCHAIN="go$(sed -n 's/^go //p' go.mod)"

# Build govulncheck against the patched x/tools in a throwaway module.
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
(
  cd "$tmp"
  go mod init vulncheck-runner >/dev/null
  go get golang.org/x/vuln/cmd/govulncheck@latest >/dev/null
  go get "$XTOOLS_FIX" >/dev/null
  go build -o "$tmp/govulncheck" golang.org/x/vuln/cmd/govulncheck
)

# Scan all non-test packages.
# shellcheck disable=SC2046
"$tmp/govulncheck" $(go list ./... | grep -v '/testing/')
