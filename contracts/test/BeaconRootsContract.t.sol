// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { console2 } from "@forge-std/console2.sol";
import { StdChains } from "@forge-std/StdChains.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { FixedPointMathLib } from "@solady/src/utils/FixedPointMathLib.sol";
import { BeaconRootsContract } from "@src/eip4788/BeaconRootsContract.sol";

/// @title BeaconRootsContractBaseTest
/// @dev This contract is a baseplate for tests that depend on the
/// BeaconRootsContract.
contract BeaconRootsContractBaseTest is SoladyTest {
    uint256 internal constant HISTORY_BUFFER_LENGTH = 8191;
    uint256 internal constant BEACON_ROOT_OFFSET = HISTORY_BUFFER_LENGTH;
    uint256 internal constant COINBASE_OFFSET =
        BEACON_ROOT_OFFSET + HISTORY_BUFFER_LENGTH;
    uint256 internal constant BLOCK_MAPPING_OFFSET =
        COINBASE_OFFSET + HISTORY_BUFFER_LENGTH;
    address internal constant SYSTEM_ADDRESS =
        0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
    address internal constant BEACON_ROOT_ADDRESS =
        0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;
    bytes4 internal constant GET_COINBASE_SELECTOR =
        bytes4(keccak256("getCoinbase(uint256)"));
    uint256 internal constant BLOCK_INTERVAL = 5;
    uint256 internal constant TIMESTAMP = 1_707_425_462;

    uint256[HISTORY_BUFFER_LENGTH] internal _timestamps;

    bytes32 internal lastBeaconRoot;
    uint256 internal snapshot;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public virtual {
        // etch the BeaconRootsContract to the BEACON_ROOT_ADDRESS
        vm.etch(
            BEACON_ROOT_ADDRESS, vm.getDeployedCode("BeaconRootsContract.sol")
        );
        // take a snapshot of the clean state
        snapshot = vm.snapshot();
        // set the initial storage of the BEACON_ROOT_ADDRESS
        setBeaconRoots(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
    }

    /// @dev Set the storage of the BeaconRootsContract by calling from the
    /// SYSTEM_ADDRESS.
    function setBeaconRoots(
        uint256 startBlock,
        uint256 startTimestamp,
        uint256 length
    )
        internal
        returns (
            uint256[] memory blockNumbers,
            uint256[] memory timestamps,
            bytes32[] memory beaconRoots,
            address[] memory coinbases
        )
    {
        // create a pseudo random number from seeds
        uint256 rando = uint256(bytes32(abi.encode(startBlock, startTimestamp)));
        blockNumbers = new uint256[](length);
        timestamps = new uint256[](length);
        beaconRoots = new bytes32[](length);
        coinbases = new address[](length);
        vm.startPrank(SYSTEM_ADDRESS);
        unchecked {
            for (uint256 i; i < length; ++i) {
                blockNumbers[i] = startBlock + i;
                rando = pseudoRandom(rando);
                timestamps[i] = startTimestamp + i * BLOCK_INTERVAL
                // a random number between 1 and BLOCK_INTERVAL such that the
                // timestamp is ever increasing
                + FixedPointMathLib.min(
                    1,
                    FixedPointMathLib.fullMulDiv(
                        rando, BLOCK_INTERVAL, type(uint256).max
                    )
                );
                _timestamps[i % HISTORY_BUFFER_LENGTH] = timestamps[i];
                rando = pseudoRandom(rando);
                beaconRoots[i] = bytes32(rando);
                rando = pseudoRandom(rando);
                coinbases[i] = address(uint160(rando));
                vm.roll(blockNumbers[i]);
                vm.warp(timestamps[i]);
                vm.coinbase(coinbases[i]);
                (bool success,) =
                    BEACON_ROOT_ADDRESS.call(abi.encode(beaconRoots[i]));
                assertTrue(success, "setStorage: set failed");
            }
            lastBeaconRoot = beaconRoots[length - 1];
        }
        vm.stopPrank();
    }

    /// @dev Call the BeaconRootsContract to retrieve the beacon root for a
    /// given timestamp.
    function callGet(uint256 _timestamp)
        internal
        view
        returns (bool success, bytes32 beaconRoot)
    {
        // BEACON_ROOT_ADDRESS.staticcall(abi.encode(timestamp))
        assembly ("memory-safe") {
            mstore(0, _timestamp)
            // `staticcall` is evaluated before `returndatasize`
            success :=
                and(
                    eq(returndatasize(), 0x20),
                    staticcall(gas(), BEACON_ROOT_ADDRESS, 0, 0x20, 0, 0x20)
                )
            beaconRoot := mload(0)
        }
    }

    /// @dev Validate the beacon roots are stored correctly in the circular
    /// buffer.
    function validateBeaconRoots(
        uint256[] memory timestamps,
        bytes32[] memory beaconRoots
    )
        internal
    {
        unchecked {
            // loop over the last `HISTORY_BUFFER_LENGTH` indices
            uint256 i = timestamps.length - 1;
            for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
                (bool success, bytes32 beaconRoot) = callGet(timestamps[i]);
                assertTrue(success, "get: failed");
                assertEq(beaconRoot, beaconRoots[i], "get: invalid beacon root");
                if (i == 0) {
                    break;
                }
                --i;
            }
        }
    }

    /// @dev Generate a pseudo random number from a seed.
    function pseudoRandom(uint256 seed) internal pure returns (uint256) {
        assembly ("memory-safe") {
            mstore(0, seed)
            seed := keccak256(0, 0x20)
        }
        return seed;
    }
}

/// @title BeaconRootsContractTest
/// @dev This contract is used for testing the BeaconRootsContract.
contract BeaconRootsContractTest is BeaconRootsContractBaseTest {
    /// @dev Test the timestamps, beacon roots, and coinbases are stored
    /// correctly in the circular buffers.
    function test_Set() public {
        testFuzz_Set(0, 1, HISTORY_BUFFER_LENGTH);
    }

    /// @dev Fuzzing test the timestamps, beacon roots, and coinbases are stored
    /// correctly in the circular buffers.
    function testFuzz_Set(
        uint64 startBlock,
        uint32 startTimestamp,
        uint256 length
    )
        public
    {
        vm.assume(startTimestamp > 0);
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 2);
        (
            ,
            uint256[] memory timestamps,
            bytes32[] memory beaconRoots,
            address[] memory coinbases
        ) = setBeaconRoots(startBlock, startTimestamp, length);
        unchecked {
            // loop over the last `HISTORY_BUFFER_LENGTH` indices
            uint256 i = length - 1;
            for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
                uint256 blockNumber = startBlock + i;
                uint256 blockIdx = blockNumber % HISTORY_BUFFER_LENGTH;
                bytes32 data = vm.load(BEACON_ROOT_ADDRESS, bytes32(blockIdx));
                assertEq(uint256(data), timestamps[i], "set: invalid timestamp");
                data = vm.load(
                    BEACON_ROOT_ADDRESS,
                    keccak256(abi.encode(timestamps[i], BLOCK_MAPPING_OFFSET))
                );
                assertEq(
                    uint256(data), blockNumber, "set: invalid block number"
                );
                data = vm.load(
                    BEACON_ROOT_ADDRESS, bytes32(blockIdx + BEACON_ROOT_OFFSET)
                );
                assertEq(data, beaconRoots[i], "set: invalid beacon root");
                data = vm.load(
                    BEACON_ROOT_ADDRESS, bytes32(blockIdx + COINBASE_OFFSET)
                );
                assertEq(
                    uint256(data),
                    uint160(coinbases[i]),
                    "set: invalid coinbase"
                );
                if (i == 0) {
                    break;
                }
                --i;
            }
        }
    }

    /// @dev Test the beacon root is retrieved correctly from the circular
    /// buffer.
    function test_Get() public {
        (bool success, bytes32 beaconRoot) = callGet(block.timestamp);
        assertTrue(success, "get: failed");
        assertEq(beaconRoot, lastBeaconRoot, "get: invalid beacon root");
    }

    /// @dev Should fail if the calldata length is invalid.
    function test_InvalidCalldataLength() public {
        bytes memory data = abi.encode(block.timestamp);
        (bool success,) =
            BEACON_ROOT_ADDRESS.staticcall(bytes.concat(data, data));
        assertFalse(success, "get: found invalid calldata length");
    }

    /// @dev Should fail if the timestamp is out of range.
    function test_GetOutOfRangeTimestamp() public {
        (bool success,) = callGet(TIMESTAMP - 1);
        assertFalse(success, "get: found out of range timestamp");
    }

    /// @dev Should fail if the timestamp is not in the circular buffer.
    function testFuzz_GetInvalidTimestamp(uint256 timestamp) public {
        timestamp = _bound(
            timestamp,
            _timestamps[0] + 1,
            _timestamps[HISTORY_BUFFER_LENGTH - 1] - 1
        );
        bool loop = true;
        // find a timestamp that is not in the circular buffer
        while (loop) {
            loop = false;
            for (uint256 i; i < HISTORY_BUFFER_LENGTH; ++i) {
                if (timestamp == _timestamps[i]) {
                    // if the timestamp is found in the circular buffer, try
                    // another one
                    loop = true;
                    timestamp = _bound(
                        _random(),
                        _timestamps[0] + 1,
                        _timestamps[HISTORY_BUFFER_LENGTH - 1] - 1
                    );
                    break;
                }
            }
        }
        (bool success,) = callGet(timestamp);
        assertFalse(success, "get: found invalid timestamp");
    }

    /// @dev Fuzzing test the beacon root is retrieved correctly from the
    /// circular buffer.
    function testFuzz_Get(
        uint64 startBlock,
        uint32 startTimestamp,
        uint256 length
    )
        public
    {
        vm.assume(startTimestamp > 0);
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 2);
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots,) =
            setBeaconRoots(startBlock, startTimestamp, length);
        // The timestamp encoded in the calldata may be in the past.
        // But the block number and timestamp in the EVM must be the latest.
        validateBeaconRoots(timestamps, beaconRoots);
    }

    /// @dev Fuzzing test the coinbase is retrieved correctly from the circular
    /// buffer.
    function testFuzz_GetCoinbase(
        uint64 startBlock,
        uint32 startTimestamp,
        uint256 length
    )
        public
    {
        vm.assume(startTimestamp > 0);
        // may wrap around the circular buffer
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH * 2);
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        (uint256[] memory blockNumbers,,, address[] memory coinbases) =
            setBeaconRoots(startBlock, startTimestamp, length);
        unchecked {
            // loop over the last `HISTORY_BUFFER_LENGTH` indices
            uint256 i = length - 1;
            for (uint256 j; j < HISTORY_BUFFER_LENGTH; ++j) {
                (bool success, bytes memory data) = BEACON_ROOT_ADDRESS
                    .staticcall(
                    abi.encodeWithSelector(
                        GET_COINBASE_SELECTOR, blockNumbers[i]
                    )
                );
                assertTrue(success, "getCoinbase: failed");
                assertEq(
                    uint256(bytes32(data)),
                    uint160(coinbases[i]),
                    "get: invalid coinbase"
                );
                if (i == 0) {
                    break;
                }
                --i;
            }
        }
    }

    /// @dev Fuzzing test the beacon root is retrieved correctly from a
    /// partially initialized buffer.
    function testFuzz_PartiallyInitializedBuffer(uint256 length) public {
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        length = _bound(length, 1, HISTORY_BUFFER_LENGTH - 1);
        // The block number starts from 1.
        (, uint256[] memory timestamps, bytes32[] memory beaconRoots,) =
            setBeaconRoots(1, TIMESTAMP, length);
        validateBeaconRoots(timestamps, beaconRoots);
    }
}

/// @title BeaconRootsContractForkTest
/// @dev This contract is used for testing the BeaconRootsContract on a fork.
contract BeaconRootsContractForkTest is
    BeaconRootsContractBaseTest,
    StdChains
{
    function setUp() public override {
        vm.createSelectFork("sepolia");
        assertGt(block.number, 1, "something went wrong");
    }

    /// @dev Test the `get` function in BeaconRootsContract on Sepolia testnet.
    function testFork_Get() public {
        console2.log("block number", block.number);
        console2.log("block timestamp", block.timestamp);
        for (uint256 i; i < 10; ++i) {
            (bool success, bytes32 beaconRoot) = callGet(block.timestamp);
            assertTrue(success, "get: failed");
            uint256 timestampIdx = block.timestamp % HISTORY_BUFFER_LENGTH;
            bytes32 timestamp =
                vm.load(BEACON_ROOT_ADDRESS, bytes32(timestampIdx));
            assertEq(
                block.timestamp, uint256(timestamp), "get: invalid timestamp"
            );
            bytes32 _beaconRoot = vm.load(
                BEACON_ROOT_ADDRESS, bytes32(timestampIdx + BEACON_ROOT_OFFSET)
            );
            assertEq(beaconRoot, _beaconRoot, "get: invalid beacon root");
            // roll back to the previous block
            vm.rollFork(block.number - 1);
        }
    }

    /// @dev Test the `set` function in BeaconRootsContract on Sepolia testnet.
    function testFork_Set() public {
        uint256 timestamp = _random();
        vm.warp(timestamp);
        bytes32 beaconRoot = bytes32(_random());
        vm.prank(SYSTEM_ADDRESS);
        (bool success,) = BEACON_ROOT_ADDRESS.call(abi.encode(beaconRoot));
        assertTrue(success, "set failed");
        (bool success2, bytes32 beaconRoot2) = callGet(timestamp);
        assertTrue(success2, "get: failed");
        assertEq(beaconRoot, beaconRoot2, "invalid beacon root");
    }
}
