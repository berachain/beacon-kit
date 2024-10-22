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
    event Deposit(bytes pubkey, bytes credentials, uint64 amount, bytes signature, uint64 index);

    /**
     * @notice Emitted when the operator of a validator is updated.
     * @param pubkey The pubkey of the validator.
     * @param newOperator The new operator address.
     * @param previousOperator The previous operator address.
     */
    event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator);

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

    /// @dev Error thrown when the input operator is zero address on the first deposit.
    error ZeroOperatorOnFirstDeposit();

    /// @dev Error thrown when the caller is not the current operator.
    error NotCurrentOperator();

    /// @dev Error thrown when the address is zero address.
    error ZeroAddress();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            VIEWS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Get the operator address for a given pubkey.
     * @dev Returns zero address if the pubkey is not registered.
     * @param pubkey The pubkey of the validator.
     * @return The operator address for the given pubkey.
     */
    function getOperator(bytes calldata pubkey) external view returns (address);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a deposit message to the Beaconchain.
     * @notice This will be used to create a new validator or to top up an existing one, increasing stake.
     * @param pubkey is the consensus public key of the validator. If subsequent deposit, its ignored.
     * @param credentials is the staking credentials of the validator. If this is the first deposit it is
     * validator operator public key, if subsequent deposit it is the depositor's public key.
     * @param amount is the amount of stake native/ERC20 token to be deposited, in Gwei.
     * @param signature is the signature used only on the first deposit.
     * @param operator is the address of the operator.
     * @dev emits the Deposit event upon successful deposit.
     */
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        uint64 amount,
        bytes calldata signature,
        address operator
    ) external payable;
}
