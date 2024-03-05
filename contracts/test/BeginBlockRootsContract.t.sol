// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@src/eip4788/extensions/BeginBlockRootsContract.sol";
import "./BeaconRootsContract.t.sol";

contract MockExternalContract {
    uint256 public timesCalled;

    function succeed() external returns (bool) {
        timesCalled = timesCalled + 1;
        return true;
    }

    function fail() external pure {
        revert("fail");
    }
}

/// @title BeginBlockRootsTest
/// @dev we inherit from the BeaconRootsContractTest to reuse the tests,
/// since we should have the same behavior as its just an extension of the
/// BeaconRootsContract with additional fallback functionality.
contract BeginBlockRootsTest is BeaconRootsContractTest {
    /// @dev The ADMIN address is the only address that can add or remove
    /// BeginBlockers.
    address internal ADMIN = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev Actions that can set a new BeginBlocker.
    bytes32 private constant SET = keccak256("SET");

    /// @dev Action that can remove BeginBlockers from the array.
    bytes32 private constant REMOVE = keccak256("REMOVE");

    /// @dev The selector for "getBeginBlockers(uint256)".
    bytes4 private constant GET_BEGIN_BLOCKERS_SELECTOR =
        bytes4(keccak256("getBeginBlockers(uint256)"));

    /// @dev The MockExternalContract.
    MockExternalContract internal mockExternalContract;

    /// @dev The BeginBlockRootsContract.
    BeginBlockRootsContract internal beginBlockRootsContract;

    /// @dev override the setUp so that we can use our version of the
    /// BeginBlockRootsContract.
    function setUp() public override {
        BEACON_ROOT_ADDRESS = address(new BeginBlockRootsContract());
        // Deploy the MockExternalContract.
        mockExternalContract = new MockExternalContract();
        // Set the ADMIN address at the correct slot.
        bytes32 value = bytes32(uint256(uint160(ADMIN)));
        bytes32 slot = bytes32(uint256(24_574));
        vm.store(BEACON_ROOT_ADDRESS, slot, value);
        // take a snapshot of the clean state.
        snapshot = vm.snapshot();
        // set the initial storage of the BEACON_ROOT_ADDRESS
        setBeaconRoots(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
        // set the BeginBlockRootsContract.
        beginBlockRootsContract = BeginBlockRootsContract(BEACON_ROOT_ADDRESS);
    }

    /// @dev Ensure that there is no selector collision affecting `set`.
    function test_SetBeaconRootShouldNeverFail() public {
        bytes32 beaconRoot = bytes32(_random());
        bytes4 selector = GET_BEGIN_BLOCKERS_SELECTOR;
        assembly {
            beaconRoot := or(selector, shr(32, beaconRoot))
        }
        _setBeaconRoot(beaconRoot);
    }

    /// @dev Test that we can set a new BeginBlocker contract and it will be
    /// set.
    function test_SetBeginBlocker() public {
        // Set the BeginBlocker as the ADMIN address.
        vm.prank(ADMIN);
        (bool success,) = BEACON_ROOT_ADDRESS.call(
            _createCRUD(
                0,
                SET,
                address(mockExternalContract),
                mockExternalContract.succeed.selector,
                address(0)
            )
        );
        assertTrue(success, "BeginBlockRootsTest: failed to set BeginBlocker");

        // Check that the BeginBlocker is set.
        BeginBlockRootsContract.BeginBlocker memory beginBlocker =
            _getBeginBlockers(0);
        assertEq(
            beginBlocker.contractAddress,
            address(mockExternalContract),
            "BeginBlockRootsTest: BeginBlocker contract address not set"
        );
        assertEq(
            beginBlocker.selector,
            mockExternalContract.succeed.selector,
            "BeginBlockRootsTest: BeginBlocker selector not set"
        );
    }

    /// @dev Should fail to get a BeginBlocker that is out of bounds.
    function test_FailGetBeginBlockersOutOfBounds() public {
        test_SetBeginBlocker();
        (bool success,) = BEACON_ROOT_ADDRESS.staticcall(
            abi.encodePacked(GET_BEGIN_BLOCKERS_SELECTOR, uint256(1))
        );
        assertFalse(success, "BeginBlockRootsTest: getBeginBlockers succeeded");
    }

    /// @dev Test that we can remove a BeginBlocker contract.
    function test_RemoveBeginBlocker() public {
        test_SuccessCallsMulti();

        // Remove the BeginBlocker.
        vm.prank(ADMIN);
        (bool success,) = BEACON_ROOT_ADDRESS.call(
            _createCRUD(0, REMOVE, address(0), bytes4(0), address(8))
        );
        assertTrue(
            success, "BeginBlockRootsTest: failed to remove BeginBlocker"
        );
    }

    function test_SuccessCallsMulti() public {
        // Set the BeginBlocker 10 times as the ADMIN address.
        vm.startPrank(ADMIN);
        bytes memory crudMsg = _createCRUD(
            0, // test insert at index 0
            SET,
            address(mockExternalContract),
            mockExternalContract.succeed.selector,
            address(0)
        );
        for (uint256 i = 0; i < 10; ++i) {
            (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
            assertTrue(
                success, "BeginBlockRootsTest: failed to set BeginBlocker"
            );
        }
        vm.stopPrank();

        // Call the BeginBlocker and should run in a loop, 10 times incrementing
        // the timesCalled.
        _setBeaconRoot(bytes32(_random()));

        assertEq(
            10,
            mockExternalContract.timesCalled(),
            "BeginBlockRootsTest: BeginBlocker not called 10 times"
        );
    }

    function test_FailCall() public {
        bytes memory crudMsg = _createCRUD(
            0,
            SET,
            address(mockExternalContract),
            mockExternalContract.fail.selector,
            address(0)
        );
        vm.prank(ADMIN);
        (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
        assertTrue(success, "BeginBlockRootsTest: failed to set BeginBlocker");

        // Call the BeginBlocker and this should revert, not updating the
        // timesCalled, but setting the beacon root.
        _setBeaconRoot(bytes32(_random()));
        assertEq(
            0,
            mockExternalContract.timesCalled(),
            "BeginBlockRootsTest: BeginBlocker called"
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

    function _setBeaconRoot(bytes32 beaconRoot) internal {
        vm.prank(SYSTEM_ADDRESS);
        (bool success,) = BEACON_ROOT_ADDRESS.call(abi.encode(beaconRoot));
        assertTrue(success, "set failed");
    }

    function _getBeginBlockers(uint256 index)
        internal
        returns (BeginBlockRootsContract.BeginBlocker memory)
    {
        (bool success, bytes memory returnData) = BEACON_ROOT_ADDRESS.staticcall(
            abi.encodePacked(GET_BEGIN_BLOCKERS_SELECTOR, index)
        );
        assertTrue(success, "BeginBlockRootsTest: failed to get BeginBlocker");
        return abi.decode(returnData, (BeginBlockRootsContract.BeginBlocker));
    }
}
