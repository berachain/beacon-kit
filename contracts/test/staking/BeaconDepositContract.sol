// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { Ownable } from "@solady/src/auth/Ownable.sol";
import { DepositContract } from "@src/staking/DepositContract.sol";

/// @notice A test contract that permissions deposits.
contract BeaconDepositContract is DepositContract, Ownable {
    error UnauthorizedDeposit();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev A flag to check if the contract has been initialized.
    bool private initialized = false;

    /// @dev depositAuth is a mapping of number of deposits an authorized address can make.
    mapping(address => uint64) private depositAuth;

    /// @dev Initializes the owner of the contract.
    function initializeOwner() external {
        require(!initialized, "Already initialized");
        _initializeOwner(0x8a73D1380345942F1cb32541F1b19C40D8e6C94B);
        initialized = true;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    function deposit(
        bytes calldata pubkey,
        bytes calldata withdrawal_credentials,
        uint64 amount,
        bytes calldata signature
    )
        public
        payable
        override
    {
        if (depositAuth[msg.sender] == 0) revert UnauthorizedDeposit();
        --depositAuth[msg.sender];

        super.deposit(pubkey, withdrawal_credentials, amount, signature);
    }

    function allowDeposit(
        address depositor,
        uint64 number
    )
        external
        onlyOwner
    {
        depositAuth[depositor] = number;
    }
}
