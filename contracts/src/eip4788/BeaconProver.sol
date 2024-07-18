// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";

/// @author Berachain Team
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
contract BeaconProver is Verifier {
    error InvalidProposer();

    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /// @dev Generalized index of the first validator struct root in the
    /// registry.
    uint256 public immutable valGIndex;

    constructor(uint256 _valGIndex) {
        valGIndex = _valGIndex;
    }

    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex
    )
        external
        view
    {
        // First check that the validator index is that of the block proposer.
        if (validatorIndex != blockHeader.proposerIndex) {
            revert InvalidProposer();
        }

        // Then verify that the block header is the valid block header for this time.
        bytes32 expectedBeaconRoot = getBlockRoot(timestamp);
        bytes32 givenBeaconRoot = SSZ.beaconHeaderHashTreeRoot(blockHeader);
        if (expectedBeaconRoot != givenBeaconRoot) {
            revert RootNotFound();
        }

        // Finally check that the validator is a validator of the beacon chain during this time.
        proveValidator(
            expectedBeaconRoot, validatorProof, validator, validatorIndex
        );
    }

    function proveValidator(
        bytes32 beaconBlockRoot,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex
    )
        internal
        view
    {
        if (validatorIndex >= VALIDATOR_REGISTRY_LIMIT) {
            revert IndexOutOfRange();
        }

        uint256 gI = valGIndex + validatorIndex;
        bytes32 validatorRoot = SSZ.validatorHashTreeRoot(validator);

        if (
            !SSZ.verifyProof(validatorProof, beaconBlockRoot, validatorRoot, gI)
        ) {
            revert InvalidProof();
        }
    }
}
