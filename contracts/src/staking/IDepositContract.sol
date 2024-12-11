// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title IDepositContract
/// @author Berachain Team.
interface IDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           ERRORS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    // Signature: 0xe8966d7a
    error NotEnoughTime();
    // Signature: 0xd92e233d
    error ZeroAddress();
    // Signature: 0x7c214f04
    error NotOperator();

    /// @dev Error thrown when the deposit amount is too small, to prevent dust deposits.
    // Signature: 0x0e1eddda
    error InsufficientDeposit();

    /// @dev Error thrown when the deposit amount is not a multiple of Gwei.
    // Signature: 0x40567b38
    error DepositNotMultipleOfGwei();

    /// @dev Error thrown when the deposit amount is too high, since it is a uint64.
    // Signature: 0x2aa66734
    error DepositValueTooHigh();

    /// @dev Error thrown when the public key length is not 48 bytes.
    // Signature: 0x9f106472
    error InvalidPubKeyLength();

    /// @dev Error thrown when the withdrawal credentials length is not 32 bytes.
    // Signature: 0xb39bca16
    error InvalidCredentialsLength();

    /// @dev Error thrown when the signature length is not 96 bytes.
    // Signature: 0x4be6321b
    error InvalidSignatureLength();

    /// @dev Error thrown when the input operator is zero address on the first deposit.
    // Signature: 0x51969a7a
    error ZeroOperatorOnFirstDeposit();

    /// @dev Error thrown when the operator is already set and caller passed non-zero operator.
    // Signature: 0xc4142b41
    error OperatorAlreadySet();

    /// @dev Error thrown when the caller is not the current operator.
    // Signature: 0x819a0d0b
    error NotNewOperator();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           EVENTS                           */
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

    /**
     * @notice Emitted when the operator change of a validator is queued.
     * @param pubkey The pubkey of the validator.
     * @param queuedOperator The new queued operator address.
     * @param currentOperator The current operator address.
     * @param queuedTimestamp The timestamp when the change was queued.
     */
    event OperatorChangeQueued(
        bytes indexed pubkey,
        address queuedOperator,
        address currentOperator,
        uint256 queuedTimestamp
    );

    /**
     * @notice Emitted when the operator change of a validator is cancelled.
     * @param pubkey The pubkey of the validator.
     */
    event OperatorChangeCancelled(bytes indexed pubkey);

    /**
     * @notice Emitted when the operator of a validator is updated.
     * @param pubkey The pubkey of the validator.
     * @param newOperator The new operator address.
     * @param previousOperator The previous operator address.
     */
    event OperatorUpdated(
        bytes indexed pubkey, address newOperator, address previousOperator
    );

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            VIEWS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Get the operator address for a given pubkey.
     * @dev Returns zero address if the pubkey is not registered.
     * @param pubkey The pubkey of the validator.
     * @return The operator address for the given pubkey.
     */
    function getOperator(bytes calldata pubkey)
        external
        view
        returns (address);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a deposit message to the Beaconchain.
     * @notice This will be used to create a new validator or to top up an existing one, increasing stake.
     * @param pubkey is the consensus public key of the validator. If subsequent deposit, its ignored.
     * @param credentials is the staking credentials of the validator. If this is the first deposit it is
     * validator operator public key, if subsequent deposit it is the depositor's public key.
     * @param signature is the signature used only on the first deposit.
     * @param operator is the address of the operator.
     * @dev emits the Deposit event upon successful deposit.
     * @dev Reverts if the operator is already set and caller passed non-zero operator.
     */
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        bytes calldata signature,
        address operator
    )
        external
        payable;

    /**
     * @notice Request to change the operator of a validator.
     * @dev Only the current operator can request a change.
     * @param pubkey The pubkey of the validator.
     * @param newOperator The new operator address.
     */
    function requestOperatorChange(
        bytes calldata pubkey,
        address newOperator
    )
        external;

    /**
     * @notice Cancel the operator change of a validator.
     * @dev Only the current operator can cancel the change.
     * @param pubkey The pubkey of the validator.
     */
    function cancelOperatorChange(bytes calldata pubkey) external;

    /**
     * @notice Accept the operator change of a validator.
     * @dev Only the new operator can accept the change.
     * @dev Reverts if the queue delay has not passed.
     * @param pubkey The pubkey of the validator.
     */
    function acceptOperatorChange(bytes calldata pubkey) external;
}
