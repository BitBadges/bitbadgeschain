// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/// @title PrecompileTransferTestContract
/// @notice Test contract for transfer and query precompile methods
/// @dev Split from PrecompileTestContract to stay under EVM size limits
contract PrecompileTransferTestContract {
    ITokenizationPrecompile constant precompile =
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    // ============ Events ============

    event TransferExecuted(
        uint256 indexed collectionId,
        address indexed recipient,
        bool success
    );

    // ============ Transfer Methods ============

    /// @notice Wrapper for transferTokens - simplified to avoid stack too deep
    function testTransfer(
        uint256 collectionId,
        address recipient,
        uint256 amount,
        uint256 tokenIdStart,
        uint256 tokenIdEnd
    ) external returns (bool) {
        // Build simple JSON for a single transfer
        string memory transferJson = string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","transfers":[{"toAddresses":["',
                TokenizationJSONHelpers.addressToString(recipient),
                '"],"balances":[{"amount":"',
                TokenizationJSONHelpers.uintToString(amount),
                '","tokenIds":[{"start":"',
                TokenizationJSONHelpers.uintToString(tokenIdStart),
                '","end":"',
                TokenizationJSONHelpers.uintToString(tokenIdEnd),
                '"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}]}]}'
            )
        );

        bool success = precompile.transferTokens(transferJson);
        emit TransferExecuted(collectionId, recipient, success);
        return success;
    }

    // ============ Query Methods (View) ============

    /// @notice Wrapper for getCollection
    function testGetCollection(
        uint256 collectionId
    ) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getCollectionJSON(collectionId);
        return precompile.getCollection(queryJson);
    }

    /// @notice Wrapper for getBalance
    function testGetBalance(
        uint256 collectionId,
        address address_
    ) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getBalanceJSON(collectionId, address_);
        return precompile.getBalance(queryJson);
    }

    /// @notice Wrapper for getBalanceAmount - simplified
    function testGetBalanceAmount(
        uint256 collectionId,
        address address_,
        uint256 tokenIdStart,
        uint256 tokenIdEnd
    ) external view returns (uint256) {
        string memory balanceJson = string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","userAddress":"',
                TokenizationJSONHelpers.addressToString(address_),
                '","tokenIds":[{"start":"',
                TokenizationJSONHelpers.uintToString(tokenIdStart),
                '","end":"',
                TokenizationJSONHelpers.uintToString(tokenIdEnd),
                '"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}'
            )
        );
        return precompile.getBalanceAmount(balanceJson);
    }

    /// @notice Wrapper for getTotalSupply - simplified
    function testGetTotalSupply(
        uint256 collectionId,
        uint256 tokenIdStart,
        uint256 tokenIdEnd
    ) external view returns (uint256) {
        string memory supplyJson = string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","tokenIds":[{"start":"',
                TokenizationJSONHelpers.uintToString(tokenIdStart),
                '","end":"',
                TokenizationJSONHelpers.uintToString(tokenIdEnd),
                '"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}'
            )
        );
        return precompile.getTotalSupply(supplyJson);
    }

    /// @notice Wrapper for getAddressList
    function testGetAddressList(
        string calldata listId
    ) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getAddressListJSON(listId);
        return precompile.getAddressList(queryJson);
    }
}
