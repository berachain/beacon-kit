// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import { Errors } from "./Errors.sol";
import { IRootFollower } from "./IRootFollower.sol";
import { RootFollower } from "./RootFollower.sol";
import { BeaconRootsContract } from "../BeaconRootsContract.sol";
import { BeaconRootsContractBaseTest } from "../BeaconRootsContract.t.sol";
import "forge-std/Test.sol";

contract RootFollowerUser is RootFollower { }

contract RootFollowerTest is BeaconRootsContractBaseTest {
    RootFollowerUser internal rootFollower;

    function setUp() public override {
        super.setUp();
        rootFollower = new RootFollowerUser();
    }

    function test_GetAndIncrementBlock() public {
        // Setup the block number
        uint256 blockNum = 1;
        vm.roll(blockNum);

        (bool success, bytes memory result) =
            BEACON_ROOT_ADDRESS.call(abi.encodeWithSelector(GET_COINBASE_SELECTOR, blockNum));
        assertEq(success, true);
        address expected = abi.decode(result, (address));

        // Get the next actionable block and assert it
        assertEq(1, rootFollower.getNextActionableBlock());

        // Get the coinbase of block 1 and assert it
        address received = rootFollower.getCoinbase(blockNum);
        assertEq(expected, received);

        // Increment the block
        vm.expectEmit(address(rootFollower));
        emit IRootFollower.AdvancedBlock(0);
        rootFollower.incrementBlock();

        // Check the last actioned block
        assertEq(1, rootFollower.getLastActionedBlock());
        assertEq(2, rootFollower.getNextActionableBlock());
    }

    function test_OutOfBuffer() public {
        // should succeed
        rootFollower.getCoinbase(1);
        vm.roll(500);

        assertEq(500 - 256, rootFollower.getNextActionableBlock());

        // Incrementing the block should fail now because out of buffer
        vm.expectRevert(Errors.AttemptedToIncrementOutOfBuffer.selector);
        rootFollower.incrementBlock();

        // Getting an out of buffer coinbase should result in a revert
        vm.expectRevert(Errors.BlockNotInBuffer.selector);
        rootFollower.getCoinbase(1);
    }

    function test_resetCount() public {
        vm.roll(1);
        // Cannot reset block to a future block
        vm.expectRevert(Errors.BlockDoesNotExist.selector);
        rootFollower.resetCount(100);

        // Set block to 500
        vm.roll(500);

        // Cannot reset to a block not in the buffer
        vm.expectRevert(Errors.BlockNotInBuffer.selector);
        rootFollower.resetCount(2);

        // Should successfully reset the block count
        rootFollower.resetCount(500 - 256);
        assertEq(500 - 256, rootFollower.getNextActionableBlock());
    }
}
