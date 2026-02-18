# Sending Tokens from EVM - Complete Guide

## Overview

The **Bank Precompile (0x0804)** is **read-only** and does NOT support sending tokens. It only provides query methods (`balances()`, `totalSupply()`, `supplyOf()`).

To send tokens from EVM, you have several options depending on your use case.

## Precompile Capabilities Summary

| Precompile | Address | Purpose | Supports Send? |
|------------|---------|---------|----------------|
| **Bank** | `0x0804` | Query balances & supply | ❌ Read-only |
| **ICS20** | `0x0802` | IBC cross-chain transfers | ✅ (IBC only, not local) |
| **ERC20 Keeper** | Dynamic | Wrap native coins as ERC20 | ✅ (via ERC20 transfer) |

## Option 1: ERC20 Wrapper (Recommended for EVM)

The **ERC20 Keeper** wraps native Cosmos coins as ERC20 tokens, allowing you to use standard ERC20 `transfer()` calls.

### How It Works

1. **Get ERC20 Address**: Each native denom has a corresponding ERC20 contract address
2. **Wrap (if needed)**: Convert native coins to ERC20 tokens
3. **Transfer**: Use standard ERC20 `transfer(address to, uint256 amount)`
4. **Unwrap (if needed)**: Convert ERC20 tokens back to native coins

### Example: Send via ERC20

```typescript
import { ethers } from "ethers";

// 1. Get ERC20 address for denom (this is chain-specific)
// The ERC20 keeper provides a method to get the ERC20 address for a denom
// You'll need to query this or have it configured

// 2. Standard ERC20 ABI
const erc20ABI = [
  "function transfer(address to, uint256 amount) external returns (bool)",
  "function balanceOf(address account) external view returns (uint256)",
  "function decimals() external view returns (uint8)"
];

// 3. Get ERC20 contract instance
const erc20Address = "0x..."; // ERC20 address for "ubadge" denom
const erc20Contract = new ethers.Contract(
  erc20Address,
  erc20ABI,
  signer
);

// 4. Send tokens (1 BADGE = 1e9 ubadge)
const recipient = "0x..."; // Recipient EVM address
const amount = ethers.parseUnits("1", 9); // 1 BADGE

const tx = await erc20Contract.transfer(recipient, amount);
await tx.wait();
```

### Solidity Example

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

contract TokenSender {
    IERC20 public token;
    
    constructor(address _tokenAddress) {
        token = IERC20(_tokenAddress);
    }
    
    function sendTokens(address to, uint256 amount) external returns (bool) {
        return token.transfer(to, amount);
    }
}
```

### Getting ERC20 Address for a Denom

The ERC20 keeper provides methods to:
- Get ERC20 address from denom
- Convert native coins to ERC20 (wrap)
- Convert ERC20 to native coins (unwrap)

**Note**: You'll need to check the ERC20 keeper documentation or query the chain to get the ERC20 address for a specific denom.

## Option 2: ICS20 Precompile (IBC Transfers Only)

The **ICS20 Precompile (0x0802)** is for **cross-chain IBC transfers**, not local sends.

### When to Use

- ✅ Sending tokens to another chain via IBC
- ❌ NOT for local (same-chain) transfers

### Example

```typescript
// ICS20 precompile is for IBC transfers only
// This requires IBC channel, receiver address on destination chain, etc.
// See ICS20 precompile documentation for details
```

## Option 3: Cosmos SDK Direct (Not from EVM)

If you're **not** in an EVM context, you can send directly via Cosmos SDK:

```typescript
import { SigningStargateClient } from "@cosmjs/stargate";

const client = await SigningStargateClient.connectWithSigner(
  rpcEndpoint,
  signer
);

const msgSend = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
    fromAddress: "bb1abc123...",
    toAddress: "bb1xyz789...",
    amount: coins(1000000000, "ubadge"), // 1 BADGE
  },
};

const result = await client.signAndBroadcast(
  signerAddress,
  [msgSend],
  "auto"
);
```

**Note**: This is NOT from EVM - it's a direct Cosmos SDK transaction.

## Option 4: Custom Precompile (Future)

You could create a custom precompile that wraps `MsgSend` or `MsgSendWithAliasRouting`:

```solidity
interface ICustomBankPrecompile {
    function send(string memory msgJson) external returns (bool);
}
```

This would accept Protobuf JSON for `cosmos.bank.v1beta1.MsgSend` or `sendmanager.MsgSendWithAliasRouting`.

**Note**: This is not currently implemented - it would require creating a new precompile.

## Comparison Table

| Method | From EVM? | Local Send? | Cross-Chain? | Complexity |
|--------|-----------|-------------|--------------|------------|
| **ERC20 Wrapper** | ✅ | ✅ | ❌ | Low (standard ERC20) |
| **ICS20 Precompile** | ✅ | ❌ | ✅ | Medium (IBC setup) |
| **Cosmos SDK Direct** | ❌ | ✅ | ✅ | Low (but not EVM) |
| **Custom Precompile** | ✅ | ✅ | ✅ | High (needs implementation) |

## Recommended Approach

For **simple (denom, amount) sends from EVM**:

1. **Use ERC20 Wrapper** (Option 1)
   - Most EVM-native approach
   - Standard ERC20 interface
   - Works with existing EVM tooling

2. **If ERC20 not available**, consider:
   - Creating a custom precompile that wraps `MsgSend`
   - Using Cosmos SDK directly (if not in EVM context)

## Getting ERC20 Address

To find the ERC20 address for a denom, you'll need to:

1. Query the ERC20 keeper module
2. Use a helper contract that calls the ERC20 keeper
3. Check chain documentation for pre-configured addresses

Example query (pseudo-code):
```solidity
// ERC20 keeper provides methods to get ERC20 address from denom
address erc20Address = erc20Keeper.getERC20Address("ubadge");
```

## Important Notes

1. **Bank Precompile is Read-Only**: The bank precompile at `0x0804` only supports queries, not sends.

2. **ICS20 is for IBC**: The ICS20 precompile at `0x0802` is for cross-chain transfers, not local sends.

3. **ERC20 Wrapper**: The ERC20 keeper wraps native coins, allowing standard ERC20 transfers from EVM.

4. **Address Conversion**: When using ERC20, recipient addresses are EVM addresses (0x...), not Cosmos bech32 addresses (bb1...).

5. **Gas Costs**: ERC20 transfers consume gas like any EVM transaction. Cosmos SDK direct sends use Cosmos gas/fees.

## Example: Complete Flow

```typescript
// 1. Get ERC20 address for denom
const erc20Address = await getERC20Address("ubadge");

// 2. Check if tokens need wrapping
const nativeBalance = await queryBankBalance(cosmosAddress, "ubadge");
if (nativeBalance > 0) {
  // Wrap native coins to ERC20 (if needed)
  await wrapCoinsToERC20("ubadge", nativeBalance);
}

// 3. Send via ERC20
const erc20 = new ethers.Contract(erc20Address, erc20ABI, signer);
await erc20.transfer(recipientEVMAddress, amount);

// 4. Recipient can unwrap if needed
// await unwrapERC20ToCoins("ubadge", amount);
```

## Summary

- **Bank Precompile (0x0804)**: Read-only queries only ❌
- **ICS20 Precompile (0x0802)**: IBC transfers only (cross-chain) ⚠️
- **ERC20 Keeper**: Wrap native coins and use standard ERC20 transfers ✅ (Recommended)
- **Cosmos SDK Direct**: Not from EVM, but works for Cosmos-native sends ✅

For simple (denom, amount) sends from EVM, **use the ERC20 wrapper approach** (Option 1).


















