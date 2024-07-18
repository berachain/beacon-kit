// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";
import { ValidatorVerifier } from "./ValidatorVerifier.sol";

contract BeaconProver is Verifier {
    error InvalidProposer();

    ValidatorVerifier public immutable validatorVerifier;

    constructor(uint256 valGIndex) {
        validatorVerifier = new ValidatorVerifier(valGIndex);
    }

    // naive, unoptimized implementation
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex,
        uint64 timestamp
    ) 
        external
    {
        // First check that the validator index is that of the block proposer.
        if (validatorIndex != blockHeader.proposerIndex) {
            revert InvalidProposer();
        }

        // Then check that the validator is a validator of the beacon chain during this time.
        validatorVerifier.proveValidator(validatorProof, validator, validatorIndex, timestamp);

        // Finally verify that the block header is the valid block header for this slot & time.
        bytes32 expectedBeaconRoot = getParentBlockRoot(timestamp);
        bytes32 givenBeaconRoot = SSZ.beaconHeaderHashTreeRoot(blockHeader);
        if (expectedBeaconRoot != givenBeaconRoot) {
            revert RootNotFound();
        }
    }
}
