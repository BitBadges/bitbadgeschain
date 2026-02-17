// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/// @title PrecompileCollectionTestContract
/// @notice Test contract for collection management precompile methods
/// @dev Split from PrecompileTestContract to stay under EVM size limits
contract PrecompileCollectionTestContract {
    ITokenizationPrecompile constant precompile =
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    // ============ Events ============

    event CollectionCreated(uint256 indexed collectionId);

    event CollectionUpdated(uint256 indexed collectionId, bool success);

    event CollectionDeleted(uint256 indexed collectionId, bool success);

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

    event AddressListsCreated(uint256 numLists, bool success);

    event VoteCast(
        uint256 indexed collectionId,
        string proposalId,
        bool success
    );

    // ============ Collection Methods ============

    /// @notice Simplified wrapper for createCollection with basic parameters
    function testCreateCollectionSimple(
        uint256 tokenIdStart,
        uint256 tokenIdEnd,
        string calldata manager
    ) external returns (uint256) {
        // Build JSON inline to avoid stack issues
        string memory createJson = _buildCreateCollectionJson(tokenIdStart, tokenIdEnd, manager);
        uint256 collectionId = precompile.createCollection(createJson);
        emit CollectionCreated(collectionId);
        return collectionId;
    }

    function _buildCreateCollectionJson(
        uint256 tokenIdStart,
        uint256 tokenIdEnd,
        string calldata manager
    ) internal pure returns (string memory) {
        return string(
            abi.encodePacked(
                '{"validTokenIds":[{"start":"',
                TokenizationJSONHelpers.uintToString(tokenIdStart),
                '","end":"',
                TokenizationJSONHelpers.uintToString(tokenIdEnd),
                '"}],"manager":"',
                manager,
                '","defaultBalances":{"autoApproveSelfInitiatedOutgoingTransfers":true,"autoApproveSelfInitiatedIncomingTransfers":true}}'
            )
        );
    }

    /// @notice Wrapper for deleteCollection
    function testDeleteCollection(uint256 collectionId) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteCollectionJSON(collectionId);
        bool success = precompile.deleteCollection(deleteJson);
        emit CollectionDeleted(collectionId, success);
        return success;
    }

    // ============ Approval Methods ============

    /// @notice Simplified wrapper for setIncomingApproval
    function testSetIncomingApprovalSimple(
        uint256 collectionId,
        string calldata approvalId,
        string calldata fromListId
    ) external returns (bool) {
        string memory msgJson = _buildIncomingApprovalJson(collectionId, approvalId, fromListId);
        bool success = precompile.setIncomingApproval(msgJson);
        emit ApprovalSet(collectionId, approvalId, true, success);
        return success;
    }

    function _buildIncomingApprovalJson(
        uint256 collectionId,
        string calldata approvalId,
        string calldata fromListId
    ) internal pure returns (string memory) {
        return string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","approval":{"approvalId":"',
                approvalId,
                '","fromListId":"',
                fromListId,
                '","initiatedByListId":"All","transferTimes":[{"start":"1","end":"18446744073709551615"}],"tokenIds":[{"start":"1","end":"18446744073709551615"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}}'
            )
        );
    }

    /// @notice Simplified wrapper for setOutgoingApproval
    function testSetOutgoingApprovalSimple(
        uint256 collectionId,
        string calldata approvalId,
        string calldata toListId
    ) external returns (bool) {
        string memory msgJson = _buildOutgoingApprovalJson(collectionId, approvalId, toListId);
        bool success = precompile.setOutgoingApproval(msgJson);
        emit ApprovalSet(collectionId, approvalId, false, success);
        return success;
    }

    function _buildOutgoingApprovalJson(
        uint256 collectionId,
        string calldata approvalId,
        string calldata toListId
    ) internal pure returns (string memory) {
        return string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","approval":{"approvalId":"',
                approvalId,
                '","toListId":"',
                toListId,
                '","initiatedByListId":"All","transferTimes":[{"start":"1","end":"18446744073709551615"}],"tokenIds":[{"start":"1","end":"18446744073709551615"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}}'
            )
        );
    }

    /// @notice Wrapper for deleteIncomingApproval
    function testDeleteIncomingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteIncomingApprovalJSON(
            collectionId,
            approvalId
        );
        bool success = precompile.deleteIncomingApproval(deleteJson);
        emit ApprovalDeleted(collectionId, approvalId, true, success);
        return success;
    }

    /// @notice Wrapper for deleteOutgoingApproval
    function testDeleteOutgoingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteOutgoingApprovalJSON(
            collectionId,
            approvalId
        );
        bool success = precompile.deleteOutgoingApproval(deleteJson);
        emit ApprovalDeleted(collectionId, approvalId, false, success);
        return success;
    }

    // ============ Address List Methods ============

    /// @notice Simplified wrapper for createAddressLists - single address version
    function testCreateAddressListSingle(
        string calldata listId,
        string calldata singleAddress,
        bool whitelist
    ) external returns (bool) {
        string memory createJson = _buildAddressListJson(listId, singleAddress, whitelist);
        bool success = precompile.createAddressLists(createJson);
        emit AddressListsCreated(1, success);
        return success;
    }

    function _buildAddressListJson(
        string calldata listId,
        string calldata singleAddress,
        bool whitelist
    ) internal pure returns (string memory) {
        return string(
            abi.encodePacked(
                '{"addressLists":[{"listId":"',
                listId,
                '","addresses":["',
                singleAddress,
                '"],"whitelist":',
                whitelist ? "true" : "false",
                "}]}"
            )
        );
    }

    // ============ Vote Methods ============

    /// @notice Wrapper for castVote - simplified
    function testCastVote(
        uint256 collectionId,
        string calldata approvalId,
        string calldata proposalId,
        uint256 yesWeight
    ) external returns (bool) {
        string memory castVoteJson = _buildCastVoteJson(collectionId, approvalId, proposalId, yesWeight);
        bool success = precompile.castVote(castVoteJson);
        emit VoteCast(collectionId, proposalId, success);
        return success;
    }

    function _buildCastVoteJson(
        uint256 collectionId,
        string calldata approvalId,
        string calldata proposalId,
        uint256 yesWeight
    ) internal pure returns (string memory) {
        return string(
            abi.encodePacked(
                '{"collectionId":"',
                TokenizationJSONHelpers.uintToString(collectionId),
                '","approvalLevel":"collection","approverAddress":"","approvalId":"',
                approvalId,
                '","proposalId":"',
                proposalId,
                '","yesWeight":"',
                TokenizationJSONHelpers.uintToString(yesWeight),
                '"}'
            )
        );
    }
}
