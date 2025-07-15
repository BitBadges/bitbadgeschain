# Extending Approvals (Advanced)

For edge cases requiring advanced functionality beyond native interfaces.

## When to Extend

Consider extending when you need:

-   Cross-approval functionality
-   Access to other blockchain data/modules
-   Custom logic not covered by native options

## Before Extending

1. **Consider Workarounds**: Many approvals can be adapted to fit native interfaces
2. **Design Alternatives**: Think creatively about using existing features
3. **Evaluate Necessity**: Ensure extension is truly required

## Implementation Options

### CosmWASM Smart Contracts

Build custom smart contracts that call into the x/badges module.

### EVM Smart Contracts (Coming Soon)

Build custom EVM contracts that call into the x/badges module.
