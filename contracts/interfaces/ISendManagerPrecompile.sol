// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title ISendManagerPrecompile
/// @notice Interface for the BitBadges sendmanager precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000001003
///      All methods use JSON string parameters matching protobuf JSON format.
///      The caller address (sender) is automatically set from msg.sender.
///      Use helper libraries to construct JSON strings from Solidity types.
///      
///      This precompile enables sending native Cosmos coins from EVM without
///      requiring ERC20 wrapping. All accounting is kept in x/bank (Cosmos side).
///      Supports both standard coins and alias denoms (e.g., badgeslp:...).
interface ISendManagerPrecompile {
    /// @notice Send native Cosmos coins from the caller to a recipient
    /// @param msgJson JSON string matching MsgSendWithAliasRouting protobuf JSON format
    ///                Example: {"to_address":"bb1...","amount":[{"denom":"ubadge","amount":"1000000000"}]}
    ///                Note: from_address is automatically set from msg.sender
    ///                Supports both standard denoms (e.g., "ubadge") and alias denoms (e.g., "badgeslp:...")
    /// @return success Whether the send succeeded
    function send(string memory msgJson) external returns (bool success);
}

