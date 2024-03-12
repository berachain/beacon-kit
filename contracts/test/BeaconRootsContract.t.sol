// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { console2 } from "@forge-std/console2.sol";
import { StdChains } from "@forge-std/StdChains.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { FixedPointMathLib } from "@solady/src/utils/FixedPointMathLib.sol";

/// @title BeaconRootsContractBaseTest
/// @dev This contract is a baseplate for tests that depend on the
/// BeaconRootsContract.
abstract contract BeaconRootsContractBaseTest is SoladyTest {
    uint256 internal constant HISTORY_BUFFER_LENGTH = 8191;
    uint256 internal constant BEACON_ROOT_OFFSET = HISTORY_BUFFER_LENGTH;
    uint256 internal constant COINBASE_OFFSET =
        BEACON_ROOT_OFFSET + HISTORY_BUFFER_LENGTH;
    uint256 internal constant BLOCK_MAPPING_OFFSET =
        COINBASE_OFFSET + HISTORY_BUFFER_LENGTH;
    address internal constant SYSTEM_ADDRESS =
        0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
    uint256 internal constant BLOCK_INTERVAL = 5;
    uint256 internal constant TIMESTAMP = 1_707_425_462;

    address internal BEACON_ROOT_ADDRESS =
        0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;
    uint256[HISTORY_BUFFER_LENGTH] internal _timestamps;

    bytes32 internal lastBeaconRoot;
    uint256 internal snapshot;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public virtual {
        bytes memory beaconRootsContractBytecode = abi.encodePacked(
            hex"3373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500"
        );
        vm.etch(BEACON_ROOT_ADDRESS, beaconRootsContractBytecode);
        BEACON_ROOT_ADDRESS =
            address(0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02);
        // take a snapshot of the clean state
        snapshot = vm.snapshot();
        // set the initial storage of the BEACON_ROOT_ADDRESS
        // setBeaconRoots(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
    }
}
