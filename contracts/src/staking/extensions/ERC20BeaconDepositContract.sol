pragma solidity 0.8.25;

import { BeaconDepositContract } from "../BeaconDepositContract.sol";
import { IStakeERC20 } from "./IStakeERC20.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits, withdrawals, and redirections of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 */
contract ERC20BeaconDepositContract is BeaconDepositContract {
    /// @notice The ERC20 token contract that is used for staking.
    /// TODO: Change this to the actual ERC20 token contract.
    address public ERC20_DEPOSIT_ASSET;

    /// @dev Validates the deposit amount and burns the staking asset from the sender.
    function _deposit(uint64 amount) internal override returns (uint64) {
        // burn the staking asset from the sender, converting the gwei to wei.
        IStakeERC20(ERC20_DEPOSIT_ASSET).burn(msg.sender, uint256(amount) * 1e9);
        return amount;
    }
}
