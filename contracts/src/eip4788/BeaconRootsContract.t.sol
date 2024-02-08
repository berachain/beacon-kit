// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../../lib/solady/test/utils/SoladyTest.sol";
import "../../lib/solady/src/utils/FixedPointMathLib.sol";
import "../eip4788/BeaconRootsContract.sol";

/// @title BeaconRootsContractTest
/// @dev This contract is used for testing the BeaconRootsContract.
contract BeaconRootsContractTest is SoladyTest {
    uint256 private constant HISTORY_BUFFER_LENGTH = 256;
    uint256 private constant BEACON_ROOT_OFFSET = HISTORY_BUFFER_LENGTH;
    uint256 private constant COINBASE_OFFSET = BEACON_ROOT_OFFSET + HISTORY_BUFFER_LENGTH;
    address private constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;
    uint256 private constant BLOCK_INTERVAL = 5;
    uint256 private constant TIMESTAMP = 1707425462;

    BeaconRootsContract internal beaconRootsContract;
    bytes32 private beaconRoot;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public {
        beaconRootsContract = new BeaconRootsContract();
        setStorage(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
    }

    function setStorage(uint256 startBlock, uint256 startTimestamp, uint256 length)
        internal
        returns (
            uint256[] memory blockNumbers,
            uint256[] memory timestamps,
            bytes32[] memory beaconRoots,
            address[] memory coinbases
        )
    {
        blockNumbers = new uint256[](length);
        timestamps = new uint256[](length);
        beaconRoots = new bytes32[](length);
        coinbases = new address[](length);
        vm.startPrank(SYSTEM_ADDRESS);
        for (uint256 i; i < length; ++i) {
            blockNumbers[i] = startBlock + i;
            timestamps[i] = startTimestamp + i * BLOCK_INTERVAL
            // a random number between 1 and BLOCK_INTERVAL such that the timestamp is ever increasing
            + FixedPointMathLib.min(1, FixedPointMathLib.fullMulDiv(_random(), BLOCK_INTERVAL, type(uint256).max));
            beaconRoots[i] = bytes32(_random());
            coinbases[i] = _randomNonZeroAddress();
            vm.roll(blockNumbers[i]);
            vm.warp(timestamps[i]);
            vm.coinbase(coinbases[i]);
            (bool success,) = address(beaconRootsContract).call(abi.encode(beaconRoots[i]));
            assertTrue(success, "setStorage: set failed");
        }
        beaconRoot = beaconRoots[length - 1];
        vm.stopPrank();
    }

    function test_Set() public {
        testFuzz_Set(0, 0);
    }

    /// @dev Test the timestamps, beacon roots, and coinbases are stored correctly in the circular buffers.
    function testFuzz_Set(uint64 startBlock, uint32 startTimestamp) public {
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots, address[] memory coinbases) =
            setStorage(startBlock, startTimestamp, HISTORY_BUFFER_LENGTH);
        for (uint256 i; i < HISTORY_BUFFER_LENGTH; ++i) {
            uint256 blockIdx = (startBlock + i) % HISTORY_BUFFER_LENGTH;
            bytes32 data = vm.load(address(beaconRootsContract), bytes32(blockIdx));
            assertEq(uint256(data), timestamps[i], "set: invalid timestamp");
            data = vm.load(address(beaconRootsContract), bytes32(blockIdx + BEACON_ROOT_OFFSET));
            assertEq(data, beaconRoots[i], "set: invalid beacon root");
            data = vm.load(address(beaconRootsContract), bytes32(blockIdx + COINBASE_OFFSET));
            assertEq(uint256(data), uint160(coinbases[i]), "set: invalid coinbase");
        }
    }

    function test_Get() public {
        (bool success, bytes memory data) = address(beaconRootsContract).call(abi.encode(block.timestamp));
        assertTrue(success, "get: failed");
        assertEq(data.length, 32, "get: invalid length");
        assertEq(bytes32(data), beaconRoot, "get: invalid beacon root");
    }

    function testFuzz_Get(uint64 startBlock, uint32 startTimestamp) public {
        vm.assume(startTimestamp > 0);
        (uint256[] memory blockNumbers, uint256[] memory timestamps, bytes32[] memory beaconRoots,) =
            setStorage(startBlock, startTimestamp, HISTORY_BUFFER_LENGTH);
        // The timestamp encoded in the calldata may be in the past.
        // But the block number and timestamp in the EVM must be the latest.
        vm.roll(blockNumbers[HISTORY_BUFFER_LENGTH - 1]);
        vm.warp(timestamps[HISTORY_BUFFER_LENGTH - 1]);
        for (uint256 i; i < HISTORY_BUFFER_LENGTH; ++i) {
            (bool success, bytes memory data) = address(beaconRootsContract).call(abi.encode(timestamps[i]));
            assertTrue(success, "get: failed");
            assertEq(data.length, 32, "get: invalid length");
            assertEq(bytes32(data), beaconRoots[i], "get: invalid beacon root");
        }
    }

    function testFuzz_GetCoinbase(uint64 startBlock, uint32 startTimestamp) public {
        vm.assume(startTimestamp > 0);
        (uint256[] memory blockNumbers,,, address[] memory coinbases) =
            setStorage(startBlock, startTimestamp, HISTORY_BUFFER_LENGTH);
        for (uint256 i; i < HISTORY_BUFFER_LENGTH; ++i) {
            (bool success, bytes memory data) =
                address(beaconRootsContract).call(abi.encodeWithSelector(GET_COINBASE_SELECTOR, blockNumbers[i]));
            assertTrue(success, "getCoinbase: failed");
            assertEq(uint256(bytes32(data)), uint160(coinbases[i]), "get: invalid coinbase");
        }
    }
}
