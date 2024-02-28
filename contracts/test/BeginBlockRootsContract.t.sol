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

    /// @dev The MockExternalContract.
    MockExternalContract internal mockExternalContract;

    /// @dev The BeginBlockRootsContract.
    BeginBlockRootsContract internal beginBlockRootsContract;

    /// @dev override the setUp so that we can use our version of the
    /// BeginBlockRootsContract.
    function setUp() public override {
        // etch the BeginBlockRootsContract to the BEACON_ROOT_ADDRESS.
        vm.etch(
            BEACON_ROOT_ADDRESS,
            vm.getDeployedCode("BeginBlockRootsContract.sol")
        );
        // Deploy the MockExternalContract.
        mockExternalContract = new MockExternalContract();
        // Set the ADMIN address at the correct slot.
        bytes32 value = bytes32(uint256(uint160(ADMIN)));
        bytes32 slot = bytes32(uint256(24_574));
        vm.store(BEACON_ROOT_ADDRESS, slot, value);
        // take a snapshot of the clean state.
        snapshot = vm.snapshot();
        // set the initial storage of the BEACON_ROOT_ADDRESS
        setStorage(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
        // set the BeginBlockRootsContract.
        beginBlockRootsContract = BeginBlockRootsContract(BEACON_ROOT_ADDRESS);
    }

    /// @dev Test that we can set a new BeginBlocker contract and it will be
    /// set.
    function test_SimpleCRUD() public {
        bytes memory crudMsg = _createCRUD(
            0,
            SET,
            address(mockExternalContract),
            mockExternalContract.succeed.selector,
            address(0)
        );

        // Set the BeginBlocker as the ADMIN address.
        vm.startPrank(ADMIN);
        (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
        assertTrue(success, "BeginBlockRootsTest: failed to set BeginBlocker");

        // Check that the BeginBlocker is set.
        (address contractAddress, bytes4 selector) =
            beginBlockRootsContract.beginBlockers(0);
        assertEq(
            address(mockExternalContract),
            contractAddress,
            "BeginBlockRootsTest: BeginBlocker contract address not set"
        );
        assertEq(
            mockExternalContract.succeed.selector,
            selector,
            "BeginBlockRootsTest: BeginBlocker selector not set"
        );

        // Remove the BeginBlocker.
        crudMsg = _createCRUD(0, REMOVE, address(0), bytes4(0), address(8));
        (success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
        assertTrue(
            success, "BeginBlockRootsTest: failed to remove BeginBlocker"
        );

        vm.stopPrank();
    }

    function test_SuccessCallsMulti() public {
        // Set the BeginBlocker 10 times as the ADMIN address.
        vm.startPrank(ADMIN);
        for (uint256 i = 0; i < 10; i++) {
            bytes memory crudMsg = _createCRUD(
                i,
                SET,
                address(mockExternalContract),
                mockExternalContract.succeed.selector,
                address(0)
            );

            (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
            assertTrue(
                success, "BeginBlockRootsTest: failed to set BeginBlocker"
            );
        }
        vm.stopPrank();

        // Call the BeginBlocker and should run in a loop, 10 times incrementing
        // the timesCalled.
        _setRandomBeaconRoot();

        assertEq(
            10,
            mockExternalContract.timesCalled(),
            "BeginBlockRootsTest: BeginBlocker not called 10 times"
        );
    }

    function test_FailCall() public {
        vm.startPrank(ADMIN);
        bytes memory crudMsg = _createCRUD(
            0,
            SET,
            address(mockExternalContract),
            mockExternalContract.fail.selector,
            address(0)
        );

        (bool success,) = BEACON_ROOT_ADDRESS.call(crudMsg);
        assertTrue(success, "BeginBlockRootsTest: failed to set BeginBlocker");
        vm.stopPrank();

        // Call the BeginBlocker and this should revert, not updating the
        // timesCalled, but setting the beacon root.
        _setRandomBeaconRoot();
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

    function _setRandomBeaconRoot() internal returns (bytes32) {
        vm.prank(SYSTEM_ADDRESS);
        bytes32 random = bytes32(_random());
        (bool success,) = BEACON_ROOT_ADDRESS.call(abi.encode(random));
        require(
            success, "BeginBlockRootsTest: failed to set random beacon root"
        );
        return random;
    }
}
