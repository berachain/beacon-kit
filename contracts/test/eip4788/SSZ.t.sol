// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { Test } from "@forge-std/Test.sol";
import { SSZ } from "@src/eip4788/SSZ.sol";

contract SSZTest is Test {
    function test_uint64HashTreeRoot() public pure {
        uint64 w = type(uint64).max;
        bytes32 expected =
            bytes32(bytes.concat(hex"FFFFFFFFFFFFFFFF", bytes24(0)));
        bytes32 actual = SSZ.uint64HashTreeRoot(w);
        assertEq(actual, expected);

        uint64 v = 0x1234567890ABCDEF;
        expected = bytes32(bytes.concat(hex"EFCDAB9078563412", bytes24(0)));
        actual = SSZ.uint64HashTreeRoot(v);
        assertEq(actual, expected);

        uint64 z = 123_456_789;
        expected = bytes32(bytes.concat(hex"15CD5B0700000000", bytes24(0)));
        actual = SSZ.uint64HashTreeRoot(z);
        assertEq(actual, expected);
    }

    function test_addressHashTreeRoot() public pure {
        address a = address(0x1234567890abCDef000000000000000000000000);
        bytes32 expected =
            bytes32(bytes.concat(hex"1234567890ABCDEF", bytes24(0)));
        bytes32 actual = SSZ.addressHashTreeRoot(a);
        assertEq(actual, expected);

        address b = address(0x0102030405060708090a0B0c0d0e0f1011121314);
        expected =
            hex"0102030405060708090a0b0c0d0e0f1011121314000000000000000000000000";
        actual = SSZ.addressHashTreeRoot(b);
        assertEq(actual, expected);
    }

    function test_ValidatorPubkeyRoot() public view {
        bytes memory pubkey =
            hex"a1c5a1e39fbe3eb23ce0a32aede615e813a8d4c36d1334bb9b36e4fb11289e0f5cce007e35ec74174de59ffa42c7f833";
        bytes32 expected =
            0x8d94e5939d9e60417b14cbe52b90124eeb99bee4a0336219fe04d693e64944d9;
        bytes32 actual = SSZ.validatorPubkeyHashTreeRoot(pubkey);
        assertEq(actual, expected);
    }
}
