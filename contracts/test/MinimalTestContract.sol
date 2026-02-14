// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationJSONHelpers.sol";

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
        // Convert UintRange arrays to JSON
        uint256[] memory tokenIdStarts = new uint256[](tokenIds.length);
        uint256[] memory tokenIdEnds = new uint256[](tokenIds.length);
        for (uint256 i = 0; i < tokenIds.length; i++) {
            tokenIdStarts[i] = tokenIds[i].start;
            tokenIdEnds[i] = tokenIds[i].end;
        }
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeArrayToJson(tokenIdStarts, tokenIdEnds);
        
        uint256[] memory ownershipStarts = new uint256[](ownershipTimes.length);
        uint256[] memory ownershipEnds = new uint256[](ownershipTimes.length);
        for (uint256 i = 0; i < ownershipTimes.length; i++) {
            ownershipStarts[i] = ownershipTimes[i].start;
            ownershipEnds[i] = ownershipTimes[i].end;
        }
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeArrayToJson(ownershipStarts, ownershipEnds);
        
        // Convert recipients array
        address[] memory recipientsArray = new address[](recipients.length);
        for (uint256 i = 0; i < recipients.length; i++) {
            recipientsArray[i] = recipients[i];
        }
        
        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            recipientsArray,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        
        bool success = precompile.transferTokens(transferJson);
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
        string memory queryJson = TokenizationJSONHelpers.getBalanceJSON(collectionId, address_);
        return precompile.getBalance(queryJson);
    }
}

