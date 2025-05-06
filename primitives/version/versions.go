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

package version

import (
	"github.com/berachain/beacon-kit/primitives/common"
)

// These are the versions of the Beacon Chain.
//
//nolint:gochecknoglobals // Kept as private to avoid modification of these variables at runtime.
var (
	// phase0 is the first version of the Beacon Chain.
	phase0 = common.Version{0x00, 0x00, 0x00, 0x00}
	// altair is the first hardfork of the Beacon Chain.
	altair = common.Version{0x01, 0x00, 0x00, 0x00}
	// bellatrix is the second hardfork of the Beacon Chain.
	bellatrix = common.Version{0x02, 0x00, 0x00, 0x00}
	// capella is the third hardfork of the Beacon Chain.
	capella = common.Version{0x03, 0x00, 0x00, 0x00}
	// deneb is the first version of the Deneb hardfork, used for genesis of Berachain mainnet.
	deneb = common.Version{0x04, 0x00, 0x00, 0x00}
	// deneb1 is the first hardfork of Deneb on Berachain mainnet.
	deneb1 = common.Version{0x04, 0x01, 0x00, 0x00}
	// electra is the first version of the Electra hardfork on Berachain mainnet.
	electra = common.Version{0x05, 0x00, 0x00, 0x00}
	// electra1 is the first hardfork of Electra on Berachain mainnet.
	// TBD if used but kept as an example.
	electra1 = common.Version{0x05, 0x01, 0x00, 0x00}
)

// Phase0 returns phase0 as a common.Version.
func Phase0() common.Version {
	return phase0
}

// Altair returns altair as a common.Version.
func Altair() common.Version {
	return altair
}

// Bellatrix returns bellatrix as a common.Version.
func Bellatrix() common.Version {
	return bellatrix
}

// Capella returns capella as a common.Version.
func Capella() common.Version {
	return capella
}

// Deneb returns deneb as a common.Version. Deneb is the genesis fork version for Berachain
// mainnet and Bepolia testnet.
func Deneb() common.Version {
	return deneb
}

// Deneb1 returns deneb1 as a common.Version.
func Deneb1() common.Version {
	return deneb1
}

// Electra returns electra as a common.Version.
func Electra() common.Version {
	return electra
}

// Electra1 returns electra1 as a common.Version.
func Electra1() common.Version {
	return electra1
}
