// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract MockLogs {
    event LogUint256(uint256 value);
    event LogBytes32(bytes32 value);
    event LogAddress(address value);
    event LogBool(bool value);
    event LogString(string value);

    event LogUint256Array(uint256[] values);
    event LogBytes32Array(bytes32[] values);
    event LogAddressArray(address[] values);
    event LogBoolArray(bool[] values);
    event LogStringArray(string[] values);

    event Log2Uint256(uint256 value1, uint256 value2);
    event AnotherLog2Uint256(uint256 value1, uint256 value2);

    struct MockStruct {
        uint256 uint256Value;
        bytes32 bytes32Value;
        address addressValue;
        bool boolValue;
        string stringValue;
    }
    event LogMockStruct(MockStruct mockStruct);
}
