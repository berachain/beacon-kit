// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package keys

// Collection prefixes.
const (
	WithdrawalQueuePrefix                  = "withdrawal_queue"
	RandaoMixPrefix                        = "randao_mix"
	SlashingsPrefix                        = "slashings"
	TotalSlashingPrefix                    = "total_slashing"
	ValidatorIndexPrefix                   = "val_idx"
	BlockRootsPrefix                       = "block_roots"
	StateRootsPrefix                       = "state_roots"
	ValidatorByIndexPrefix                 = "val_idx_to_pk"
	ValidatorPubkeyToIndexPrefix           = "val_pk_to_idx"
	ValidatorConsAddrToIndexPrefix         = "val_cons_addr_to_idx"
	ValidatorEffectiveBalanceToIndexPrefix = "val_eff_bal_to_idx"
	LatestBeaconBlockHeaderPrefix          = "latest_beacon_block_header"
	SlotPrefix                             = "slot"
	BalancesPrefix                         = "balances"
	Eth1BlockHashPrefix                    = "eth1_block_hash"
	Eth1DataPrefix                         = "eth1_data"
	Eth1DepositIndexPrefix                 = "eth1_deposit_idx"
	LatestExecutionPayloadHeaderPrefix     = "latest_execution_payload_header"
	GenesisValidatorsRootPrefix            = "genesis_validators_root"
	NextWithdrawalIndexPrefix              = "next_withdrawal_index"
	NextWithdrawalValidatorIndexPrefix     = "next_withdrawal_val_idx"
	ForkPrefix                             = "fork"
)
