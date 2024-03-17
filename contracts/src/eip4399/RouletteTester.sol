// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title RouletteTester
/// @dev This contract is used during integration testing of
///      EIP-4399 in BeaconKit.
///      DO NOT USE THIS FOR GENERATING RANDOMNESS IN PRODUCTION
//       FOR ANYTHING IMPORTANT.
/// @author https://eips.ethereum.org/EIPS/eip-4399
/// @author ocnc@berachain.com
contract RouletteTester {
    address private owner;
    uint256 public minimumBet = 1 ether;
    uint256 public houseEdgePercent = 2; // Represents the 2% house odds

    event BetResult(
        address indexed bettor,
        uint256 amountBet,
        bool indexed onRed,
        bool win,
        uint256 winningAmount
    );
    event Deposit(address indexed sender, uint256 amount);
    event Withdrawal(address indexed to, uint256 amount);

    constructor() {
        owner = msg.sender; // Set the contract creator as the owner
    }

    modifier onlyOwner() {
        require(
            msg.sender == owner,
            "This function is restricted to the contract's owner."
        );
        _;
    }

    modifier meetsMinimumBet(uint256 _amount) {
        require(
            _amount >= minimumBet, "Bet does not meet the minimum requirement."
        );
        _;
    }

    function betOnRed() external payable meetsMinimumBet(msg.value) {
        resolveBet(true);
    }

    function betOnBlack() external payable meetsMinimumBet(msg.value) {
        resolveBet(false);
    }

    function resolveBet(bool red) internal {
        uint256 randomNumber = uint256(block.prevrandao) % 100;
        uint256 betAmount = msg.value;

        if (randomNumber < houseEdgePercent) {
            emit BetResult(msg.sender, betAmount, red, false, 0);
        } else {
            bool isRedWin = (randomNumber % 2 == 0);
            if (red == isRedWin) {
                uint256 winningAmount = betAmount * 2;
                payable(msg.sender).transfer(winningAmount);
                emit BetResult(msg.sender, betAmount, red, true, winningAmount);
            } else {
                emit BetResult(msg.sender, betAmount, red, false, 0);
            }
        }
    }

    function depositFunds() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    function withdrawFunds(uint256 _amount) external onlyOwner {
        require(address(this).balance >= _amount, "Insufficient balance.");
        payable(owner).transfer(_amount);
        emit Withdrawal(owner, _amount);
    }

    function setMinimumBet(uint256 _newMinimumBet) external onlyOwner {
        minimumBet = _newMinimumBet;
    }

    function modifyHouseOdds(uint256 _newHouseEdgePercent) external onlyOwner {
        houseEdgePercent = _newHouseEdgePercent;
    }

    fallback() external payable {
        revert(
            "Please use the betOnRed or betOnBlack functions to place a bet."
        );
    }

    // Receive function to accept plain Ether transactions
    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }
}
