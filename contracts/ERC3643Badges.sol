// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IERC3643.sol";
import "./interfaces/IBadgesPrecompile.sol";

/**
 * @title ERC3643Badges
 * @dev ERC-3643 compliant contract that uses the badges precompile for token transfers
 */
contract ERC3643Badges is IERC3643 {
    // Badges precompile address: 0x0000000000000000000000000000000000001001
    address public constant BADGES_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000001001;
    
    // Collection ID for this token instance
    uint256 public immutable collectionId;
    
    // Token ID range (1-1 for this proof of concept)
    UintRange private constant TOKEN_IDS = UintRange({start: 1, end: 1});
    
    // Ownership times "All" = 1 to MaxUint64
    UintRange private constant OWNERSHIP_TIMES = UintRange({
        start: 1,
        end: type(uint256).max
    });
    
    // Reference to the badges precompile
    IBadgesPrecompile private constant badgesPrecompile = IBadgesPrecompile(BADGES_PRECOMPILE_ADDRESS);
    
    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to another (`to`).
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);
    
    /**
     * @dev Constructor sets the collection ID
     * @param _collectionId The collection ID to use for transfers
     */
    constructor(uint256 _collectionId) {
        collectionId = _collectionId;
    }
    
    /**
     * @dev Transfer tokens using the badges precompile
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
        tokenIds[0] = TOKEN_IDS;
        
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = OWNERSHIP_TIMES;
        
        // Call badges precompile (msg.sender is automatically used as from)
        (bool success, bytes memory returnData) = BADGES_PRECOMPILE_ADDRESS.call(
            abi.encodeWithSelector(
                IBadgesPrecompile.transferTokens.selector,
                collectionId,
                toAddresses,
                amount,
                tokenIds,
                ownershipTimes
            )
        );
        
        require(success, "ERC3643: transfer failed");
        
        // Decode return value
        bool result = abi.decode(returnData, (bool));
        require(result, "ERC3643: transfer returned false");
        
        // Emit Transfer event
        emit Transfer(msg.sender, to, amount);
        
        return true;
    }
    
    /**
     * @dev Get balance of an account
     * @param account The account to query
     * @return balance The balance for the account
     */
    function balanceOf(address account) external view override returns (uint256) {
        require(account != address(0), "ERC3643: balance query for zero address");
        
        // Create arrays from constants for the precompile call
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TOKEN_IDS;
        
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = OWNERSHIP_TIMES;
        
        // Call badges precompile to get balance amount
        (bool success, bytes memory returnData) = BADGES_PRECOMPILE_ADDRESS.staticcall(
            abi.encodeWithSelector(
                IBadgesPrecompile.getBalanceAmount.selector,
                collectionId,
                account,
                tokenIds,
                ownershipTimes
            )
        );
        
        require(success, "ERC3643: balance query failed");
        
        // Decode return value
        uint256 balance = abi.decode(returnData, (uint256));
        return balance;
    }
    
    /**
     * @dev Get total supply
     * @return supply The total supply
     */
    function totalSupply() external view override returns (uint256) {
        // Create arrays from constants for the precompile call
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TOKEN_IDS;
        
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = OWNERSHIP_TIMES;
        
        // Call badges precompile to get total supply
        (bool success, bytes memory returnData) = BADGES_PRECOMPILE_ADDRESS.staticcall(
            abi.encodeWithSelector(
                IBadgesPrecompile.getTotalSupply.selector,
                collectionId,
                tokenIds,
                ownershipTimes
            )
        );
        
        require(success, "ERC3643: total supply query failed");
        
        // Decode return value
        uint256 supply = abi.decode(returnData, (uint256));
        return supply;
    }
}

