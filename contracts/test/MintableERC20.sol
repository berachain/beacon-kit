pragma solidity ^0.8.25;

import { ERC20 } from "@solady/src/tokens/ERC20.sol";

contract MintableERC20 is ERC20 {
    constructor() ERC20() { }

    event Mint(address indexed to, uint256 amount);

    function name() public view virtual override returns (string memory) {
        return "Token";
    }

    function symbol() public view virtual override returns (string memory) {
        return "TK";
    }

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
        emit Mint(to, amount);
    }
}
