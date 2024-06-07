pragma solidity ^0.8.25;

import { Script } from "@forge-std/Script.sol";
import { MintableERC20 } from "../test/MintableERC20.sol";

contract DeployAndCallERC20 is Script {
    function run() public {
        address dropAddress = address(12);
        uint256 quantity = 50_000;

        vm.startBroadcast();
        MintableERC20 drop = new MintableERC20();

        for (uint256 i = 0; i < 1000; i++) {
            quantity += 50_000;
            drop.mint(dropAddress, quantity);
        }

        vm.stopBroadcast();
    }
}
