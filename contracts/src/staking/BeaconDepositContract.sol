// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

pragma solidity 0.8.25;

import { SSZ } from "../eip4788/SSZ.sol";
import { IBeaconDepositContract } from "./IBeaconDepositContract.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 * @dev It is based on the Ethereum 2.0 specification.
 * @dev From https://github.com/ethereum/consensus-specs/blob/dev/solidity_deposit_contract/deposit_contract.sol
 */
contract BeaconDepositContract is IBeaconDepositContract {
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

    /// @dev The length of the withdrawal credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 private constant WITHDRAWAL_CREDENTIALS_LENGTH = 32;

    /// @dev The maximum depth of the deposit contract's Merkle tree.
    uint256 constant DEPOSIT_CONTRACT_TREE_DEPTH = 32;

    /// @dev The maximum number of deposits that can be made to the contract.
    uint256 constant MAX_DEPOSIT_COUNT = 2 ** DEPOSIT_CONTRACT_TREE_DEPTH - 1;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    uint64 public depositCount;

    /// @dev branch stores the Merkle tree branch for the deposit contract.
    bytes32[DEPOSIT_CONTRACT_TREE_DEPTH] branch;

    /// @dev zeroHashes stores the zero hashes for the deposit contract.
    bytes32[DEPOSIT_CONTRACT_TREE_DEPTH] zeroHashes;

    constructor() {
        // Compute hashes in empty sparse Merkle tree
        unchecked {
            uint256 depthMinusOne = DEPOSIT_CONTRACT_TREE_DEPTH - 1;
            for (uint256 height; height < depthMinusOne; ++height) {
                zeroHashes[height + 1] = sha256(
                    abi.encodePacked(zeroHashes[height], zeroHashes[height])
                );
            }
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata withdrawal_credentials,
        bytes calldata signature,
        bytes32 deposit_data_root
    )
        external
        payable
    {
        if (pubkey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (withdrawal_credentials.length != WITHDRAWAL_CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (signature.length != SIGNATURE_LENGTH) {
            revert InvalidSignatureLength();
        }

        uint64 amountInGwei = _deposit();

        emit Deposit(
            pubkey,
            withdrawal_credentials,
            amountInGwei,
            signature,
            depositCount
        );

        bytes32 amount = SSZ.toLittleEndian(amountInGwei);
        // Compute deposit data root (`DepositData` hash tree root)
        bytes32 pubkey_root = sha256(abi.encodePacked(pubkey, bytes16(0)));
        bytes32 signature_root = sha256(
            abi.encodePacked(
                sha256(signature[:64]),
                sha256(abi.encodePacked(signature[64:], bytes32(0)))
            )
        );
        bytes32 node = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(pubkey_root, withdrawal_credentials)),
                sha256(abi.encodePacked(amount, signature_root))
            )
        );

        if (node != deposit_data_root) {
            revert InvalidDepositDataRoot();
        }

        if (depositCount >= MAX_DEPOSIT_COUNT) {
            revert MerkleTreeFull();
        }

        // Add deposit data root to Merkle tree (update a single `branch` node)
        unchecked {
            ++depositCount;
        }
        uint256 size = depositCount;
        unchecked {
            for (uint256 height; height < DEPOSIT_CONTRACT_TREE_DEPTH; ++height)
            {
                if ((size & 1) == 1) {
                    branch[height] = node;
                    return;
                }
                node = sha256(abi.encodePacked(branch[height], node));
                size >>= 1;
            }
        }
        // As the loop should always end prematurely with the `return` statement,
        // this code should be unreachable. We assert `false` just to be safe.
        assert(false);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            READS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconDepositContract
    function getDepositRoot() external view override returns (bytes32) {
        //slither-disable-next-line uninitialized-local
        bytes32 node;
        uint256 size = depositCount;
        unchecked {
            for (uint256 height; height < DEPOSIT_CONTRACT_TREE_DEPTH; ++height)
            {
                if ((size & 1) == 1) {
                    node = sha256(abi.encodePacked(branch[height], node));
                } else {
                    node = sha256(abi.encodePacked(node, zeroHashes[height]));
                }
                size >>= 1;
            }
        }
        return sha256(abi.encodePacked(node, SSZ.toLittleEndian(depositCount)));
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           INTERNAL                         */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Validates the deposit amount and sends the native asset to the zero address.
    function _deposit() internal virtual returns (uint64) {
        if (msg.value % 1 gwei != 0) {
            revert DepositNotMultipleOfGwei();
        }

        uint256 amountInGwei = msg.value / 1 gwei;
        if (amountInGwei < MIN_DEPOSIT_AMOUNT_IN_GWEI) {
            revert DepositValueTooLow();
        }

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
