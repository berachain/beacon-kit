// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

interface IBeaconVerifier {
    /// @notice Emitted when the zero validator pubkey Generalized Index is
    /// changed.
    event ZeroValidatorPubkeyGIndexChanged(
        uint256 newZeroValidatorPubkeyGIndex
    );
    /// @notice Emitted when the execution number Generalized Index is changed.
    event ExecutionNumberGIndexChanged(uint256 newExecutionNumberGIndex);
    /// @notice Emitted when the execution fee recipient Generalized Index is
    /// changed.
    event ExecutionFeeRecipientGIndexChanged(
        uint256 newExecutionFeeRecipientGIndex
    );

    /// @notice Generalized Index of the pubkey of the first validator
    /// (validator index of 0) in the registry of the beacon state in the
    /// beacon block.
    /// @dev In the Deneb beacon chain fork, this should be 3254554418216960.
    function zeroValidatorPubkeyGIndex() external view returns (uint256);

    /// @notice Generalized Index of the block number in the latest execution
    /// payload header in the beacon state in the beacon block.
    /// @dev In the Deneb beacon chain fork, this should be 5894.
    function executionNumberGIndex() external view returns (uint256);

    /// @notice Generalized Index of the fee recipient in the latest execution
    /// payload header in the beacon state in the beacon block.
    /// @dev In the Deneb beacon chain fork, this should be 5889.
    function executionFeeRecipientGIndex() external view returns (uint256);

    /// @notice Get the parent beacon block root from the given timestamp.
    function getParentBeaconBlockRootAt(uint64 timestamp)
        external
        view
        returns (bytes32);

    /// @notice Get the parent beacon block root at `block.timestamp`.
    function getParentBeaconBlockRoot() external view returns (bytes32);

    /// @notice Verifies the proposer within the beacon block at the given
    /// timestamp. Reverts if proof invalid.
    /// @param timestamp `uint64` timestamp of the parent beacon block.
<<<<<<< HEAD
    /// @param validatorPubkeyProof `bytes32[]` proof of the validator pubkey.
    /// @param validatorPubkey `ValidatorPubkey` to verify.
    /// @param proposerIndex `uint64` validator index of the proposer of the
    /// parent beacon block.
=======
    /// @param proposerIndex `uint64` validator index of the proposer of the
    /// parent beacon block.
    /// @param proposerPubkey `bytes` proposer validator pubkey to verify.
    /// @param proposerPubkeyProof `bytes32[]` proof of the proposer validator
    /// pubkey.
>>>>>>> main
    function verifyBeaconBlockProposer(
        uint64 timestamp,
        uint64 proposerIndex,
        bytes calldata proposerPubkey,
        bytes32[] calldata proposerPubkeyProof
    )
        external
        view;

    /// @notice Verifies the execution number in the parent beacon block at the
    /// given timestamp. Reverts if proof invalid.
    /// @param timestamp `uint64` timestamp of the parent beacon block.
<<<<<<< HEAD
    /// @param executionNumberProof `bytes32[]` proof of the execution number.
    /// @param blockNumber `uint64` execution number of the parent beacon block.
=======
    /// @param executionNumber `uint64` execution number of the parent beacon
    /// block to verify.
    /// @param executionNumberProof `bytes32[]` proof of the execution number.
>>>>>>> main
    function verifyExecutionNumber(
        uint64 timestamp,
        uint64 executionNumber,
        bytes32[] calldata executionNumberProof
    )
        external
        view;

    /// @notice Verifies the coinbase (fee recipient) in the parent beacon
    /// block at the given timestamp. Reverts if proof invalid.
    /// @param timestamp `uint64` timestamp of the parent beacon block.
    /// @param coinbase `address` fee recipient of the parent beacon block to
    /// verify.
    /// @param coinbaseProof `bytes32[]` proof of the fee recipient.
    function verifyCoinbase(
        uint64 timestamp,
        address coinbase,
        bytes32[] calldata coinbaseProof
    )
        external
        view;

    /// @notice Verifies the coinbase (fee recipient) in the parent beacon
    /// block at the given timestamp. Reverts if proof invalid.
    /// @param timestamp `uint64` timestamp of the parent beacon block.
    /// @param coinbaseProof `bytes32[]` proof of the coinbase.
    /// @param coinbase `address` fee recipient of the parent beacon block.
    function verifyCoinbase(
        uint64 timestamp,
        bytes32[] calldata coinbaseProof,
        address coinbase
    )
        external
        view;
}
