// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package suite

import (
	"path/filepath"
	"runtime"

	"github.com/berachain/beacon-kit/mod/errors"
)

// GetRelativePathToKurtosis returns the relative path to the kurtosis folder
// at the project root.
func GetRelativePathToKurtosis() (string, error) {
	// Get the current file path
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get current file path")
	}

	// Get the directory of the current file
	currentDir := filepath.Dir(currentFile)

	// Define the target directory (kurtosis folder at project root)
	targetDir := filepath.Join(
		filepath.Dir(filepath.Dir(filepath.Dir(currentDir))), "kurtosis",
	)

	// Calculate the relative path
	relPath, err := filepath.Rel(currentDir, targetDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to calculate relative path")
	}

	return relPath, nil
}
