// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title TokenizationErrors
 * @notice Custom error types for tokenization precompile operations
 * @dev Provides gas-efficient, descriptive error types for better error handling.
 *      These errors match common failure modes in the tokenization system.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/TokenizationErrors.sol";
 * 
 * function transfer(uint256 collectionId, address to, uint256 amount) external {
 *     if (collectionId == 0) {
 *         revert TokenizationErrors.InvalidCollectionId(collectionId);
 *     }
 *     // ... transfer logic
 * }
 * ```
 */
library TokenizationErrors {
    // ============ Collection Errors ============

    /// @notice Thrown when a collection is not found
    /// @param collectionId The collection ID that was not found
    error CollectionNotFound(uint256 collectionId);

    /// @notice Thrown when a collection ID is invalid (e.g., zero)
    /// @param collectionId The invalid collection ID
    error InvalidCollectionId(uint256 collectionId);

    /// @notice Thrown when attempting to access an archived collection
    /// @param collectionId The archived collection ID
    error CollectionArchived(uint256 collectionId);

    // ============ Balance Errors ============

    /// @notice Thrown when insufficient balance for a transfer
    /// @param collectionId The collection ID
    /// @param required The required amount
    /// @param available The available amount
    error InsufficientBalance(uint256 collectionId, uint256 required, uint256 available);

    /// @notice Thrown when balance query fails
    /// @param collectionId The collection ID
    /// @param address_ The address queried
    error BalanceQueryFailed(uint256 collectionId, address address_);

    // ============ Transfer Errors ============

    /// @notice Thrown when a transfer fails
    /// @param collectionId The collection ID
    /// @param reason The reason for failure
    error TransferFailed(uint256 collectionId, string reason);

    /// @notice Thrown when transfer is not allowed due to permissions
    /// @param collectionId The collection ID
    /// @param from The sender address
    /// @param to The recipient address
    error TransferNotAllowed(uint256 collectionId, address from, address to);

    /// @notice Thrown when transfer amount is zero
    /// @param collectionId The collection ID
    error TransferAmountZero(uint256 collectionId);

    // ============ Token ID Errors ============

    /// @notice Thrown when a token ID is invalid
    /// @param collectionId The collection ID
    /// @param tokenId The invalid token ID
    error InvalidTokenId(uint256 collectionId, uint256 tokenId);

    /// @notice Thrown when token ID range is invalid (start > end)
    /// @param collectionId The collection ID
    /// @param start The start value
    /// @param end The end value
    error InvalidTokenIdRange(uint256 collectionId, uint256 start, uint256 end);

    // ============ Approval Errors ============

    /// @notice Thrown when an approval is denied
    /// @param collectionId The collection ID
    /// @param approvalId The approval ID
    error ApprovalDenied(uint256 collectionId, string approvalId);

    /// @notice Thrown when an approval is not found
    /// @param collectionId The collection ID
    /// @param approvalId The approval ID
    error ApprovalNotFound(uint256 collectionId, string approvalId);

    /// @notice Thrown when approval criteria are not met
    /// @param collectionId The collection ID
    /// @param approvalId The approval ID
    /// @param reason The reason criteria were not met
    error ApprovalCriteriaNotMet(uint256 collectionId, string approvalId, string reason);

    // ============ Dynamic Store Errors ============

    /// @notice Thrown when a dynamic store is not found
    /// @param storeId The store ID that was not found
    error DynamicStoreNotFound(uint256 storeId);

    /// @notice Thrown when a dynamic store ID is invalid (e.g., zero)
    /// @param storeId The invalid store ID
    error InvalidDynamicStoreId(uint256 storeId);

    // ============ Address List Errors ============

    /// @notice Thrown when an address list is not found
    /// @param listId The list ID that was not found
    error AddressListNotFound(string listId);

    /// @notice Thrown when an address is not in a required list
    /// @param listId The list ID
    /// @param address_ The address that is not in the list
    error AddressNotInList(string listId, address address_);

    /// @notice Thrown when an address is in a forbidden list
    /// @param listId The list ID
    /// @param address_ The address that is in the forbidden list
    error AddressInForbiddenList(string listId, address address_);

    // ============ Permission Errors ============

    /// @notice Thrown when caller lacks required permission
    /// @param collectionId The collection ID
    /// @param permission The required permission
    error PermissionDenied(uint256 collectionId, string permission);

    /// @notice Thrown when caller is not the collection creator
    /// @param collectionId The collection ID
    /// @param caller The caller address
    error NotCollectionCreator(uint256 collectionId, address caller);

    /// @notice Thrown when caller is not the collection manager
    /// @param collectionId The collection ID
    /// @param caller The caller address
    error NotCollectionManager(uint256 collectionId, address caller);

    // ============ Validation Errors ============

    /// @notice Thrown when an address is invalid (e.g., zero address)
    /// @param address_ The invalid address
    error InvalidAddress(address address_);

    /// @notice Thrown when a string parameter is empty
    /// @param parameterName The name of the parameter
    error EmptyString(string parameterName);

    /// @notice Thrown when an array is empty when it should not be
    /// @param parameterName The name of the parameter
    error EmptyArray(string parameterName);

    /// @notice Thrown when array lengths do not match
    /// @param array1Name The name of the first array
    /// @param array2Name The name of the second array
    error ArrayLengthMismatch(string array1Name, string array2Name);

    // ============ Ownership Time Errors ============

    /// @notice Thrown when ownership time range is invalid
    /// @param collectionId The collection ID
    /// @param start The start time
    /// @param end The end time
    error InvalidOwnershipTimeRange(uint256 collectionId, uint256 start, uint256 end);

    /// @notice Thrown when ownership time has expired
    /// @param collectionId The collection ID
    /// @param expirationTime The expiration time
    error OwnershipTimeExpired(uint256 collectionId, uint256 expirationTime);

    // ============ Query Errors ============

    /// @notice Thrown when a query fails
    /// @param queryType The type of query that failed
    /// @param reason The reason for failure
    error QueryFailed(string queryType, string reason);

    // ============ Helper Functions ============

    /**
     * @notice Validate that a collection ID is not zero
     * @param collectionId The collection ID to validate
     */
    function requireValidCollectionId(uint256 collectionId) internal pure {
        if (collectionId == 0) {
            revert InvalidCollectionId(collectionId);
        }
    }

    /**
     * @notice Validate that an address is not zero
     * @param address_ The address to validate
     */
    function requireValidAddress(address address_) internal pure {
        if (address_ == address(0)) {
            revert InvalidAddress(address_);
        }
    }

    /**
     * @notice Validate that a string is not empty
     * @param str The string to validate
     * @param parameterName The name of the parameter for error reporting
     */
    function requireNonEmptyString(string memory str, string memory parameterName) internal pure {
        if (bytes(str).length == 0) {
            revert EmptyString(parameterName);
        }
    }

    /**
     * @notice Validate that an array is not empty
     * @param arr The array to validate
     * @param parameterName The name of the parameter for error reporting
     */
    function requireNonEmptyArray(uint256[] memory arr, string memory parameterName) internal pure {
        if (arr.length == 0) {
            revert EmptyArray(parameterName);
        }
    }

    /**
     * @notice Validate that two arrays have the same length
     * @param arr1 The first array
     * @param arr1Name The name of the first array
     * @param arr2 The second array
     * @param arr2Name The name of the second array
     */
    function requireSameLength(
        uint256[] memory arr1,
        string memory arr1Name,
        uint256[] memory arr2,
        string memory arr2Name
    ) internal pure {
        if (arr1.length != arr2.length) {
            revert ArrayLengthMismatch(arr1Name, arr2Name);
        }
    }

    /**
     * @notice Validate that a token ID range is valid (start <= end)
     * @param collectionId The collection ID
     * @param start The start value
     * @param end The end value
     */
    function requireValidTokenIdRange(
        uint256 collectionId,
        uint256 start,
        uint256 end
    ) internal pure {
        if (start > end) {
            revert InvalidTokenIdRange(collectionId, start, end);
        }
    }
}

