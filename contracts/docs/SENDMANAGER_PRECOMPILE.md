# SendManager Precompile - Frontend Integration Guide

## Overview

The SendManager Precompile enables Solidity smart contracts to send native Cosmos coins from EVM without requiring ERC20 wrapping. All accounting is kept in `x/bank` (Cosmos side), making it perfect for dual wallet support (Cosmos and EVM).

**Key Features:**
- ✅ Send native Cosmos coins directly from EVM
- ✅ No ERC20 wrapping required
- ✅ Supports both standard coins and alias denoms (e.g., `badgeslp:...`)
- ✅ All accounting in `x/bank` (Cosmos side)
- ✅ Perfect for dual wallet support

## Precompile Address

```
0x0000000000000000000000000000000000001003
```

**Address Space Convention:**
- **0x0800-0x0806**: Reserved for default Cosmos precompiles
- **0x1001+**: Reserved for custom BitBadges precompiles
  - 0x1001: Tokenization precompile
  - 0x1002: Gamm precompile
  - 0x1003: SendManager precompile ✅ (This document)
  - 0x1004+: Available for future precompiles

## Registration & Enablement

The sendmanager precompile is **registered** in `app/evm.go` and **enabled** in the genesis configuration (`config.yml`). It's also enabled via the upgrade handler for existing chains.

## ABI Interface

```solidity
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
```

## Protobuf JSON Format

The `send` method accepts a JSON string matching the `MsgSendWithAliasRouting` protobuf format:

```json
{
  "to_address": "bb1xyz789...",
  "amount": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    }
  ]
}
```

**Important Notes:**
- `from_address` is **automatically set** from `msg.sender` (cannot be spoofed)
- `to_address` must be a valid Cosmos bech32 address (e.g., `bb1...`)
- `amount` is an array of coins (can send multiple denoms in one transaction)
- Supports both standard denoms (e.g., `"ubadge"`) and alias denoms (e.g., `"badgeslp:..."`)

## Frontend Integration

### TypeScript/JavaScript Example

```typescript
import { ethers } from "ethers";

// Precompile address
const SENDMANAGER_PRECOMPILE_ADDRESS = "0x0000000000000000000000000000000000001003";

// ABI
const sendManagerABI = [
  "function send(string memory msgJson) external returns (bool success)"
];

// Create contract instance
const sendManager = new ethers.Contract(
  SENDMANAGER_PRECOMPILE_ADDRESS,
  sendManagerABI,
  signer
);

// Build message JSON
const msgJson = JSON.stringify({
  to_address: "bb1xyz789...", // Recipient Cosmos bech32 address
  amount: [
    {
      denom: "ubadge",
      amount: "1000000000" // 1 BADGE (1e9 ubadge)
    }
  ]
});

// Send transaction
const tx = await sendManager.send(msgJson);
await tx.wait();

console.log("Send successful!");
```

### Multiple Coins Example

```typescript
const msgJson = JSON.stringify({
  to_address: "bb1xyz789...",
  amount: [
    {
      denom: "ubadge",
      amount: "1000000000" // 1 BADGE
    },
    {
      denom: "ustake",
      amount: "500000000" // 0.5 STAKE
    }
  ]
});

const tx = await sendManager.send(msgJson);
```

### Alias Denom Example

```typescript
// Send alias denom (e.g., badgeslp:...)
const msgJson = JSON.stringify({
  to_address: "bb1xyz789...",
  amount: [
    {
      denom: "badgeslp:1:ubadge", // Alias denom format
      amount: "1000000000"
    }
  ]
});

const tx = await sendManager.send(msgJson);
```

## Solidity Example

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/ISendManagerPrecompile.sol";

contract TokenSender {
    ISendManagerPrecompile public constant SENDMANAGER = 
        ISendManagerPrecompile(0x0000000000000000000000000000000000001003);
    
    function sendTokens(
        string memory toAddress,
        string memory denom,
        uint256 amount
    ) external returns (bool) {
        // Build JSON string
        string memory msgJson = string(abi.encodePacked(
            '{"to_address":"',
            toAddress,
            '","amount":[{"denom":"',
            denom,
            '","amount":"',
            _uintToString(amount),
            '"}]}'
        ));
        
        return SENDMANAGER.send(msgJson);
    }
    
    function _uintToString(uint256 v) private pure returns (string memory) {
        if (v == 0) {
            return "0";
        }
        uint256 j = v;
        uint256 len;
        while (j != 0) {
            len++;
            j /= 10;
        }
        bytes memory bstr = new bytes(len);
        uint256 k = len;
        while (v != 0) {
            k = k-1;
            uint8 temp = (48 + uint8(v - v / 10 * 10));
            bytes1 b1 = bytes1(temp);
            bstr[k] = b1;
            v /= 10;
        }
        return string(bstr);
    }
}
```

## Gas Costs

- **Base Gas**: 30,000 gas per `send` transaction
- Additional gas may be consumed based on the number of coins and complexity of alias denom routing

## Error Handling

The precompile uses structured error codes:

- **ErrorCodeInvalidInput (1)**: Invalid input parameters (e.g., invalid JSON, invalid address)
- **ErrorCodeSendFailed (2)**: Send operation failed (e.g., insufficient balance)
- **ErrorCodeInsufficientBalance (3)**: Insufficient balance for the requested send
- **ErrorCodeInternalError (4)**: Internal error
- **ErrorCodeUnauthorized (5)**: Unauthorized operation

## Security Considerations

1. **Sender Authentication**: The `from_address` is **always** set from `msg.sender` and cannot be spoofed. Any `from_address` in the JSON is ignored.

2. **Address Validation**: All addresses are validated using Cosmos SDK address validation.

3. **Coin Validation**: All coins are validated before sending (denom format, amount > 0, etc.).

4. **Alias Denom Routing**: Alias denoms are routed through the sendmanager module, which handles routing to the appropriate module (e.g., tokenization for `badgeslp:` prefixes).

## Comparison with Other Methods

| Method | From EVM? | ERC20 Wrapping? | Alias Denoms? | Accounting Location |
|--------|-----------|-----------------|---------------|---------------------|
| **SendManager Precompile** | ✅ | ❌ | ✅ | x/bank (Cosmos) |
| ERC20 Wrapper | ✅ | ✅ | ❌ | ERC20 contract |
| Cosmos SDK Direct | ❌ | ❌ | ✅ | x/bank (Cosmos) |

## Use Cases

1. **Dual Wallet Support**: Send tokens from EVM wallets to Cosmos wallets (and vice versa) without wrapping
2. **Cross-Chain DApps**: Build DApps that work seamlessly with both EVM and Cosmos wallets
3. **Alias Denom Support**: Send alias denoms (e.g., `badgeslp:...`) directly from EVM
4. **Multi-Denom Transfers**: Send multiple denominations in a single transaction

## Important Notes

1. **Recipient Address Format**: Recipient addresses must be Cosmos bech32 addresses (e.g., `bb1...`), not EVM addresses (0x...).

2. **Amount Format**: Amounts are strings in the smallest unit (e.g., `"1000000000"` for 1 BADGE if BADGE has 9 decimals).

3. **Alias Denoms**: Alias denoms are supported and routed through the sendmanager module. The format depends on the alias denom type (e.g., `badgeslp:1:ubadge`).

4. **Gas Costs**: Gas is paid in the native EVM gas token, not in the coins being sent.

5. **Event Emission**: The precompile emits events for monitoring and indexing:
   - `precompile_send`: Emitted on successful send
   - `precompile_usage`: Emitted for all precompile method calls

## Testing

```typescript
import { ethers } from "hardhat";

describe("SendManager Precompile", () => {
  it("should send tokens", async () => {
    const sendManager = await ethers.getContractAt(
      "ISendManagerPrecompile",
      "0x0000000000000000000000000000000000001003"
    );

    const [signer] = await ethers.getSigners();
    const recipient = "bb1xyz789..."; // Cosmos bech32 address
    
    const msgJson = JSON.stringify({
      to_address: recipient,
      amount: [
        {
          denom: "ubadge",
          amount: "1000000000" // 1 BADGE
        }
      ]
    });

    const tx = await sendManager.send(msgJson);
    await tx.wait();

    // Verify success
    expect(tx).to.not.be.null;
  });
});
```

## Summary

The SendManager Precompile provides a simple, secure way to send native Cosmos coins from EVM without ERC20 wrapping. It's perfect for dual wallet support and maintains all accounting in `x/bank` on the Cosmos side.

**Key Benefits:**
- ✅ No ERC20 wrapping required
- ✅ Supports alias denoms
- ✅ All accounting in x/bank (Cosmos side)
- ✅ Perfect for dual wallet support
- ✅ Simple JSON-based API












