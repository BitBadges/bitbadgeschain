// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/ITokenizationPrecompile.sol";
import "../interfaces/IERC3643.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/**
 * @title ERC3643Template
 * @notice Standardized ERC-3643 security token template for BitBadges
 * @dev Maps ERC-3643 standard functions to the BitBadges tokenization precompile.
 *
 * Key design decisions:
 * - Token ID 1 is used for all balance operations (fungible-like behavior)
 * - Ownership times span 1 to uint64.max (perpetual ownership)
 * - Auto-scan mode for approvals (no prioritized approvals)
 * - Compliance checks are performed in the contract before transfers
 * - Forceful operations (recovery, freeze) use Dynamic Stores and can be
 *   managed through the BitBadges module directly or via admin functions
 *
 * To use this template:
 * 1. Deploy the contract with token name and symbol
 * 2. Call initializeCollection() with desired total supply
 * 3. Register investors via registerIdentity() before they can hold tokens
 * 4. Transfers automatically check compliance via canTransfer()
 *
 * ERC-3643 Standard Mapping:
 * - Identity Registry → Dynamic Stores (kycRegistryId, accreditedRegistryId)
 * - Compliance Module → canTransfer() function with configurable rules
 * - Token Operations → BitBadges precompile with ID 1, full time range
 */
contract ERC3643Template is IERC3643 {
    // ============ Constants ============

    /// @notice BitBadges tokenization precompile address
    ITokenizationPrecompile constant PRECOMPILE =
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    /// @notice Token ID used for all operations (fungible-like behavior)
    uint256 private constant TOKEN_ID = 1;

    /// @notice Maximum uint64 value for ownership time range
    uint64 private constant MAX_TIME = type(uint64).max;

    // ============ State Variables ============

    /// @notice BitBadges collection ID (set on initialization)
    uint256 public collectionId;

    /// @notice Token metadata
    string public name;
    string public symbol;
    uint8 public constant decimals = 0;

    /// @notice Identity Registry via Dynamic Stores
    uint256 public kycRegistryId;
    uint256 public accreditedRegistryId;
    uint256 public frozenRegistryId;

    /// @notice Admin address (can manage compliance)
    address public admin;

    /// @notice Compliance agents (can manage investor identities)
    mapping(address => bool) public isComplianceAgent;

    // ============ Implementation-Specific Events ============
    // Note: Transfer, IdentityRegistered, IdentityRemoved, AddressFrozen, AddressUnfrozen
    // are inherited from IERC3643 interface

    event CollectionInitialized(uint256 indexed collectionId);
    event ComplianceAgentAdded(address indexed agent);
    event ComplianceAgentRemoved(address indexed agent);

    // ============ Modifiers ============

    modifier onlyAdmin() {
        require(msg.sender == admin, "ERC3643: not admin");
        _;
    }

    modifier onlyComplianceAgent() {
        require(
            isComplianceAgent[msg.sender] || msg.sender == admin,
            "ERC3643: not compliance agent"
        );
        _;
    }

    modifier whenInitialized() {
        require(collectionId != 0, "ERC3643: not initialized");
        _;
    }

    // ============ Constructor ============

    /**
     * @notice Deploy a new ERC-3643 security token
     * @param _name Token name
     * @param _symbol Token symbol
     */
    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
        admin = msg.sender;
        isComplianceAgent[msg.sender] = true;

        // Create Dynamic Stores for compliance registries
        kycRegistryId = PRECOMPILE.createDynamicStore(
            TokenizationJSONHelpers.createDynamicStoreJSON(
                false,
                "",
                string(abi.encodePacked('{"type":"kyc","token":"', _symbol, '"}'))
            )
        );

        accreditedRegistryId = PRECOMPILE.createDynamicStore(
            TokenizationJSONHelpers.createDynamicStoreJSON(
                false,
                "",
                '{"type":"accredited"}'
            )
        );

        frozenRegistryId = PRECOMPILE.createDynamicStore(
            TokenizationJSONHelpers.createDynamicStoreJSON(
                false,
                "",
                '{"type":"frozen"}'
            )
        );
    }

    // ============ Initialization ============

    /**
     * @notice Initialize the underlying BitBadges collection
     * @dev Must be called after deployment. Mints all tokens to the admin.
     * @param _totalSupply Total number of tokens to create
     */
    function initializeCollection(uint256 _totalSupply) external onlyAdmin {
        require(collectionId == 0, "ERC3643: already initialized");

        // Build the balance JSON for initial supply (all to admin)
        string memory balanceJson = string(abi.encodePacked(
            '[{"amount":"', TokenizationJSONHelpers.uintToString(_totalSupply),
            '","ownershipTimes":', _ownershipTimesJson(),
            ',"tokenIds":', _tokenIdsJson(), '}]'
        ));

        // Default balances with auto-approve for self-initiated transfers
        string memory defaultBalancesJson = string(abi.encodePacked(
            '{"balances":', balanceJson,
            ',"autoApproveSelfInitiatedOutgoingTransfers":true',
            ',"autoApproveSelfInitiatedIncomingTransfers":true',
            ',"autoApproveAllIncomingTransfers":false',
            ',"outgoingApprovals":[],"incomingApprovals":[],"userPermissions":{}}'
        ));

        // Create collection with ERC-3643 standard tag
        string[] memory standards = new string[](1);
        standards[0] = "ERC-3643";

        string memory createJson = TokenizationJSONHelpers.createCollectionJSON(
            _tokenIdsJson(),
            TokenizationJSONHelpers.addressToString(address(this)),
            TokenizationJSONHelpers.collectionMetadataToJson(
                "",
                string(abi.encodePacked(
                    '{"name":"', name,
                    '","symbol":"', symbol,
                    '","standard":"ERC-3643"}'
                ))
            ),
            defaultBalancesJson,
            "{}",
            TokenizationJSONHelpers.stringArrayToJson(standards),
            "",
            false
        );

        collectionId = PRECOMPILE.createCollection(createJson);
        emit CollectionInitialized(collectionId);
    }

    // ============ ERC-3643 Core Functions ============

    /**
     * @notice Transfer tokens with compliance checks
     * @param to Recipient address
     * @param amount Amount to transfer
     * @return success True if transfer succeeded
     */
    function transfer(address to, uint256 amount)
        external
        override
        whenInitialized
        returns (bool)
    {
        require(to != address(0), "ERC3643: transfer to zero address");
        require(amount > 0, "ERC3643: zero amount");

        // Check compliance
        (bool allowed, string memory reason) = canTransfer(msg.sender, to, amount);
        require(allowed, reason);

        // Execute transfer via precompile
        address[] memory recipients = new address[](1);
        recipients[0] = to;

        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            recipients,
            amount,
            _tokenIdsJson(),
            _ownershipTimesJson()
        );

        bool success = PRECOMPILE.transferTokens(transferJson);
        require(success, "ERC3643: transfer failed");

        emit Transfer(msg.sender, to, amount);
        return true;
    }

    /**
     * @notice Get token balance of an address
     * @param account Address to query
     * @return balance Token balance
     */
    function balanceOf(address account)
        external
        view
        override
        returns (uint256)
    {
        if (collectionId == 0) return 0;

        string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
            collectionId,
            account,
            1,                        // Single token ID (fungible tokens use ID 1)
            block.timestamp * 1000    // Current time in milliseconds
        );
        return PRECOMPILE.getBalanceAmount(balanceJson);
    }

    /**
     * @notice Get total token supply
     * @return supply Total supply
     */
    function totalSupply() external view override returns (uint256) {
        if (collectionId == 0) return 0;

        string memory supplyJson = TokenizationJSONHelpers.getTotalSupplyJSON(
            collectionId,
            1,                        // Single token ID (fungible tokens use ID 1)
            block.timestamp * 1000    // Current time in milliseconds
        );
        return PRECOMPILE.getTotalSupply(supplyJson);
    }

    // ============ ERC-3643 Compliance Functions ============

    /**
     * @notice Check if a transfer is compliant
     * @param from Sender address
     * @param to Recipient address
     * @param amount Amount to transfer (unused in basic check, available for extensions)
     * @return allowed True if transfer is allowed
     * @return reason Reason if transfer is not allowed
     */
    function canTransfer(address from, address to, uint256 amount)
        public
        view
        returns (bool allowed, string memory reason)
    {
        // Silence unused variable warning - amount can be used in extensions
        amount;

        // Check KYC status
        if (!isKYCVerified(from)) {
            return (false, "ERC3643: sender not KYC verified");
        }
        if (!isKYCVerified(to)) {
            return (false, "ERC3643: recipient not KYC verified");
        }

        // Check frozen status
        if (isFrozen(from)) {
            return (false, "ERC3643: sender frozen");
        }
        if (isFrozen(to)) {
            return (false, "ERC3643: recipient frozen");
        }

        return (true, "");
    }

    // ============ Identity Registry Functions ============

    /**
     * @notice Register an investor's identity
     * @param investor Address of the investor
     * @param _isAccredited Whether the investor is accredited
     */
    function registerIdentity(address investor, bool _isAccredited)
        external
        onlyComplianceAgent
    {
        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                kycRegistryId,
                investor,
                true
            )
        );

        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                accreditedRegistryId,
                investor,
                _isAccredited
            )
        );

        emit IdentityRegistered(investor, _isAccredited);
    }

    /**
     * @notice Remove an investor's identity registration
     * @param investor Address of the investor
     */
    function removeIdentity(address investor) external onlyComplianceAgent {
        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                kycRegistryId,
                investor,
                false
            )
        );

        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                accreditedRegistryId,
                investor,
                false
            )
        );

        emit IdentityRemoved(investor);
    }

    /**
     * @notice Check if an address is KYC verified
     * @param investor Address to check
     * @return verified True if KYC verified
     */
    function isKYCVerified(address investor) public view returns (bool) {
        bytes memory result = PRECOMPILE.getDynamicStoreValue(
            TokenizationJSONHelpers.getDynamicStoreValueJSON(kycRegistryId, investor)
        );
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    /**
     * @notice Check if an address is accredited
     * @param investor Address to check
     * @return accredited True if accredited
     */
    function isAccredited(address investor) public view returns (bool) {
        bytes memory result = PRECOMPILE.getDynamicStoreValue(
            TokenizationJSONHelpers.getDynamicStoreValueJSON(accreditedRegistryId, investor)
        );
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    /**
     * @notice Freeze an address (prevent all transfers)
     * @param investor Address to freeze
     */
    function freezeAddress(address investor) external onlyComplianceAgent {
        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                frozenRegistryId,
                investor,
                true
            )
        );
        emit AddressFrozen(investor);
    }

    /**
     * @notice Unfreeze an address
     * @param investor Address to unfreeze
     */
    function unfreezeAddress(address investor) external onlyComplianceAgent {
        PRECOMPILE.setDynamicStoreValue(
            TokenizationJSONHelpers.setDynamicStoreValueJSON(
                frozenRegistryId,
                investor,
                false
            )
        );
        emit AddressUnfrozen(investor);
    }

    /**
     * @notice Check if an address is frozen
     * @param investor Address to check
     * @return frozen True if frozen
     */
    function isFrozen(address investor) public view returns (bool) {
        bytes memory result = PRECOMPILE.getDynamicStoreValue(
            TokenizationJSONHelpers.getDynamicStoreValueJSON(frozenRegistryId, investor)
        );
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    /**
     * @notice Get full compliance status of an investor
     * @param investor Address to check
     * @return kyc KYC verified status
     * @return accredited Accredited investor status
     * @return frozen Frozen status
     */
    function getInvestorStatus(address investor)
        external
        view
        returns (bool kyc, bool accredited, bool frozen)
    {
        return (isKYCVerified(investor), isAccredited(investor), isFrozen(investor));
    }

    // ============ Admin Functions ============

    /**
     * @notice Add a compliance agent
     * @param agent Address to add
     */
    function addComplianceAgent(address agent) external onlyAdmin {
        isComplianceAgent[agent] = true;
        emit ComplianceAgentAdded(agent);
    }

    /**
     * @notice Remove a compliance agent
     * @param agent Address to remove
     */
    function removeComplianceAgent(address agent) external onlyAdmin {
        isComplianceAgent[agent] = false;
        emit ComplianceAgentRemoved(agent);
    }

    /**
     * @notice Transfer admin role
     * @param newAdmin New admin address
     */
    function transferAdmin(address newAdmin) external onlyAdmin {
        require(newAdmin != address(0), "ERC3643: invalid address");
        admin = newAdmin;
    }

    /**
     * @notice Set collection ID (for existing collections)
     * @param _collectionId Collection ID to use
     */
    function setCollectionId(uint256 _collectionId) external onlyAdmin {
        collectionId = _collectionId;
    }

    /**
     * @notice Set registry IDs (for existing registries)
     */
    function setRegistryIds(
        uint256 _kycRegistryId,
        uint256 _accreditedRegistryId,
        uint256 _frozenRegistryId
    ) external onlyAdmin {
        kycRegistryId = _kycRegistryId;
        accreditedRegistryId = _accreditedRegistryId;
        frozenRegistryId = _frozenRegistryId;
    }

    // ============ Internal Helpers ============

    /**
     * @notice Get JSON for token ID range (always ID 1)
     */
    function _tokenIdsJson() internal pure returns (string memory) {
        return '[{"start":"1","end":"1"}]';
    }

    /**
     * @notice Get JSON for ownership time range (1 to uint64.max)
     */
    function _ownershipTimesJson() internal pure returns (string memory) {
        return '[{"start":"1","end":"18446744073709551615"}]';
    }

}
