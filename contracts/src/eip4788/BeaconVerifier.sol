// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Ownable } from "@solady/src/auth/Ownable.sol";

import { SSZ } from "./SSZ.sol";

import { Verifier } from "./Verifier.sol";

import { IBeaconVerifier } from "./interfaces/IBeaconVerifier.sol";

/// @author Berachain Team
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
/// eigenlayer approach: 0x124363b6D0866118A8b6899F2674856618E0Ea4c
/// madlabman approach: 0x5793a71D3eF074f71dCC21216Dbfd5C0e780132c
contract BeaconVerifier is Verifier, Ownable, IBeaconVerifier {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    uint256 public zeroValidatorsGIndex;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       ADMIN FUNCTIONS                      */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    constructor(uint256 _zeroValGIndex) {
        zeroValidatorsGIndex = _zeroValGIndex;

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

    function setZeroValidatorsGIndex(uint256 _zeroValGIndex) external onlyOwner {
        zeroValidatorsGIndex = _zeroValGIndex;
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

    function getBeaconHeaderHTREigenlayer(SSZ.BeaconBlockHeader calldata blockHeader)
        public
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
        return SSZ.merkleizeSha256Eigenlayer(nodes);
    }

    function getValidatorHTREigenlayer(SSZ.Validator calldata validator)
        public
        view
        returns (bytes32)
    {
        bytes32 pubkeyRoot;
        assembly {
            // Dynamic data types such as bytes are stored at the specified offset.
            let offset := mload(validator)
            // Call sha256 precompile (0x02) with the pubkey pointer
            let result := staticcall(gas(), 0x02, add(offset, 32), 0x40, 0x00, 0x20)
            // Precompile returns no data on OutOfGas error.
            if eq(result, 0) { revert(0, 0) }
            pubkeyRoot := mload(0x00)
        }
        bytes32[] memory nodes = new bytes32[](8);
        nodes[0] = pubkeyRoot;
        nodes[1] = validator.withdrawalCredentials;
        nodes[2] = SSZ.toLittleEndian(validator.effectiveBalance);
        nodes[3] = SSZ.toLittleEndian(validator.slashed);
        nodes[4] = SSZ.toLittleEndian(validator.activationEligibilityEpoch);
        nodes[5] = SSZ.toLittleEndian(validator.activationEpoch);
        nodes[6] = SSZ.toLittleEndian(validator.exitEpoch);
        nodes[7] = SSZ.toLittleEndian(validator.withdrawableEpoch);
        return SSZ.merkleizeSha256Eigenlayer(nodes);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           PROOFS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconVerifier
    /// @dev gas used with eigenlayer ~600568
    /// @dev gas used with madlabman ~84783
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
        // bytes32 givenBeaconRoot = getBeaconHeaderHTREigenlayer(blockHeader);
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

        uint256 gIndex = zeroValidatorsGIndex + validatorIndex;
        bytes32 validatorRoot = SSZ.validatorHashTreeRoot(validator);
        // bytes32 validatorRoot = getValidatorHTREigenlayer(validator);

        if (
            !SSZ.verifyProof(validatorProof, beaconBlockRoot, validatorRoot, gIndex)
        ) {
            revert InvalidProof();
        }

        // // Convert validatorProof to bytes memory and verify with eigenlayer.
        // bytes memory proof = new bytes(validatorProof.length * 32);
        // for (uint256 i = 0; i < validatorProof.length; i++) {
        //     bytes32 node = validatorProof[i];
        //     for (uint256 j = 0; j < 32; j++) {
        //         proof[i * 32 + j] = node[j];
        //     }
        // }
        // if (
        //     !SSZ.verifyInclusionSha256Eigenlayer(proof, beaconBlockRoot, validatorRoot, gIndex)
        // ) {
        //     revert InvalidProof();
        // }
    }
}
