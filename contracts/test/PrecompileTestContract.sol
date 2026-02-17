// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/// @title PrecompileTestContract
/// @notice Test contract for wrapping precompile calls and emitting events
/// @dev This contract is used for E2E testing of the tokenization precompile
contract PrecompileTestContract {
    ITokenizationPrecompile constant precompile = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    // ============ Events ============
    
    event TransferExecuted(
        uint256 indexed collectionId,
        address indexed recipient,
        bool success
    );
    
    event CollectionCreated(
        uint256 indexed collectionId
    );
    
    event CollectionUpdated(
        uint256 indexed collectionId,
        bool success
    );
    
    event CollectionDeleted(
        uint256 indexed collectionId,
        bool success
    );
    
    event ApprovalSet(
        uint256 indexed collectionId,
        string approvalId,
        bool isIncoming,
        bool success
    );
    
    event ApprovalDeleted(
        uint256 indexed collectionId,
        string approvalId,
        bool isIncoming,
        bool success
    );
    
    event DynamicStoreCreated(
        uint256 indexed storeId
    );
    
    event DynamicStoreUpdated(
        uint256 indexed storeId,
        bool success
    );
    
    event DynamicStoreDeleted(
        uint256 indexed storeId,
        bool success
    );
    
    event DynamicStoreValueSet(
        uint256 indexed storeId,
        address indexed address_,
        bool value,
        bool success
    );
    
    event AddressListsCreated(
        uint256 numLists,
        bool success
    );
    
    event VoteCast(
        uint256 indexed collectionId,
        string proposalId,
        bool success
    );
    
    // ============ Transfer Methods ============
    
    /// @notice Wrapper for transferTokens
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
    
    // ============ Approval Methods ============
    
    /// @notice Wrapper for setIncomingApproval
    function testSetIncomingApproval(
        uint256 collectionId,
        UserIncomingApproval calldata approval
    ) external returns (bool) {
        string memory approvalJson = _userIncomingApprovalToJson(approval);
        string memory msgJson = string(abi.encodePacked(
            '{"collectionId":"', TokenizationJSONHelpers.uintToString(collectionId),
            '","approval":', approvalJson, '}'
        ));
        bool success = precompile.setIncomingApproval(msgJson);
        emit ApprovalSet(collectionId, approval.approvalId, true, success);
        return success;
    }
    
    /// @notice Wrapper for setOutgoingApproval
    function testSetOutgoingApproval(
        uint256 collectionId,
        UserOutgoingApproval calldata approval
    ) external returns (bool) {
        string memory approvalJson = _userOutgoingApprovalToJson(approval);
        string memory msgJson = string(abi.encodePacked(
            '{"collectionId":"', TokenizationJSONHelpers.uintToString(collectionId),
            '","approval":', approvalJson, '}'
        ));
        bool success = precompile.setOutgoingApproval(msgJson);
        emit ApprovalSet(collectionId, approval.approvalId, false, success);
        return success;
    }
    
    /// @notice Wrapper for deleteIncomingApproval
    function testDeleteIncomingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteIncomingApprovalJSON(collectionId, approvalId);
        bool success = precompile.deleteIncomingApproval(deleteJson);
        emit ApprovalDeleted(collectionId, approvalId, true, success);
        return success;
    }
    
    /// @notice Wrapper for deleteOutgoingApproval
    function testDeleteOutgoingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteOutgoingApprovalJSON(collectionId, approvalId);
        bool success = precompile.deleteOutgoingApproval(deleteJson);
        emit ApprovalDeleted(collectionId, approvalId, false, success);
        return success;
    }
    
    // ============ Collection Methods ============
    
    /// @notice Wrapper for createCollection
    function testCreateCollection(
        MsgCreateCollection calldata msg_
    ) external returns (uint256) {
        string memory createJson = _msgCreateCollectionToJson(msg_);
        uint256 collectionId = precompile.createCollection(createJson);
        emit CollectionCreated(collectionId);
        return collectionId;
    }
    
    /// @notice Wrapper for updateCollection
    function testUpdateCollection(
        MsgUpdateCollection calldata msg_
    ) external returns (bool) {
        string memory updateJson = _msgUpdateCollectionToJson(msg_);
        uint256 resultCollectionId = precompile.updateCollection(updateJson);
        bool success = (resultCollectionId == msg_.collectionId);
        emit CollectionUpdated(msg_.collectionId, success);
        return success;
    }
    
    /// @notice Wrapper for universalUpdateCollection
    function testUniversalUpdateCollection(
        MsgUniversalUpdateCollection calldata msg_
    ) external returns (bool) {
        string memory updateJson = _msgUniversalUpdateCollectionToJson(msg_);
        uint256 resultCollectionId = precompile.universalUpdateCollection(updateJson);
        bool success = (resultCollectionId == msg_.collectionId);
        emit CollectionUpdated(msg_.collectionId, success);
        return success;
    }
    
    /// @notice Wrapper for deleteCollection
    function testDeleteCollection(
        uint256 collectionId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteCollectionJSON(collectionId);
        bool success = precompile.deleteCollection(deleteJson);
        emit CollectionDeleted(collectionId, success);
        return success;
    }
    
    // ============ Dynamic Store Methods ============
    
    /// @notice Wrapper for createDynamicStore
    function testCreateDynamicStore(
        bool defaultValue,
        string calldata uri,
        string calldata customData
    ) external returns (uint256) {
        string memory createJson = TokenizationJSONHelpers.createDynamicStoreJSON(defaultValue, uri, customData);
        uint256 storeId = precompile.createDynamicStore(createJson);
        emit DynamicStoreCreated(storeId);
        return storeId;
    }
    
    /// @notice Wrapper for updateDynamicStore
    function testUpdateDynamicStore(
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string calldata uri,
        string calldata customData
    ) external returns (bool) {
        string memory updateJson = TokenizationJSONHelpers.updateDynamicStoreJSON(
            storeId,
            defaultValue,
            globalEnabled,
            uri,
            customData
        );
        bool success = precompile.updateDynamicStore(updateJson);
        emit DynamicStoreUpdated(storeId, success);
        return success;
    }
    
    /// @notice Wrapper for deleteDynamicStore
    function testDeleteDynamicStore(
        uint256 storeId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteDynamicStoreJSON(storeId);
        bool success = precompile.deleteDynamicStore(deleteJson);
        emit DynamicStoreDeleted(storeId, success);
        return success;
    }
    
    /// @notice Wrapper for setDynamicStoreValue
    function testSetDynamicStoreValue(
        uint256 storeId,
        address address_,
        bool value
    ) external returns (bool) {
        string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            storeId,
            address_,
            value
        );
        bool success = precompile.setDynamicStoreValue(setValueJson);
        emit DynamicStoreValueSet(storeId, address_, value, success);
        return success;
    }
    
    // ============ Address List Methods ============
    
    /// @notice Wrapper for createAddressLists
    function testCreateAddressLists(
        AddressListInput[] calldata addressLists
    ) external returns (bool) {
        string memory createJson = _addressListsToJson(addressLists);
        bool success = precompile.createAddressLists(createJson);
        emit AddressListsCreated(addressLists.length, success);
        return success;
    }
    
    // ============ Vote Methods ============
    
    /// @notice Wrapper for castVote
    function testCastVote(
        uint256 collectionId,
        string calldata approvalLevel,
        string calldata approverAddress,
        string calldata approvalId,
        string calldata proposalId,
        uint256 yesWeight
    ) external returns (bool) {
        string memory castVoteJson = string(abi.encodePacked(
            '{"collectionId":"', TokenizationJSONHelpers.uintToString(collectionId),
            '","approvalLevel":"', _escapeJson(approvalLevel),
            '","approverAddress":"', _escapeJson(approverAddress),
            '","approvalId":"', _escapeJson(approvalId),
            '","proposalId":"', _escapeJson(proposalId),
            '","yesWeight":"', TokenizationJSONHelpers.uintToString(yesWeight), '"}'
        ));
        bool success = precompile.castVote(castVoteJson);
        emit VoteCast(collectionId, proposalId, success);
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
    
    /// @notice Wrapper for getBalanceAmount (single tokenId and ownershipTime)
    function testGetBalanceAmount(
        uint256 collectionId,
        address address_,
        uint256 tokenId,
        uint256 ownershipTime
    ) external view returns (uint256) {
        string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
            collectionId,
            address_,
            tokenId,
            ownershipTime
        );
        return precompile.getBalanceAmount(balanceJson);
    }

    /// @notice Wrapper for getTotalSupply (single tokenId and ownershipTime)
    function testGetTotalSupply(
        uint256 collectionId,
        uint256 tokenId,
        uint256 ownershipTime
    ) external view returns (uint256) {
        string memory supplyJson = TokenizationJSONHelpers.getTotalSupplyJSON(
            collectionId,
            tokenId,
            ownershipTime
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
    
    /// @notice Wrapper for getDynamicStore
    function testGetDynamicStore(
        uint256 storeId
    ) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getDynamicStoreJSON(storeId);
        return precompile.getDynamicStore(queryJson);
    }
    
    /// @notice Wrapper for getDynamicStoreValue
    function testGetDynamicStoreValue(
        uint256 storeId,
        address address_
    ) external view returns (bool) {
        // Note: getDynamicStoreValue returns bytes, but we need bool
        // This is a simplified wrapper - full implementation would decode bytes
        string memory queryJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(storeId, address_);
        bytes memory result = precompile.getDynamicStoreValue(queryJson);
        // For now, return false if empty, true if not empty (simplified)
        return result.length > 0;
    }
    
    // ============ Internal Helper Functions ============
    
    /// @notice Convert UserIncomingApproval struct to JSON
    function _userIncomingApprovalToJson(UserIncomingApproval calldata approval) internal pure returns (string memory) {
        string memory result = "{";
        
        // fromListId
        if (bytes(approval.fromListId).length > 0) {
            result = string(abi.encodePacked(result, '"fromListId":"', _escapeJson(approval.fromListId), '"'));
        }
        
        // initiatedByListId
        if (bytes(approval.initiatedByListId).length > 0) {
            if (bytes(approval.fromListId).length > 0) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"initiatedByListId":"', _escapeJson(approval.initiatedByListId), '"'));
        }
        
        // transferTimes
        if (approval.transferTimes.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"transferTimes":', _uintRangeArrayToJson(approval.transferTimes)));
        }
        
        // tokenIds
        if (approval.tokenIds.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"tokenIds":', _uintRangeArrayToJson(approval.tokenIds)));
        }
        
        // ownershipTimes
        if (approval.ownershipTimes.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"ownershipTimes":', _uintRangeArrayToJson(approval.ownershipTimes)));
        }
        
        // uri
        if (bytes(approval.uri).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"uri":"', _escapeJson(approval.uri), '"'));
        }
        
        // customData
        if (bytes(approval.customData).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"customData":"', _escapeJson(approval.customData), '"'));
        }
        
        // approvalId
        if (bytes(approval.approvalId).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"approvalId":"', _escapeJson(approval.approvalId), '"'));
        }
        
        // Note: approvalCriteria and version are complex - simplified for now
        // In practice, these would need full JSON construction
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }
    
    /// @notice Convert UserOutgoingApproval struct to JSON
    function _userOutgoingApprovalToJson(UserOutgoingApproval calldata approval) internal pure returns (string memory) {
        string memory result = "{";
        
        // toListId
        if (bytes(approval.toListId).length > 0) {
            result = string(abi.encodePacked(result, '"toListId":"', _escapeJson(approval.toListId), '"'));
        }
        
        // initiatedByListId
        if (bytes(approval.initiatedByListId).length > 0) {
            if (bytes(approval.toListId).length > 0) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"initiatedByListId":"', _escapeJson(approval.initiatedByListId), '"'));
        }
        
        // transferTimes
        if (approval.transferTimes.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"transferTimes":', _uintRangeArrayToJson(approval.transferTimes)));
        }
        
        // tokenIds
        if (approval.tokenIds.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"tokenIds":', _uintRangeArrayToJson(approval.tokenIds)));
        }
        
        // ownershipTimes
        if (approval.ownershipTimes.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"ownershipTimes":', _uintRangeArrayToJson(approval.ownershipTimes)));
        }
        
        // uri
        if (bytes(approval.uri).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"uri":"', _escapeJson(approval.uri), '"'));
        }
        
        // customData
        if (bytes(approval.customData).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"customData":"', _escapeJson(approval.customData), '"'));
        }
        
        // approvalId
        if (bytes(approval.approvalId).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"approvalId":"', _escapeJson(approval.approvalId), '"'));
        }
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }
    
    /// @notice Convert UintRange array to JSON
    function _uintRangeArrayToJson(UintRange[] calldata ranges) internal pure returns (string memory) {
        if (ranges.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < ranges.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                '{"start":"', TokenizationJSONHelpers.uintToString(ranges[i].start),
                '","end":"', TokenizationJSONHelpers.uintToString(ranges[i].end), '"}'
            ));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }
    
    /// @notice Convert MsgCreateCollection to JSON (simplified - handles common fields)
    function _msgCreateCollectionToJson(MsgCreateCollection calldata msg_) internal pure returns (string memory) {
        string memory result = "{";
        
        // validTokenIds
        if (msg_.validTokenIds.length > 0) {
            result = string(abi.encodePacked(result, '"validTokenIds":', _uintRangeArrayToJson(msg_.validTokenIds)));
        }
        
        // manager
        if (bytes(msg_.manager).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"manager":"', _escapeJson(msg_.manager), '"'));
        }
        
        // collectionMetadata
        if (bytes(msg_.collectionMetadata.uri).length > 0 || bytes(msg_.collectionMetadata.customData).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"collectionMetadata":', 
                TokenizationJSONHelpers.collectionMetadataToJson(msg_.collectionMetadata.uri, msg_.collectionMetadata.customData)));
        }
        
        // defaultBalances (simplified - just auto-approve flags)
        if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
        result = string(abi.encodePacked(result, '"defaultBalances":',
            TokenizationJSONHelpers.simpleUserBalanceStoreToJson(
                msg_.defaultBalances.autoApproveSelfInitiatedOutgoingTransfers,
                msg_.defaultBalances.autoApproveSelfInitiatedIncomingTransfers,
                msg_.defaultBalances.autoApproveAllIncomingTransfers
            )));
        
        // standards
        if (msg_.standards.length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"standards":', TokenizationJSONHelpers.stringArrayToJson(msg_.standards)));
        }
        
        // customData
        if (bytes(msg_.customData).length > 0) {
            if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"customData":"', _escapeJson(msg_.customData), '"'));
        }
        
        // isArchived
        if (bytes(result).length > 1) result = string(abi.encodePacked(result, ","));
        result = string(abi.encodePacked(result, '"isArchived":', msg_.isArchived ? 'true' : 'false'));
        
        // Note: collectionPermissions, tokenMetadata, collectionApprovals are complex - simplified
        // In practice, these would need full JSON construction
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }
    
    /// @notice Convert MsgUpdateCollection to JSON (simplified)
    function _msgUpdateCollectionToJson(MsgUpdateCollection calldata msg_) internal pure returns (string memory) {
        string memory result = "{";
        
        // collectionId
        result = string(abi.encodePacked(result, '"collectionId":"', TokenizationJSONHelpers.uintToString(msg_.collectionId), '"'));
        
        // Update flags and fields (simplified - only include if update flag is true)
        if (msg_.updateValidTokenIds && msg_.validTokenIds.length > 0) {
            result = string(abi.encodePacked(result, ',"validTokenIds":', _uintRangeArrayToJson(msg_.validTokenIds)));
        }
        if (msg_.updateManager && bytes(msg_.manager).length > 0) {
            result = string(abi.encodePacked(result, ',"manager":"', _escapeJson(msg_.manager), '"'));
        }
        if (msg_.updateCollectionMetadata) {
            result = string(abi.encodePacked(result, ',"collectionMetadata":',
                TokenizationJSONHelpers.collectionMetadataToJson(msg_.collectionMetadata.uri, msg_.collectionMetadata.customData)));
        }
        if (msg_.updateStandards && msg_.standards.length > 0) {
            result = string(abi.encodePacked(result, ',"standards":', TokenizationJSONHelpers.stringArrayToJson(msg_.standards)));
        }
        if (msg_.updateCustomData && bytes(msg_.customData).length > 0) {
            result = string(abi.encodePacked(result, ',"customData":"', _escapeJson(msg_.customData), '"'));
        }
        if (msg_.updateIsArchived) {
            result = string(abi.encodePacked(result, ',"isArchived":', msg_.isArchived ? 'true' : 'false'));
        }
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }
    
    /// @notice Convert MsgUniversalUpdateCollection to JSON (simplified)
    function _msgUniversalUpdateCollectionToJson(MsgUniversalUpdateCollection calldata msg_) internal pure returns (string memory) {
        // Similar to updateCollection but with all fields
        string memory result = "{";
        result = string(abi.encodePacked(result, '"collectionId":"', TokenizationJSONHelpers.uintToString(msg_.collectionId), '"'));
        
        // Add other fields as needed (simplified for now)
        // In practice, this would need full JSON construction for all fields
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }
    
    /// @notice Convert AddressListInput array to JSON
    function _addressListsToJson(AddressListInput[] calldata addressLists) internal pure returns (string memory) {
        string memory result = '{"addressLists":[';
        for (uint256 i = 0; i < addressLists.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, "{"));
            
            // listId
            result = string(abi.encodePacked(result, '"listId":"', _escapeJson(addressLists[i].listId), '"'));
            
            // addresses
            if (addressLists[i].addresses.length > 0) {
                result = string(abi.encodePacked(result, ',"addresses":['));
                for (uint256 j = 0; j < addressLists[i].addresses.length; j++) {
                    if (j > 0) result = string(abi.encodePacked(result, ","));
                    result = string(abi.encodePacked(result, '"', _escapeJson(addressLists[i].addresses[j]), '"'));
                }
                result = string(abi.encodePacked(result, "]"));
            }
            
            // whitelist
            result = string(abi.encodePacked(result, ',"whitelist":', addressLists[i].whitelist ? 'true' : 'false'));
            
            // uri
            if (bytes(addressLists[i].uri).length > 0) {
                result = string(abi.encodePacked(result, ',"uri":"', _escapeJson(addressLists[i].uri), '"'));
            }
            
            // customData
            if (bytes(addressLists[i].customData).length > 0) {
                result = string(abi.encodePacked(result, ',"customData":"', _escapeJson(addressLists[i].customData), '"'));
            }
            
            result = string(abi.encodePacked(result, "}"));
        }
        result = string(abi.encodePacked(result, "]}"));
        return result;
    }
    
    /// @notice Simple JSON string escaping (basic)
    function _escapeJson(string calldata str) private pure returns (string memory) {
        // Use the helper library's escape function via a workaround
        // Since we can't call internal functions from calldata, we'll do basic escaping here
        bytes memory strBytes = bytes(str);
        bytes memory result = new bytes(strBytes.length * 2);
        uint256 resultIndex = 0;
        
        for (uint256 i = 0; i < strBytes.length; i++) {
            bytes1 char = strBytes[i];
            if (char == 0x22) { // "
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x22; // "
            } else if (char == 0x5C) { // \
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x5C; // \
            } else if (char >= 0x20) {
                result[resultIndex++] = char;
            }
        }
        
        bytes memory finalResult = new bytes(resultIndex);
        for (uint256 i = 0; i < resultIndex; i++) {
            finalResult[i] = result[i];
        }
        return string(finalResult);
    }
}

