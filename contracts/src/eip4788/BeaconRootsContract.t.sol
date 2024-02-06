// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../../lib/solady/test/utils/SoladyTest.sol";
import "../eip4788/BeaconRootsContract.sol";

/// @title BeaconRootsContractTest
/// @dev This contract is used for testing the BeaconRootsContract.
contract BeaconRootsContractTest is SoladyTest {
    BeaconRootsContract beaconRootsContract;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public {
        beaconRootsContract = new BeaconRootsContract();
    }
}
