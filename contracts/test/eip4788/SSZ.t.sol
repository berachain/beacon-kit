// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Test } from "@forge-std/Test.sol";
import { SSZ } from "@src/eip4788/SSZ.sol";

contract SSZTest is Test {
    function test_toLittleEndian() public pure {
        uint256 v = 0x1234567890ABCDEF;
        bytes32 expected =
            bytes32(bytes.concat(hex"EFCDAB9078563412", bytes24(0)));
        bytes32 actual = SSZ.toLittleEndian(v);
        assertEq(actual, expected);
    }

    function test_log2floor() public pure {
        uint64 v = 31;
        uint64 expected = 4;
        uint64 actual = uint64(SSZ.log2(v));
        assertEq(actual, expected);
    }

    function test_concatGIndices() public pure {
        uint64 expected = 3230;
        uint64 actual = SSZ.concatGindices(SSZ.concatGindices(12, 25), 30);
        assertEq(actual, expected);
    }

    function test_ValidatorPubkeyRoot() public view {
        bytes memory pubkey = hex"b68a04945cecb9157b40f31ecd5fba35dec6c30a4b3a772234ecb477170c198713bb9400c2142e7b3f88af9abb3362bb";
        bytes32 expected = hex"f5a5fd42d16a20302798ef6ed309979b43003d2320d9f0e8ea9831a92759fb4b";
        bytes32 actual = SSZ.validatorPubkeyHashTreeRoot(pubkey);
        assertEq(actual, expected);
    }

    /// Slot 7172576 withdrawal.at(0)
    function test_withdrawalRoot() public pure {
        SSZ.Withdrawal memory w = SSZ.Withdrawal({
            index: 15_213_404,
            validatorIndex: 429_156,
            _address: 0xB9D7934878B5FB9610B3fE8A5e441e8fad7E293f,
            amount: 15_428_006
        });
        bytes32 expected =
            0x900838206a9d83fec95bd54289eb52a8500cbb4a198d000f9f9c2c0662bb8fa2;
        bytes32 actual = SSZ.withdrawalHashTreeRoot(w);
        assertEq(actual, expected);
    }

    // Slot 7172576 validator.at(1025214)
    function test_ValidatorRootExitedSlashed() public view {
        SSZ.Validator memory v = SSZ.Validator({
            pubkey: hex"91760f8a17729cfcb68bfc621438e5d9dfa831cd648e7b2b7d33540a7cbfda1257e4405e67cd8d3260351ab3ff71b213",
            withdrawalCredentials: 0x01000000000000000000000006676e8584342cc8b6052cfdf381c3a281f00ac8,
            effectiveBalance: 30_000_000_000,
            slashed: true,
            activationEligibilityEpoch: 242_529,
            activationEpoch: 242_551,
            exitEpoch: 242_556,
            withdrawableEpoch: 250_743
        });

        bytes32 expected =
            0xe4674dc5c27e7d3049fcd298745c00d3e314f03d33c877f64bf071d3b77eb942;
        bytes32 actual = SSZ.validatorHashTreeRoot(v);
        assertEq(actual, expected);
    }

    // Slot 7172576 validator.at(44444)
    function test_ValidatorRootActive() public view {
        SSZ.Validator memory v = SSZ.Validator({
            pubkey: hex"8fb78536e82bcec34e98fff85c907f0a8e6f4b1ccdbf1e8ace26b59eb5a06d16f34e50837f6c490e2ad6a255db8d543b",
            withdrawalCredentials: 0x0023b9d00bf66e7f8071208a85afde59b3148dea046ee3db5d79244880734881,
            effectiveBalance: 32_000_000_000,
            slashed: false,
            activationEligibilityEpoch: 2593,
            activationEpoch: 5890,
            exitEpoch: type(uint64).max,
            withdrawableEpoch: type(uint64).max
        });

        bytes32 expected =
            0x60fb91184416404ddfc62bef6df9e9a52c910751daddd47ea426aabaf19dfa09;
        bytes32 actual = SSZ.validatorHashTreeRoot(v);
        assertEq(actual, expected);
    }

    function test_BeaconBlockHeaderRoot() public view {
        // Can be obtained via /eth/v1/beacon/headers/{block_id}.
        SSZ.BeaconBlockHeader memory h = SSZ.BeaconBlockHeader({
            slot: 7_472_518,
            proposerIndex: 152_834,
            parentRoot: 0x4916af1ff31b06f1b27125d2d20cd26e123c425a4b34ebd414e5f0120537e78d,
            stateRoot: 0x76ca64f3732754bc02c7966271fb6356a9464fe5fce85be8e7abc403c8c7b56b,
            bodyRoot: 0x6d858c959f1c95f411dba526c4ae9ab8b2690f8b1e59ed1b79ad963ab798b01a
        });

        bytes32 expected =
            0x26631ee28ab4dd44a39c3756e03714d6a35a256560de5e2885caef9c3efd5516;
        bytes32 actual = SSZ.beaconHeaderHashTreeRoot(h);
        assertEq(actual, expected);
    }
}
