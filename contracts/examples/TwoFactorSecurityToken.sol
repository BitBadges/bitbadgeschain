// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/**
 * @title TwoFactorSecurityToken
 * @notice Security token requiring 2FA token ownership for transfers
 * @dev Demonstrates using a separate BitBadges collection as a 2FA mechanism
 *
 * How it works:
 * 1. Admin issues time-limited 2FA tokens to users (e.g., valid for 24 hours)
 * 2. Before any sensitive operation, we check if user owns a valid 2FA token
 * 3. 2FA tokens use ownershipTimes to auto-expire
 * 4. Users must request new 2FA tokens periodically (like session tokens)
 *
 * Use cases:
 * - High-value transfers requiring additional verification
 * - Admin operations requiring multi-step auth
 * - Time-boxed trading sessions
 */
contract TwoFactorSecurityToken {
    ITokenizationPrecompile constant TOKENIZATION = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    // ============ State Variables ============

    // Main security token collection
    uint256 public securityTokenCollectionId;

    // 2FA token collection (separate collection for auth tokens)
    uint256 public twoFactorCollectionId;

    // Compliance registry
    uint256 public kycRegistryId;

    // Token metadata
    string public name;
    string public symbol;

    // 2FA configuration
    uint256 public twoFactorValidityPeriod = 24 hours;  // How long 2FA tokens are valid
    uint256 public twoFactorCooldown = 1 hours;         // Minimum time between 2FA requests

    // Roles
    address public issuer;
    address public twoFactorAuthority;  // Can issue 2FA tokens

    // 2FA request tracking
    mapping(address => uint256) public lastTwoFactorIssued;
    mapping(address => uint256) public twoFactorNonce;  // Increments with each 2FA token

    // Token ranges
    UintRange[] private _securityTokenIds;
    UintRange[] private _perpetualOwnership;

    // ============ Events ============

    event TwoFactorIssued(address indexed user, uint256 validUntil, uint256 nonce);
    event TwoFactorRevoked(address indexed user, uint256 nonce);
    event TransferWithTwoFactor(address indexed from, address indexed to, uint256 amount);
    event TwoFactorRequired(address indexed user, string operation);

    // ============ Errors ============

    error TwoFactorRequired();
    error TwoFactorExpired();
    error TwoFactorCooldownActive();
    error NotAuthorized();

    // ============ Modifiers ============

    modifier onlyIssuer() {
        if (msg.sender != issuer) revert NotAuthorized();
        _;
    }

    modifier onlyTwoFactorAuthority() {
        if (msg.sender != twoFactorAuthority && msg.sender != issuer) revert NotAuthorized();
        _;
    }

    /**
     * @notice Require valid 2FA token ownership at current time
     * @dev Checks the 2FA collection for a valid, non-expired token
     */
    modifier requiresTwoFactor() {
        if (!hasValidTwoFactor(msg.sender)) {
            emit TwoFactorRequired(msg.sender, "operation");
            revert TwoFactorRequired();
        }
        _;
    }

    // ============ Constructor ============

    /**
     * @notice Deploy the 2FA-protected security token
     * @param _name Token name
     * @param _symbol Token symbol
     * @param _twoFactorAuthority Address authorized to issue 2FA tokens
     */
    constructor(
        string memory _name,
        string memory _symbol,
        address _twoFactorAuthority
    ) {
        issuer = msg.sender;
        twoFactorAuthority = _twoFactorAuthority;
        name = _name;
        symbol = _symbol;

        // Security token uses token ID 1
        _securityTokenIds = new UintRange[](1);
        _securityTokenIds[0] = UintRange(1, 1);

        // Perpetual ownership for security tokens
        _perpetualOwnership = new UintRange[](1);
        _perpetualOwnership[0] = UintRange(1, type(uint64).max);

        // Create KYC registry
        string memory kycJson = TokenizationJSONHelpers.createDynamicStoreJSON(
            false,
            "ipfs://kyc-registry",
            "{\"type\":\"kyc-verified\"}"
        );
        kycRegistryId = TOKENIZATION.createDynamicStore(kycJson);
    }

    /**
     * @notice Initialize both token collections
     * @param totalSupply Total security tokens to create
     */
    function initialize(uint256 totalSupply) external onlyIssuer {
        require(securityTokenCollectionId == 0, "Already initialized");

        // ===== Create Security Token Collection =====
        // Build initial balance JSON (for minting totalSupply to creator)
        string memory balanceJson = string(abi.encodePacked(
            '[{"amount":"', TokenizationJSONHelpers.uintToString(totalSupply),
            '","ownershipTimes":', TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max),
            ',"tokenIds":', TokenizationJSONHelpers.uintRangeToJson(1, 1), '}]'
        ));
        
        // Build defaultBalances JSON with initial balances
        string memory defaultBalancesJson = string(abi.encodePacked(
            '{"balances":', balanceJson,
            ',"autoApproveSelfInitiatedOutgoingTransfers":true',
            ',"autoApproveSelfInitiatedIncomingTransfers":true',
            ',"autoApproveAllIncomingTransfers":false',
            ',"outgoingApprovals":[],"incomingApprovals":[],"userPermissions":{}}'
        ));
        
        string memory validTokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
        string memory collectionMetadataJson = TokenizationJSONHelpers.collectionMetadataToJson(
            "ipfs://2fa-security-token",
            string(abi.encodePacked(
                "{\"name\":\"", name,
                "\",\"symbol\":\"", symbol,
                "\",\"requires2FA\":true}"
            ))
        );
        
        string[] memory standards = new string[](2);
        standards[0] = "ERC-3643";
        standards[1] = "2FA-Protected";
        string memory standardsJson = TokenizationJSONHelpers.stringArrayToJson(standards);
        
        string memory createJson = TokenizationJSONHelpers.createCollectionJSON(
            validTokenIdsJson,
            _addressToString(address(this)),
            collectionMetadataJson,
            defaultBalancesJson,
            "{}",
            standardsJson,
            "",
            false
        );
        
        securityTokenCollectionId = TOKENIZATION.createCollection(createJson);

        // ===== Create 2FA Token Collection =====
        // 2FA tokens use incrementing token IDs (nonce-based)
        string memory twoFactorTokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
        
        // 2FA tokens should not be transferable
        string memory twoFactorDefaultsJson = TokenizationJSONHelpers.simpleUserBalanceStoreToJson(
            false,  // autoApproveSelfInitiatedOutgoingTransfers
            true,   // autoApproveSelfInitiatedIncomingTransfers
            false   // autoApproveAllIncomingTransfers
        );
        
        string memory twoFactorMetadataJson = TokenizationJSONHelpers.collectionMetadataToJson(
            "ipfs://2fa-tokens",
            "{\"type\":\"2FA\",\"transferable\":false}"
        );
        
        string[] memory twoFactorStandards = new string[](1);
        twoFactorStandards[0] = "2FA-Token";
        string memory twoFactorStandardsJson = TokenizationJSONHelpers.stringArrayToJson(twoFactorStandards);
        
        string memory twoFactorCreateJson = TokenizationJSONHelpers.createCollectionJSON(
            twoFactorTokenIdsJson,
            _addressToString(address(this)),
            twoFactorMetadataJson,
            twoFactorDefaultsJson,
            "{}",
            twoFactorStandardsJson,
            "",
            false
        );
        
        twoFactorCollectionId = TOKENIZATION.createCollection(twoFactorCreateJson);
    }

    // ============ 2FA Token Management ============

    /**
     * @notice Issue a 2FA token to a user
     * @dev Creates a time-limited ownership token
     * @param user Address to receive the 2FA token
     */
    function issueTwoFactor(address user) external onlyTwoFactorAuthority {
        // Check cooldown
        if (block.timestamp < lastTwoFactorIssued[user] + twoFactorCooldown) {
            revert TwoFactorCooldownActive();
        }

        // Increment nonce for this user (each 2FA session gets unique token ID)
        uint256 nonce = ++twoFactorNonce[user];
        uint256 validUntil = block.timestamp + twoFactorValidityPeriod;

        // Create time-limited ownership
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(nonce, nonce);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
            block.timestamp, 
            validUntil
        );

        address[] memory recipients = new address[](1);
        recipients[0] = user;

        // Issue the 2FA token (1 token with time-limited ownership)
        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            twoFactorCollectionId,
            recipients,
            1,
            tokenIdsJson,
            ownershipTimesJson
        );
        TOKENIZATION.transferTokens(transferJson);

        lastTwoFactorIssued[user] = block.timestamp;

        emit TwoFactorIssued(user, validUntil, nonce);
    }

    /**
     * @notice Batch issue 2FA tokens to multiple users
     * @param users Array of addresses to receive 2FA tokens
     */
    function batchIssueTwoFactor(address[] calldata users) external onlyTwoFactorAuthority {
        for (uint256 i = 0; i < users.length; i++) {
            address user = users[i];

            // Skip if cooldown active
            if (block.timestamp < lastTwoFactorIssued[user] + twoFactorCooldown) {
                continue;
            }

            uint256 nonce = ++twoFactorNonce[user];
            uint256 validUntil = block.timestamp + twoFactorValidityPeriod;

            string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(nonce, nonce);
            string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
                block.timestamp,
                validUntil
            );

            address[] memory recipients = new address[](1);
            recipients[0] = user;

            string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
                twoFactorCollectionId,
                recipients,
                1,
                tokenIdsJson,
                ownershipTimesJson
            );
            TOKENIZATION.transferTokens(transferJson);

            lastTwoFactorIssued[user] = block.timestamp;

            emit TwoFactorIssued(user, validUntil, nonce);
        }
    }

    /**
     * @notice Check if user has a valid (non-expired) 2FA token
     * @param user Address to check
     * @return bool True if user has valid 2FA
     */
    function hasValidTwoFactor(address user) public view returns (bool) {
        uint256 currentNonce = twoFactorNonce[user];
        if (currentNonce == 0) return false;

        // Check ownership of the current nonce token at the CURRENT TIME
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(currentNonce, currentNonce);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
            block.timestamp,
            block.timestamp
        );

        string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
            twoFactorCollectionId,
            user,
            tokenIdsJson,
            ownershipTimesJson
        );
        uint256 balance = TOKENIZATION.getBalanceAmount(balanceJson);

        return balance > 0;
    }

    /**
     * @notice Get 2FA token expiration time for a user
     * @param user Address to check
     * @return validUntil Expiration timestamp (0 if no valid token)
     */
    function getTwoFactorExpiration(address user) external view returns (uint256 validUntil) {
        uint256 lastIssued = lastTwoFactorIssued[user];
        if (lastIssued == 0) return 0;

        uint256 expiration = lastIssued + twoFactorValidityPeriod;
        if (block.timestamp >= expiration) return 0;

        return expiration;
    }

    // ============ Security Token Operations (2FA Protected) ============

    /**
     * @notice Transfer security tokens (requires valid 2FA)
     * @param to Recipient address
     * @param amount Amount to transfer
     */
    function transfer(address to, uint256 amount) external requiresTwoFactor returns (bool) {
        require(isKYCVerified(to), "Recipient not KYC verified");

        address[] memory recipients = new address[](1);
        recipients[0] = to;

        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
        
        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            securityTokenCollectionId,
            recipients,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        bool success = TOKENIZATION.transferTokens(transferJson);

        if (success) {
            emit TransferWithTwoFactor(msg.sender, to, amount);
        }

        return success;
    }

    /**
     * @notice High-value transfer with additional checks
     * @param to Recipient
     * @param amount Amount (must be > threshold for this function)
     */
    function highValueTransfer(
        address to,
        uint256 amount
    ) external requiresTwoFactor returns (bool) {
        require(isKYCVerified(to), "Recipient not KYC verified");

        // Additional check: recipient must also have valid 2FA for high-value receives
        require(hasValidTwoFactor(to), "Recipient needs 2FA for high-value transfer");

        address[] memory recipients = new address[](1);
        recipients[0] = to;

        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
        
        string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
            securityTokenCollectionId,
            recipients,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        bool success = TOKENIZATION.transferTokens(transferJson);

        if (success) {
            emit TransferWithTwoFactor(msg.sender, to, amount);
        }

        return success;
    }

    // ============ KYC Management ============

    function setKYCStatus(address user, bool verified) external onlyIssuer {
        string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            kycRegistryId,
            user,
            verified
        );
        TOKENIZATION.setDynamicStoreValue(setValueJson);
    }

    function isKYCVerified(address user) public view returns (bool) {
        string memory getValueJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(
            kycRegistryId,
            user
        );
        bytes memory result = TOKENIZATION.getDynamicStoreValue(getValueJson);
        if (result.length == 0) return false;
        return abi.decode(result, (bool));
    }

    // ============ Configuration ============

    /**
     * @notice Update 2FA validity period
     * @param newPeriod New validity period in seconds
     */
    function setTwoFactorValidityPeriod(uint256 newPeriod) external onlyIssuer {
        require(newPeriod >= 1 hours && newPeriod <= 7 days, "Invalid period");
        twoFactorValidityPeriod = newPeriod;
    }

    /**
     * @notice Update 2FA cooldown period
     * @param newCooldown New cooldown in seconds
     */
    function setTwoFactorCooldown(uint256 newCooldown) external onlyIssuer {
        require(newCooldown <= twoFactorValidityPeriod, "Cooldown > validity");
        twoFactorCooldown = newCooldown;
    }

    /**
     * @notice Update the 2FA authority
     * @param newAuthority New authority address
     */
    function setTwoFactorAuthority(address newAuthority) external onlyIssuer {
        require(newAuthority != address(0), "Invalid address");
        twoFactorAuthority = newAuthority;
    }

    // ============ View Functions ============

    function balanceOf(address account) external view returns (uint256) {
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
        
        string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
            securityTokenCollectionId,
            account,
            tokenIdsJson,
            ownershipTimesJson
        );
        return TOKENIZATION.getBalanceAmount(balanceJson);
    }

    function totalSupply() external view returns (uint256) {
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
        
        string memory supplyJson = TokenizationJSONHelpers.getTotalSupplyJSON(
            securityTokenCollectionId,
            tokenIdsJson,
            ownershipTimesJson
        );
        return TOKENIZATION.getTotalSupply(supplyJson);
    }

    /**
     * @notice Get user's 2FA status
     */
    function getTwoFactorStatus(address user) external view returns (
        bool hasValid2FA,
        uint256 currentNonce,
        uint256 lastIssued,
        uint256 expiresAt
    ) {
        currentNonce = twoFactorNonce[user];
        lastIssued = lastTwoFactorIssued[user];
        hasValid2FA = hasValidTwoFactor(user);

        if (lastIssued > 0) {
            uint256 expiration = lastIssued + twoFactorValidityPeriod;
            expiresAt = block.timestamp < expiration ? expiration : 0;
        }
    }

    /**
     * @notice Check if transfer would succeed
     */
    function canTransfer(
        address from,
        address to,
        uint256 /* amount */
    ) external view returns (bool, string memory) {
        if (!hasValidTwoFactor(from)) {
            return (false, "Sender needs valid 2FA");
        }
        if (!isKYCVerified(to)) {
            return (false, "Recipient not KYC verified");
        }
        return (true, "");
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
