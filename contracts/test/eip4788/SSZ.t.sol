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
        bytes memory pubkey = hex"a1c5a1e39fbe3eb23ce0a32aede615e813a8d4c36d1334bb9b36e4fb11289e0f5cce007e35ec74174de59ffa42c7f833";
        bytes32 expected =
            0x8d94e5939d9e60417b14cbe52b90124eeb99bee4a0336219fe04d693e64944d9;
        bytes32 actual = SSZ.validatorPubkeyHashTreeRoot(pubkey);
        assertEq(actual, expected);
    }
}
