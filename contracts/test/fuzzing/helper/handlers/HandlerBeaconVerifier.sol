// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "../../util/PreconditionsBase.sol";
import "../../properties/Properties.sol";
import "@openzeppelin-contracts/contracts/utils/Strings.sol";
import {Vm} from "@forge-std/Vm.sol";
import {BeaconRoots} from "../../libraries/BeaconRoots.sol";

/**
 * @title HandlerBeaconVerifier
 * @author 0xScourgedev, Rappie
 * @notice Fuzz handlers for BeaconVerifier
 */
contract HandlerBeaconVerifier is PreconditionsBase, Properties {
    using Strings for int64;
    using Strings for uint256;

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         STRUCTS                                           //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    struct SetProposerIndexGIndexParams {
        uint256 _proposerIndexGIndex;
    }

    struct SetZeroValidatorPubkeyGIndexParams {
        uint256 _zeroValidatorPubkeyGIndex;
    }

    struct VerifyValidatorPubkeyParams {
        bytes32 beaconBlockRoot;
        uint64 proposerIndex;
        bytes proposerPubkey;
        bytes32[] proposerPubkeyProof;
    }

    struct VerifyProposerIndexParams {
        bytes32 beaconBlockRoot;
        uint64 proposerIndex;
        bytes32[] proposerIndexProof;
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                      PRECONDITIONS                                        //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @notice This is a preconditions stub for setProposerIndexGIndex
     */
    function setProposerIndexGIndexPreconditions(uint64 _proposerIndexGIndex)
        internal
        returns (SetProposerIndexGIndexParams memory)
    {
        return SetProposerIndexGIndexParams({_proposerIndexGIndex: _proposerIndexGIndex});
    }

    /**
     * @notice This is a preconditions stub for setZeroValidatorPubkeyGIndex
     */
    function setZeroValidatorPubkeyGIndexPreconditions(uint64 _zeroValidatorPubkeyGIndex)
        internal
        returns (SetZeroValidatorPubkeyGIndexParams memory)
    {
        return SetZeroValidatorPubkeyGIndexParams({_zeroValidatorPubkeyGIndex: _zeroValidatorPubkeyGIndex});
    }

    /**
     * @notice Preconditions for the verifyValidatorPubkey function
     * @custom:precondition The valLen (validators length) is clamped between 1 and MAX_VALIDATORS_LENGTH
     * @custom:precondition Generate the beacon block proposer proof with the given seed and valLen from the go script
     * @custom:precondition If proofMod is true, shift the proposerPubkeyProof at the given shiftIndex
     * @custom:precondition Set the mock block root with what is returned from the proof generation script
     */
    function verifyValidatorPubkeyPreconditions(int64 seed, uint256 valLen, bool proofMod, uint16 shiftIndex)
        internal
        returns (VerifyValidatorPubkeyParams memory)
    {
        valLen = fl.clamp(valLen, 1, MAX_VALIDATORS_LENGTH);
        (
            bytes32 beaconBlockRoot,
            uint64 proposerIndex,
            bytes memory proposerPubkey,
            bytes32[] memory proposerPubkeyProof
        ) = generateValidatorPubkeyProof(seed, valLen);

        if (proofMod) {
            shiftIndex = uint16(fl.clamp(shiftIndex, 0, proposerPubkeyProof.length - 1));
            fl.log("OLD: ", proposerPubkeyProof[shiftIndex]);
            bytes32 newProof;
            if (proposerPubkeyProof[shiftIndex] == bytes32(0)) {
                newProof = bytes32(uint256(1));
            } else {
                newProof = proposerPubkeyProof[shiftIndex] >> 1;
            }
            proposerPubkeyProof[shiftIndex] = newProof;
            fl.log("NEW: ", proposerPubkeyProof[shiftIndex]);
        }

        setMockBlockRoot(beaconBlockRoot);

        return VerifyValidatorPubkeyParams({
            beaconBlockRoot: beaconBlockRoot,
            proposerIndex: proposerIndex,
            proposerPubkey: proposerPubkey,
            proposerPubkeyProof: proposerPubkeyProof
        });
    }

    /**
     * @notice Preconditions for the verifyProposerIndex function
     * @custom:precondition The valLen (validators length) is clamped between 1 and MAX_VALIDATORS_LENGTH
     * @custom:precondition Generate the execution number proof with the given seed and valLen from the go script
     * @custom:precondition If proofMod is true, shift the proposerIndexProof at the given shiftIndex
     * @custom:precondition Set the mock block root with what is returned from the proof generation script
     */
    function verifyProposerIndexPreconditions(int64 seed, uint256 valLen, bool proofMod, uint16 shiftIndex)
        internal
        returns (VerifyProposerIndexParams memory)
    {
        valLen = fl.clamp(valLen, 1, MAX_VALIDATORS_LENGTH);
        (bytes32 beaconBlockRoot, uint64 proposerIndex, bytes32[] memory proposerIndexProof) =
            generateProposerIndexProof(seed, valLen);

        if (proofMod) {
            shiftIndex = uint16(fl.clamp(shiftIndex, 0, proposerIndexProof.length - 1));
            fl.log("OLD: ", proposerIndexProof[shiftIndex]);
            bytes32 newProof;
            if (proposerIndexProof[shiftIndex] == bytes32(0)) {
                newProof = bytes32(uint256(1));
            } else {
                newProof = proposerIndexProof[shiftIndex] >> 1;
            }
            proposerIndexProof[shiftIndex] = newProof;
            fl.log("NEW: ", proposerIndexProof[shiftIndex]);
        }

        setMockBlockRoot(beaconBlockRoot);

        return VerifyProposerIndexParams({
            beaconBlockRoot: beaconBlockRoot,
            proposerIndex: proposerIndex,
            proposerIndexProof: proposerIndexProof
        });
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         HANDLERS                                          //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @notice This handler calls the setProposerIndexGIndex function
     * @param _proposerIndexGIndex The execution number GIndex to set to
     */
    function handler_setProposerIndexGIndex(uint64 _proposerIndexGIndex) public setCurrentActor {
        SetProposerIndexGIndexParams memory params = setProposerIndexGIndexPreconditions(_proposerIndexGIndex);

        (bool success, bytes memory returnData) = _setProposerIndexGIndexCall(params._proposerIndexGIndex);

        setProposerIndexGIndexPostconditions(success, returnData);
    }

    /**
     * @notice This handler calls the setZeroValidatorPubkeyGIndex function
     * @param _zeroValidatorPubkeyGIndex The zero validator pubkey GIndex to set to
     */
    function handler_setZeroValidatorPubkeyGIndex(uint64 _zeroValidatorPubkeyGIndex) public setCurrentActor {
        SetZeroValidatorPubkeyGIndexParams memory params =
            setZeroValidatorPubkeyGIndexPreconditions(_zeroValidatorPubkeyGIndex);

        (bool success, bytes memory returnData) = _setZeroValidatorPubkeyGIndexCall(params._zeroValidatorPubkeyGIndex);

        setZeroValidatorPubkeyGIndexPostconditions(success, returnData);
    }

    /**
     * @notice This handler calls the verifyValidatorPubkey function
     * @param seed The seed used for deterministic randomization in generating the proof
     * @param valLen The length of the validators array
     * @param proofMod Whether to modify the proof or not
     * @param shiftIndex The index to shift the proof at
     */
    function handler_verifyValidatorPubkey(int64 seed, uint256 valLen, bool proofMod, uint16 shiftIndex)
        public
        setCurrentActor
    {
        VerifyValidatorPubkeyParams memory params =
            verifyValidatorPubkeyPreconditions(seed, valLen, proofMod, shiftIndex);

        _before();

        (bool success, bytes memory returnData) = _verifyValidatorPubkeyCall(
            params.beaconBlockRoot, params.proposerPubkeyProof, params.proposerPubkey, params.proposerIndex
        );

        verifyValidatorPubkeyPostconditions(success, returnData, proofMod);
    }

    /**
     * @notice This handler calls the verifyProposerIndex function
     * @param seed The seed used for deterministic randomization in generating the proof
     * @param valLen The length of the validators array
     * @param proofMod Whether to modify the proof or not
     * @param shiftIndex The index to shift the proof at
     */
    function handler_verifyProposerIndex(int64 seed, uint256 valLen, bool proofMod, uint16 shiftIndex)
        public
        setCurrentActor
    {
        VerifyProposerIndexParams memory params = verifyProposerIndexPreconditions(seed, valLen, proofMod, shiftIndex);

        _before();

        (bool success, bytes memory returnData) =
            _verifyProposerIndexCall(params.beaconBlockRoot, params.proposerIndex, params.proposerIndexProof);

        verifyProposerIndexPostconditions(success, returnData, proofMod);
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                     POSTCONDITIONS                                        //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Postconditions for the setProposerIndexGIndex function
     * @custom:invariant REV-02: setProposerIndexGIndex never reverts
     */
    function setProposerIndexGIndexPostconditions(bool success, bytes memory returnData) internal {
        if (success) {} else {
            invariant_REV_02(bytes4(returnData));
        }
    }

    /**
     * @notice Postconditions for the setZeroValidatorPubkeyGIndex function
     * @custom:invariant REV-03: setZeroValidatorPubkeyGIndex never reverts
     */
    function setZeroValidatorPubkeyGIndexPostconditions(bool success, bytes memory returnData) internal {
        if (success) {} else {
            invariant_REV_01(bytes4(returnData));
        }
    }

    /**
     * @notice Postconditions for the verifyValidatorPubkey function
     * @custom:invariant PROOF-01: If the zeroValidatorPubkeyGIndex is different, the proof should never be valid
     * @custom:invariant PROOF-03: If the proof for verifyValidatorPubkey was not modified post-generation,
     * then the proof should always be valid
     *
     * @custom:invariant PROOF-05: If the zeroValidatorPubkeyGIndex is the same and the proof for
     * verifyValidatorPubkey was modified post-generation, then the proof should never be valid
     */
    function verifyValidatorPubkeyPostconditions(bool success, bytes memory returnData, bool proofMod) internal {
        if (success) {
            _after();
            invariant_PROOF_01();
            invariant_PROOF_03(proofMod);
        } else {
            invariant_PROOF_05(proofMod);
        }
    }

    /**
     * @notice Postconditions for the verifyProposerIndex function
     * @custom:invariant PROOF-02: If the proposerIndexGIndex is different, the proof should never be valid
     * @custom:invariant PROOF-04: If the proof for verifyProposerIndex was not modified post-generation,
     * then the proof should always be valid
     *
     * @custom:invariant PROOF-06: If the proposerIndexGIndex is the same and the proof for verifyProposerIndex
     * was modified post-generation, then the proof should never be valid
     */
    function verifyProposerIndexPostconditions(bool success, bytes memory returnData, bool proofMod) internal {
        if (success) {
            _after();
            invariant_PROOF_02();
            invariant_PROOF_04(proofMod);
        } else {
            invariant_PROOF_06(proofMod);
        }
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         HELPER                                            //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Generate the beacon block proposer proof from the go script with ffi using the given seed and valLen
     */
    function generateValidatorPubkeyProof(int64 seed, uint256 valLen)
        internal
        returns (bytes32, uint64, bytes memory, bytes32[] memory)
    {
        string[] memory inputs = new string[](6);
        inputs[0] = "go";
        inputs[1] = "run";
        inputs[2] = "./test/fuzzing/script/proof_gen";

        inputs[3] = string.concat("-seed=", seed.toStringSigned());
        inputs[4] = string.concat("-sel=", "0");
        inputs[5] = string.concat("-valLen=", valLen.toString());

        bytes memory result = vm.ffi(inputs);

        return abi.decode(result, (bytes32, uint64, bytes, bytes32[]));
    }

    /**
     * @notice Generate the execution number proof from the go script with ffi using the given seed and valLen
     */
    function generateProposerIndexProof(int64 seed, uint256 valLen)
        internal
        returns (bytes32, uint64, bytes32[] memory)
    {
        string[] memory inputs = new string[](6);
        inputs[0] = "go";
        inputs[1] = "run";
        inputs[2] = "./test/fuzzing/script/proof_gen";

        inputs[3] = string.concat("-seed=", seed.toStringSigned());
        inputs[4] = string.concat("-sel=", "1");
        inputs[5] = string.concat("-valLen=", valLen.toString());

        bytes memory result = vm.ffi(inputs);

        return abi.decode(result, (bytes32, uint64, bytes32[]));
    }

    /**
     * @notice Set the mock block root with the given beaconBlockRoot
     */
    function setMockBlockRoot(bytes32 beaconBlockRoot) internal {
        BeaconRootMock(payable(BeaconRoots.ADDRESS)).setBeaconRoot(beaconBlockRoot);
    }
}
