// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { FixedPointMathLib } from "@solady/src/utils/FixedPointMathLib.sol";

import { IPOLErrors } from "./interfaces/IPOLErrors.sol";

import { BeaconRoots } from "../libraries/BeaconRoots.sol";
import { SSZ } from "@src/eip4788/SSZ.sol";
import { Utils } from "../libraries/Utils.sol";

/// @title BeaconRootsHelper
/// @author Berachain Team
/// @notice A contract that follows the beacon chain roots (using the EIP-4788 Beacon Roots contract) to
/// maintain a buffer of processed timestamps. The buffer mirrors the buffer used in the Beacon Roots contract
/// for consistency, but is not 1:1 since we are not guaranteed updates to this buffer every block.
/// @notice The contract also allows verifying the pubkey of the proposer of a beacon block.
abstract contract BeaconRootsHelper is IPOLErrors {
    using BeaconRoots for uint64;
    using Utils for bytes4;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                         CONSTANTS                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The length of the history buffer in the EIP-4788 Beacon Roots contract.
    uint64 private constant HISTORY_BUFFER_LENGTH = 8191;
    /// @dev The validator registry limit in the Beacon spec.
    uint64 private constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;
    /// @dev The validator pubkey offset of Generalized Index based on the fields of the validator SSZ container in the
    /// Beacon spec.
    uint8 private constant VALIDATOR_PUBKEY_OFFSET = 8;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          EVENTS                            */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Emitted when the zero validator pubkey Generalized Index is changed.
    event ZeroValidatorPubkeyGIndexChanged(uint64 newZeroValidatorPubkeyGIndex);

    /// @notice Emitted when the proposer index Generalized Index is changed.
    event ProposerIndexGIndexChanged(uint64 newProposerIndexGIndex);

    /// @notice Emitted when a timestamp is successfully processed in the history buffer.
    event TimestampProcessed(uint64 timestamp);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Mapping `timestamp_idx` (timestamp % HISTORY_BUFFER_LENGTH) in the history buffer to the
    /// processed timestamp.
    /// @dev Ensure that the value at `timestamp_idx` is the desired timestamp to process.
    uint64[HISTORY_BUFFER_LENGTH] private _processedTimestampsBuffer;

    /// @notice Generalized Index of the pubkey of the first validator (validator index of 0) in the registry of the
    /// beacon state in the beacon block.
    /// @dev In the Deneb beacon fork on Berachain, this should be 3254554418216960.
    uint64 public zeroValidatorPubkeyGIndex;

    /// @notice Generalized Index of the proposer index in the beacon block.
    /// @dev In the Deneb beacon fork on Berachain, this should be 9.
    uint64 public proposerIndexGIndex;

    /// @dev Storage gap for future upgrades.
    uint256[50] private _gap;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            ADMIN                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev This action should be permissioned to prevent unauthorized actors from modifying inappropriately.
    function setZeroValidatorPubkeyGIndex(uint64 _zeroValidatorPubkeyGIndex) public virtual {
        zeroValidatorPubkeyGIndex = _zeroValidatorPubkeyGIndex;
        emit ZeroValidatorPubkeyGIndexChanged(_zeroValidatorPubkeyGIndex);
    }

    /// @dev This action should be permissioned to prevent unauthorized actors from modifying inappropriately.
    function setProposerIndexGIndex(uint64 _proposerIndexGIndex) public virtual {
        proposerIndexGIndex = _proposerIndexGIndex;
        emit ProposerIndexGIndexChanged(_proposerIndexGIndex);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                  EXTERNAL VIEW FUNCTIONS                   */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Returns whether a timestamp is in the history buffer and can still be processed.
    /// @param timestamp The timestamp to search for.
    /// @return actionable Whether the timestamp is in the history buffer and it has not yet been processed.
    function isTimestampActionable(uint64 timestamp) external view returns (bool) {
        // First check if the timestamp is in the Beacon Roots history buffer.
        if (!timestamp.isParentBlockRootAt()) return false;

        // If we know the timestamp is in the buffer, return if the timestamp has not been processed.
        uint64 timestampIndex = timestamp % HISTORY_BUFFER_LENGTH;
        return _processedTimestampsBuffer[timestampIndex] != timestamp;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     INTERNAL FUNCTIONS                     */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Processes the timestamp in the history buffer.
    /// @param timestamp The timestamp to process in the history buffer.
    /// @dev Reverts if the timestamp is not available in the Beacon Roots history buffer.
    /// @dev Reverts if the timestamp is available in the Beacon Roots history buffer but has already been processed.
    /// @return parentBeaconBlockRoot The parent beacon block root at the given timestamp.
    function _processTimestampInBuffer(uint64 timestamp) internal returns (bytes32 parentBeaconBlockRoot) {
        // First enforce the timestamp is in the Beacon Roots history buffer, reverting if not found.
        parentBeaconBlockRoot = timestamp.getParentBlockRootAt();

        // Mark the in buffer timestamp as processed if it has not been processed yet.
        uint64 timestampIndex = timestamp % HISTORY_BUFFER_LENGTH;
        if (timestamp == _processedTimestampsBuffer[timestampIndex]) TimestampAlreadyProcessed.selector.revertWith();
        _processedTimestampsBuffer[timestampIndex] = timestamp;

        // Emit the event that the timestamp has been processed.
        emit TimestampProcessed(timestamp);
    }

    /// @notice Verifies the proposer index is in the beacon block, reverting if the proof is invalid.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param proposerIndexProof `bytes32[]` proof of the proposer index in the beacon block.
    /// @param proposerIndex `uint64` proposer index to verify.
    function _verifyProposerIndexInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata proposerIndexProof,
        uint64 proposerIndex
    )
        internal
        view
    {
        bytes32 proposerIndexRoot = SSZ.uint64HashTreeRoot(proposerIndex);

        if (!SSZ.verifyProof(proposerIndexProof, beaconBlockRoot, proposerIndexRoot, proposerIndexGIndex)) {
            InvalidProof.selector.revertWith();
        }
    }

    /// @notice Verifies the validator pubkey is in the registry of beacon state at the given validator index,
    /// reverting if the proof is invalid.
    /// @param beaconBlockRoot `bytes32` root of the beacon block.
    /// @param validatorPubkeyProof `bytes32[]` proof of the validator pubkey in the beacon block.
    /// @param validatorPubkey `bytes` 40 byte validator pubkey to verify.
    /// @param validatorIndex `uint64` validator index in the validator registry of the beacon state.
    function _verifyValidatorPubkeyInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata validatorPubkeyProof,
        bytes calldata validatorPubkey,
        uint64 validatorIndex
    )
        internal
        view
    {
        if (validatorIndex >= VALIDATOR_REGISTRY_LIMIT) {
            IndexOutOfRange.selector.revertWith();
        }

        bytes32 validatorPubkeyRoot = SSZ.validatorPubkeyHashTreeRoot(validatorPubkey);
        uint256 gIndex = zeroValidatorPubkeyGIndex + (VALIDATOR_PUBKEY_OFFSET * validatorIndex);

        if (!SSZ.verifyProof(validatorPubkeyProof, beaconBlockRoot, validatorPubkeyRoot, gIndex)) {
            InvalidProof.selector.revertWith();
        }
    }
}
