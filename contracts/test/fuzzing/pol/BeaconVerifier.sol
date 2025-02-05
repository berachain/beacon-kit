// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import {BeaconRootsHelper} from "./BeaconRootsHelper.sol";

/// @author Berachain Team
contract BeaconVerifier is BeaconRootsHelper {
    constructor(uint64 _zeroValidatorPubkeyGIndex, uint64 _proposerIndexGIndex) {
        super.setZeroValidatorPubkeyGIndex(_zeroValidatorPubkeyGIndex);
        super.setProposerIndexGIndex(_proposerIndexGIndex);
    }

    function setZeroValidatorPubkeyGIndex(uint64 _zeroValidatorPubkeyGIndex) public override {
        super.setZeroValidatorPubkeyGIndex(_zeroValidatorPubkeyGIndex);
    }

    function setProposerIndexGIndex(uint64 _proposerIndexGIndex) public override {
        super.setProposerIndexGIndex(_proposerIndexGIndex);
    }

    function verifyProposerIndexInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata proposerIndexProof,
        uint64 proposerIndex
    ) public view {
        return super._verifyProposerIndexInBeaconBlock(beaconBlockRoot, proposerIndexProof, proposerIndex);
    }

    function verifyValidatorPubkeyInBeaconBlock(
        bytes32 beaconBlockRoot,
        bytes32[] calldata validatorPubkeyProof,
        bytes calldata validatorPubkey,
        uint64 validatorIndex
    ) public view {
        return super._verifyValidatorPubkeyInBeaconBlock(
            beaconBlockRoot, validatorPubkeyProof, validatorPubkey, validatorIndex
        );
    }
}
