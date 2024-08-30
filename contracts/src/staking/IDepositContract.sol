// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IDepositContract
/// @author Berachain Team.
/// @dev This contract is used to create validator, deposit and withdraw stake from the Beacon chain.
interface IDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        EVENTS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when a deposit is made, which could mean a new validator or a top up of an existing one.
     * @param pubkey the public key of the validator who is being deposited for if not a new validator.
     * @param credentials the public key of the operator if new validator or the depositor if top up.
     * @param amount the amount of stake being deposited, in Gwei.
     * @param signature the signature of the deposit message, only checked when creating a new validator.
     * @param index the index of the deposit.
     */
    event Deposit(
        bytes pubkey,
        bytes credentials,
        uint64 amount,
        bytes signature,
        uint64 index
    );

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ERRORS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Error thrown when the deposit amount is too small, to prevent dust deposits.
    error InsufficientDeposit();

    /// @dev Error thrown when the deposit amount is not a multiple of Gwei.
    error DepositNotMultipleOfGwei();

    /// @dev Error thrown when the deposit amount is too high, since it is a uint64.
    error DepositValueTooHigh();

    /// @dev Error thrown when the public key length is not 48 bytes.
    error InvalidPubKeyLength();

    /// @dev Error thrown when the withdrawal credentials length is not 32 bytes.
    error InvalidCredentialsLength();

    /// @dev Error thrown when the signature length is not 96 bytes.
    error InvalidSignatureLength();
}
