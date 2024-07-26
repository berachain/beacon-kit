// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Ownable } from "@solady/src/auth/Ownable.sol";

import { SSZ } from "./SSZ.sol";

import { Verifier } from "./Verifier.sol";

import { IBeaconVerifier } from "./interfaces/IBeaconVerifier.sol";

/// @author Berachain Team
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
/// 0x3945f611Fe77A51C7F3e1f84709C1a2fDcDfAC5B
contract BeaconVerifier is Verifier, Ownable, IBeaconVerifier {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    uint256 public validatorsGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(uint256 _valGIndex) {
        validatorsGIndex = _valGIndex;

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

    function setValidatorsGIndex(uint256 _valGIndex) external onlyOwner {
        validatorsGIndex = _valGIndex;
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
    /*                       TEMP HELPERS                         */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    function getBeaconHeaderHTR(SSZ.BeaconBlockHeader calldata blockHeader)
        external
        pure
        returns (bytes32)
    {
        bytes32[] memory nodes = new bytes32[](8);
        nodes[0] = SSZ.toLittleEndian(blockHeader.slot);
        nodes[1] = SSZ.toLittleEndian(blockHeader.proposerIndex);
        nodes[2] = blockHeader.parentRoot;
        nodes[3] = blockHeader.stateRoot;
        nodes[4] = blockHeader.bodyRoot;
        nodes[5] = bytes32(0);
        nodes[6] = bytes32(0);
        nodes[7] = bytes32(0);
        return SSZ.merkleizeSha256(nodes);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           PROOFS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    function proveBlockProposer(
        SSZ.BeaconBlockHeader calldata blockHeader,
        uint64 timestamp,
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator
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

        // Finally check that the validator is a validator of the beacon chain during this time.
        proveValidatorInBlock(
            expectedBeaconRoot,
            validatorProof,
            validator,
            blockHeader.proposerIndex
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

        uint256 gI = validatorsGIndex + validatorIndex;
        bytes32 validatorRoot = SSZ.validatorHashTreeRoot(validator);

        if (
            !SSZ.verifyProof(validatorProof, beaconBlockRoot, validatorRoot, gI)
        ) {
            revert InvalidProof();
        }
    }
}
