# Extended Examples

Extended examples beyond the basic getting started guide. See `contracts/examples/` for complete contract implementations.

## Table of Contents

- [ERC-3643 Compliant Token](#erc-3643-compliant-token)
- [Time-Limited Access Tokens](#time-limited-access-tokens)
- [Multi-Registry Compliance](#multi-registry-compliance)
- [Voting Token](#voting-token)
- [Subscription Token](#subscription-token)

## ERC-3643 Compliant Token

A security token with KYC/AML compliance using dynamic stores.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../libraries/TokenizationWrappers.sol";
import "../libraries/TokenizationHelpers.sol";
import "../libraries/TokenizationBuilders.sol";

contract ERC3643Token {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public collectionId;
    uint256 public kycRegistryId;
    address public complianceOfficer;
    
    mapping(address => bool) public frozen;
    
    event Transfer(address indexed from, address indexed to, uint256 amount);
    event Freeze(address indexed account);
    event Unfreeze(address indexed account);
    
    constructor(
        string memory name,
        string memory symbol,
        address _complianceOfficer
    ) {
        complianceOfficer = _complianceOfficer;
        
        // Create KYC registry
        kycRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION,
            false,  // Default: not KYC'd
            string(abi.encodePacked("ipfs://kyc-", name)),
            "KYC Registry"
        );
        
        // Create collection
        TokenizationBuilders.CollectionBuilder memory builder = 
            TokenizationBuilders.newCollection();
        
        // Note: Use uint64.max for ownership times/token IDs (BitBadges internal limit)
        builder = builder.withValidTokenIdRange(1, type(uint64).max);
        builder = builder.withManager(TokenizationHelpers.addressToCosmosString(msg.sender));
        
        TokenizationTypes.CollectionMetadata memory metadata = 
            TokenizationHelpers.createCollectionMetadata(
                string(abi.encodePacked("ipfs://", name)),
                string(abi.encodePacked('{"name":"', name, '","symbol":"', symbol, '"}'))
            );
        builder = builder.withMetadata(metadata);
        
        string[] memory standards = new string[](1);
        standards[0] = "ERC-3643";
        builder = builder.withStandards(standards);
        
        // Auto-approve self-initiated transfers
        builder = builder.withDefaultBalancesFromFlags(true, true, false);
        
        collectionId = TOKENIZATION.createCollection(builder.build());
    }
    
    modifier onlyComplianceOfficer() {
        require(msg.sender == complianceOfficer, "Not compliance officer");
        _;
    }
    
    function setKYCStatus(address user, bool isKYCd) external onlyComplianceOfficer {
        TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION,
            kycRegistryId,
            user,
            isKYCd
        );
    }
    
    function freeze(address account) external onlyComplianceOfficer {
        frozen[account] = true;
        emit Freeze(account);
    }
    
    function unfreeze(address account) external onlyComplianceOfficer {
        frozen[account] = false;
        emit Unfreeze(account);
    }
    
    function transfer(address to, uint256 amount) external {
        require(!frozen[msg.sender], "Sender is frozen");
        require(!frozen[to], "Recipient is frozen");
        
        // In production, check KYC status here
        // (simplified for example)
        
        TokenizationWrappers.transferSingleToken(
            TOKENIZATION,
            collectionId,
            to,
            amount,
            1  // Assuming single token ID
        );
        
        emit Transfer(msg.sender, to, amount);
    }
    
    function balanceOf(address account) external view returns (uint256) {
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
        
        return TokenizationWrappers.getBalanceAmount(
            TOKENIZATION,
            collectionId,
            account,
            tokenIds,
            ownershipTimes
        );
    }
}
```

## Time-Limited Access Tokens

Tokens that expire after a certain time period.

```solidity
contract TimeLimitedAccess {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public collectionId;
    mapping(address => uint256) public accessExpiry;
    
    function issueAccessToken(
        address user,
        uint256 duration
    ) external {
        uint256 expirationTime = block.timestamp + duration;
        accessExpiry[user] = expirationTime;
        
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createUintRange(
            block.timestamp,
            expirationTime
        );
        
        address[] memory recipients = new address[](1);
        recipients[0] = user;
        
        TokenizationWrappers.transferTokens(
            TOKENIZATION,
            collectionId,
            recipients,
            1,
            tokenIds,
            ownershipTimes
        );
    }
    
    function hasAccess(address user) external view returns (bool) {
        if (accessExpiry[user] == 0) return false;
        if (block.timestamp > accessExpiry[user]) return false;
        
        // Also check actual balance (token may have been transferred)
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createTimePoint(block.timestamp);
        
        uint256 balance = TokenizationWrappers.getBalanceAmount(
            TOKENIZATION,
            collectionId,
            user,
            tokenIds,
            ownershipTimes
        );
        
        return balance > 0;
    }
}
```

## Multi-Registry Compliance

Using multiple dynamic stores for different compliance requirements.

```solidity
contract MultiComplianceToken {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000001001);
    
    uint256 public collectionId;
    uint256 public kycRegistryId;
    uint256 public accreditedInvestorRegistryId;
    uint256 public jurisdictionRegistryId;
    
    address public complianceOfficer;
    
    constructor() {
        complianceOfficer = msg.sender;
        
        // Create all registries
        kycRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://kyc", "KYC"
        );
        accreditedInvestorRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://accredited", "Accredited Investor"
        );
        jurisdictionRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://jurisdiction", "Jurisdiction"
        );
        
        // Create collection...
    }
    
    function setComplianceStatus(
        address user,
        bool kyc,
        bool accredited,
        string memory jurisdiction
    ) external {
        require(msg.sender == complianceOfficer, "Not authorized");
        
        TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION, kycRegistryId, user, kyc
        );
        TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION, accreditedInvestorRegistryId, user, accredited
        );
        // Jurisdiction would need a different approach (not boolean)
    }
    
    function canTransfer(address from, address to) external view returns (bool) {
        // Check all compliance requirements
        // In production, decode registry values properly
        return true; // Simplified
    }
}
```

## Voting Token

Tokens used for voting with time-bound ownership.

```solidity
contract VotingToken {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public collectionId;
    
    struct Proposal {
        string description;
        uint256 deadline;
        mapping(address => bool) hasVoted;
        uint256 yesVotes;
        uint256 noVotes;
    }
    
    mapping(uint256 => Proposal) public proposals;
    uint256 public proposalCount;
    
    function createProposal(
        string memory description,
        uint256 votingDuration
    ) external returns (uint256) {
        uint256 proposalId = proposalCount++;
        proposals[proposalId].description = description;
        proposals[proposalId].deadline = block.timestamp + votingDuration;
        return proposalId;
    }
    
    function vote(uint256 proposalId, bool support) external {
        Proposal storage proposal = proposals[proposalId];
        require(block.timestamp <= proposal.deadline, "Voting closed");
        require(!proposal.hasVoted[msg.sender], "Already voted");
        
        // Check voting power (balance at proposal creation time)
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createTimePoint(block.timestamp);
        
        uint256 votingPower = TokenizationWrappers.getBalanceAmount(
            TOKENIZATION,
            collectionId,
            msg.sender,
            tokenIds,
            ownershipTimes
        );
        
        require(votingPower > 0, "No voting power");
        
        proposal.hasVoted[msg.sender] = true;
        if (support) {
            proposal.yesVotes += votingPower;
        } else {
            proposal.noVotes += votingPower;
        }
    }
}
```

## Subscription Token

Tokens that represent active subscriptions with recurring ownership times.

```solidity
contract SubscriptionToken {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public collectionId;
    
    struct Subscription {
        address subscriber;
        uint256 startTime;
        uint256 duration;
        bool active;
    }
    
    mapping(address => Subscription) public subscriptions;
    
    function subscribe(uint256 duration) external payable {
        require(msg.value >= getSubscriptionPrice(duration), "Insufficient payment");
        require(!subscriptions[msg.sender].active, "Already subscribed");
        
        uint256 startTime = block.timestamp;
        uint256 endTime = startTime + duration;
        
        subscriptions[msg.sender] = Subscription({
            subscriber: msg.sender,
            startTime: startTime,
            duration: duration,
            active: true
        });
        
        // Issue subscription token with time-bound ownership
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createUintRange(startTime, endTime);
        
        address[] memory recipients = new address[](1);
        recipients[0] = msg.sender;
        
        TokenizationWrappers.transferTokens(
            TOKENIZATION,
            collectionId,
            recipients,
            1,
            tokenIds,
            ownershipTimes
        );
    }
    
    function isSubscriptionActive(address user) external view returns (bool) {
        Subscription memory sub = subscriptions[user];
        if (!sub.active) return false;
        if (block.timestamp > sub.startTime + sub.duration) return false;
        
        // Check actual token ownership
        TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
        
        TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createTimePoint(block.timestamp);
        
        uint256 balance = TokenizationWrappers.getBalanceAmount(
            TOKENIZATION,
            collectionId,
            user,
            tokenIds,
            ownershipTimes
        );
        
        return balance > 0;
    }
    
    function getSubscriptionPrice(uint256 duration) public pure returns (uint256) {
        // Pricing logic
        return duration * 1e15; // Example: 0.001 ETH per second
    }
}
```

## Additional Resources

- See `contracts/examples/` for complete, production-ready examples:
  - `PrivateEquityToken.sol` - Private equity fund tokens
  - `RealEstateSecurityToken.sol` - Real estate tokenization
  - `CarbonCreditToken.sol` - Carbon credit tracking
  - `TwoFactorSecurityToken.sol` - 2FA security tokens

- Review [Patterns](./PATTERNS.md) for more common patterns
- Check [Best Practices](./BEST_PRACTICES.md) for security considerations


















