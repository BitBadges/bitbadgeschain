// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct UintRange {
    uint256 start;
    uint256 end;
}

// Simplified approval structure for precompile
struct UserApproval {
    string approvalId;
    string listId; // toListId for outgoing, fromListId for incoming
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
}

interface IBadgesPrecompile {
    // ============ Transactions ============
    
    function transferTokens(
        uint256 collectionId,
        address[] calldata toAddresses,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool);
    
    function setIncomingApproval(
        uint256 collectionId,
        UserApproval calldata approval
    ) external returns (bool);
    
    function setOutgoingApproval(
        uint256 collectionId,
        UserApproval calldata approval
    ) external returns (bool);
    
    // ============ Queries ============
    
    function getCollection(
        uint256 collectionId
    ) external view returns (bytes memory);
    
    function getBalance(
        uint256 collectionId,
        address userAddress
    ) external view returns (bytes memory);
    
    function getAddressList(
        string calldata listId
    ) external view returns (bytes memory);
    
    function getApprovalTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata amountTrackerId,
        string calldata trackerType,
        address approvedAddress,
        string calldata approvalId
    ) external view returns (bytes memory);
    
    function getChallengeTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata challengeTrackerId,
        uint256 leafIndex,
        string calldata approvalId
    ) external view returns (bytes memory);
    
    function getETHSignatureTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata challengeTrackerId,
        string calldata signature
    ) external view returns (bytes memory);
    
    function getDynamicStore(
        uint256 storeId
    ) external view returns (bytes memory);
    
    function getDynamicStoreValue(
        uint256 storeId,
        address userAddress
    ) external view returns (bytes memory);
    
    function getWrappableBalances(
        string calldata denom,
        address userAddress
    ) external view returns (uint256);
    
    function isAddressReservedProtocol(
        address addr
    ) external view returns (bool);
    
    function getAllReservedProtocolAddresses(
    ) external view returns (address[] memory);
    
    function getVote(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata proposalId,
        address voterAddress
    ) external view returns (bytes memory);
    
    function getVotes(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata proposalId
    ) external view returns (bytes memory);
    
    function params(
    ) external view returns (bytes memory);
    
    function getBalanceAmount(
        uint256 collectionId,
        address userAddress,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
    
    function getTotalSupply(
        uint256 collectionId,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
}
