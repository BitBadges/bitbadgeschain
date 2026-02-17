// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/TokenizationTypes.sol";

/**
 * @title TokenizationHelpers
 * @notice Helper library for constructing and validating BitBadges tokenization types
 * @dev Provides utility functions for creating structs, validating inputs, and building default values
 */
library TokenizationHelpers {
    /**
     * @notice Maximum uint64 value - used for "forever" ownership times
     * @dev BitBadges uses uint64 internally for timestamps and IDs. Using uint256.max will cause errors.
     */
    uint64 public constant MAX_UINT64 = type(uint64).max;  // 18446744073709551615

    /**
     * @notice Creates a UintRange struct
     * @param start The start value (inclusive)
     * @param end The end value (inclusive)
     * @return range The constructed UintRange
     */
    function createUintRange(uint256 start, uint256 end) internal pure returns (TokenizationTypes.UintRange memory range) {
        require(start <= end, "TokenizationHelpers: start must be <= end");
        return TokenizationTypes.UintRange({start: start, end: end});
    }

    /**
     * @notice Creates an array of UintRange structs
     * @param starts Array of start values
     * @param ends Array of end values
     * @return ranges Array of constructed UintRange structs
     */
    function createUintRangeArray(
        uint256[] memory starts,
        uint256[] memory ends
    ) internal pure returns (TokenizationTypes.UintRange[] memory ranges) {
        require(starts.length == ends.length, "TokenizationHelpers: starts and ends arrays must have same length");
        ranges = new TokenizationTypes.UintRange[](starts.length);
        for (uint256 i = 0; i < starts.length; i++) {
            ranges[i] = createUintRange(starts[i], ends[i]);
        }
    }

    /**
     * @notice Creates a Balance struct
     * @param amount The amount of tokens
     * @param ownershipTimes Array of ownership time ranges
     * @param tokenIds Array of token ID ranges
     * @return balance The constructed Balance
     */
    function createBalance(
        uint256 amount,
        TokenizationTypes.UintRange[] memory ownershipTimes,
        TokenizationTypes.UintRange[] memory tokenIds
    ) internal pure returns (TokenizationTypes.Balance memory balance) {
        return TokenizationTypes.Balance({
            amount: amount,
            ownershipTimes: ownershipTimes,
            tokenIds: tokenIds
        });
    }

    /**
     * @notice Creates a CollectionMetadata struct
     * @param uri The URI for the collection metadata
     * @param customData Custom data string
     * @return metadata The constructed CollectionMetadata
     */
    function createCollectionMetadata(
        string memory uri,
        string memory customData
    ) internal pure returns (TokenizationTypes.CollectionMetadata memory metadata) {
        return TokenizationTypes.CollectionMetadata({
            uri: uri,
            customData: customData
        });
    }

    /**
     * @notice Creates a TokenMetadata struct
     * @param uri The URI for the token metadata
     * @param customData Custom data string
     * @param tokenIds Array of token ID ranges this metadata applies to
     * @return metadata The constructed TokenMetadata
     */
    function createTokenMetadata(
        string memory uri,
        string memory customData,
        TokenizationTypes.UintRange[] memory tokenIds
    ) internal pure returns (TokenizationTypes.TokenMetadata memory metadata) {
        return TokenizationTypes.TokenMetadata({
            uri: uri,
            customData: customData,
            tokenIds: tokenIds
        });
    }

    /**
     * @notice Creates an ActionPermission struct
     * @param permittedTimes Array of time ranges when action is permanently permitted
     * @param forbiddenTimes Array of time ranges when action is permanently forbidden
     * @return permission The constructed ActionPermission
     */
    function createActionPermission(
        TokenizationTypes.UintRange[] memory permittedTimes,
        TokenizationTypes.UintRange[] memory forbiddenTimes
    ) internal pure returns (TokenizationTypes.ActionPermission memory permission) {
        return TokenizationTypes.ActionPermission({
            permanentlyPermittedTimes: permittedTimes,
            permanentlyForbiddenTimes: forbiddenTimes
        });
    }

    /**
     * @notice Creates an empty UserBalanceStore with default values
     * @return store The constructed UserBalanceStore with empty arrays and false booleans
     */
    function createEmptyUserBalanceStore() internal pure returns (TokenizationTypes.UserBalanceStore memory store) {
        return TokenizationTypes.UserBalanceStore({
            balances: new TokenizationTypes.Balance[](0),
            outgoingApprovals: new TokenizationTypes.UserOutgoingApproval[](0),
            incomingApprovals: new TokenizationTypes.UserIncomingApproval[](0),
            autoApproveSelfInitiatedOutgoingTransfers: false,
            autoApproveSelfInitiatedIncomingTransfers: false,
            autoApproveAllIncomingTransfers: false,
            userPermissions: createEmptyUserPermissions()
        });
    }

    /**
     * @notice Creates an empty UserPermissions struct
     * @return permissions The constructed UserPermissions with empty arrays
     */
    function createEmptyUserPermissions() internal pure returns (TokenizationTypes.UserPermissions memory permissions) {
        return TokenizationTypes.UserPermissions({
            canUpdateOutgoingApprovals: new TokenizationTypes.UserOutgoingApprovalPermission[](0),
            canUpdateIncomingApprovals: new TokenizationTypes.UserIncomingApprovalPermission[](0),
            canUpdateAutoApproveSelfInitiatedOutgoingTransfers: new TokenizationTypes.ActionPermission[](0),
            canUpdateAutoApproveSelfInitiatedIncomingTransfers: new TokenizationTypes.ActionPermission[](0),
            canUpdateAutoApproveAllIncomingTransfers: new TokenizationTypes.ActionPermission[](0)
        });
    }

    /**
     * @notice Creates an empty CollectionPermissions struct
     * @return permissions The constructed CollectionPermissions with empty arrays
     */
    function createEmptyCollectionPermissions() internal pure returns (TokenizationTypes.CollectionPermissions memory permissions) {
        return TokenizationTypes.CollectionPermissions({
            canDeleteCollection: new TokenizationTypes.ActionPermission[](0),
            canArchiveCollection: new TokenizationTypes.ActionPermission[](0),
            canUpdateStandards: new TokenizationTypes.ActionPermission[](0),
            canUpdateCustomData: new TokenizationTypes.ActionPermission[](0),
            canUpdateManager: new TokenizationTypes.ActionPermission[](0),
            canUpdateCollectionMetadata: new TokenizationTypes.ActionPermission[](0),
            canUpdateValidTokenIds: new TokenizationTypes.TokenIdsActionPermission[](0),
            canUpdateTokenMetadata: new TokenizationTypes.TokenIdsActionPermission[](0),
            canUpdateCollectionApprovals: new TokenizationTypes.CollectionApprovalPermission[](0),
            canAddMoreAliasPaths: new TokenizationTypes.ActionPermission[](0),
            canAddMoreCosmosCoinWrapperPaths: new TokenizationTypes.ActionPermission[](0)
        });
    }

    /**
     * @notice Validates that a UintRange is valid (start <= end)
     * @param range The UintRange to validate
     * @return valid True if the range is valid
     */
    function validateUintRange(TokenizationTypes.UintRange memory range) internal pure returns (bool valid) {
        return range.start <= range.end;
    }

    /**
     * @notice Validates an array of UintRanges
     * @param ranges Array of UintRanges to validate
     * @return valid True if all ranges are valid
     */
    function validateUintRangeArray(TokenizationTypes.UintRange[] memory ranges) internal pure returns (bool valid) {
        for (uint256 i = 0; i < ranges.length; i++) {
            if (!validateUintRange(ranges[i])) {
                return false;
            }
        }
        return true;
    }

    /**
     * @notice Validates that a Balance has valid ranges
     * @param balance The Balance to validate
     * @return valid True if the balance is valid
     */
    function validateBalance(TokenizationTypes.Balance memory balance) internal pure returns (bool valid) {
        return validateUintRangeArray(balance.ownershipTimes) && validateUintRangeArray(balance.tokenIds);
    }

    /**
     * @notice Creates a full ownership time range (from 1 to max uint64)
     * @dev Uses uint64 max because BitBadges internally uses uint64 for timestamps
     * @return range The full ownership time range
     */
    function createFullOwnershipTimeRange() internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(1, MAX_UINT64);
    }

    /**
     * @notice Creates a single token ID range
     * @param tokenId The token ID (both start and end will be this value)
     * @return range The UintRange for the single token ID
     */
    function createSingleTokenIdRange(uint256 tokenId) internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(tokenId, tokenId);
    }

    /**
     * @notice Creates token ID ranges for a consecutive sequence
     * @param startTokenId The first token ID
     * @param endTokenId The last token ID (inclusive)
     * @return range The UintRange for the token ID sequence
     */
    function createTokenIdSequence(
        uint256 startTokenId,
        uint256 endTokenId
    ) internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(startTokenId, endTokenId);
    }

    /**
     * @notice Creates an ownership time range from current time to a future time
     * @param startTime The start time (typically block.timestamp)
     * @param duration The duration in seconds
     * @return range The ownership time range
     */
    function createOwnershipTimeRange(
        uint256 startTime,
        uint256 duration
    ) internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(startTime, startTime + duration);
    }

    /**
     * @notice Creates an ownership time range from current time to expiration
     * @param expirationTime The expiration timestamp
     * @return range The ownership time range from now to expiration
     */
    function createOwnershipTimeRangeToExpiration(
        uint256 expirationTime
    ) internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(block.timestamp, expirationTime);
    }

    /**
     * @notice Creates a time range for a specific time point (single timestamp)
     * @param timestamp The timestamp (both start and end)
     * @return range The UintRange for the single timestamp
     */
    function createTimePoint(
        uint256 timestamp
    ) internal pure returns (TokenizationTypes.UintRange memory range) {
        return createUintRange(timestamp, timestamp);
    }

    /**
     * @notice Creates a time range for current time (single point)
     * @return range The UintRange for current block timestamp
     */
    function createCurrentTimePoint() internal view returns (TokenizationTypes.UintRange memory range) {
        return createTimePoint(block.timestamp);
    }

    /**
     * @notice Creates an empty ActionPermission (no restrictions)
     * @return permission An ActionPermission with empty arrays
     */
    function createEmptyActionPermission() internal pure returns (TokenizationTypes.ActionPermission memory permission) {
        return TokenizationTypes.ActionPermission({
            permanentlyPermittedTimes: new TokenizationTypes.UintRange[](0),
            permanentlyForbiddenTimes: new TokenizationTypes.UintRange[](0)
        });
    }

    /**
     * @notice Creates an ActionPermission that is always permitted
     * @return permission An ActionPermission with full time range permitted
     */
    function createAlwaysPermittedActionPermission() internal pure returns (TokenizationTypes.ActionPermission memory permission) {
        TokenizationTypes.UintRange[] memory permittedTimes = new TokenizationTypes.UintRange[](1);
        permittedTimes[0] = createFullOwnershipTimeRange();
        return TokenizationTypes.ActionPermission({
            permanentlyPermittedTimes: permittedTimes,
            permanentlyForbiddenTimes: new TokenizationTypes.UintRange[](0)
        });
    }

    /**
     * @notice Creates a Balance with full ownership time range
     * @param amount The amount of tokens
     * @param tokenIds Array of token ID ranges
     * @return balance The constructed Balance with full ownership time
     */
    function createBalanceWithFullOwnership(
        uint256 amount,
        TokenizationTypes.UintRange[] memory tokenIds
    ) internal pure returns (TokenizationTypes.Balance memory balance) {
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = createFullOwnershipTimeRange();
        return createBalance(amount, ownershipTimes, tokenIds);
    }

    /**
     * @notice Creates a Balance for a single token ID with full ownership
     * @param amount The amount of tokens
     * @param tokenId The token ID
     * @return balance The constructed Balance
     */
    function createBalanceForSingleToken(
        uint256 amount,
        uint256 tokenId
    ) internal pure returns (TokenizationTypes.Balance memory balance) {
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = createSingleTokenIdRange(tokenId);
        return createBalanceWithFullOwnership(amount, tokenIds);
    }

    /**
     * @notice Converts an address to a Cosmos address string (hex format with 0x prefix)
     * @param addr The EVM address
     * @return cosmosAddr The address as a hex string (for use in manager fields, etc.)
     */
    function addressToCosmosString(address addr) internal pure returns (string memory cosmosAddr) {
        bytes memory data = abi.encodePacked(addr);
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < data.length; i++) {
            str[2 + i * 2] = alphabet[uint8(data[i] >> 4)];
            str[3 + i * 2] = alphabet[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }

    /**
     * @notice Creates a simple CollectionMetadata with only URI
     * @param uri The metadata URI
     * @return metadata The constructed CollectionMetadata
     */
    function createCollectionMetadataFromURI(string memory uri) internal pure returns (TokenizationTypes.CollectionMetadata memory metadata) {
        return createCollectionMetadata(uri, "");
    }

    /**
     * @notice Creates a simple TokenMetadata with only URI
     * @param uri The metadata URI
     * @param tokenIds Array of token ID ranges
     * @return metadata The constructed TokenMetadata
     */
    function createTokenMetadataFromURI(
        string memory uri,
        TokenizationTypes.UintRange[] memory tokenIds
    ) internal pure returns (TokenizationTypes.TokenMetadata memory metadata) {
        return createTokenMetadata(uri, "", tokenIds);
    }
}

