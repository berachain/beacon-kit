// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Test } from "@forge-std/Test.sol";
import { stdJson } from "@forge-std/StdJson.sol";

import { SSZ } from "@src/eip4788/SSZ.sol";
import { WithdrawalsVerifier } from "@src/eip4788/WithdrawalsVerifier.sol";

contract WithdrawalsVerifierTest is Test {
    using stdJson for string;

    struct ProofJson {
        bytes32[] withdrawalProof;
        SSZ.Withdrawal withdrawal;
        uint8 withdrawalIndex;
        bytes32 blockRoot;
    }

    uint256 internal constant DENEB_ZERO_WITHDRAWAL_GINDEX = 206_272;

    WithdrawalsVerifier public verifier;
    ProofJson public proofJson;

    function setUp() public {
        string memory json =
            vm.readFile("./test/eip4788/fixtures/withdrawal_proof.json");
        bytes memory data = json.parseRaw("$");
        proofJson = abi.decode(data, (ProofJson));
        verifier = new WithdrawalsVerifier(DENEB_ZERO_WITHDRAWAL_GINDEX);
    }

    function test_SubmitWithdrawal() public {
        uint64 ts = 31_337;

        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );

        // forgefmt: disable-next-item
        verifier.submitWithdrawal(
            proofJson.withdrawalProof,
            proofJson.withdrawal,
            proofJson.withdrawalIndex,
            ts
        );
    }

    function test_SubmitWithdrawalWrongIndex() public {
        uint64 ts = 31_337;

        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );
        vm.expectRevert(bytes4(keccak256("IndexOutOfRange()")));
        // forgefmt: disable-next-item
        verifier.submitWithdrawal(
            proofJson.withdrawalProof,
            proofJson.withdrawal,
            (1 << 4) + 1, // MAX_WITHDRAWALS + 1
            ts
        );
    }

    function test_SubmitWithdrawalInvalidProof() public {
        uint64 ts = 31_337;
        vm.mockCall(
            verifier.BEACON_ROOTS(),
            abi.encode(ts),
            abi.encode(proofJson.blockRoot)
        );
        vm.expectRevert(bytes4(keccak256("InvalidProof()")));
        // forgefmt: disable-next-item
        verifier.submitWithdrawal(
            proofJson.withdrawalProof,
            proofJson.withdrawal,
            proofJson.withdrawalIndex-1,
            ts
        );
    }
}
