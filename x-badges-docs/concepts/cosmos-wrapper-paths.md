# Cosmos Wrapper Paths

Cosmos Wrapper Paths enable 1:1 wrapping between BitBadges badges and native Cosmos SDK coin asset types, making badges IBC-compatible. These paths automatically mint and burn badges when transferring to/from specific wrapper addresses. These transfers to/from are handled within the badges module, so you can set up customizable logic for how these transfers are handled.

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

### Symbol

-   **Display symbol** - Human-readable symbol for the wrapped asset

### Balances

-   **Custom conversion rates** - Defines which badges and ownership times participate in wrapping and how many badges are wrapped for each native coin unit

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
    ]
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
    ]
}
```

This creates a system where:

-   `utoken` is the base unit (smallest denomination)
-   `mtoken` = 1,000 `utoken` (milli-token)
-   `TOKEN` = 1,000,000 `utoken` (full token, default display)

## Use Cases

### IBC Transfers

-   **Cross-chain transfers** - Send wrapped badges to other Cosmos chains
-   **DeFi integration** - Use wrapped badges in Cosmos DeFi protocols
-   **Liquidity provision** - Add wrapped badges to AMM pools

### Multi-Chain Ecosystems

-   **Ecosystem bridges** - Connect BitBadges to broader Cosmos ecosystem
-   **Shared liquidity** - Participate in cross-chain liquidity pools
-   **Governance tokens** - Use wrapped badges in governance across chains

### Trading and Exchange

-   **DEX compatibility** - Trade on Cosmos-native decentralized exchanges
-   **Price discovery** - Enable market-driven price discovery
-   **Arbitrage opportunities** - Cross-chain arbitrage possibilities

### Featured Use Case: List on Osmosis

With BitBadges' existing relayer infrastructure and IBC-compatible wrapped denominations, listing wrapped badges on Osmosis is streamlined:

-   **IBC Relayer Ready** - BitBadges already has relayer infrastructure set up for seamless cross-chain transfers
-   **Native IBC Compatibility** - Wrapped badges become native SDK coins that work seamlessly with IBC protocols
-   **Automatic Liquidity** - Create liquidity pools on Osmosis DEX with wrapped badge assets
-   **Streamlined Process** - The technical infrastructure eliminates common barriers to cross-chain trading
-   **Enhanced Discoverability** - Badges gain exposure to the broader Cosmos DeFi ecosystem

## Conversion Process

### Badge to Coin (Wrapping)

1. User transfers badges to the wrapper address
2. System burns the badges from user's balance
3. System mints equivalent native coins
4. Coins are credited to the user's account

### Coin to Badge (Unwrapping)

1. User transfers coins to the wrapper address
2. System burns the native coins
3. System mints equivalent badges
4. Badges are credited to the user's balance

Cosmos Wrapper Paths provide seamless interoperability between BitBadges and the broader Cosmos ecosystem while maintaining the unique properties of both badge and coin systems.
