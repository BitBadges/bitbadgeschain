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
    function createUintRange(uint256 start, uint256 end) internal pure returns (UintRange memory range) {
        require(start <= end, "TokenizationHelpers: start must be <= end");
        return UintRange({start: start, end: end});
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
    ) internal pure returns (UintRange[] memory ranges) {
        require(starts.length == ends.length, "TokenizationHelpers: starts and ends arrays must have same length");
        ranges = new UintRange[](starts.length);
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
        UintRange[] memory ownershipTimes,
        UintRange[] memory tokenIds
    ) internal pure returns (Balance memory balance) {
        return Balance({
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
    ) internal pure returns (CollectionMetadata memory metadata) {
        return CollectionMetadata({
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
        UintRange[] memory tokenIds
    ) internal pure returns (TokenMetadata memory metadata) {
        return TokenMetadata({
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
        UintRange[] memory permittedTimes,
        UintRange[] memory forbiddenTimes
    ) internal pure returns (ActionPermission memory permission) {
        return ActionPermission({
            permanentlyPermittedTimes: permittedTimes,
            permanentlyForbiddenTimes: forbiddenTimes
        });
    }

    /**
     * @notice Creates an empty UserBalanceStore with default values
     * @return store The constructed UserBalanceStore with empty arrays and false booleans
     */
    function createEmptyUserBalanceStore() internal pure returns (UserBalanceStore memory store) {
        return UserBalanceStore({
            balances: new Balance[](0),
            outgoingApprovals: new UserOutgoingApproval[](0),
            incomingApprovals: new UserIncomingApproval[](0),
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
    function createEmptyUserPermissions() internal pure returns (UserPermissions memory permissions) {
        return UserPermissions({
            canUpdateOutgoingApprovals: new UserOutgoingApprovalPermission[](0),
            canUpdateIncomingApprovals: new UserIncomingApprovalPermission[](0),
            canUpdateAutoApproveSelfInitiatedOutgoingTransfers: new ActionPermission[](0),
            canUpdateAutoApproveSelfInitiatedIncomingTransfers: new ActionPermission[](0),
            canUpdateAutoApproveAllIncomingTransfers: new ActionPermission[](0)
        });
    }

    /**
     * @notice Creates an empty CollectionPermissions struct
     * @return permissions The constructed CollectionPermissions with empty arrays
     */
    function createEmptyCollectionPermissions() internal pure returns (CollectionPermissions memory permissions) {
        return CollectionPermissions({
            canDeleteCollection: new ActionPermission[](0),
            canArchiveCollection: new ActionPermission[](0),
            canUpdateStandards: new ActionPermission[](0),
            canUpdateCustomData: new ActionPermission[](0),
            canUpdateManager: new ActionPermission[](0),
            canUpdateCollectionMetadata: new ActionPermission[](0),
            canUpdateValidTokenIds: new TokenIdsActionPermission[](0),
            canUpdateTokenMetadata: new TokenIdsActionPermission[](0),
            canUpdateCollectionApprovals: new CollectionApprovalPermission[](0),
            canAddMoreAliasPaths: new ActionPermission[](0),
            canAddMoreCosmosCoinWrapperPaths: new ActionPermission[](0)
        });
    }

    /**
     * @notice Validates that a UintRange is valid (start <= end)
     * @param range The UintRange to validate
     * @return valid True if the range is valid
     */
    function validateUintRange(UintRange memory range) internal pure returns (bool valid) {
        return range.start <= range.end;
    }

    /**
     * @notice Validates an array of UintRanges
     * @param ranges Array of UintRanges to validate
     * @return valid True if all ranges are valid
     */
    function validateUintRangeArray(UintRange[] memory ranges) internal pure returns (bool valid) {
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
    function validateBalance(Balance memory balance) internal pure returns (bool valid) {
        return validateUintRangeArray(balance.ownershipTimes) && validateUintRangeArray(balance.tokenIds);
    }

    /**
     * @notice Creates a full ownership time range (from 1 to max uint64)
     * @dev Uses uint64 max because BitBadges internally uses uint64 for timestamps
     * @return range The full ownership time range
     */
    function createFullOwnershipTimeRange() internal pure returns (UintRange memory range) {
        return createUintRange(1, MAX_UINT64);
    }

    /**
     * @notice Creates a single token ID range
     * @param tokenId The token ID (both start and end will be this value)
     * @return range The UintRange for the single token ID
     */
    function createSingleTokenIdRange(uint256 tokenId) internal pure returns (UintRange memory range) {
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
    ) internal pure returns (UintRange memory range) {
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
    ) internal pure returns (UintRange memory range) {
        return createUintRange(startTime, startTime + duration);
    }

    /**
     * @notice Creates an ownership time range from current time to expiration
     * @param expirationTime The expiration timestamp
     * @return range The ownership time range from now to expiration
     */
    function createOwnershipTimeRangeToExpiration(
        uint256 expirationTime
    ) internal pure returns (UintRange memory range) {
        return createUintRange(block.timestamp, expirationTime);
    }

    /**
     * @notice Creates a time range for a specific time point (single timestamp)
     * @param timestamp The timestamp (both start and end)
     * @return range The UintRange for the single timestamp
     */
    function createTimePoint(
        uint256 timestamp
    ) internal pure returns (UintRange memory range) {
        return createUintRange(timestamp, timestamp);
    }

    /**
     * @notice Creates a time range for current time (single point)
     * @return range The UintRange for current block timestamp
     */
    function createCurrentTimePoint() internal view returns (UintRange memory range) {
        return createTimePoint(block.timestamp);
    }

    /**
     * @notice Creates an empty ActionPermission (no restrictions)
     * @return permission An ActionPermission with empty arrays
     */
    function createEmptyActionPermission() internal pure returns (ActionPermission memory permission) {
        return ActionPermission({
            permanentlyPermittedTimes: new UintRange[](0),
            permanentlyForbiddenTimes: new UintRange[](0)
        });
    }

    /**
     * @notice Creates an ActionPermission that is always permitted
     * @return permission An ActionPermission with full time range permitted
     */
    function createAlwaysPermittedActionPermission() internal pure returns (ActionPermission memory permission) {
        UintRange[] memory permittedTimes = new UintRange[](1);
        permittedTimes[0] = createFullOwnershipTimeRange();
        return ActionPermission({
            permanentlyPermittedTimes: permittedTimes,
            permanentlyForbiddenTimes: new UintRange[](0)
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
        UintRange[] memory tokenIds
    ) internal pure returns (Balance memory balance) {
        UintRange[] memory ownershipTimes = new UintRange[](1);
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
    ) internal pure returns (Balance memory balance) {
        UintRange[] memory tokenIds = new UintRange[](1);
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
    function createCollectionMetadataFromURI(string memory uri) internal pure returns (CollectionMetadata memory metadata) {
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
        UintRange[] memory tokenIds
    ) internal pure returns (TokenMetadata memory metadata) {
        return createTokenMetadata(uri, "", tokenIds);
    }

    // ============ EVM Query Challenge Helpers (v25) ============

    /**
     * @notice Creates an EVMQueryChallenge struct
     * @param contractAddress The EVM contract address to call (hex string with 0x prefix)
     * @param callData The calldata for the static call (hex string)
     * @param expectedResult The expected result to compare against (hex string)
     * @param comparisonOperator The comparison operator ("eq", "ne", "gt", "gte", "lt", "lte")
     * @param gasLimit Gas limit for the call (default 100000, max 500000)
     * @return challenge The constructed EVMQueryChallenge
     * @dev Use this for building EVM query challenges in approvals or invariants.
     *      Placeholders available in callData: $initiator, $sender, $recipient, $collectionId
     *
     * Example:
     * ```solidity
     * // Check that a contract returns true for the sender
     * EVMQueryChallenge memory challenge = TokenizationHelpers.createEVMQueryChallenge(
     *     "0x1234...",           // contract address
     *     "0x70a08231...",       // balanceOf(address) calldata with $sender placeholder
     *     "0x0000...0001",       // expected: at least 1
     *     "gte",                 // greater than or equal
     *     100000                 // gas limit
     * );
     * ```
     */
    function createEVMQueryChallenge(
        string memory contractAddress,
        string memory callData,
        string memory expectedResult,
        string memory comparisonOperator,
        uint256 gasLimit
    ) internal pure returns (EVMQueryChallenge memory challenge) {
        return EVMQueryChallenge({
            contractAddress: contractAddress,
            callData: callData,
            expectedResult: expectedResult,
            comparisonOperator: comparisonOperator,
            gasLimit: gasLimit,
            uri: "",
            customData: ""
        });
    }

    /**
     * @notice Creates an EVMQueryChallenge with metadata
     * @param contractAddress The EVM contract address to call
     * @param callData The calldata for the static call
     * @param expectedResult The expected result to compare against
     * @param comparisonOperator The comparison operator
     * @param gasLimit Gas limit for the call
     * @param uri Metadata URI for the challenge
     * @param customData Custom data for the challenge
     * @return challenge The constructed EVMQueryChallenge
     */
    function createEVMQueryChallengeWithMetadata(
        string memory contractAddress,
        string memory callData,
        string memory expectedResult,
        string memory comparisonOperator,
        uint256 gasLimit,
        string memory uri,
        string memory customData
    ) internal pure returns (EVMQueryChallenge memory challenge) {
        return EVMQueryChallenge({
            contractAddress: contractAddress,
            callData: callData,
            expectedResult: expectedResult,
            comparisonOperator: comparisonOperator,
            gasLimit: gasLimit,
            uri: uri,
            customData: customData
        });
    }

    /**
     * @notice Creates an equality check EVM query challenge
     * @param contractAddress The EVM contract address to call
     * @param callData The calldata for the static call
     * @param expectedResult The expected result (must equal)
     * @return challenge The constructed EVMQueryChallenge with "eq" operator
     */
    function createEVMQueryChallengeEq(
        string memory contractAddress,
        string memory callData,
        string memory expectedResult
    ) internal pure returns (EVMQueryChallenge memory challenge) {
        return createEVMQueryChallenge(contractAddress, callData, expectedResult, "eq", 100000);
    }

    /**
     * @notice Creates a greater-than-or-equal check EVM query challenge
     * @param contractAddress The EVM contract address to call
     * @param callData The calldata for the static call
     * @param minValue The minimum expected result
     * @return challenge The constructed EVMQueryChallenge with "gte" operator
     */
    function createEVMQueryChallengeGte(
        string memory contractAddress,
        string memory callData,
        string memory minValue
    ) internal pure returns (EVMQueryChallenge memory challenge) {
        return createEVMQueryChallenge(contractAddress, callData, minValue, "gte", 100000);
    }

    /**
     * @notice Creates a less-than-or-equal check EVM query challenge
     * @param contractAddress The EVM contract address to call
     * @param callData The calldata for the static call
     * @param maxValue The maximum expected result
     * @return challenge The constructed EVMQueryChallenge with "lte" operator
     */
    function createEVMQueryChallengeLte(
        string memory contractAddress,
        string memory callData,
        string memory maxValue
    ) internal pure returns (EVMQueryChallenge memory challenge) {
        return createEVMQueryChallenge(contractAddress, callData, maxValue, "lte", 100000);
    }

    // ============ Collection Invariants Helpers (v25) ============

    /**
     * @notice Creates an empty CollectionInvariants struct
     * @return invariants CollectionInvariants with all fields set to defaults
     */
    function createEmptyCollectionInvariants() internal pure returns (CollectionInvariants memory invariants) {
        return CollectionInvariants({
            noCustomOwnershipTimes: false,
            maxSupplyPerId: 0,
            cosmosCoinBackedPath: CosmosCoinBackedPath({
                addr: "",
                conversion: Conversion({
                    sideA: ConversionSideAWithDenom({amount: 0, denom: ""}),
                    sideB: new Balance[](0)
                })
            }),
            noForcefulPostMintTransfers: false,
            disablePoolCreation: false,
            evmQueryChallenges: new EVMQueryChallenge[](0)
        });
    }

    /**
     * @notice Creates CollectionInvariants with EVM query challenges
     * @param evmQueryChallenges Array of EVM query challenges to enforce post-transfer
     * @return invariants CollectionInvariants with the specified EVM query challenges
     * @dev EVM query challenges in invariants are executed after every transfer.
     *      If any challenge fails, the transfer is reverted.
     *
     * Example - Max holder count invariant:
     * ```solidity
     * EVMQueryChallenge[] memory challenges = new EVMQueryChallenge[](1);
     * challenges[0] = TokenizationHelpers.createEVMQueryChallenge(
     *     maxHoldersCheckerAddress,
     *     abi.encodeWithSignature("checkMaxHolders(uint256,uint256)", collectionId, 100),
     *     bytes32(uint256(1)),  // expect pass (1)
     *     "eq",
     *     200000
     * );
     * CollectionInvariants memory invariants = TokenizationHelpers.createInvariantsWithEVMChallenges(challenges);
     * ```
     */
    function createInvariantsWithEVMChallenges(
        EVMQueryChallenge[] memory evmQueryChallenges
    ) internal pure returns (CollectionInvariants memory invariants) {
        invariants = createEmptyCollectionInvariants();
        invariants.evmQueryChallenges = evmQueryChallenges;
        return invariants;
    }

    /**
     * @notice Creates CollectionInvariants with max supply per token ID
     * @param maxSupplyPerId Maximum supply allowed per token ID
     * @return invariants CollectionInvariants with max supply set
     */
    function createInvariantsWithMaxSupply(
        uint256 maxSupplyPerId
    ) internal pure returns (CollectionInvariants memory invariants) {
        invariants = createEmptyCollectionInvariants();
        invariants.maxSupplyPerId = maxSupplyPerId;
        return invariants;
    }

    /**
     * @notice Creates CollectionInvariants with multiple settings
     * @param noCustomOwnershipTimes Disallow custom ownership times
     * @param maxSupplyPerId Maximum supply per token ID (0 for unlimited)
     * @param noForcefulPostMintTransfers Prevent forceful transfers after minting
     * @param disablePoolCreation Disable liquidity pool creation
     * @param evmQueryChallenges Array of EVM query challenges
     * @return invariants The configured CollectionInvariants
     */
    function createCollectionInvariants(
        bool noCustomOwnershipTimes,
        uint256 maxSupplyPerId,
        bool noForcefulPostMintTransfers,
        bool disablePoolCreation,
        EVMQueryChallenge[] memory evmQueryChallenges
    ) internal pure returns (CollectionInvariants memory invariants) {
        return CollectionInvariants({
            noCustomOwnershipTimes: noCustomOwnershipTimes,
            maxSupplyPerId: maxSupplyPerId,
            cosmosCoinBackedPath: CosmosCoinBackedPath({
                addr: "",
                conversion: Conversion({
                    sideA: ConversionSideAWithDenom({amount: 0, denom: ""}),
                    sideB: new Balance[](0)
                })
            }),
            noForcefulPostMintTransfers: noForcefulPostMintTransfers,
            disablePoolCreation: disablePoolCreation,
            evmQueryChallenges: evmQueryChallenges
        });
    }

    // ============ Calldata Placeholder Helpers ============

    /**
     * @notice Placeholder constants for EVM query challenge calldata
     * @dev These are replaced at runtime with actual values:
     *      $initiator - The address that initiated the transfer
     *      $sender - The sender (from) address
     *      $recipient - The recipient (to) address
     *      $collectionId - The collection ID (as uint256)
     */
    string constant PLACEHOLDER_INITIATOR = "$initiator";
    string constant PLACEHOLDER_SENDER = "$sender";
    string constant PLACEHOLDER_RECIPIENT = "$recipient";
    string constant PLACEHOLDER_COLLECTION_ID = "$collectionId";
}

