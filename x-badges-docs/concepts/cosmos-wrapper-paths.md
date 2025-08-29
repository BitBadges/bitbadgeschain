# Cosmos Wrapper Paths

Cosmos Wrapper Paths enable 1:1 wrapping between BitBadges tokens and native Cosmos SDK coin asset types, making tokens IBC-compatible. These paths automatically mint and burn tokens when transferring to/from specific wrapper addresses. These transfers to/from are handled within the badges module, so you can set up customizable logic for how these transfers are handled.

> **Important**: Since wrapper addresses are uncontrollable (no private keys), approval design requires careful consideration. You must override the wrapper address's user-level approvals where necessary using collection approvals to ensure wrapping/unwrapping functions properly.

### Auto-Generating Wrapper Addresses

You can programmatically generate wrapper addresses using the `bitbadgesjs-sdk` npm package. Note that the address is generated based on the denom set. It is just the custom denom, not the full `badges:collectionId:denom` format.

```typescript
import { generateAliasAddressForDenom } from 'bitbadgesjs-sdk';

const denom = 'utoken1';
const wrapperAddress = generateAliasAddressForDenom(denom);
console.log('Wrapper Address:', wrapperAddress);
```

## Proto Definition

```protobuf
message CosmosCoinWrapperPathAddObject {
  string denom = 1;
  repeated Balance balances = 2;
  string symbol = 3;
  repeated DenomUnit denomUnits = 4;
  bool allowOverrideWithAnyValidToken = 5;
}

message DenomUnit {
  string decimals = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
  string symbol = 2;
  bool isDefaultDisplay = 3;
}
```

## Configuration Fields

### Denom

-   **Base denomination** - The fundamental unit name for the wrapped coin. For the full Cosmos denomination, it will be "badges:collectionId:denom"
-   **{id} placeholder support** - You can use `{id}` in the denom to dynamically replace it with the actual badge ID during transfers if allowOverrideWithAnyValidToken is true

### Symbol

-   **Display symbol** - Human-readable symbol for the wrapped asset

### Balances

-   **Custom conversion rates** - Defines which tokens and ownership times participate in wrapping and how many tokens are wrapped for each native coin unit

### Denomination Units

Multiple denomination units allow for different display formats:

#### Decimals

-   **Precision level** - Number of decimal places for this unit
-   **Conversion factor** - How this unit relates to the base denomination

#### Symbol

-   **Unit symbol** - Symbol for this specific denomination unit
-   **Different from base** - Can differ from the main symbol
-   **Context-specific** - Used in appropriate contexts (micro, milli, etc.)

#### Default Display

-   **Primary unit** - Which unit is shown by default in interfaces
-   **Only one default** - Only one unit can be marked as default display. If none are marked as default, the base level with 0 decimals is shown by default.
-   **User experience** - Determines what users see first

### Allow Override With Any Valid Token

-   **Dynamic badge ID handling** - When `true`, allows the wrapper to accept any SINGLE valid badge ID from the collection's `validBadgeIds` range
-   **Override conversion balances** - Temporarily overrides the `balances` field's badge ID ranges with the actual token ID being transferred
-   **Validation required** - The transferred token ID must be within the collection's valid badge ID range
-   **Use case** - Useful when you want a single wrapper path to handle multiple badge IDs dynamically

#### Example with {id} Placeholder

```json
{
    "denom": "utoken{id}",
    "symbol": "TOKEN:{id}",
    "balances": [
        {
            "amount": "1",
            "badgeIds": [{ "start": "1", "end": "1" }],
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }]
        }
    ],
    "denomUnits": [
        {
            "decimals": "6",
            "symbol": "TOKEN",
            "isDefaultDisplay": true
        }
    ],
    "allowOverrideWithAnyValidToken": true
}
```

When transferring badge ID 5, the final denomination becomes `utoken5`.

### Allow Override With Any Valid Token

This feature provides flexibility in handling different badge IDs within a single wrapper path.

#### How It Works

1. **Validation Check** - Verifies the transferred token ID is within the collection's `validBadgeIds` range
2. **Dynamic Override** - Temporarily replaces the wrapper's badge ID ranges with the actual token ID being transferred
3. **Conversion** - Uses the overridden badge ID for the conversion calculation

#### Example with Override Enabled

```json
{
    "denom": "utoken",
    "symbol": "TOKEN",
    "balances": [
        {
            "amount": "1",
            "badgeIds": [{ "start": "1", "end": "1" }], // These technically don't matter since we can override the badge IDs during conversion
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }]
        }
    ],
    "denomUnits": [
        {
            "decimals": "6",
            "symbol": "TOKEN",
            "isDefaultDisplay": true
        }
    ],
    "allowOverrideWithAnyValidToken": true
}
```

This wrapper can now handle any badge ID from 1-100 (assuming the collection's `validBadgeIds` includes that range), dynamically overriding the badge ID during conversion.

## Usage Examples

### Basic Wrapper Path

```json
{
    "denom": "utoken1",
    "symbol": "TOKEN1",
    "balances": [
        {
            "amount": "1",
            "badgeIds": [{ "start": "1", "end": "1" }],
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }]
        }
    ],
    "denomUnits": [
        {
            "decimals": "6",
            "symbol": "TOKEN1",
            "isDefaultDisplay": true
        }
    ],
    "allowOverrideWithAnyValidToken": false
}
```

### Multi-Unit Display System

```json
{
    "denom": "utoken",
    "symbol": "TOKEN",
    "balances": [
        {
            "amount": "1",
            "badgeIds": [{ "start": "1", "end": "100" }],
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }]
        }
    ],
    "denomUnits": [
        {
            "decimals": "3",
            "symbol": "mtoken",
            "isDefaultDisplay": false
        },
        {
            "decimals": "6",
            "symbol": "TOKEN",
            "isDefaultDisplay": true
        }
    ],
    "allowOverrideWithAnyValidToken": false
}
```

This creates a system where:

-   `utoken` is the base unit (smallest denomination)
-   `mtoken` = 1,000 `utoken` (milli-token)
-   `TOKEN` = 1,000,000 `utoken` (full token, default display)

## Use Cases

### IBC Transfers

-   **Cross-chain transfers** - Send wrapped tokens to other Cosmos chains
-   **DeFi integration** - Use wrapped tokens in Cosmos DeFi protocols
-   **Liquidity provision** - Add wrapped tokens to AMM pools

### Multi-Chain Ecosystems

-   **Ecosystem bridges** - Connect BitBadges to broader Cosmos ecosystem
-   **Shared liquidity** - Participate in cross-chain liquidity pools
-   **Governance tokens** - Use wrapped tokens in governance across chains

### Trading and Exchange

-   **DEX compatibility** - Trade on Cosmos-native decentralized exchanges
-   **Price discovery** - Enable market-driven price discovery
-   **Arbitrage opportunities** - Cross-chain arbitrage possibilities

### Featured Use Case: List on Osmosis

With BitBadges' existing relayer infrastructure and IBC-compatible wrapped denominations, listing wrapped tokens on Osmosis is streamlined:

-   **IBC Relayer Ready** - BitBadges already has relayer infrastructure set up for seamless cross-chain transfers
-   **Native IBC Compatibility** - Wrapped tokens become native SDK coins that work seamlessly with IBC protocols
-   **Automatic Liquidity** - Create liquidity pools on Osmosis DEX with wrapped token assets
-   **Streamlined Process** - The technical infrastructure eliminates common barriers to cross-chain trading
-   **Enhanced Discoverability** - Tokens gain exposure to the broader Cosmos DeFi ecosystem

## Conversion Process

### Token to Coin (Wrapping)

1. User transfers tokens to the wrapper address
2. System processes the denom (replaces {id} if present, validates override if enabled)
3. System burns the tokens from user's balance
4. System mints equivalent native coins
5. Coins are credited to the user's account

### Coin to Token (Unwrapping)

1. User transfers coins to the wrapper address
2. System processes the denom (replaces {id} if present, validates override if enabled)
3. System burns the native coins
4. System mints equivalent tokens
5. Tokens are credited to the user's balance

Cosmos Wrapper Paths provide seamless interoperability between BitBadges and the broader Cosmos ecosystem while maintaining the unique properties of both token and coin systems.
