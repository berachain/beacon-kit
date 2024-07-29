// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Ownable } from "@solady/src/auth/Ownable.sol";

import { IBeaconVerifier } from "./interfaces/IBeaconVerifier.sol";

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";

/// @author Berachain Team
contract BeaconVerifier is Verifier, Ownable, IBeaconVerifier {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;
    uint8 internal constant VALIDATOR_PUBKEY_OFFSET = 8;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    uint256 public zeroValidatorPubkeyGIndex;
    /// @inheritdoc IBeaconVerifier
    uint256 public executionNumberGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(
        uint256 _zeroValidatorPubkeyGIndex,
        uint256 _executionNumberGIndex
    ) {
        zeroValidatorPubkeyGIndex = _zeroValidatorPubkeyGIndex;
        executionNumberGIndex = _executionNumberGIndex;

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
        emit ZeroValidatorPubkeyGIndexChanged(_zeroValidatorPubkeyGIndex);
    }

    function setExecutionNumberGIndex(uint256 _executionNumberGIndex)
        external
        onlyOwner
    {
        executionNumberGIndex = _executionNumberGIndex;
        emit ExecutionNumberGIndexChanged(_executionNumberGIndex);
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
    /*                         VERIFIERS                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used ~75381
    function verifyBeaconBlockProposer(
        uint64 timestamp,
        bytes32[] calldata validatorPubkeyProof,
        bytes calldata validatorPubkey,
        uint64 proposerIndex
    )
        external
        view
    {
        proveValidatorPubkeyInBeaconBlock(
            getParentBlockRoot(timestamp),
            validatorPubkeyProof,
            validatorPubkey,
            proposerIndex
        );
    }

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used ~41647
    function verifyExecutionNumber(
        uint64 timestamp,
        bytes32[] calldata executionNumberProof,
        uint64 blockNumber
    )
        external
        view
    {
        proveExecutionNumberInBeaconBlock(
            getParentBlockRoot(uint64(timestamp)),
            executionNumberProof,
            blockNumber
        );
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           PROOFS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Verifies the validator pubkey is in the registry of beacon state.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param validatorPubkeyProof `bytes32[]` proof of the validator.
    /// @param validatorPubkey `ValidatorPubkey` to verify.
    /// @param validatorIndex `uint64` index of the validator.
    function proveValidatorPubkeyInBeaconBlock(
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

        bytes32 validatorPubkeyRoot =
            SSZ.validatorPubkeyHashTreeRoot(validatorPubkey);
        uint256 gIndex = zeroValidatorPubkeyGIndex
            + (VALIDATOR_PUBKEY_OFFSET * validatorIndex);

        if (
            !SSZ.verifyProof(
                validatorPubkeyProof,
                beaconBlockRoot,
                validatorPubkeyRoot,
                gIndex
            )
        ) revert InvalidProof();
    }

    /// @notice Verifies the execution number in the beacon block.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param executionNumberProof `bytes32[]` proof of the execution number.
    /// @param blockNumber `uint64` execution number of the block.
    function proveExecutionNumberInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata executionNumberProof,
        uint64 blockNumber
    )
        internal
        view
    {
        bytes32 executionNumberRoot = SSZ.toLittleEndian(uint256(blockNumber));

        if (
            !SSZ.verifyProof(
                executionNumberProof,
                beaconBlockRoot,
                executionNumberRoot,
                executionNumberGIndex
            )
        ) revert InvalidProof();
    }
}
