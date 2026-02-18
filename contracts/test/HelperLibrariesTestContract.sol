// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationWrappers.sol";
import "../libraries/TokenizationBuilders.sol";
import "../libraries/TokenizationHelpers.sol";
import "../libraries/TokenizationJSONHelpers.sol";
import "../libraries/TokenizationErrors.sol";

/**
 * @title HelperLibrariesTestContract
 * @notice Comprehensive test contract for all helper libraries
 * @dev Tests all helper libraries E2E to verify:
 *      - Proper data population
 *      - Correct JSON construction
 *      - Correct underlying message calls
 *      - Type safety and error handling
 */
contract HelperLibrariesTestContract {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    // ============ Events for Test Verification ============
    
    event TestResult(
        string testName,
        bool success,
        bytes returnData
    );
    
    event JSONConstructed(
        string methodName,
        string json
    );
    
    // ============ TokenizationWrappers Tests ============
    
    /**
     * @notice Test transferTokens wrapper
     */
    function testTransferTokensWrapper(
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        TokenizationTypes.UintRange[] memory tokenIds,
        TokenizationTypes.UintRange[] memory ownershipTimes
    ) external returns (bool) {
        bool success = TokenizationWrappers.transferTokens(
            TOKENIZATION,
            collectionId,
            toAddresses,
            amount,
            tokenIds,
            ownershipTimes
        );
        emit TestResult("transferTokensWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test transferSingleToken convenience wrapper
     */
    function testTransferSingleTokenWrapper(
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) external returns (bool) {
        bool success = TokenizationWrappers.transferSingleToken(
            TOKENIZATION,
            collectionId,
            to,
            amount,
            tokenId
        );
        emit TestResult("transferSingleTokenWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test transferTokensWithFullOwnership wrapper
     */
    function testTransferTokensWithFullOwnershipWrapper(
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        TokenizationTypes.UintRange[] memory tokenIds
    ) external returns (bool) {
        bool success = TokenizationWrappers.transferTokensWithFullOwnership(
            TOKENIZATION,
            collectionId,
            toAddresses,
            amount,
            tokenIds
        );
        emit TestResult("transferTokensWithFullOwnershipWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test createDynamicStore wrapper
     */
    function testCreateDynamicStoreWrapper(
        bool defaultValue,
        string memory uri,
        string memory customData
    ) external returns (uint256) {
        uint256 storeId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION,
            defaultValue,
            uri,
            customData
        );
        emit TestResult("createDynamicStoreWrapper", true, abi.encode(storeId));
        return storeId;
    }
    
    /**
     * @notice Test setDynamicStoreValue wrapper
     */
    function testSetDynamicStoreValueWrapper(
        uint256 storeId,
        address address_,
        bool value
    ) external returns (bool) {
        bool success = TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION,
            storeId,
            address_,
            value
        );
        emit TestResult("setDynamicStoreValueWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test updateDynamicStore wrapper
     */
    function testUpdateDynamicStoreWrapper(
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string memory uri,
        string memory customData
    ) external returns (bool) {
        bool success = TokenizationWrappers.updateDynamicStore(
            TOKENIZATION,
            storeId,
            defaultValue,
            globalEnabled,
            uri,
            customData
        );
        emit TestResult("updateDynamicStoreWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test deleteDynamicStore wrapper
     */
    function testDeleteDynamicStoreWrapper(
        uint256 storeId
    ) external returns (bool) {
        bool success = TokenizationWrappers.deleteDynamicStore(
            TOKENIZATION,
            storeId
        );
        emit TestResult("deleteDynamicStoreWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test deleteCollection wrapper
     */
    function testDeleteCollectionWrapper(
        uint256 collectionId
    ) external returns (bool) {
        bool success = TokenizationWrappers.deleteCollection(
            TOKENIZATION,
            collectionId
        );
        emit TestResult("deleteCollectionWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test deleteIncomingApproval wrapper
     */
    function testDeleteIncomingApprovalWrapper(
        uint256 collectionId,
        string memory approvalId
    ) external returns (bool) {
        bool success = TokenizationWrappers.deleteIncomingApproval(
            TOKENIZATION,
            collectionId,
            approvalId
        );
        emit TestResult("deleteIncomingApprovalWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test deleteOutgoingApproval wrapper
     */
    function testDeleteOutgoingApprovalWrapper(
        uint256 collectionId,
        string memory approvalId
    ) external returns (bool) {
        bool success = TokenizationWrappers.deleteOutgoingApproval(
            TOKENIZATION,
            collectionId,
            approvalId
        );
        emit TestResult("deleteOutgoingApprovalWrapper", success, abi.encode(success));
        return success;
    }
    
    /**
     * @notice Test getCollection wrapper
     */
    function testGetCollectionWrapper(
        uint256 collectionId
    ) external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.getCollection(
            TOKENIZATION,
            collectionId
        );
        emit TestResult("getCollectionWrapper", result.length > 0, result);
        return result;
    }
    
    /**
     * @notice Test getBalance wrapper
     */
    function testGetBalanceWrapper(
        uint256 collectionId,
        address userAddress
    ) external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.getBalance(
            TOKENIZATION,
            collectionId,
            userAddress
        );
        emit TestResult("getBalanceWrapper", result.length > 0, result);
        return result;
    }
    
    /**
     * @notice Test getBalanceAmount wrapper
     */
    function testGetBalanceAmountWrapper(
        uint256 collectionId,
        address userAddress,
        TokenizationTypes.UintRange[] memory tokenIds,
        TokenizationTypes.UintRange[] memory ownershipTimes
    ) external view returns (uint256) {
        uint256 amount = TokenizationWrappers.getBalanceAmount(
            TOKENIZATION,
            collectionId,
            userAddress,
            tokenIds,
            ownershipTimes
        );
        emit TestResult("getBalanceAmountWrapper", true, abi.encode(amount));
        return amount;
    }
    
    /**
     * @notice Test getTotalSupply wrapper
     */
    function testGetTotalSupplyWrapper(
        uint256 collectionId,
        TokenizationTypes.UintRange[] memory tokenIds,
        TokenizationTypes.UintRange[] memory ownershipTimes
    ) external view returns (uint256) {
        uint256 amount = TokenizationWrappers.getTotalSupply(
            TOKENIZATION,
            collectionId,
            tokenIds,
            ownershipTimes
        );
        emit TestResult("getTotalSupplyWrapper", true, abi.encode(amount));
        return amount;
    }
    
    /**
     * @notice Test getDynamicStore wrapper
     */
    function testGetDynamicStoreWrapper(
        uint256 storeId
    ) external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.getDynamicStore(
            TOKENIZATION,
            storeId
        );
        emit TestResult("getDynamicStoreWrapper", result.length > 0, result);
        return result;
    }
    
    /**
     * @notice Test getDynamicStoreValue wrapper
     */
    function testGetDynamicStoreValueWrapper(
        uint256 storeId,
        address userAddress
    ) external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.getDynamicStoreValue(
            TOKENIZATION,
            storeId,
            userAddress
        );
        emit TestResult("getDynamicStoreValueWrapper", result.length > 0, result);
        return result;
    }
    
    /**
     * @notice Test getAddressList wrapper
     */
    function testGetAddressListWrapper(
        string memory listId
    ) external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.getAddressList(
            TOKENIZATION,
            listId
        );
        emit TestResult("getAddressListWrapper", result.length > 0, result);
        return result;
    }
    
    /**
     * @notice Test isAddressReservedProtocol wrapper
     */
    function testIsAddressReservedProtocolWrapper(
        address addr
    ) external view returns (bool) {
        bool isReserved = TokenizationWrappers.isAddressReservedProtocol(
            TOKENIZATION,
            addr
        );
        emit TestResult("isAddressReservedProtocolWrapper", true, abi.encode(isReserved));
        return isReserved;
    }
    
    /**
     * @notice Test getAllReservedProtocolAddresses wrapper
     */
    function testGetAllReservedProtocolAddressesWrapper() external view returns (address[] memory) {
        address[] memory addresses = TokenizationWrappers.getAllReservedProtocolAddresses(
            TOKENIZATION
        );
        emit TestResult("getAllReservedProtocolAddressesWrapper", true, abi.encode(addresses));
        return addresses;
    }
    
    /**
     * @notice Test params wrapper
     */
    function testParamsWrapper() external view returns (bytes memory) {
        bytes memory result = TokenizationWrappers.params(TOKENIZATION);
        emit TestResult("paramsWrapper", result.length > 0, result);
        return result;
    }
    
    // ============ TokenizationBuilders Tests ============
    
    /**
     * @notice Test CollectionBuilder
     */
    function testCollectionBuilder(
        uint256 tokenIdStart,
        uint256 tokenIdEnd,
        string memory manager,
        string memory uri,
        string memory customData
    ) external returns (string memory) {
        TokenizationBuilders.CollectionBuilder memory builder = 
            TokenizationBuilders.newCollection();
        
        builder = builder.withValidTokenIdRange(tokenIdStart, tokenIdEnd);
        builder = builder.withManager(manager);
        
        TokenizationTypes.CollectionMetadata memory metadata = 
            TokenizationHelpers.createCollectionMetadata(uri, customData);
        builder = builder.withMetadata(metadata);
        
        string[] memory standards = new string[](1);
        standards[0] = "ERC-3643";
        builder = builder.withStandards(standards);
        
        builder = builder.withDefaultBalancesFromFlags(true, true, false);
        
        string memory json = builder.build();
        emit JSONConstructed("CollectionBuilder", json);
        emit TestResult("CollectionBuilder", bytes(json).length > 0, abi.encode(json));
        return json;
    }
    
    /**
     * @notice Test TransferBuilder
     */
    function testTransferBuilder(
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) external returns (string memory) {
        TokenizationBuilders.TransferBuilder memory builder = 
            TokenizationBuilders.newTransfer();
        
        builder = builder.withRecipient(to);
        builder = builder.withAmount(amount);
        builder = builder.withTokenId(tokenId);
        builder = builder.withFullOwnershipTime();
        
        string memory json = builder.buildTransfer(collectionId);
        emit JSONConstructed("TransferBuilder", json);
        emit TestResult("TransferBuilder", bytes(json).length > 0, abi.encode(json));
        return json;
    }
    
    // ============ TokenizationHelpers Tests ============
    
    /**
     * @notice Test helper functions
     */
    function testTokenizationHelpers() external pure returns (bool) {
        // Test createUintRange
        TokenizationTypes.UintRange memory range = TokenizationHelpers.createUintRange(1, 100);
        require(range.start == 1 && range.end == 100, "createUintRange failed");
        
        // Test createFullOwnershipTimeRange
        TokenizationTypes.UintRange memory fullRange = TokenizationHelpers.createFullOwnershipTimeRange();
        require(fullRange.start == 1 && fullRange.end == type(uint64).max, "createFullOwnershipTimeRange failed");
        
        // Test createSingleTokenIdRange
        TokenizationTypes.UintRange memory singleRange = TokenizationHelpers.createSingleTokenIdRange(5);
        require(singleRange.start == 5 && singleRange.end == 5, "createSingleTokenIdRange failed");
        
        // Test createOwnershipTimeRange
        uint256 startTime = 1000;
        uint256 duration = 3600;
        TokenizationTypes.UintRange memory timeRange = TokenizationHelpers.createOwnershipTimeRange(startTime, duration);
        require(timeRange.start == startTime && timeRange.end == startTime + duration, "createOwnershipTimeRange failed");
        
        // Test createCollectionMetadata
        TokenizationTypes.CollectionMetadata memory metadata = TokenizationHelpers.createCollectionMetadata("ipfs://test", "custom");
        require(keccak256(bytes(metadata.uri)) == keccak256(bytes("ipfs://test")), "createCollectionMetadata uri failed");
        require(keccak256(bytes(metadata.customData)) == keccak256(bytes("custom")), "createCollectionMetadata customData failed");
        
        // Test addressToCosmosString
        address testAddr = address(0x1234567890123456789012345678901234567890);
        string memory cosmosAddr = TokenizationHelpers.addressToCosmosString(testAddr);
        require(bytes(cosmosAddr).length > 0, "addressToCosmosString failed");
        
        return true;
    }
    
    // ============ TokenizationJSONHelpers Tests ============
    
    /**
     * @notice Test JSON construction for transferTokens
     */
    function testTransferTokensJSON(
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        uint256[] memory tokenIdStarts,
        uint256[] memory tokenIdEnds,
        uint256[] memory ownershipStarts,
        uint256[] memory ownershipEnds
    ) external pure returns (string memory) {
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeArrayToJson(tokenIdStarts, tokenIdEnds);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeArrayToJson(ownershipStarts, ownershipEnds);
        
        string memory json = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            toAddresses,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        
        emit JSONConstructed("transferTokensJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction for getCollection
     */
    function testGetCollectionJSON(
        uint256 collectionId
    ) external pure returns (string memory) {
        string memory json = TokenizationJSONHelpers.getCollectionJSON(collectionId);
        emit JSONConstructed("getCollectionJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction for getBalance
     */
    function testGetBalanceJSON(
        uint256 collectionId,
        address userAddress
    ) external pure returns (string memory) {
        string memory json = TokenizationJSONHelpers.getBalanceJSON(collectionId, userAddress);
        emit JSONConstructed("getBalanceJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction for getBalanceAmount (single tokenId and ownershipTime)
     */
    function testGetBalanceAmountJSON(
        uint256 collectionId,
        address userAddress,
        uint256 tokenId,
        uint256 ownershipTime
    ) external pure returns (string memory) {
        string memory json = TokenizationJSONHelpers.getBalanceAmountJSON(
            collectionId,
            userAddress,
            tokenId,
            ownershipTime
        );

        emit JSONConstructed("getBalanceAmountJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction for createDynamicStore
     */
    function testCreateDynamicStoreJSON(
        bool defaultValue,
        string memory uri,
        string memory customData
    ) external pure returns (string memory) {
        string memory json = TokenizationJSONHelpers.createDynamicStoreJSON(defaultValue, uri, customData);
        emit JSONConstructed("createDynamicStoreJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction for setDynamicStoreValue
     */
    function testSetDynamicStoreValueJSON(
        uint256 storeId,
        address address_,
        bool value
    ) external pure returns (string memory) {
        string memory json = TokenizationJSONHelpers.setDynamicStoreValueJSON(storeId, address_, value);
        emit JSONConstructed("setDynamicStoreValueJSON", json);
        return json;
    }
    
    /**
     * @notice Test JSON construction helpers
     */
    function testJSONHelpers() external pure returns (bool) {
        // Test uintRangeArrayToJson
        uint256[] memory starts = new uint256[](2);
        uint256[] memory ends = new uint256[](2);
        starts[0] = 1;
        ends[0] = 10;
        starts[1] = 20;
        ends[1] = 30;
        string memory rangesJson = TokenizationJSONHelpers.uintRangeArrayToJson(starts, ends);
        require(bytes(rangesJson).length > 0, "uintRangeArrayToJson failed");
        
        // Test uintRangeToJson
        string memory singleRangeJson = TokenizationJSONHelpers.uintRangeToJson(1, 100);
        require(bytes(singleRangeJson).length > 0, "uintRangeToJson failed");
        
        // Test stringArrayToJson
        string[] memory strings = new string[](2);
        strings[0] = "standard1";
        strings[1] = "standard2";
        string memory stringsJson = TokenizationJSONHelpers.stringArrayToJson(strings);
        require(bytes(stringsJson).length > 0, "stringArrayToJson failed");
        
        // Test collectionMetadataToJson
        string memory metadataJson = TokenizationJSONHelpers.collectionMetadataToJson("ipfs://test", "custom");
        require(bytes(metadataJson).length > 0, "collectionMetadataToJson failed");
        
        return true;
    }
    
    // ============ TokenizationErrors Tests ============
    
    /**
     * @notice Test error validation helpers
     */
    function testTokenizationErrors(
        uint256 collectionId,
        address address_
    ) external pure returns (bool) {
        // Test requireValidCollectionId
        if (collectionId == 0) {
            // This should revert with InvalidCollectionId error
            TokenizationErrors.requireValidCollectionId(collectionId);
        }
        
        // Test requireValidAddress
        if (address_ == address(0)) {
            // This should revert with InvalidAddress error
            TokenizationErrors.requireValidAddress(address_);
        }
        
        // Test requireNonEmptyString
        string memory emptyStr = "";
        // This should revert with EmptyString error
        TokenizationErrors.requireNonEmptyString(emptyStr, "testParam");
        
        return true;
    }
    
    /**
     * @notice Test error validation with valid inputs (should not revert)
     */
    function testTokenizationErrorsValid(
        uint256 collectionId,
        address address_,
        string memory nonEmptyStr
    ) external pure returns (bool) {
        require(collectionId > 0, "CollectionId must be > 0 for this test");
        require(address_ != address(0), "Address must not be zero for this test");
        require(bytes(nonEmptyStr).length > 0, "String must not be empty for this test");
        
        // These should not revert with valid inputs
        TokenizationErrors.requireValidCollectionId(collectionId);
        TokenizationErrors.requireValidAddress(address_);
        TokenizationErrors.requireNonEmptyString(nonEmptyStr, "testParam");
        
        return true;
    }
}





















