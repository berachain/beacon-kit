// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Test } from "@forge-std/Test.sol";
import { stdJson } from "@forge-std/StdJson.sol";
import { Vm } from "@forge-std/Vm.sol";

import { SSZ } from "@src/eip4788/SSZ.sol";
import { BeaconVerifier } from "@src/eip4788/BeaconVerifier.sol";

contract BeaconVerifierTest is Test {
    using stdJson for string;

    struct BlockProposerProofJson {
        bytes32 beaconBlockRoot;
        uint64 proposerIndex;
        bytes proposerPubkey;
        bytes32[] proposerPubkeyProof;
    }

    struct CoinbaseProofJson {
        bytes32 beaconBlockRoot;
        address coinbase;
        bytes32[] coinbaseProof;
    }

    struct ExecutionNumberProofJson {
        bytes32 beaconBlockRoot;
        uint64 executionNumber;
        bytes32[] executionNumberProof;
    }

    uint256 constant DENEB_ZERO_VALIDATOR_PUBKEY_GINDEX = 3_254_554_418_216_960;
    uint256 constant DENEB_EXECUTION_NUMBER_GINDEX = 5894;
    uint256 constant DENEB_EXECUTION_FEE_RECIPIENT_GINDEX = 5889;

    uint64 timestamp = 31_337;
    BeaconVerifier public verifier;
    BlockProposerProofJson public blockProposerProofJson;
    CoinbaseProofJson public coinbaseProofJson;
    ExecutionNumberProofJson public executionNumberProofJson;

    function setUp() public {
        string memory root = vm.projectRoot();

        string memory blockProposerPath = string.concat(
            root, "/test/eip4788/fixtures/block_proposer_proof.json"
        );
        string memory blockProposerJson = vm.readFile(blockProposerPath);
        bytes memory blockProposerData = blockProposerJson.parseRaw("$");
        blockProposerProofJson =
            abi.decode(blockProposerData, (BlockProposerProofJson));

        string memory coinbasePath =
            string.concat(root, "/test/eip4788/fixtures/coinbase_proof.json");
        string memory coinbaseJson = vm.readFile(coinbasePath);
        bytes memory coinbaseData = coinbaseJson.parseRaw("$");
        coinbaseProofJson = abi.decode(coinbaseData, (CoinbaseProofJson));

        string memory executionNumberPath = string.concat(
            root, "/test/eip4788/fixtures/execution_number_proof.json"
        );
        string memory executionNumberJson = vm.readFile(executionNumberPath);
        bytes memory executionNumberData = executionNumberJson.parseRaw("$");
        executionNumberProofJson =
            abi.decode(executionNumberData, (ExecutionNumberProofJson));

        verifier = new BeaconVerifier(
            DENEB_ZERO_VALIDATOR_PUBKEY_GINDEX,
            DENEB_EXECUTION_NUMBER_GINDEX,
            DENEB_EXECUTION_FEE_RECIPIENT_GINDEX
        );
    }

    function test_verifyBeaconBlockProposer() public {
        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(timestamp),
            abi.encode(blockProposerProofJson.beaconBlockRoot)
        );

        verifier.verifyBeaconBlockProposer(
            timestamp,
            blockProposerProofJson.proposerIndex,
            blockProposerProofJson.proposerPubkey,
            blockProposerProofJson.proposerPubkeyProof
        );
    }

    function test_verifyCoinbase() public {
        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(timestamp),
            abi.encode(coinbaseProofJson.beaconBlockRoot)
        );

        verifier.verifyCoinbase(
            timestamp,
            coinbaseProofJson.coinbase,
            coinbaseProofJson.coinbaseProof
        );
    }

    function test_verifyExecutionNumber() public {
        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(timestamp),
            abi.encode(executionNumberProofJson.beaconBlockRoot)
        );

        verifier.verifyExecutionNumber(
            timestamp,
            executionNumberProofJson.executionNumber,
            executionNumberProofJson.executionNumberProof
        );
    }
}
