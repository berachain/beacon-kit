// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Script, console} from "forge-std/Script.sol";
import {SimplePoLDistributor} from "../src/brip0004/MockPoL.sol";
import {ValidatorRegistry} from "../src/brip0004/MockValidatorRegistry.sol";

contract GetBytecodeScript is Script {
    function run() external {
        // Deploy ValidatorRegistry first
        ValidatorRegistry validatorRegistry = new ValidatorRegistry();
        bytes memory registryRuntimeCode = address(validatorRegistry).code;
        
        // Deploy SimplePoLDistributor
        SimplePoLDistributor polDistributor = new SimplePoLDistributor();
        bytes memory polRuntimeCode = address(polDistributor).code;
        
        console.log("=== ValidatorRegistry Bytecode ===");
        console.log("Contract Address:", address(validatorRegistry));
        console.log("Expected Genesis Address: 0x4200000000000000000000000000000000000043");
        console.log("Runtime Code Length:", registryRuntimeCode.length);
        console.log("Runtime Bytecode (hex):");
        console.logBytes(registryRuntimeCode);
        
        console.log("\n=== SimplePoLDistributor Bytecode ===");
        console.log("Contract Address:", address(polDistributor));
        console.log("Expected Genesis Address: 0x4200000000000000000000000000000000000042");
        console.log("Runtime Code Length:", polRuntimeCode.length);
        console.log("System Address:", 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE);
        console.log("Runtime Bytecode (hex):");
        console.logBytes(polRuntimeCode);
        
        // Log function selectors for verification
        bytes4 distributeForSelector = SimplePoLDistributor.distributeFor.selector;
        bytes4 recordActivitySelector = ValidatorRegistry.recordValidatorActivity.selector;
        console.log("\n=== Function Selectors ===");
        console.log("distributeFor selector:", vm.toString(distributeForSelector));
        console.log("recordValidatorActivity selector:", vm.toString(recordActivitySelector));
    }
}