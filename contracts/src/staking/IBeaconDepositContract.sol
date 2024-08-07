// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IBeaconDepositContract
/// @author Berachain Team.
/// @dev This contract is used to create validator, deposit and withdraw stake from the Beaconchain.
interface IBeaconDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        EVENTS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Emitted when an operator address change is cancelled.
     * @param pubkey The public key of the validator.
     */
    event OperatorChangeCancelled(bytes pubkey);

    /**
     * @notice Emitted when an operator address change is queued.
     * @param pubkey The public key of the validator.
     * @param newOperator The new operator address.
     */
    event OperatorChangeQueued(bytes pubkey, address newOperator);

    /**
     * @notice Emitted when an operator is set for a given public key.
     * @param pubkey The public key of the validator.
     * @param newOperator The new operator address.
     * @param oldOperator The old operator address.
     */
    event OperatorSet(bytes pubkey, address newOperator, address oldOperator);

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

    /// @dev Error thrown when the deposit is not authorized.
    error UnauthorizedDeposit();

    error NotOperator();

    error NotNewOperator();

    error NotEnoughTimePassed();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        WRITES                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a deposit message to the Beaconchain.
     * @notice This will be used to create a new validator or to top up an existing one, increasing stake.
     * @param pubkey is the consensus public key of the validator. If subsequent deposit, its ignored.
     * @param credentials is the staking credentials of the validator. If this is the first deposit it is
     * validator operator public key, if subsequent deposit it is the depositor's public key.
     * @param signature is the signature used only on the first deposit.
     */
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        bytes calldata signature
    )
        external
        payable;

    /**
     * @notice Permits an address to use the deposit contract for N number of deposits.
     * @param depositor the address to allow deposits for.
     * @param number the number of deposits to allow.
     */
    function allowDeposit(address depositor, uint64 number) external;

    /**
     * @notice Request a change of operator address for a given public key.
     * @dev This will queue the change and require confirmation.
     * @dev Only the current operator can request a change.
     * @param pubkey The public key of the validator.
     * @param newOperator The new operator address.
     */
    function requestOperatorChange(
        bytes calldata pubkey,
        address newOperator
    )
        external;

    /**
     * @notice Cancel a pending operator address change request.
     * @dev Only the current operator can cancel a change.
     * @param pubkey The public key of the validator.
     */
    function cancelOperatorChange(bytes calldata pubkey) external;

    /**
     * @notice Confirm a pending operator address change request.
     * @dev Only the new operator can confirm the change.
     * @param pubkey The public key of the validator.
     */
    function acceptOperatorChange(bytes calldata pubkey) external;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        WRITES                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Get the operator address for a given public key.
     * @dev This is guaranteed to return a non-zero address if validator exists.
     * @param pubkey the public key to get the operator for.
     * @return the operator address.
     */
    function getOperator(bytes calldata pubkey)
        external
        view
        returns (address);
}
