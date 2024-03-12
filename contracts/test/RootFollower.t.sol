// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { Errors } from "@src/eip4788/extensions/Errors.sol";
import { IRootFollower } from "@src/eip4788/extensions/IRootFollower.sol";
import { RootFollower } from "@src/eip4788/extensions/RootFollower.sol";
import { BeaconRootsContractBaseTest } from "./BeaconRootsContract.t.sol";

contract RootFollowerUser is RootFollower { }

contract RootFollowerTest is BeaconRootsContractBaseTest {
    RootFollowerUser internal rootFollower;

    function setUp() public override {
        // etch the BeaconRootsContract to the BEACON_ROOT_ADDRESS
        bytes memory beaconRootsContractBytecode = abi.encodePacked(
            hex"3373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500"
        );
        vm.etch(BEACON_ROOT_ADDRESS, beaconRootsContractBytecode);
        // set the initial storage of the BEACON_ROOT_ADDRESS
        // setBeaconRoots(0, TIMESTAMP, HISTORY_BUFFER_LENGTH);
        rootFollower = new RootFollowerUser();
    }

    function test_GetAndIncrementBlock() public {
        // Setup the block number
        uint256 blockNum = 1;
        vm.roll(blockNum);

        // Get the next actionable block and assert it
        assertEq(1, rootFollower.getNextActionableBlock());

        // Increment the block
        vm.expectEmit(address(rootFollower));
        emit IRootFollower.AdvancedBlock(0);
        rootFollower.incrementBlock();

        // Check the last actioned block
        assertEq(1, rootFollower.getLastActionedBlock());
        assertEq(2, rootFollower.getNextActionableBlock());
    }

    function test_resetCount() public {
        vm.roll(1);
        // Cannot reset block to a future block
        vm.expectRevert(Errors.BlockDoesNotExist.selector);
        rootFollower.resetCount(100);

        // Set block to 10000
        vm.roll(10_000);

        // Cannot reset to a block not in the buffer
        vm.expectRevert(Errors.BlockNotInBuffer.selector);
        rootFollower.resetCount(2);

        // Should successfully reset the block count
        rootFollower.resetCount(10_000 - 8191);
        assertEq(10_000 - 8191, rootFollower.getNextActionableBlock());
    }
}
