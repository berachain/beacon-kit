// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Test } from "@forge-std/Test.sol";
import { stdJson } from "@forge-std/StdJson.sol";

import { SSZ } from "@src/eip4788/SSZ.sol";
import { ValidatorVerifier } from "@src/eip4788/ValidatorVerifier.sol";

contract ValidatorVerifierTest is Test {
    using stdJson for string;

    struct ProofJson {
        bytes32[] validatorProof;
        SSZ.Validator validator;
        uint64 validatorIndex;
        bytes32 blockRoot;
    }

    uint256 internal constant DENEB_ZERO_VALIDATOR_GINDEX = 798_245_441_765_376;

    ValidatorVerifier public verifier;
    ProofJson public proofJson;

    function setUp() public {
        string memory json =
            vm.readFile("./test/eip4788/fixtures/validator_proof.json");
        bytes memory data = json.parseRaw("$");
        proofJson = abi.decode(data, (ProofJson));
    }

    function test_ProveValidator() public {
        uint64 ts = 31_337;

        verifier = new ValidatorVerifier(DENEB_ZERO_VALIDATOR_GINDEX);

        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );

        verifier.proveValidator(
            proofJson.validatorProof,
            proofJson.validator,
            proofJson.validatorIndex,
            ts
        );
    }

    function test_ProveValidatorWrongIndex() public {
        uint64 ts = 31_337;

        verifier = new ValidatorVerifier(DENEB_ZERO_VALIDATOR_GINDEX);

        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );
        vm.expectRevert(bytes4(keccak256("IndexOutOfRange()")));
        verifier.proveValidator(
            proofJson.validatorProof,
            proofJson.validator,
            (1 << 40) + 1, // VALIDATOR_REGISTRY_LIMIT + 1
            ts
        );
    }

    function test_ProveValidatorInvalidProof() public {
        uint64 ts = 31_337;

        verifier = new ValidatorVerifier(DENEB_ZERO_VALIDATOR_GINDEX);

        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );
        vm.expectRevert(bytes4(keccak256("InvalidProof()")));
        verifier.proveValidator(
            proofJson.validatorProof,
            proofJson.validator,
            proofJson.validatorIndex - 1,
            ts
        );
    }

    function test_ProveValidator_OnFork() public {
        string memory forkUrl = vm.envOr("FORK_URL", string(""));
        vm.skip(_isEmptyString(forkUrl));
        vm.createSelectFork(forkUrl);
        _checkChainId(5); // Only works on Goerli for now.

        // Timestamp of the block which parent root is a `proofJson.blockRoot`.
        uint64 ts = 1_705_602_156;

        // Move to a block at ts 1705605588 (doesn't matter, but a root at `ts`
        // should be still available to work.
        vm.rollFork(10_395_866);

        verifier = new ValidatorVerifier(DENEB_ZERO_VALIDATOR_GINDEX);

        verifier.proveValidator(
            proofJson.validatorProof,
            proofJson.validator,
            proofJson.validatorIndex,
            ts
        );
    }

    function _checkChainId(uint256 chainId) internal view {
        if (chainId != block.chainid) {
            revert("wrong chain id");
        }
    }

    function _isEmptyString(string memory str) internal pure returns (bool) {
        return bytes(str).length == 0;
    }
}
