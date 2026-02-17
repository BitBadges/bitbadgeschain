// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/**
 * @title CarbonCreditToken
 * @notice ERC-3643 compliant token for verified carbon credits
 * @dev Uses BitBadges precompile with full mapping to underlying module
 *
 * This example demonstrates:
 * - Using BitBadges collections for carbon credit vintages (different token IDs per vintage)
 * - Dynamic Stores for carbon registry verification
 * - Ownership time ranges for credit expiration
 * - Retirement tracking via burning to a sink address
 */
contract CarbonCreditToken {
    ITokenizationPrecompile constant TOKENIZATION = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    // ============ State Variables ============

    uint256 public collectionId;

    // Registry stores
    uint256 public verifiedBuyersRegistryId;    // Verified carbon credit buyers
    uint256 public verifiedSellersRegistryId;   // Verified project developers
    uint256 public retiredCreditsRegistryId;    // Track retired credits per address

    // Metadata
    string public name;
    string public symbol;
    string public standard;     // e.g., "VCS", "Gold Standard", "ACR"
    string public projectId;    // Registry project identifier

    // Roles
    address public registryOperator;
    address public verifier;

    // Vintage tracking (token ID = vintage year)
    mapping(uint256 => VintageInfo) public vintages;

    struct VintageInfo {
        uint256 year;
        uint256 totalIssued;
        uint256 totalRetired;
        uint256 expirationTime;  // Credits expire and can't be used after this
        string verificationUri;
    }

    // Retirement tracking
    mapping(address => uint256) public retiredByAddress;
    uint256 public totalRetired;

    // Sink address for retired credits
    address public constant RETIREMENT_SINK = address(0xdead);

    // ============ Events ============

    event CreditIssued(uint256 indexed vintage, address indexed to, uint256 amount);
    event CreditTransferred(address indexed from, address indexed to, uint256 indexed vintage, uint256 amount);
    event CreditRetired(address indexed retiree, uint256 indexed vintage, uint256 amount, string reason);
    event VintageCreated(uint256 indexed year, uint256 expirationTime);
    event BuyerVerified(address indexed buyer);
    event SellerVerified(address indexed seller);

    // ============ Modifiers ============

    modifier onlyOperator() {
        require(msg.sender == registryOperator, "Only registry operator");
        _;
    }

    modifier onlyVerifier() {
        require(msg.sender == verifier || msg.sender == registryOperator, "Only verifier");
        _;
    }

    // ============ Constructor ============

    /**
     * @notice Deploy a carbon credit token for a specific project
     * @param _name Project name
     * @param _symbol Token symbol
     * @param _standard Carbon standard (VCS, Gold Standard, etc.)
     * @param _projectId Registry project ID
     * @param _verifier Address of the verification agent
     */
    constructor(
        string memory _name,
        string memory _symbol,
        string memory _standard,
        string memory _projectId,
        address _verifier
    ) {
        registryOperator = msg.sender;
        verifier = _verifier;
        name = _name;
        symbol = _symbol;
        standard = _standard;
        projectId = _projectId;

        // Create verification registries
        string memory buyersJson = TokenizationJSONHelpers.createDynamicStoreJSON(
            false,
            "ipfs://carbon-buyers-registry",
            "{\"type\":\"verified-buyers\"}"
        );
        verifiedBuyersRegistryId = TOKENIZATION.createDynamicStore(buyersJson);

        string memory sellersJson = TokenizationJSONHelpers.createDynamicStoreJSON(
            false,
            "ipfs://carbon-sellers-registry",
            "{\"type\":\"verified-sellers\"}"
        );
        verifiedSellersRegistryId = TOKENIZATION.createDynamicStore(sellersJson);

        string memory retiredJson = TokenizationJSONHelpers.createDynamicStoreJSON(
            false,
            "ipfs://retired-credits-tracker",
            "{\"type\":\"retirement-tracker\"}"
        );
        retiredCreditsRegistryId = TOKENIZATION.createDynamicStore(retiredJson);
    }

    /**
     * @notice Initialize the BitBadges collection
     */
    function initializeCollection() external onlyOperator {
        require(collectionId == 0, "Already initialized");

        // Build JSON components for collection creation
        string memory validTokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(2020, 2100);
        
        string memory defaultBalancesJson = TokenizationJSONHelpers.simpleUserBalanceStoreToJson(
            true,  // autoApproveSelfInitiatedOutgoingTransfers
            true,  // autoApproveSelfInitiatedIncomingTransfers
            false  // autoApproveAllIncomingTransfers
        );
        
        string memory collectionMetadataJson = TokenizationJSONHelpers.collectionMetadataToJson(
            "ipfs://carbon-credit-collection",
            string(abi.encodePacked(
                "{\"name\":\"", name,
                "\",\"standard\":\"", standard,
                "\",\"projectId\":\"", projectId,
                "\",\"compliance\":\"ERC-3643\"}"
            ))
        );
        
        string[] memory standards = new string[](3);
        standards[0] = "ERC-3643";
        standards[1] = "Carbon Credit";
        standards[2] = standard;
        string memory standardsJson = TokenizationJSONHelpers.stringArrayToJson(standards);
        
        string memory createJson = TokenizationJSONHelpers.createCollectionJSON(
            validTokenIdsJson,
            _addressToString(address(this)),  // manager
            collectionMetadataJson,
            defaultBalancesJson,
            "{}",  // collectionPermissions (empty)
            standardsJson,
            "",    // customData (already in metadata)
            false  // isArchived
        );
        
        collectionId = TOKENIZATION.createCollection(createJson);
    }

    // ============ Vintage Management ============

    /**
     * @notice Create a new vintage year for carbon credits
     * @param year Vintage year (e.g., 2024)
     * @param expirationTime Unix timestamp when credits expire
     * @param verificationUri IPFS URI to verification documents
     */
    function createVintage(
        uint256 year,
        uint256 expirationTime,
        string calldata verificationUri
    ) external onlyOperator {
        require(year >= 2020 && year <= 2100, "Invalid vintage year");
        require(vintages[year].year == 0, "Vintage exists");
        require(expirationTime > block.timestamp, "Expiration must be future");

        vintages[year] = VintageInfo({
            year: year,
            totalIssued: 0,
            totalRetired: 0,
            expirationTime: expirationTime,
            verificationUri: verificationUri
        });

        emit VintageCreated(year, expirationTime);
    }

    /**
     * @notice Issue carbon credits for a vintage to a verified seller
     * @param vintage Vintage year
     * @param to Recipient (must be verified seller/project developer)
     * @param amount Number of credits (1 credit = 1 tonne CO2e)
     */
    function issueCredits(
        uint256 vintage,
        address to,
        uint256 amount
    ) external onlyOperator {
        require(vintages[vintage].year != 0, "Vintage not found");
        require(isVerifiedSeller(to), "Recipient not verified seller");
        require(block.timestamp < vintages[vintage].expirationTime, "Vintage expired");

        vintages[vintage].totalIssued += amount;

        // Build token ID and ownership time ranges
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(vintage, vintage);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
            block.timestamp, 
            vintages[vintage].expirationTime
        );

        address[] memory recipients = new address[](1);
        recipients[0] = to;

        // Transfer from mint escrow (collection manager)
        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            recipients,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        TOKENIZATION.transferTokens(transferJson);

        emit CreditIssued(vintage, to, amount);
    }

    // ============ Transfer Functions ============

    /**
     * @notice Check if transfer is compliant
     * @param from Sender
     * @param to Recipient
     * @param vintage Vintage year
     */
    function canTransfer(
        address from,
        address to,
        uint256 vintage
    ) public view returns (bool) {
        // Check vintage hasn't expired
        if (block.timestamp >= vintages[vintage].expirationTime) return false;

        // Seller must be verified (for selling) OR buyer verified (for buying)
        if (!isVerifiedSeller(from) && !isVerifiedBuyer(from)) return false;
        if (!isVerifiedBuyer(to) && to != RETIREMENT_SINK) return false;

        return true;
    }

    /**
     * @notice Transfer carbon credits
     * @param to Recipient address
     * @param vintage Vintage year
     * @param amount Number of credits
     */
    function transfer(
        address to,
        uint256 vintage,
        uint256 amount
    ) external returns (bool) {
        require(canTransfer(msg.sender, to, vintage), "Transfer not compliant");

        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(vintage, vintage);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
            block.timestamp, 
            vintages[vintage].expirationTime
        );

        address[] memory recipients = new address[](1);
        recipients[0] = to;

        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            recipients,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        bool success = TOKENIZATION.transferTokens(transferJson);

        if (success) {
            emit CreditTransferred(msg.sender, to, vintage, amount);
        }
        return success;
    }

    // ============ Retirement Functions ============

    /**
     * @notice Retire carbon credits (permanent removal from circulation)
     * @param vintage Vintage year
     * @param amount Number of credits to retire
     * @param reason Retirement reason (e.g., "Offsetting 2024 emissions")
     */
    function retire(
        uint256 vintage,
        uint256 amount,
        string calldata reason
    ) external returns (bool) {
        require(vintages[vintage].year != 0, "Invalid vintage");
        require(block.timestamp < vintages[vintage].expirationTime, "Credits expired");

        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(vintage, vintage);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
            block.timestamp, 
            vintages[vintage].expirationTime
        );

        // Transfer to retirement sink (effectively burning)
        address[] memory recipients = new address[](1);
        recipients[0] = RETIREMENT_SINK;

        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            recipients,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        bool success = TOKENIZATION.transferTokens(transferJson);

        if (success) {
            vintages[vintage].totalRetired += amount;
            retiredByAddress[msg.sender] += amount;
            totalRetired += amount;

            emit CreditRetired(msg.sender, vintage, amount, reason);
        }
        return success;
    }

    // ============ Verification Functions ============

    /**
     * @notice Verify a buyer
     */
    function verifyBuyer(address buyer) external onlyVerifier {
        string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            verifiedBuyersRegistryId,
            buyer,
            true
        );
        TOKENIZATION.setDynamicStoreValue(setValueJson);
        emit BuyerVerified(buyer);
    }

    /**
     * @notice Verify a seller/project developer
     */
    function verifySeller(address seller) external onlyVerifier {
        string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            verifiedSellersRegistryId,
            seller,
            true
        );
        TOKENIZATION.setDynamicStoreValue(setValueJson);
        emit SellerVerified(seller);
    }

    /**
     * @notice Check if address is verified buyer
     */
    function isVerifiedBuyer(address account) public view returns (bool) {
        string memory getValueJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(
            verifiedBuyersRegistryId,
            account
        );
        bytes memory result = TOKENIZATION.getDynamicStoreValue(getValueJson);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    /**
     * @notice Check if address is verified seller
     */
    function isVerifiedSeller(address account) public view returns (bool) {
        string memory getValueJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(
            verifiedSellersRegistryId,
            account
        );
        bytes memory result = TOKENIZATION.getDynamicStoreValue(getValueJson);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    // ============ View Functions ============

    /**
     * @notice Get balance for a specific vintage
     */
    function balanceOf(address account, uint256 vintage) external view returns (uint256) {
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(vintage, vintage);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);

        string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
            collectionId,
            account,
            tokenIdsJson,
            ownershipTimesJson
        );
        return TOKENIZATION.getBalanceAmount(balanceJson);
    }

    /**
     * @notice Get total active credits for a vintage
     */
    function activeSupply(uint256 vintage) external view returns (uint256) {
        return vintages[vintage].totalIssued - vintages[vintage].totalRetired;
    }

    /**
     * @notice Get retirement certificate data
     */
    function getRetirementInfo(address account) external view returns (
        uint256 totalCreditsRetired
    ) {
        return retiredByAddress[account];
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
