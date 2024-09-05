// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { Ownable } from "@solady/src/auth/Ownable.sol";
import { DepositContract } from "@src/staking/DepositContract.sol";

/// @notice A test contract that permissions deposits.
contract PermissionedDepositContract is DepositContract, Ownable {
    error UnauthorizedDeposit();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev depositAuth is a mapping of number of deposits an authorized
    /// address can make.
    mapping(address => uint64) public depositAuth;

    /// @dev Initializes the owner of the contract.
    constructor(address owner) {
        _initializeOwner(owner);
    }

    /// @dev Override to return true to make `_initializeOwner` prevent
    /// double-initialization.
    function _guardInitializeOwner()
        internal
        pure
        override
        returns (bool guard)
    {
        return true;
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
