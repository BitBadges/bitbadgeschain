// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IERC3643.sol";
import "./interfaces/ITokenizationPrecompile.sol";
import "./libraries/TokenizationWrappers.sol";
import "./libraries/TokenizationHelpers.sol";
import "./types/TokenizationTypes.sol";

/**
 * @title ERC3643Tokenization
 * @dev Minimal ERC-3643 style token backed by a BitBadges collection.
 *
 * All precompile methods take a single JSON string argument; the calls are
 * built through TokenizationWrappers so the encoding always matches the
 * precompile ABI. Token ID 1 is used for every balance, and balances are read
 * at the current block time.
 */
contract ERC3643Tokenization is IERC3643 {
    using TokenizationWrappers for ITokenizationPrecompile;

    // Tokenization precompile address: 0x0000000000000000000000000000000000001001
    address public constant TOKENIZATION_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000001001;

    // Collection ID for this token instance
    uint256 public immutable collectionId;

    // Single token ID (fungible-like behaviour)
    uint256 private constant TOKEN_ID = 1;

    // Reference to the tokenization precompile
    ITokenizationPrecompile private constant tokenizationPrecompile =
        ITokenizationPrecompile(TOKENIZATION_PRECOMPILE_ADDRESS);

    /**
     * @dev Constructor sets the collection ID
     * @param _collectionId The collection ID to use for transfers
     */
    constructor(uint256 _collectionId) {
        collectionId = _collectionId;
    }

    /**
     * @dev Transfer tokens using the tokenization precompile (msg.sender is the sender)
     * @param to The recipient address
     * @param amount The amount to transfer
     * @return success Whether the transfer succeeded
     */
    function transfer(address to, uint256 amount) external override returns (bool) {
        require(to != address(0), "ERC3643: transfer to zero address");
        require(amount > 0, "ERC3643: transfer amount must be greater than zero");

        address[] memory toAddresses = new address[](1);
        toAddresses[0] = to;

        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(TOKEN_ID);

        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();

        bool success = tokenizationPrecompile.transferTokens(collectionId, toAddresses, amount, tokenIds, ownershipTimes);
        require(success, "ERC3643: transfer failed");

        emit Transfer(msg.sender, to, amount);
        return true;
    }

    /**
     * @dev Get balance of an account at the current block time
     * @param account The account to query
     * @return balance The balance for the account
     */
    function balanceOf(address account) external view override returns (uint256) {
        require(account != address(0), "ERC3643: balance query for zero address");
        return tokenizationPrecompile.getBalanceAmount(collectionId, account, TOKEN_ID, _nowMillis());
    }

    /**
     * @dev Get total supply at the current block time
     * @return supply The total supply
     */
    function totalSupply() external view override returns (uint256) {
        return tokenizationPrecompile.getTotalSupply(collectionId, TOKEN_ID, _nowMillis());
    }

    // Ownership times are millisecond timestamps.
    function _nowMillis() private view returns (uint256) {
        return block.timestamp * 1000;
    }
}
