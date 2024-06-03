pragma solidity 0.8.25;

interface IStakeERC20 {
    /**
     * @notice Burns the specified amount of tokens from the specified account.
     * @param account The address of the account to burn from.
     * @param amount The amount of tokens to burn.
     */
    function burn(address account, uint256 amount) external;
}
