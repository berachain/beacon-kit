// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../../lib/solady/test/utils/SoladyTest.sol";
import "../../lib/solady/src/utils/FixedPointMathLib.sol";
import "./BeaconRootsContract.t.sol";
import "./extensions/BeginBlockRootsContract.sol";

contract MockExternalContract {
    function succeed() external pure returns (bool) {
        return true;
    }

    function fail() external pure {
        revert("fail");
    }
}

/// @title BeginBlockRootsTest
/// @notice we inherit from the BeaconRootsContractTest to reuse the tests, since we should have
/// the same behavior
// as its just an extension of the BeaconRootsContract with additional fallback functionality.
contract BeginBlockRootsTest is BeaconRootsContractTest {
    /// @dev The ADMIN address is the only address that can add or remove BeginBlockers.
    address internal ADMIN = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev Actions that can set a new BeginBlocker.
    bytes32 private constant SET = keccak256("SET");

    /// @dev The MockExternalContract.
    MockExternalContract internal mockExternalContract;

    /// @dev The BeginBlockRootsContract.
    BeginBlockRootsContract internal beginBlockRootsContract;

    /// @dev override the setUp so that we can use our version of the BeginBlockRootsContract.
    function setUp() public override {
        // etch the BeginBlockRootsContract to the BEACON_ROOT_ADDRESS.
        vm.etch(BEACON_ROOT_ADDRESS, vm.getDeployedCode("BeginBlockRootsContract.sol"));
        // Deploy the MockExternalContract.
        mockExternalContract = new MockExternalContract();
        // take a snapshot of the clean state.
        snapshot = vm.snapshot();
        // set the initial storage of the BEACON_ROOT_ADDRESS
        setStorage(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
        // set the BeginBlockRootsContract.
        beginBlockRootsContract = BeginBlockRootsContract(BEACON_ROOT_ADDRESS);
    }

    /// @dev Test that we can set a new BeginBlocker contract and it will be set.
    function test_SetBeginBlocker() public {
        bytes memory crudMsg = _createCRUD(
            0,
            SET,
            address(mockExternalContract),
            mockExternalContract.succeed.selector,
            address(0)
        );

        // Set the BeginBlocker as the ADMIN address.
        vm.prank(ADMIN);
        (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
        assertTrue(success, "BeginBlockRootsTest: failed to set BeginBlocker");

        // Check that the BeginBlocker is set.
        (address contractAddress, bytes4 selector) = beginBlockRootsContract.beginBlockers(0);
        assertEq(
            contractAddress,
            address(mockExternalContract),
            "BeginBlockRootsTest: contractAddress is not set"
        );
        assertEq(
            selector,
            mockExternalContract.succeed.selector,
            "BeginBlockRootsTest: selector is not set"
        );
    }

    /// @dev Create a BeginBlockCRUD message.
    function _createCRUD(
        uint256 i,
        bytes32 action,
        address contractAddress,
        bytes4 selector,
        address admin
    )
        internal
        pure
        returns (bytes memory)
    {
        return abi.encode(i, action, contractAddress, selector, admin);
    }
}
