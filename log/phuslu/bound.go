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

package phuslu

import (
	"fmt"
	"io"
)

// MaxLogLineBytes caps the size of one serialized log line, marker included, so a single oversized value cannot expand it into megabytes.
// phuslu writes one entry per Write, so bounding each Write bounds each line.
const MaxLogLineBytes = 64 * 1024

// truncatingWriter caps every write at MaxLogLineBytes, reserving room to append a marker that names the original total size, and keeps
// the log stream newline-terminated.
type truncatingWriter struct {
	out io.Writer
}

func (w truncatingWriter) Write(p []byte) (int, error) {
	if len(p) <= MaxLogLineBytes {
		return w.out.Write(p)
	}
	// Reserve room for the marker so the emitted line, marker included, never exceeds MaxLogLineBytes.
	marker := fmt.Sprintf(" ...[truncated, total size %d bytes]\n", len(p))
	keep := MaxLogLineBytes - len(marker)
	if keep < 0 {
		keep = 0
	}
	line := make([]byte, 0, keep+len(marker))
	line = append(line, p[:keep]...)
	line = append(line, marker...)
	n, err := w.out.Write(line)
	if err != nil {
		return 0, err
	}
	if n < len(line) {
		// A short write with a nil error violates the io.Writer contract; surface it rather than lose the line silently.
		return 0, io.ErrShortWrite
	}
	// Report the full length as consumed so the logger does not treat the deliberate truncation as a short write.
	return len(p), nil
}
