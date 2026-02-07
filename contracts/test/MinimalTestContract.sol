// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";

/// @title MinimalTestContract
/// @notice Minimal test contract for testing precompile integration
/// @dev This is a test-only contract with minimal functionality to stay under EVM size limits
contract MinimalTestContract {
    ITokenizationPrecompile constant precompile = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    event TransferExecuted(
        uint256 indexed collectionId,
        address indexed recipient,
        bool success
    );
    
    /// @notice Wrapper for transferTokens - the most essential method for testing
    function testTransfer(
        uint256 collectionId,
        address[] calldata recipients,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool) {
        bool success = precompile.transferTokens(
            collectionId,
            recipients,
            amount,
            tokenIds,
            ownershipTimes
        );
        if (recipients.length > 0) {
            emit TransferExecuted(collectionId, recipients[0], success);
        }
        return success;
    }
    
    /// @notice Simple wrapper for getBalance - useful for testing queries
    function testGetBalance(
        uint256 collectionId,
        address address_
    ) external view returns (bytes memory) {
        return precompile.getBalance(collectionId, address_);
    }
}

