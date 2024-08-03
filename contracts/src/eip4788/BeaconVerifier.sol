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
    /// @inheritdoc IBeaconVerifier
    uint256 public executionFeeRecipientGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(
        uint256 _zeroValidatorPubkeyGIndex,
        uint256 _executionNumberGIndex,
        uint256 _executionFeeRecipientGIndex
    ) {
        zeroValidatorPubkeyGIndex = _zeroValidatorPubkeyGIndex;
        executionNumberGIndex = _executionNumberGIndex;
        executionFeeRecipientGIndex = _executionFeeRecipientGIndex;

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

    function setExecutionFeeRecipientGIndex(
        uint256 _executionFeeRecipientGIndex
    )
        external
        onlyOwner
    {
        executionFeeRecipientGIndex = _executionFeeRecipientGIndex;
        emit ExecutionFeeRecipientGIndexChanged(_executionFeeRecipientGIndex);
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
        uint64 proposerIndex,
        bytes calldata proposerPubkey,
        bytes32[] calldata proposerPubkeyProof
    )
        external
        view
    {
        proveValidatorPubkeyInBeaconBlock(
            getParentBlockRoot(timestamp),
            proposerPubkeyProof,
            proposerPubkey,
            proposerIndex
        );
    }

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used ~41647
    function verifyExecutionNumber(
        uint64 timestamp,
        uint64 executionNumber,
        bytes32[] calldata executionNumberProof
    )
        external
        view
    {
        proveExecutionNumberInBeaconBlock(
<<<<<<< HEAD
            getParentBlockRoot(timestamp), executionNumberProof, blockNumber
=======
            getParentBlockRoot(timestamp), executionNumberProof, executionNumber
>>>>>>> main
        );
    }

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used ~41784
    function verifyCoinbase(
        uint64 timestamp,
<<<<<<< HEAD
        bytes32[] calldata coinbaseProof,
        address coinbase
=======
        address coinbase,
        bytes32[] calldata coinbaseProof
>>>>>>> main
    )
        external
        view
    {
        proveExecutionFeeRecipientInBeaconBlock(
            getParentBlockRoot(timestamp), coinbaseProof, coinbase
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

    /// @notice Verifies the block number in the latest execution payload header
    /// in the beacon state in the beacon block.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param executionNumberProof `bytes32[]` proof of the execution number.
    /// @param blockNumber `uint64` execution number of the beacon block.
    function proveExecutionNumberInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata executionNumberProof,
        uint64 blockNumber
    )
        internal
        view
    {
        bytes32 executionNumberRoot = SSZ.uint64HashTreeRoot(blockNumber);

        if (
            !SSZ.verifyProof(
                executionNumberProof,
                beaconBlockRoot,
                executionNumberRoot,
                executionNumberGIndex
            )
        ) revert InvalidProof();
    }

    /// @notice Verifies the coinbase (fee recipient) in the latest execution
    /// payload header in the beacon state in the beacon block.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param coinbaseProof `bytes32[]` proof of the coinbase.
<<<<<<< HEAD
    /// @param coinbase `address` to verify.
=======
    /// @param coinbase `address` fee recipient of the beacon block.
>>>>>>> main
    function proveExecutionFeeRecipientInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata coinbaseProof,
        address coinbase
    )
        internal
        view
    {
        bytes32 coinbaseRoot = SSZ.addressHashTreeRoot(coinbase);

        if (
            !SSZ.verifyProof(
                coinbaseProof,
                beaconBlockRoot,
                coinbaseRoot,
                executionFeeRecipientGIndex
            )
        ) revert InvalidProof();
    }
}
