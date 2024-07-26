// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "../SSZ.sol";

interface IBeaconVerifier {
    /// @dev Generalized index of the first validator struct root in the
    /// registry.
    function zeroValidatorsGIndex() external view returns (uint256);

    /// @notice Get the parent beacon block root from a block's timestamp.
    /// @param timestamp `uint64` timestamp of the block.
    function getParentBeaconBlockRootAt(uint64 timestamp)
        external
        view
        returns (bytes32);

    /// @notice Get the parent beacon block root at `block.timestamp`.
    function getParentBeaconBlockRoot() external view returns (bytes32);

    /// @notice Verifies the proposer of a beacon block.
    /// @param blockHeader `BeaconBlockHeader` to verify.
    /// @param timestamp `uint64` timestamp of the block.
    /// @param validatorProof `bytes32[]` proof of the validator.
    /// @param validator `Validator` to verify.
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator
    )
        external
        view;
}
