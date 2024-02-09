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
    uint256[HISTORY_BUFFER_LENGTH] private _timestamps;

    BeaconRootsContract internal beaconRootsContract;
    bytes32 private beaconRoot;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public {
        beaconRootsContract = new BeaconRootsContract();
        setStorage(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
    }

    /// @dev Set the storage of the BeaconRootsContract by calling from the SYSTEM_ADDRESS.
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
            _timestamps[i % HISTORY_BUFFER_LENGTH] = timestamps[i];
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

    function callGet(uint256 timestamp) internal returns (bool success, bytes memory data) {
        (success, data) = address(beaconRootsContract).call(abi.encode(timestamp));
    }

    function validateBeaconRoots(uint256[] memory timestamps, bytes32[] memory beaconRoots) internal {
        // loop over the last `HISTORY_BUFFER_LENGTH` indices
        uint256 i = timestamps.length - 1;
        for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
            (bool success, bytes memory data) = callGet(timestamps[i]);
            assertTrue(success, "get: failed");
            assertEq(data.length, 32, "get: invalid length");
            assertEq(bytes32(data), beaconRoots[i], "get: invalid beacon root");
            if (i == 0) {
                break;
            }
            --i;
        }
    }

    /// @dev Test the timestamps, beacon roots, and coinbases are stored correctly in the circular buffers.
    function test_Set() public {
        testFuzz_Set(0, 0, HISTORY_BUFFER_LENGTH);
    }

    /// @dev Fuzzing test the timestamps, beacon roots, and coinbases are stored correctly in the circular buffers.
    function testFuzz_Set(uint64 startBlock, uint32 startTimestamp, uint256 length) public {
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 4);
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots, address[] memory coinbases) =
            setStorage(startBlock, startTimestamp, length);
        // loop over the last `HISTORY_BUFFER_LENGTH` indices
        uint256 i = length - 1;
        for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
            uint256 blockIdx = (startBlock + i) % HISTORY_BUFFER_LENGTH;
            bytes32 data = vm.load(address(beaconRootsContract), bytes32(blockIdx));
            assertEq(uint256(data), timestamps[i], "set: invalid timestamp");
            data = vm.load(address(beaconRootsContract), bytes32(blockIdx + BEACON_ROOT_OFFSET));
            assertEq(data, beaconRoots[i], "set: invalid beacon root");
            data = vm.load(address(beaconRootsContract), bytes32(blockIdx + COINBASE_OFFSET));
            assertEq(uint256(data), uint160(coinbases[i]), "set: invalid coinbase");
            if (i == 0) {
                break;
            }
            --i;
        }
    }

    /// @dev Test the beacon root is retrieved correctly from the circular buffer.
    function test_Get() public {
        (bool success, bytes memory data) = callGet(block.timestamp);
        assertTrue(success, "get: failed");
        assertEq(data.length, 32, "get: invalid length");
        assertEq(bytes32(data), beaconRoot, "get: invalid beacon root");
    }

    function test_InvalidCalldataLength() public {
        bytes memory data = abi.encode(block.timestamp);
        (bool success,) = address(beaconRootsContract).call(bytes.concat(data, data));
        assertFalse(success, "get: found invalid calldata length");
    }

    function test_GetOutOfRangeTimestamp() public {
        (bool success,) = callGet(TIMESTAMP - 1);
        assertFalse(success, "get: found out of range timestamp");
    }

    function test_GetInvalidTimestamp(uint256 timestamp) public {
        timestamp = _bound(timestamp, _timestamps[0] + 1, _timestamps[HISTORY_BUFFER_LENGTH - 1] - 1);
        for (uint256 i; i < HISTORY_BUFFER_LENGTH; ++i) {
            vm.assume(timestamp != _timestamps[i]);
        }
        (bool success,) = callGet(timestamp);
        assertFalse(success, "get: found invalid timestamp");
    }

    /// @dev Fuzzing test the beacon root is retrieved correctly from the circular buffer.
    function testFuzz_Get(uint64 startBlock, uint32 startTimestamp, uint256 length) public {
        vm.assume(startTimestamp > 0);
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 4);
        beaconRootsContract = new BeaconRootsContract();
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots,) = setStorage(startBlock, startTimestamp, length);
        // The timestamp encoded in the calldata may be in the past.
        // But the block number and timestamp in the EVM must be the latest.
        validateBeaconRoots(timestamps, beaconRoots);
    }

    /// @dev Fuzzing test the coinbase is retrieved correctly from the circular buffer.
    function testFuzz_GetCoinbase(uint64 startBlock, uint32 startTimestamp, uint256 length) public {
        vm.assume(startTimestamp > 0);
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 4);
        beaconRootsContract = new BeaconRootsContract();
        (uint256[] memory blockNumbers,,, address[] memory coinbases) = setStorage(startBlock, startTimestamp, length);
        // loop over the last `HISTORY_BUFFER_LENGTH` indices
        uint256 i = length - 1;
        for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
            (bool success, bytes memory data) =
                address(beaconRootsContract).call(abi.encodeWithSelector(GET_COINBASE_SELECTOR, blockNumbers[i]));
            assertTrue(success, "getCoinbase: failed");
            assertEq(uint256(bytes32(data)), uint160(coinbases[i]), "get: invalid coinbase");
            if (i == 0) {
                break;
            }
            --i;
        }
    }

    /// @dev Fuzzing test the beacon root is retrieved correctly from a partially initialized buffer.
    function testFuzz_PartiallyInitializedBuffer(uint256 length) public {
        beaconRootsContract = new BeaconRootsContract();
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH - 1);
        // The block number starts from 1.
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots,) = setStorage(1, TIMESTAMP, length);
        validateBeaconRoots(timestamps, beaconRoots);
    }
}
