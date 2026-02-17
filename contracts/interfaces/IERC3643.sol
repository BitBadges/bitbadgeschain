// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title IERC3643
 * @dev ERC-3643 security token interface for BitBadges
 *
 * Core functions map to BitBadges precompile with:
 * - Token ID 1 (fungible-like behavior)
 * - Full time range (1 to uint64.max)
 * - Auto-scan approval mode
 */
interface IERC3643 {
    // ============ Events ============

    /// @notice Emitted when tokens are transferred
    event Transfer(address indexed from, address indexed to, uint256 value);

    /// @notice Emitted when an identity is registered
    event IdentityRegistered(address indexed investor, bool accredited);

    /// @notice Emitted when an identity is removed
    event IdentityRemoved(address indexed investor);

    /// @notice Emitted when an address is frozen
    event AddressFrozen(address indexed investor);

    /// @notice Emitted when an address is unfrozen
    event AddressUnfrozen(address indexed investor);

    // ============ Core Functions ============

    /// @notice Transfer tokens to a recipient
    /// @param to Recipient address
    /// @param amount Amount to transfer
    /// @return success True if transfer succeeded
    function transfer(address to, uint256 amount) external returns (bool);

    /// @notice Get token balance of an address
    /// @param account Address to query
    /// @return balance Token balance
    function balanceOf(address account) external view returns (uint256);

    /// @notice Get total token supply
    /// @return supply Total supply
    function totalSupply() external view returns (uint256);
}

