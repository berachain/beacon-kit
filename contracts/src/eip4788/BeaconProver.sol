// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Ownable } from "@solady/src/auth/Ownable.sol";

import { SSZ } from "./libraries/SSZ.sol";

import { Verifier } from "./Verifier.sol";

import { IBeaconProver } from "./interfaces/IBeaconProver.sol";

/// @author Berachain Team
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
contract BeaconProver is Verifier, Ownable, IBeaconProver {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconProver
    uint256 public valGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(uint256 _valGIndex) {
        valGIndex = _valGIndex;

        _initializeOwner(msg.sender);
    }

    function _guardInitializeOwner()
        internal
        pure
        override
        returns (bool guard)
    {
        return true;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     BEACON ROOT VIEWS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconProver
    function getBeaconBlockRootAt(uint64 timestamp)
        external
        view
        returns (bytes32)
    {
        return getBeaconBlockRoot(timestamp);
    }

    /// @inheritdoc IBeaconProver
    function getCurrentBeaconBlockRoot() external view returns (bytes32) {
        return getBeaconBlockRoot(uint64(block.timestamp));
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           PROOFS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconProver
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex
    )
        public
        view
    {
        // First check that the validator index is that of the block proposer.
        if (validatorIndex != blockHeader.proposerIndex) {
            revert InvalidProposer();
        }

        // Then verify that the block header is the valid block header for this time.
        bytes32 expectedBeaconRoot = getBeaconBlockRoot(timestamp);
        bytes32 givenBeaconRoot = SSZ.beaconHeaderHashTreeRoot(blockHeader);
        if (expectedBeaconRoot != givenBeaconRoot) {
            revert RootNotFound();
        }

        // Finally check that the validator is a validator of the beacon chain during this time.
        proveValidatorInBlock(
            expectedBeaconRoot, validatorProof, validator, validatorIndex
        );
    }

    /// @notice Verifies the validator is in the registry of beacon state.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param validatorProof `bytes32[]` proof of the validator.
    /// @param validator `Validator` to verify.
    /// @param validatorIndex `uint64` index of the validator.
    function proveValidatorInBlock(
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
