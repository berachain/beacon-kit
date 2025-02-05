// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract BeaconRootMock {
    bytes32 public beaconRoot;

    constructor() payable { }

    function setBeaconRoot(bytes32 _beaconRoot) public {
        beaconRoot = _beaconRoot;
    }

    fallback(bytes calldata) external payable returns (bytes memory) {
        return abi.encode(beaconRoot);
    }
}
