// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "../libraries/SSZ.sol";

interface IBeaconVerifier {
    error InvalidProposer();

    /// @dev Generalized index of the first validator struct root in the
    /// registry.
    function valGIndex() external view returns (uint256);

    /// @notice Get the parent block root from a block's timestamp.
    /// @param timestamp `uint64` timestamp of the block.
    function getParentBlockRootAt(uint64 timestamp)
        external
        view
        returns (bytes32);

    /// @notice Get the parent block root at `block.timestamp`.
    function getParentBlockRoot() external view returns (bytes32);

    /// @notice Verifies the proposer of a beacon block.
    /// @param blockHeader `BeaconBlockHeader` to verify.
    /// @param timestamp `uint64` timestamp of the block.
    /// @param validatorProof `bytes32[]` proof of the validator.
    /// @param validator `Validator` to verify.
    /// @param validatorIndex `uint64` index of the validator. Must be the same
    /// as the `proposerIndex` in the block header.
    ///
    /// TODO: Change the validator to just take in bytes pubKey and the
    /// validatorProof correspondingly.
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex
    )
        external
        view;
}
