// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title GammErrors
 * @notice Custom error types and validation helpers for gamm precompile
 * @dev Provides custom errors and require functions for common validations
 */
library GammErrors {
    // ============ Custom Errors ============

    /// @notice Pool ID is invalid (zero)
    error InvalidPoolId(uint64 poolId);

    /// @notice Coin is invalid
    error InvalidCoin(string message);

    /// @notice Swap route is invalid
    error InvalidRoute(string message);

    /// @notice Affiliate is invalid
    error InvalidAffiliate(string message);

    /// @notice Pool not found
    error PoolNotFound(uint64 poolId);

    /// @notice Transaction failed
    error TransactionFailed(string message);

    /// @notice Query failed
    error QueryFailed(string message);

    /// @notice IBC transfer info is invalid
    error InvalidIBCTransferInfo(string message);

    // ============ Validation Helpers ============

    /**
     * @notice Require that pool ID is valid (non-zero)
     */
    function requireValidPoolId(uint64 poolId) internal pure {
        if (poolId == 0) {
            revert InvalidPoolId(poolId);
        }
    }

    /**
     * @notice Require that coin is valid (non-zero amount, non-empty denom)
     */
    function requireValidCoin(string memory denom, uint256 amount) internal pure {
        if (bytes(denom).length == 0) {
            revert InvalidCoin("denom cannot be empty");
        }
        if (amount == 0) {
            revert InvalidCoin("amount cannot be zero");
        }
    }

    /**
     * @notice Require that swap route is valid
     */
    function requireValidRoute(uint64 poolId, string memory tokenOutDenom) internal pure {
        if (poolId == 0) {
            revert InvalidRoute("poolId cannot be zero");
        }
        if (bytes(tokenOutDenom).length == 0) {
            revert InvalidRoute("tokenOutDenom cannot be empty");
        }
    }

    /**
     * @notice Require that affiliate is valid
     */
    function requireValidAffiliate(address address_, uint256 basisPointsFee) internal pure {
        if (address_ == address(0)) {
            revert InvalidAffiliate("address cannot be zero");
        }
        if (basisPointsFee > 10000) {
            revert InvalidAffiliate("basisPointsFee cannot exceed 10000 (100%)");
        }
    }

    /**
     * @notice Require that IBC transfer info is valid
     */
    function requireValidIBCTransferInfo(
        string memory sourceChannel,
        string memory receiver
    ) internal pure {
        if (bytes(sourceChannel).length == 0) {
            revert InvalidIBCTransferInfo("sourceChannel cannot be empty");
        }
        if (bytes(receiver).length == 0) {
            revert InvalidIBCTransferInfo("receiver cannot be empty");
        }
    }

    /**
     * @notice Require that transaction succeeded
     * @dev Use this after calling transaction methods to check return values
     */
    function requireTransactionSuccess(bool success, string memory message) internal pure {
        if (!success) {
            revert TransactionFailed(message);
        }
    }

    /**
     * @notice Require that query result is valid
     * @dev Use this after calling query methods to check return values
     */
    function requireQuerySuccess(bytes memory result, string memory message) internal pure {
        if (result.length == 0) {
            revert QueryFailed(message);
        }
    }
}















