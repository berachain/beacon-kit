// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";

/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
contract WithdrawalsVerifier is Verifier {
    uint64 internal constant MAX_WITHDRAWALS = 1 << 4;

    // Generalized index of the first withdrawal struct root in the withdrawals.
    uint256 public immutable gIndex;

    /// @notice Emitted when a withdrawal is submitted
    event WithdrawalSubmitted(uint64 indexed validatorIndex, uint64 amount);

    constructor(uint256 _gIndex) {
        gIndex = _gIndex;
    }

    function submitWithdrawal(
        bytes32[] calldata withdrawalProof,
        SSZ.Withdrawal memory withdrawal,
        uint8 withdrawalIndex,
        uint64 ts
    )
        public
    {
        if (withdrawalIndex >= MAX_WITHDRAWALS) {
            revert IndexOutOfRange();
        }

        uint256 gI = gIndex + withdrawalIndex;
        bytes32 withdrawalRoot = SSZ.withdrawalHashTreeRoot(withdrawal);
        bytes32 blockRoot = getParentBlockRoot(ts);

        if (!SSZ.verifyProof(withdrawalProof, blockRoot, withdrawalRoot, gI)) {
            revert InvalidProof();
        }

        emit WithdrawalSubmitted(withdrawal.validatorIndex, withdrawal.amount);
    }
}
