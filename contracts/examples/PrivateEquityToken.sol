// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";

/**
 * @title PrivateEquityToken
 * @notice ERC-3643 compliant token for private equity fund interests
 * @dev Demonstrates advanced features:
 * - Lock-up periods using ownership time ranges
 * - Maximum ownership limits per investor
 * - Qualified purchaser requirements
 * - Transfer approval workflows via collection approvals
 * - Dividend distribution tracking
 */
contract PrivateEquityToken {
    ITokenizationPrecompile constant TOKENIZATION = ITokenizationPrecompile(0x0000000000000000000000000000000000000808);

    // ============ State Variables ============

    uint256 public collectionId;

    // Compliance registries (Dynamic Stores)
    uint256 public qualifiedPurchaserRegistryId;   // QP status (>$5M investments)
    uint256 public accreditedInvestorRegistryId;   // Accredited investor
    uint256 public blacklistRegistryId;            // Prohibited investors
    uint256 public lockUpRegistryId;               // Custom lock-up status per investor

    // Fund metadata
    string public name;
    string public symbol;
    string public fundType;             // e.g., "Buyout", "Venture", "Growth"
    uint256 public fundSize;            // Target fund size in USD
    uint256 public minimumInvestment;   // Minimum investment amount

    // Fund timing
    uint256 public fundingPeriodEnd;    // End of commitment period
    uint256 public lockUpPeriod;        // Standard lock-up in seconds
    uint256 public fundTermination;     // Fund termination date

    // Ownership limits (basis points, 10000 = 100%)
    uint256 public maxOwnershipBps = 2500;  // 25% max per investor

    // Roles
    address public generalPartner;
    address public fundAdmin;
    address public transferAgent;

    // Investor tracking
    mapping(address => InvestorInfo) public investors;
    address[] public investorList;

    struct InvestorInfo {
        uint256 commitmentAmount;
        uint256 paidInAmount;
        uint256 lockUpEnd;
        bool isLP;              // Limited Partner status
        string investorClass;   // e.g., "Class A", "Class B"
    }

    // Dividend tracking
    uint256 public totalDistributed;
    mapping(address => uint256) public claimedDistributions;

    // Token ID structure
    UintRange[] private _tokenIds;

    // ============ Events ============

    event InvestorOnboarded(address indexed investor, uint256 commitment, string investorClass);
    event CapitalCall(uint256 callNumber, uint256 percentage);
    event DistributionDeclared(uint256 amount, uint256 perToken);
    event TransferRequested(address indexed from, address indexed to, uint256 amount, bytes32 requestId);
    event TransferApproved(bytes32 indexed requestId, address indexed approver);
    event TransferExecuted(address indexed from, address indexed to, uint256 amount);
    event LockUpExtended(address indexed investor, uint256 newLockUpEnd);

    // ============ Modifiers ============

    modifier onlyGP() {
        require(msg.sender == generalPartner, "Only GP");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == fundAdmin || msg.sender == generalPartner, "Only admin");
        _;
    }

    modifier onlyTransferAgent() {
        require(
            msg.sender == transferAgent ||
            msg.sender == fundAdmin ||
            msg.sender == generalPartner,
            "Only transfer agent"
        );
        _;
    }

    // ============ Constructor ============

    /**
     * @notice Deploy a private equity fund token
     * @param _name Fund name
     * @param _symbol Token symbol
     * @param _fundType Type of PE fund
     * @param _fundSize Target fund size
     * @param _minimumInvestment Minimum commitment
     * @param _lockUpPeriod Lock-up period in seconds
     * @param _fundAdmin Fund administrator address
     * @param _transferAgent Transfer agent address
     */
    constructor(
        string memory _name,
        string memory _symbol,
        string memory _fundType,
        uint256 _fundSize,
        uint256 _minimumInvestment,
        uint256 _lockUpPeriod,
        address _fundAdmin,
        address _transferAgent
    ) {
        generalPartner = msg.sender;
        fundAdmin = _fundAdmin;
        transferAgent = _transferAgent;

        name = _name;
        symbol = _symbol;
        fundType = _fundType;
        fundSize = _fundSize;
        minimumInvestment = _minimumInvestment;
        lockUpPeriod = _lockUpPeriod;

        // Default fund periods (can be updated)
        fundingPeriodEnd = block.timestamp + 365 days;
        fundTermination = block.timestamp + 10 * 365 days;

        // Token ID 1 for LP interests
        _tokenIds = new UintRange[](1);
        _tokenIds[0] = UintRange(1, 1);

        // Create compliance registries
        qualifiedPurchaserRegistryId = TOKENIZATION.createDynamicStore(
            false,
            "ipfs://qp-registry",
            "{\"type\":\"qualified-purchaser\",\"threshold\":5000000}"
        );

        accreditedInvestorRegistryId = TOKENIZATION.createDynamicStore(
            false,
            "ipfs://accredited-registry",
            "{\"type\":\"accredited-investor\"}"
        );

        blacklistRegistryId = TOKENIZATION.createDynamicStore(
            false,
            "ipfs://blacklist-registry",
            "{\"type\":\"prohibited-persons\"}"
        );

        lockUpRegistryId = TOKENIZATION.createDynamicStore(
            true,  // Default: in lock-up
            "ipfs://lockup-registry",
            "{\"type\":\"lock-up-status\"}"
        );
    }

    /**
     * @notice Initialize the fund token collection
     */
    function initializeCollection() external onlyGP {
        require(collectionId == 0, "Already initialized");

        // Ownership times based on fund lifecycle
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = UintRange(block.timestamp, fundTermination);

        TokenizationTypes.Balance[] memory initialBalances = new TokenizationTypes.Balance[](1);
        initialBalances[0] = TokenizationTypes.Balance({
            amount: fundSize / 1000,  // 1 token = $1000 of commitment
            ownershipTimes: ownershipTimes,
            tokenIds: _tokenIds
        });

        TokenizationTypes.UserBalanceStore memory defaultBalances;
        defaultBalances.balances = initialBalances;
        // Require explicit approvals for all transfers
        defaultBalances.autoApproveSelfInitiatedOutgoingTransfers = false;
        defaultBalances.autoApproveSelfInitiatedIncomingTransfers = false;

        TokenizationTypes.MsgCreateCollection memory createMsg;
        createMsg.defaultBalances = defaultBalances;
        createMsg.validTokenIds = _tokenIds;
        createMsg.manager = _addressToString(address(this));
        createMsg.collectionMetadata = TokenizationTypes.CollectionMetadata({
            uri: "ipfs://pe-fund-metadata",
            customData: string(abi.encodePacked(
                "{\"name\":\"", name,
                "\",\"fundType\":\"", fundType,
                "\",\"compliance\":\"ERC-3643\",\"reg\":\"Reg D 506(c)\"}"
            ))
        });

        string[] memory standards = new string[](3);
        standards[0] = "ERC-3643";
        standards[1] = "Private Equity";
        standards[2] = "Reg D 506(c)";
        createMsg.standards = standards;

        collectionId = TOKENIZATION.createCollection(createMsg);
    }

    // ============ Investor Onboarding ============

    /**
     * @notice Onboard a new LP investor
     * @param investor Investor address
     * @param commitmentAmount Capital commitment in USD
     * @param investorClass LP class (A, B, etc.)
     */
    function onboardInvestor(
        address investor,
        uint256 commitmentAmount,
        string calldata investorClass
    ) external onlyAdmin {
        require(commitmentAmount >= minimumInvestment, "Below minimum");
        require(block.timestamp < fundingPeriodEnd, "Funding period ended");
        require(!isBlacklisted(investor), "Investor blacklisted");
        require(
            isQualifiedPurchaser(investor) || isAccreditedInvestor(investor),
            "Not qualified"
        );

        // Check ownership limit
        uint256 newOwnership = (commitmentAmount * 10000) / fundSize;
        require(newOwnership <= maxOwnershipBps, "Exceeds ownership limit");

        // Set lock-up end for this investor
        uint256 investorLockUpEnd = block.timestamp + lockUpPeriod;

        investors[investor] = InvestorInfo({
            commitmentAmount: commitmentAmount,
            paidInAmount: 0,
            lockUpEnd: investorLockUpEnd,
            isLP: true,
            investorClass: investorClass
        });

        investorList.push(investor);

        // Mark as in lock-up
        TOKENIZATION.setDynamicStoreValue(lockUpRegistryId, investor, true);

        emit InvestorOnboarded(investor, commitmentAmount, investorClass);
    }

    /**
     * @notice Process capital call (issue tokens for paid-in capital)
     * @param investor LP address
     * @param amount Amount called and paid
     */
    function processCapitalCall(
        address investor,
        uint256 amount
    ) external onlyAdmin {
        require(investors[investor].isLP, "Not an LP");
        require(
            investors[investor].paidInAmount + amount <= investors[investor].commitmentAmount,
            "Exceeds commitment"
        );

        investors[investor].paidInAmount += amount;

        // Issue tokens proportional to paid-in capital
        uint256 tokensToIssue = amount / 1000;  // 1 token = $1000

        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = UintRange(block.timestamp, fundTermination);

        address[] memory recipients = new address[](1);
        recipients[0] = investor;

        TOKENIZATION.transferTokens(collectionId, recipients, tokensToIssue, _tokenIds, ownershipTimes);
    }

    // ============ Transfer Functions with Approval Workflow ============

    /**
     * @notice Check comprehensive transfer eligibility
     */
    function canTransfer(
        address from,
        address to,
        uint256 amount
    ) public view returns (bool, string memory) {
        // Check blacklist
        if (isBlacklisted(from)) return (false, "Sender blacklisted");
        if (isBlacklisted(to)) return (false, "Recipient blacklisted");

        // Check lock-up
        if (isInLockUp(from)) return (false, "Sender in lock-up");

        // Check recipient qualification
        if (!isQualifiedPurchaser(to) && !isAccreditedInvestor(to)) {
            return (false, "Recipient not qualified");
        }

        // Check ownership limits for recipient
        uint256 recipientBalance = this.balanceOf(to);
        uint256 totalAfterTransfer = recipientBalance + amount;
        uint256 totalSupply = this.totalSupply();

        if (totalSupply > 0) {
            uint256 newOwnershipBps = (totalAfterTransfer * 10000) / totalSupply;
            if (newOwnershipBps > maxOwnershipBps) {
                return (false, "Would exceed ownership limit");
            }
        }

        return (true, "");
    }

    /**
     * @notice Transfer LP interests (requires transfer agent approval)
     * @dev In practice, secondary transfers typically require GP consent
     */
    function transferWithApproval(
        address to,
        uint256 amount
    ) external returns (bool) {
        (bool eligible, string memory reason) = canTransfer(msg.sender, to, amount);
        require(eligible, reason);

        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = UintRange(block.timestamp, fundTermination);

        address[] memory recipients = new address[](1);
        recipients[0] = to;

        bool success = TOKENIZATION.transferTokens(collectionId, recipients, amount, _tokenIds, ownershipTimes);

        if (success) {
            // Update investor records
            if (!investors[to].isLP) {
                investors[to] = InvestorInfo({
                    commitmentAmount: 0,  // Secondary purchaser
                    paidInAmount: 0,
                    lockUpEnd: block.timestamp + lockUpPeriod,
                    isLP: true,
                    investorClass: "Secondary"
                });
                investorList.push(to);
                TOKENIZATION.setDynamicStoreValue(lockUpRegistryId, to, true);
            }

            emit TransferExecuted(msg.sender, to, amount);
        }
        return success;
    }

    // ============ Lock-Up Management ============

    /**
     * @notice Release an investor from lock-up
     */
    function releaseLockUp(address investor) external onlyAdmin {
        require(investors[investor].isLP, "Not an LP");
        require(block.timestamp >= investors[investor].lockUpEnd, "Lock-up not expired");

        TOKENIZATION.setDynamicStoreValue(lockUpRegistryId, investor, false);
    }

    /**
     * @notice Extend lock-up for an investor (e.g., as condition of co-investment)
     */
    function extendLockUp(address investor, uint256 newLockUpEnd) external onlyGP {
        require(newLockUpEnd > investors[investor].lockUpEnd, "Must extend");

        investors[investor].lockUpEnd = newLockUpEnd;
        TOKENIZATION.setDynamicStoreValue(lockUpRegistryId, investor, true);

        emit LockUpExtended(investor, newLockUpEnd);
    }

    /**
     * @notice Check if investor is in lock-up
     */
    function isInLockUp(address investor) public view returns (bool) {
        // Check both time-based and manual lock-up
        if (block.timestamp < investors[investor].lockUpEnd) return true;

        bytes memory result = TOKENIZATION.getDynamicStoreValue(lockUpRegistryId, investor);
        if (result.length == 0) return true;  // Default: locked
        return abi.decode(result, (bool));
    }

    // ============ Compliance Registry Functions ============

    function setQualifiedPurchaser(address investor, bool status) external onlyAdmin {
        TOKENIZATION.setDynamicStoreValue(qualifiedPurchaserRegistryId, investor, status);
    }

    function setAccreditedInvestor(address investor, bool status) external onlyAdmin {
        TOKENIZATION.setDynamicStoreValue(accreditedInvestorRegistryId, investor, status);
    }

    function setBlacklisted(address investor, bool status) external onlyAdmin {
        TOKENIZATION.setDynamicStoreValue(blacklistRegistryId, investor, status);
    }

    function isQualifiedPurchaser(address investor) public view returns (bool) {
        bytes memory result = TOKENIZATION.getDynamicStoreValue(qualifiedPurchaserRegistryId, investor);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    function isAccreditedInvestor(address investor) public view returns (bool) {
        bytes memory result = TOKENIZATION.getDynamicStoreValue(accreditedInvestorRegistryId, investor);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    function isBlacklisted(address investor) public view returns (bool) {
        bytes memory result = TOKENIZATION.getDynamicStoreValue(blacklistRegistryId, investor);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    // ============ View Functions ============

    function balanceOf(address account) external view returns (uint256) {
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = UintRange(1, type(uint256).max);
        return TOKENIZATION.getBalanceAmount(collectionId, account, _tokenIds, ownershipTimes);
    }

    function totalSupply() external view returns (uint256) {
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = UintRange(1, type(uint256).max);
        return TOKENIZATION.getTotalSupply(collectionId, _tokenIds, ownershipTimes);
    }

    function getInvestorInfo(address investor) external view returns (
        uint256 commitment,
        uint256 paidIn,
        uint256 lockUpEnd,
        bool isLP,
        string memory investorClass,
        bool inLockUp
    ) {
        InvestorInfo memory info = investors[investor];
        return (
            info.commitmentAmount,
            info.paidInAmount,
            info.lockUpEnd,
            info.isLP,
            info.investorClass,
            isInLockUp(investor)
        );
    }

    function getInvestorCount() external view returns (uint256) {
        return investorList.length;
    }

    // ============ Fund Management ============

    function updateMaxOwnership(uint256 newMaxBps) external onlyGP {
        require(newMaxBps <= 10000, "Invalid bps");
        maxOwnershipBps = newMaxBps;
    }

    function updateFundDates(
        uint256 _fundingPeriodEnd,
        uint256 _fundTermination
    ) external onlyGP {
        fundingPeriodEnd = _fundingPeriodEnd;
        fundTermination = _fundTermination;
    }

    // ============ Internal ============

    function _addressToString(address addr) internal pure returns (string memory) {
        bytes memory alphabet = "0123456789abcdef";
        bytes memory data = abi.encodePacked(addr);
        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < data.length; i++) {
            str[2 + i * 2] = alphabet[uint8(data[i] >> 4)];
            str[3 + i * 2] = alphabet[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }
}
