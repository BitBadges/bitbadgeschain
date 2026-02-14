# Chain ID Registration Guide

This document outlines the process for claiming and registering EVM chain IDs for BitBadges mainnet and testnet.

## Current Chain IDs

- **Mainnet**: `50024` (BitBadges Mainnet)
- **Testnet**: `50025` (BitBadges Testnet)

These chain IDs are defined in `app/params/constants.go` and configured in the genesis files (`config.yml` for mainnet, `config.testnet.yml` for testnet).

## Required Registries

To properly claim and register EVM chain IDs, you need to submit pull requests to the following registries:

### 1. Primary Registry: ethereum-lists/chains

**Repository**: [ethereum-lists/chains](https://github.com/ethereum-lists/chains)

This is the **authoritative source** for EVM chain metadata and powers:
- [chainid.network](https://chainid.network/)
- [ChainList](https://chainlist.org/)
- MetaMask and other wallet providers
- Various DeFi tools and explorers

#### Steps to Register:

1. **Fork the repository** on GitHub
2. **Create JSON files** in the `_data/chains` directory:
   - Filename format: `e<chainId>.json` (e.g., `e50024.json` for mainnet, `e50025.json` for testnet)
   - Follow the [CAIP-2](https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-2.md) naming convention

3. **Required JSON structure** (example for mainnet):

```json
{
  "name": "BitBadges",
  "chain": "BitBadges",
  "rpc": [
    "https://rpc.bitbadges.io",
    "https://bitbadges-rpc.publicnode.com"
  ],
  "faucets": [],
  "nativeCurrency": {
    "name": "Badge",
    "symbol": "BADGE",
    "decimals": 18
  },
  "infoURL": "https://bitbadges.io",
  "shortName": "bitbadges",
  "chainId": 50024,
  "networkId": 50024,
  "explorers": [
    {
      "name": "BitBadges Explorer",
      "url": "https://explorer.bitbadges.io",
      "standard": "EIP3091"
    }
  ],
  "testnet": false
}
```

4. **For testnet** (`e50025.json`):

```json
{
  "name": "BitBadges Testnet",
  "chain": "BitBadges",
  "rpc": [
    "https://testnet-rpc.bitbadges.io",
    "https://bitbadges-testnet-rpc.publicnode.com"
  ],
  "faucets": [
    "https://faucet.bitbadges.io"
  ],
  "nativeCurrency": {
    "name": "Badge",
    "symbol": "BADGE",
    "decimals": 18
  },
  "infoURL": "https://bitbadges.io",
  "shortName": "bitbadges-testnet",
  "chainId": 50025,
  "networkId": 50025,
  "explorers": [
    {
      "name": "BitBadges Testnet Explorer",
      "url": "https://testnet-explorer.bitbadges.io",
      "standard": "EIP3091"
    }
  ],
  "testnet": true
}
```

5. **Validation**:
   - Run `./gradlew run` to validate the JSON structure
   - Ensure Prettier formatting is applied
   - Verify that `shortName` and `name` are unique (no conflicts with existing chains)

6. **Submit Pull Request**:
   - Create a PR with a clear title: "Add BitBadges Mainnet (50024) and Testnet (50025)"
   - Include links to your chain's documentation and explorer
   - Wait for maintainer review and approval

#### Important Constraints:

- **Uniqueness**: Your `shortName` and chain `name` must be unique across all registered chains
- **No Reuse**: Chain IDs cannot be reused - the first PR claiming an ID gets it assigned
- **Permanent**: Once assigned, a chain ID can only be reused if the old chain is deprecated
- **Validation Required**: All JSON files must pass validation checks before merging

### 2. Optional: Cosmos Chain Registry

**Repository**: [cosmos/chain-registry](https://github.com/cosmos/chain-registry)

If BitBadges has Cosmos SDK integration (which it does), you may also want to register in the Cosmos chain registry for cross-chain interoperability.

#### Steps:

1. Create a directory: `chain-registry/bitbadges/`
2. Add `chain.json` with Cosmos-specific metadata
3. Include EVM chain ID information in the metadata

## Verification

After registration:

1. **Check chainid.network**: Verify your chain appears at `https://chainid.network/`
2. **Check ChainList**: Verify your chain appears at `https://chainlist.org/`
3. **Test MetaMask**: Try adding your network to MetaMask using the chain ID
4. **Test WalletConnect**: Verify wallet connections work correctly

## Updating Chain Information

If you need to update chain information (RPC endpoints, explorers, etc.):

1. Submit a new PR to `ethereum-lists/chains` with updated JSON files
2. Update the corresponding files in this repository
3. Update any documentation that references the old information

## Resources

- [ethereum-lists/chains README](https://github.com/ethereum-lists/chains#readme)
- [ChainList Documentation](https://github.com/ethereum-lists/chains#chainlist)
- [EIP-155: Simple replay attack protection](https://eips.ethereum.org/EIPS/eip-155)
- [CAIP-2: Blockchain ID Specification](https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-2.md)

## Notes

- Chain IDs are **permanent** once claimed - choose carefully
- The registration process may take several days for review
- Ensure all RPC endpoints and explorers are operational before submitting
- Keep your chain information up-to-date in the registry

