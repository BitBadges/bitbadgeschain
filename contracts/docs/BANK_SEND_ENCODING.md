# Bank Send Transaction Encoding

## Overview

The Cosmos bank precompile at `0x0804` is **read-only** and does NOT support sending tokens. However, you can send tokens via:

1. **Cosmos SDK directly** - Using `cosmos.bank.v1beta1.MsgSend`
2. **SendManager module** - Using `sendmanager.MsgSendWithAliasRouting` (supports alias denoms)

This document shows how to encode send transactions in both formats.

## Cosmos SDK MsgSend Format

### Protobuf JSON Format

```json
{
  "from_address": "bb1abc123...",
  "to_address": "bb1xyz789...",
  "amount": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    }
  ]
}
```

### TypeScript Encoding Example

```typescript
import { ethers } from "ethers";

// Cosmos SDK MsgSend message
const msgSend = {
  from_address: "bb1abc123...",  // Sender Cosmos bech32 address
  to_address: "bb1xyz789...",    // Recipient Cosmos bech32 address
  amount: [
    {
      denom: "ubadge",           // Token denomination
      amount: "1000000000"        // Amount as string (1 BADGE = 1e9 ubadge)
    }
  ]
};

// Convert to JSON string
const msgJson = JSON.stringify(msgSend);
console.log(msgJson);
// Output: {"from_address":"bb1abc123...","to_address":"bb1xyz789...","amount":[{"denom":"ubadge","amount":"1000000000"}]}
```

### Multiple Coins Example

```json
{
  "from_address": "bb1abc123...",
  "to_address": "bb1xyz789...",
  "amount": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    },
    {
      "denom": "ustake",
      "amount": "500000000"
    }
  ]
}
```

## SendManager MsgSendWithAliasRouting Format

The SendManager module's `MsgSendWithAliasRouting` supports both standard coins and alias denoms (e.g., `badgeslp:...`).

### Protobuf JSON Format

```json
{
  "from_address": "bb1abc123...",
  "to_address": "bb1xyz789...",
  "amount": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    }
  ]
}
```

**Note**: The format is identical to `MsgSend`, but it routes through the sendmanager module which handles alias denom routing.

### Alias Denom Example

```json
{
  "from_address": "bb1abc123...",
  "to_address": "bb1xyz789...",
  "amount": [
    {
      "denom": "badgeslp:pool123",
      "amount": "1000000"
    }
  ]
}
```

## If Bank Precompile Had a Send Method (Hypothetical)

If the bank precompile supported sending (it doesn't), the encoding would follow the same pattern as other BitBadges precompiles:

### Solidity Interface (Hypothetical)

```solidity
interface IBankPrecompile {
    function send(string memory msgJson) external returns (bool success);
}
```

### Function Call Encoding

```typescript
import { ethers } from "ethers";

const BANK_PRECOMPILE_ADDRESS = "0x0000000000000000000000000000000000000804"; // Hypothetical
const BANK_PRECOMPILE_ABI = [
  "function send(string memory msgJson) external returns (bool success)"
];

// Prepare message
const msgSend = {
  // from_address would be automatically set from msg.sender
  to_address: "bb1xyz789...",
  amount: [
    {
      denom: "ubadge",
      amount: "1000000000"
    }
  ]
};

const msgJson = JSON.stringify(msgSend);

// Create contract instance
const bankPrecompile = new ethers.Contract(
  BANK_PRECOMPILE_ADDRESS,
  BANK_PRECOMPILE_ABI,
  signer
);

// Encode function call
const txData = bankPrecompile.interface.encodeFunctionData("send", [msgJson]);

// The encoded data would look like:
// 0x66792ba1... (function selector) + ... (encoded JSON string)
```

### Encoded Transaction Data Structure

```
0x66792ba1 (function selector for send(string))
+ offset to string data (32 bytes)
+ length of string (32 bytes)
+ string data (padded to 32-byte words)
```

**Example encoding breakdown:**
```
0x66792ba1                                    // Function selector: send(string)
0000000000000000000000000000000000000000000000000000000000000020  // Offset to string (32 bytes)
000000000000000000000000000000000000000000000000000000000000006f  // String length (111 bytes)
7b22746f5f61646472657373223a2262623178797a...  // JSON string (padded)
```

## Actual Sending Methods

### Method 1: Cosmos SDK Transaction

```typescript
import { SigningStargateClient } from "@cosmjs/stargate";
import { coins } from "@cosmjs/proto-signing";

// Create client
const client = await SigningStargateClient.connectWithSigner(
  rpcEndpoint,
  signer
);

// Build MsgSend
const msgSend = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
    fromAddress: "bb1abc123...",
    toAddress: "bb1xyz789...",
    amount: coins(1000000000, "ubadge"), // 1 BADGE
  },
};

// Send transaction
const result = await client.signAndBroadcast(
  signerAddress,
  [msgSend],
  "auto" // fee
);
```

### Method 2: Using SendManager (Supports Alias Denoms)

```typescript
import { SigningStargateClient } from "@cosmjs/stargate";

const client = await SigningStargateClient.connectWithSigner(
  rpcEndpoint,
  signer
);

// Build MsgSendWithAliasRouting
const msgSend = {
  typeUrl: "/sendmanager.v1.MsgSendWithAliasRouting",
  value: {
    fromAddress: "bb1abc123...",
    toAddress: "bb1xyz789...",
    amount: coins(1000000000, "ubadge"),
  },
};

const result = await client.signAndBroadcast(
  signerAddress,
  [msgSend],
  "auto"
);
```

### Method 3: Via EVM (If Using ERC20 Wrapper)

If tokens are wrapped as ERC20, use standard ERC20 transfer:

```typescript
import { ethers } from "ethers";

const erc20ABI = [
  "function transfer(address to, uint256 amount) external returns (bool)"
];

const erc20Contract = new ethers.Contract(
  erc20TokenAddress,
  erc20ABI,
  signer
);

const tx = await erc20Contract.transfer(
  recipientAddress,
  ethers.parseUnits("1", 9) // 1 BADGE (9 decimals)
);
```

## Complete Encoding Example

Here's a complete example showing the JSON encoding for a send transaction:

```typescript
// Helper function to build MsgSend JSON
function buildMsgSendJSON(
  fromAddress: string,
  toAddress: string,
  denom: string,
  amount: string
): string {
  const msgSend = {
    from_address: fromAddress,
    to_address: toAddress,
    amount: [
      {
        denom: denom,
        amount: amount
      }
    ]
  };
  
  return JSON.stringify(msgSend);
}

// Usage
const msgJson = buildMsgSendJSON(
  "bb1abc123...",
  "bb1xyz789...",
  "ubadge",
  "1000000000" // 1 BADGE
);

console.log(msgJson);
// Output: {"from_address":"bb1abc123...","to_address":"bb1xyz789...","amount":[{"denom":"ubadge","amount":"1000000000"}]}
```

## Key Points

1. **Bank Precompile is Read-Only**: The precompile at `0x0804` only supports queries, not sends.

2. **JSON Format**: Both `MsgSend` and `MsgSendWithAliasRouting` use the same protobuf JSON format.

3. **Amount as String**: Always use string format for amounts (e.g., `"1000000000"` not `1000000000`).

4. **Multiple Coins**: The `amount` field is an array, so you can send multiple coin types in one transaction.

5. **Address Format**: Use Cosmos bech32 addresses (e.g., `bb1...`), not EVM addresses.

6. **Alias Denoms**: Use `MsgSendWithAliasRouting` if you need to send alias denoms like `badgeslp:pool123`.

## References

- [Cosmos SDK Bank Module](https://docs.cosmos.network/v0.47/modules/bank)
- [SendManager Module](../x/sendmanager/README.md)
- [Bank Precompile Documentation](./BANK_PRECOMPILE.md)

