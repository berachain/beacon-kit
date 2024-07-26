// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Ownable } from "@solady/src/auth/Ownable.sol";

import { IBeaconVerifier } from "./interfaces/IBeaconVerifier.sol";

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";

/// @author Berachain Team
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
contract BeaconVerifier is Verifier, Ownable, IBeaconVerifier {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    uint256 public zeroValidatorPubkeyGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(uint256 _zeroValidatorPubkeyGIndex) {
        zeroValidatorPubkeyGIndex = _zeroValidatorPubkeyGIndex;

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

    function setZeroValidatorPubkeyGIndex(uint256 _zeroValidatorPubkeyGIndex)
        external
        onlyOwner
    {
        zeroValidatorPubkeyGIndex = _zeroValidatorPubkeyGIndex;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     BEACON ROOT VIEWS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    function getParentBeaconBlockRootAt(uint64 timestamp)
        external
        view
        returns (bytes32)
    {
        return getParentBlockRoot(timestamp);
    }

    /// @inheritdoc IBeaconVerifier
    function getParentBeaconBlockRoot() external view returns (bytes32) {
        return getParentBlockRoot(uint64(block.timestamp));
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           PROOFS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used ~81281
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorPubkeyProof,
        bytes calldata validatorPubkey
    )
        public
        view
    {
        // First verify that the block header is the valid block header for this time.
        bytes32 expectedBeaconRoot = getParentBlockRoot(timestamp);
        bytes32 givenBeaconRoot = SSZ.beaconHeaderHashTreeRoot(blockHeader);
        if (expectedBeaconRoot != givenBeaconRoot) {
            revert RootNotFound();
        }

        // Then check that the validator is a validator of the beacon chain during this time.
        proveValidatorPubkeyInBlock(
            expectedBeaconRoot,
            validatorPubkeyProof,
            validatorPubkey,
            blockHeader.proposerIndex
        );
    }

    /// @notice Verifies the validator pubkey is in the registry of beacon state.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param validatorPubkeyProof `bytes32[]` proof of the validator.
    /// @param validatorPubkey `ValidatorPubkey` to verify.
    /// @param validatorIndex `uint64` index of the validator.
    function proveValidatorPubkeyInBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata validatorPubkeyProof,
        bytes calldata validatorPubkey,
        uint64 validatorIndex
    )
        internal
        view
    {
        if (validatorIndex >= VALIDATOR_REGISTRY_LIMIT) {
            revert IndexOutOfRange();
        }

        uint256 gIndex = zeroValidatorPubkeyGIndex + (8 * validatorIndex);
        bytes32 validatorPubkeyRoot =
            SSZ.validatorPubkeyHashTreeRoot(validatorPubkey);

        if (
            !SSZ.verifyProof(
                validatorPubkeyProof,
                beaconBlockRoot,
                validatorPubkeyRoot,
                gIndex
            )
        ) {
            revert InvalidProof();
        }
    }
}
