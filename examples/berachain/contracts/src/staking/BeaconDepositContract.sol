// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import { IBeaconDepositContract } from "./IBeaconDepositContract.sol";

import { Ownable } from "@solady/src/auth/Ownable.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 */
contract BeaconDepositContract is IBeaconDepositContract, Ownable {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The minimum amount of stake that can be deposited to prevent dust.
    /// @dev This is 32 ether in Gwei since our deposit contract denominates in Gwei. 32e9 * 1e9 = 32e18.
    uint64 private constant MIN_DEPOSIT_AMOUNT_IN_GWEI = 32e9;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 private constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 private constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 private constant CREDENTIALS_LENGTH = 32;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev A flag to check if the contract has been initialized.
    bool private initialized = false;

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    uint64 public depositCount;

    /// @dev depositAuth is a mapping of number of deposits an authorized address can make.
    mapping(address => uint64) private depositAuth;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Initializes the owner of the contract.
    function initializeOwner() external {
        require(!initialized, "Already initialized");
        _initializeOwner(0x8a73D1380345942F1cb32541F1b19C40D8e6C94B);
        initialized = true;
    }

    /// @inheritdoc IBeaconDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        uint64 amount,
        bytes calldata signature
    )
        external
        payable
    {
        if (depositAuth[msg.sender] == 0) {
            revert UnauthorizedDeposit();
        }

        if (pubkey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (credentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (signature.length != SIGNATURE_LENGTH) {
            revert InvalidSignatureLength();
        }

        uint64 amountInGwei = _deposit(amount);

        if (amountInGwei < MIN_DEPOSIT_AMOUNT_IN_GWEI) {
            revert InsufficientDeposit();
        }

        --depositAuth[msg.sender];

        unchecked {
            // slither-disable-next-line reentrancy-benign,reentrancy-events
            emit Deposit(
                pubkey, credentials, amountInGwei, signature, depositCount++
            );
        }
    }

    /// @inheritdoc IBeaconDepositContract
    function allowDeposit(
        address depositor,
        uint64 number
    )
        external
        onlyOwner
    {
        depositAuth[depositor] = number;
    }

    /// @dev Validates the deposit amount and sends the native asset to the zero address.
    function _deposit(uint64) internal virtual returns (uint64) {
        if (msg.value % 1 gwei != 0) {
            revert DepositNotMultipleOfGwei();
        }

        uint256 amountInGwei = msg.value / 1 gwei;
        if (amountInGwei > type(uint64).max) {
            revert DepositValueTooHigh();
        }

        _safeTransferETH(address(0), msg.value);

        return uint64(amountInGwei);
    }

    /**
     * @notice Safely transfers ETH to the given address.
     * @dev From the Solady library.
     * @param to The address to transfer the ETH to.
     * @param amount The amount of ETH to transfer.
     */
    function _safeTransferETH(address to, uint256 amount) private {
        /// @solidity memory-safe-assembly
        assembly {
            if iszero(
                call(gas(), to, amount, codesize(), 0x00, codesize(), 0x00)
            ) {
                mstore(0x00, 0xb12d13eb) // `ETHTransferFailed()`.
                revert(0x1c, 0x04)
            }
        }
    }
}
