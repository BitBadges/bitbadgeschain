// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";

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
    
    // ============ Approval Methods ============
    
    /// @notice Wrapper for setIncomingApproval
    function testSetIncomingApproval(
        uint256 collectionId,
        UserIncomingApproval calldata approval
    ) external returns (bool) {
        bool success = precompile.setIncomingApproval(collectionId, approval);
        emit ApprovalSet(collectionId, approval.approvalId, true, success);
        return success;
    }
    
    /// @notice Wrapper for setOutgoingApproval
    function testSetOutgoingApproval(
        uint256 collectionId,
        UserOutgoingApproval calldata approval
    ) external returns (bool) {
        bool success = precompile.setOutgoingApproval(collectionId, approval);
        emit ApprovalSet(collectionId, approval.approvalId, false, success);
        return success;
    }
    
    /// @notice Wrapper for deleteIncomingApproval
    function testDeleteIncomingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        bool success = precompile.deleteIncomingApproval(collectionId, approvalId);
        emit ApprovalDeleted(collectionId, approvalId, true, success);
        return success;
    }
    
    /// @notice Wrapper for deleteOutgoingApproval
    function testDeleteOutgoingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        bool success = precompile.deleteOutgoingApproval(collectionId, approvalId);
        emit ApprovalDeleted(collectionId, approvalId, false, success);
        return success;
    }
    
    // ============ Collection Methods ============
    
    /// @notice Wrapper for createCollection
    function testCreateCollection(
        MsgCreateCollection calldata msg_
    ) external returns (uint256) {
        uint256 collectionId = precompile.createCollection(msg_);
        emit CollectionCreated(collectionId);
        return collectionId;
    }
    
    /// @notice Wrapper for updateCollection
    function testUpdateCollection(
        MsgUpdateCollection calldata msg_
    ) external returns (bool) {
        uint256 resultCollectionId = precompile.updateCollection(msg_);
        bool success = (resultCollectionId == msg_.collectionId);
        emit CollectionUpdated(msg_.collectionId, success);
        return success;
    }
    
    /// @notice Wrapper for universalUpdateCollection
    function testUniversalUpdateCollection(
        MsgUniversalUpdateCollection calldata msg_
    ) external returns (bool) {
        uint256 resultCollectionId = precompile.universalUpdateCollection(msg_);
        bool success = (resultCollectionId == msg_.collectionId);
        emit CollectionUpdated(msg_.collectionId, success);
        return success;
    }
    
    /// @notice Wrapper for deleteCollection
    function testDeleteCollection(
        uint256 collectionId
    ) external returns (bool) {
        bool success = precompile.deleteCollection(collectionId);
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
        uint256 storeId = precompile.createDynamicStore(defaultValue, uri, customData);
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
        bool success = precompile.updateDynamicStore(
            storeId,
            defaultValue,
            globalEnabled,
            uri,
            customData
        );
        emit DynamicStoreUpdated(storeId, success);
        return success;
    }
    
    /// @notice Wrapper for deleteDynamicStore
    function testDeleteDynamicStore(
        uint256 storeId
    ) external returns (bool) {
        bool success = precompile.deleteDynamicStore(storeId);
        emit DynamicStoreDeleted(storeId, success);
        return success;
    }
    
    /// @notice Wrapper for setDynamicStoreValue
    function testSetDynamicStoreValue(
        uint256 storeId,
        address address_,
        bool value
    ) external returns (bool) {
        bool success = precompile.setDynamicStoreValue(storeId, address_, value);
        emit DynamicStoreValueSet(storeId, address_, value, success);
        return success;
    }
    
    // ============ Address List Methods ============
    
    /// @notice Wrapper for createAddressLists
    function testCreateAddressLists(
        AddressListInput[] calldata addressLists
    ) external returns (bool) {
        bool success = precompile.createAddressLists(addressLists);
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
        bool success = precompile.castVote(
            collectionId,
            approvalLevel,
            approverAddress,
            approvalId,
            proposalId,
            yesWeight
        );
        emit VoteCast(collectionId, proposalId, success);
        return success;
    }
    
    // ============ Query Methods (View) ============
    
    /// @notice Wrapper for getCollection
    function testGetCollection(
        uint256 collectionId
    ) external view returns (bytes memory) {
        return precompile.getCollection(collectionId);
    }
    
    /// @notice Wrapper for getBalance
    function testGetBalance(
        uint256 collectionId,
        address address_
    ) external view returns (bytes memory) {
        return precompile.getBalance(collectionId, address_);
    }
    
    /// @notice Wrapper for getBalanceAmount
    function testGetBalanceAmount(
        uint256 collectionId,
        address address_,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256) {
        return precompile.getBalanceAmount(collectionId, address_, tokenIds, ownershipTimes);
    }
    
    /// @notice Wrapper for getTotalSupply
    function testGetTotalSupply(
        uint256 collectionId,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256) {
        return precompile.getTotalSupply(collectionId, tokenIds, ownershipTimes);
    }
    
    /// @notice Wrapper for getAddressList
    function testGetAddressList(
        string calldata listId
    ) external view returns (bytes memory) {
        return precompile.getAddressList(listId);
    }
    
    /// @notice Wrapper for getDynamicStore
    function testGetDynamicStore(
        uint256 storeId
    ) external view returns (bytes memory) {
        return precompile.getDynamicStore(storeId);
    }
    
    /// @notice Wrapper for getDynamicStoreValue
    function testGetDynamicStoreValue(
        uint256 storeId,
        address address_
    ) external view returns (bool) {
        // Note: getDynamicStoreValue returns bytes, but we need bool
        // This is a simplified wrapper - full implementation would decode bytes
        bytes memory result = precompile.getDynamicStoreValue(storeId, address_);
        // For now, return false if empty, true if not empty (simplified)
        return result.length > 0;
    }
}

