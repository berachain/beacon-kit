// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package phuslu_test

import (
	"bytes"
	"testing"

	"github.com/berachain/beacon-kit/log/phuslu"
)

// commitSig is a small element; a large slice of these stands in for an oversized logged value.
type commitSig struct {
	BlockIDFlag int
	Address     []byte
}

// msgInfo nests a large slice so that logging the whole struct as a single field produces an oversized log line.
type msgInfo struct {
	Height     int
	Round      int
	Signatures []commitSig
	PeerID     string
}

func newLoggerForStyle(out *bytes.Buffer, style string) *phuslu.Logger {
	return phuslu.NewLogger(out, &phuslu.Config{TimeFormat: "RFC3339", LogLevel: "info", Style: style})
}

// TestLogLineIsBounded asserts that long log lines stay bounded and are marked as truncated in both output styles.
func TestLogLineIsBounded(t *testing.T) {
	t.Parallel()

	attack := msgInfo{Signatures: make([]commitSig, 250_000), PeerID: "1cc8ac47fe215c7da976e41cf04ce1e8b6c8f33c"}

	for _, style := range []string{phuslu.StyleJSON, phuslu.StylePretty} {
		t.Run(style, func(t *testing.T) {
			var out bytes.Buffer
			newLoggerForStyle(&out, style).Error("Peer sent us invalid msg", "msg", attack)

			b := out.Bytes()
			if len(b) > phuslu.MaxLogLineBytes {
				t.Fatalf("log line not bounded: %d bytes", len(b))
			}
			if !bytes.HasSuffix(b, []byte("\n")) {
				t.Fatal("truncated line must remain newline-terminated")
			}
			if !bytes.Contains(b, []byte("truncated, total size")) {
				t.Fatal("truncated line must carry a marker with the original total size")
			}
		})
	}
}
