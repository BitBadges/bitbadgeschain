// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title IERC3643
 * @dev Minimal ERC-3643 interface for proof of concept
 */
interface IERC3643 {
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function totalSupply() external view returns (uint256);
}

