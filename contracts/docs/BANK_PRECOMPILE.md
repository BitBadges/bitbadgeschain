# Bank Precompile - Frontend Integration Guide

## Overview

The Bank Precompile is a **read-only** precompile provided by the Cosmos EVM module that enables querying token balances and supply information directly from Solidity smart contracts. It's located at address `0x0000000000000000000000000000000000000804`.

**Important**: This is the **Cosmos default bank precompile** (0x0804), not a custom BitBadges precompile. It provides query functionality only - it does NOT support sending tokens.

## Precompile Address

```
0x0000000000000000000000000000000000000804
```

**Address Space Convention:**
- **0x0800-0x0806**: Reserved for default Cosmos precompiles (from cosmos/evm/x/vm/types/precompiles.go)
  - 0x0800: Staking precompile
  - 0x0801: Distribution precompile
  - 0x0802: ICS20 (IBC) precompile
  - 0x0803: Vesting precompile
  - 0x0804: Bank precompile ✅ (This document)
  - 0x0805: Governance precompile
  - 0x0806: Slashing precompile
- **0x1001+**: Reserved for custom BitBadges precompiles
  - 0x1001: Tokenization precompile
  - 0x1002: Gamm precompile

## Registration & Enablement

The bank precompile is **automatically registered** via `DefaultStaticPrecompiles` in `app/evm.go` and is **enabled** in the genesis configuration (`config.yml`). No additional setup is required.

## ABI Interface

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title IBankPrecompile
/// @notice Interface for the Cosmos bank precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000000804
///      All methods are read-only (view functions)
interface IBankPrecompile {
    /// @notice Query all token balances for an account
    /// @param account The account address to query
    /// @return balances Array of Balance structs containing denom and amount
    function balances(address account) external view returns (Balance[] memory balances);

    /// @notice Query the total supply of all tokens on the chain
    /// @return supply Array of Balance structs containing denom and total supply
    function totalSupply() external view returns (Balance[] memory supply);

    /// @notice Query the total supply of a specific token denomination
    /// @param denom The token denomination (e.g., "ubadge", "ustake")
    /// @return amount The total supply amount
    function supplyOf(string memory denom) external view returns (uint256 amount);
}

/// @notice Balance struct returned by the precompile
struct Balance {
    string denom;    // Token denomination (e.g., "ubadge")
    uint256 amount;  // Amount in smallest unit (e.g., 1 BADGE = 1e9 ubadge)
}
```

## ABI JSON

```json
{
  "_format": "hh-sol-artifact-1",
  "contractName": "IBankPrecompile",
  "abi": [
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balances",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "denom",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "amount",
              "type": "uint256"
            }
          ],
          "internalType": "struct IBankPrecompile.Balance[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "denom",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "amount",
              "type": "uint256"
            }
          ],
          "internalType": "struct IBankPrecompile.Balance[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "denom",
          "type": "string"
        }
      ],
      "name": "supplyOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]
}
```

## Frontend Integration

### Step 1: Install Dependencies

```bash
npm install ethers
# or
yarn add ethers
```

### Step 2: Setup Constants

```typescript
// constants/bankPrecompile.ts
export const BANK_PRECOMPILE_ADDRESS = "0x0000000000000000000000000000000000000804";

export const BANK_PRECOMPILE_ABI = [
  {
    name: "balances",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "account", type: "address" }],
    outputs: [
      {
        name: "",
        type: "tuple[]",
        components: [
          { name: "denom", type: "string" },
          { name: "amount", type: "uint256" }
        ]
      }
    ]
  },
  {
    name: "totalSupply",
    type: "function",
    stateMutability: "view",
    inputs: [],
    outputs: [
      {
        name: "",
        type: "tuple[]",
        components: [
          { name: "denom", type: "string" },
          { name: "amount", type: "uint256" }
        ]
      }
    ]
  },
  {
    name: "supplyOf",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "denom", type: "string" }],
    outputs: [{ name: "", type: "uint256" }]
  }
];
```

### Step 3: Create Helper Functions

```typescript
// utils/bankPrecompile.ts
import { ethers } from "ethers";
import { BANK_PRECOMPILE_ADDRESS, BANK_PRECOMPILE_ABI } from "../constants/bankPrecompile";

export interface Balance {
  denom: string;
  amount: bigint;
}

/**
 * Get all token balances for an account
 */
export async function getBalances(
  provider: ethers.Provider,
  accountAddress: string
): Promise<Balance[]> {
  const bankContract = new ethers.Contract(
    BANK_PRECOMPILE_ADDRESS,
    BANK_PRECOMPILE_ABI,
    provider
  );

  const balances = await bankContract.balances(accountAddress);
  return balances.map((b: any) => ({
    denom: b.denom,
    amount: BigInt(b.amount.toString())
  }));
}

/**
 * Get total supply of all tokens
 */
export async function getTotalSupply(
  provider: ethers.Provider
): Promise<Balance[]> {
  const bankContract = new ethers.Contract(
    BANK_PRECOMPILE_ADDRESS,
    BANK_PRECOMPILE_ABI,
    provider
  );

  const supply = await bankContract.totalSupply();
  return supply.map((s: any) => ({
    denom: s.denom,
    amount: BigInt(s.amount.toString())
  }));
}

/**
 * Get total supply of a specific token denomination
 */
export async function getSupplyOf(
  provider: ethers.Provider,
  denom: string
): Promise<bigint> {
  const bankContract = new ethers.Contract(
    BANK_PRECOMPILE_ADDRESS,
    BANK_PRECOMPILE_ABI,
    provider
  );

  const supply = await bankContract.supplyOf(denom);
  return BigInt(supply.toString());
}

/**
 * Get balance of a specific denomination for an account
 */
export async function getBalance(
  provider: ethers.Provider,
  accountAddress: string,
  denom: string
): Promise<bigint | null> {
  const balances = await getBalances(provider, accountAddress);
  const balance = balances.find(b => b.denom === denom);
  return balance ? balance.amount : null;
}
```

### Step 4: Usage Examples

#### Display User Balances

```typescript
import { ethers } from "ethers";
import { getBalances } from "./utils/bankPrecompile";

const provider = new ethers.JsonRpcProvider("YOUR_RPC_URL");
const userAddress = "0x742d35Cc6634C0532925a3b844Bc9e7595f25e4";

async function displayBalances() {
  const balances = await getBalances(provider, userAddress);
  
  console.log(`Balances for ${userAddress}:`);
  balances.forEach(balance => {
    // Convert from smallest unit (e.g., ubadge) to display unit (e.g., BADGE)
    const displayAmount = Number(balance.amount) / 1e9;
    console.log(`  ${balance.denom}: ${displayAmount}`);
  });
}
```

#### Get Specific Token Balance

```typescript
import { getBalance } from "./utils/bankPrecompile";

async function getBadgeBalance(userAddress: string) {
  const provider = new ethers.JsonRpcProvider("YOUR_RPC_URL");
  const balance = await getBalance(provider, userAddress, "ubadge");
  
  if (balance !== null) {
    const badgeAmount = Number(balance) / 1e9; // Convert ubadge to BADGE
    console.log(`BADGE balance: ${badgeAmount}`);
  } else {
    console.log("No BADGE balance found");
  }
}
```

#### Monitor Total Supply

```typescript
import { getTotalSupply, getSupplyOf } from "./utils/bankPrecompile";

async function monitorSupply() {
  const provider = new ethers.JsonRpcProvider("YOUR_RPC_URL");
  
  // Get all token supplies
  const allSupply = await getTotalSupply(provider);
  console.log("Total supplies:");
  allSupply.forEach(supply => {
    const displayAmount = Number(supply.amount) / 1e9;
    console.log(`  ${supply.denom}: ${displayAmount}`);
  });
  
  // Get specific token supply
  const badgeSupply = await getSupplyOf(provider, "ubadge");
  const badgeAmount = Number(badgeSupply) / 1e9;
  console.log(`Total BADGE supply: ${badgeAmount}`);
}
```

### Step 5: React Component Example

```typescript
import { useState, useEffect } from "react";
import { ethers } from "ethers";
import { getBalances, Balance } from "./utils/bankPrecompile";

function BalanceDisplay({ userAddress }: { userAddress: string }) {
  const [balances, setBalances] = useState<Balance[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchBalances() {
      const provider = new ethers.JsonRpcProvider("YOUR_RPC_URL");
      try {
        const userBalances = await getBalances(provider, userAddress);
        setBalances(userBalances);
      } catch (error) {
        console.error("Failed to fetch balances:", error);
      } finally {
        setLoading(false);
      }
    }

    fetchBalances();
  }, [userAddress]);

  if (loading) return <div>Loading balances...</div>;

  return (
    <div>
      <h3>Balances for {userAddress}</h3>
      <ul>
        {balances.map((balance, index) => {
          const displayAmount = Number(balance.amount) / 1e9;
          return (
            <li key={index}>
              {balance.denom}: {displayAmount}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
```

## Complete Frontend Client

```typescript
// bankPrecompileClient.ts
import { ethers } from "ethers";
import { BANK_PRECOMPILE_ADDRESS, BANK_PRECOMPILE_ABI } from "./constants";

export class BankPrecompileClient {
  private contract: ethers.Contract;

  constructor(provider: ethers.Provider) {
    this.contract = new ethers.Contract(
      BANK_PRECOMPILE_ADDRESS,
      BANK_PRECOMPILE_ABI,
      provider
    );
  }

  /**
   * Get all balances for an account
   */
  async balances(account: string): Promise<Balance[]> {
    const result = await this.contract.balances(account);
    return result.map((b: any) => ({
      denom: b.denom,
      amount: BigInt(b.amount.toString())
    }));
  }

  /**
   * Get total supply of all tokens
   */
  async totalSupply(): Promise<Balance[]> {
    const result = await this.contract.totalSupply();
    return result.map((s: any) => ({
      denom: s.denom,
      amount: BigInt(s.amount.toString())
    }));
  }

  /**
   * Get total supply of a specific denomination
   */
  async supplyOf(denom: string): Promise<bigint> {
    const result = await this.contract.supplyOf(denom);
    return BigInt(result.toString());
  }

  /**
   * Get balance of a specific denomination for an account
   */
  async getBalance(account: string, denom: string): Promise<bigint | null> {
    const balances = await this.balances(account);
    const balance = balances.find(b => b.denom === denom);
    return balance ? balance.amount : null;
  }
}

// Usage
const provider = new ethers.JsonRpcProvider("YOUR_RPC_URL");
const bankClient = new BankPrecompileClient(provider);

// Get user balances
const balances = await bankClient.balances("0x742d35Cc6634C0532925a3b844Bc9e7595f25e4");

// Get total supply
const supply = await bankClient.totalSupply();

// Get specific denom supply
const badgeSupply = await bankClient.supplyOf("ubadge");
```

## Important Notes

1. **Read-Only**: The bank precompile is **read-only** - it only provides query functionality. It does NOT support sending tokens.

2. **Address**: Use `0x0000000000000000000000000000000000000804` (0x0804), NOT 0x1003.

3. **Denominations**: Denominations are returned as strings (e.g., "ubadge", "ustake"). The base denomination is typically the smallest unit.

4. **Amounts**: Amounts are returned as `uint256` in the smallest unit. For example:
   - 1 BADGE = 1,000,000,000 ubadge (1e9)
   - 1 STAKE = 1,000,000,000 ustake (1e9)

5. **Gas Costs**: These are view functions, so they don't consume gas when called from frontend (read-only calls).

6. **Error Handling**: If an account has no balance for a denomination, it won't appear in the `balances()` array. Use `getBalance()` helper to check for specific denominations.

## Sending Tokens

**The bank precompile does NOT support sending tokens.** To send tokens, you have a few options:

1. **Use Cosmos SDK directly**: Send tokens via the Cosmos SDK (e.g., using `cosmos.bank.v1beta1.MsgSend`)

2. **Use ERC20 tokens**: If tokens are wrapped as ERC20, use standard ERC20 transfer methods

3. **Use other precompiles**: Some custom precompiles may support token transfers

## Testing

```typescript
import { ethers } from "hardhat";

describe("Bank Precompile", () => {
  it("should query balances", async () => {
    const bankPrecompile = await ethers.getContractAt(
      "IBankPrecompile",
      "0x0000000000000000000000000000000000000804"
    );

    const [signer] = await ethers.getSigners();
    const balances = await bankPrecompile.balances(signer.address);
    
    expect(balances.length).to.be.greaterThan(0);
    console.log("Balances:", balances);
  });

  it("should query total supply", async () => {
    const bankPrecompile = await ethers.getContractAt(
      "IBankPrecompile",
      "0x0000000000000000000000000000000000000804"
    );

    const supply = await bankPrecompile.totalSupply();
    expect(supply.length).to.be.greaterThan(0);
  });

  it("should query specific denom supply", async () => {
    const bankPrecompile = await ethers.getContractAt(
      "IBankPrecompile",
      "0x0000000000000000000000000000000000000804"
    );

    const supply = await bankPrecompile.supplyOf("ubadge");
    expect(supply).to.be.greaterThan(0);
  });
});
```

## References

- [Cosmos EVM Bank Precompile Documentation](https://evm.cosmos.network/docs/next/documentation/smart-contracts/precompiles/bank)
- [Cosmos EVM Precompiles Overview](https://evm.cosmos.network/docs/next/documentation/smart-contracts/precompiles/overview)
