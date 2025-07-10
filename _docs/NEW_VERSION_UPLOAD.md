# New Version Upload Guide

This document outlines the process for releasing a new version of the BitBadges chain.

## Version Naming Convention

We use "v8", "v9", etc. naming conventions incremented by 1 each time.

## Release Flow

### 1. Tag and Push to Origin

Tag and push the current codebase as a git tag with the new version number:

```bash
git tag v9
git push origin v9
```

### 2. Build All Binaries

Run the build command to generate all platform binaries:

```bash
make build-all
```

### 3. Generate Release Information

Create a `release-info/v9` directory in the root and generate the following files:

#### 3.1 Release Notes

Create a `RELEASE_NOTES.md` file in the `release-info/v9` directory using the exact format below. Replace the version numbers, dates, and content as appropriate:

**NOTE: Release notes should use MAINNET information only (not testnet).**

**IMPORTANT: Use find-and-replace to update all version references:**

-   Replace all instances of `v8` with `v9` (or the current version)
-   Replace all instances of `[INSERT_BLOCK_HEIGHT]` with the actual mainnet block height
-   Replace all instances of `[INSERT_DATE_TIME]` with the actual mainnet date/time
-   Replace all instances of `[INSERT_SUMMARY]` with the actual summary
-   Replace all instances of `[INSERT_CHANGE_1]`, `[INSERT_CHANGE_2]`, etc. with actual changes

```markdown
[v8](https://github.com/BitBadges/bitbadgeschain/releases/tag/v8) [Latest](https://github.com/BitBadges/bitbadgeschain/releases/latest)

🔧 BitBadges Chain Upgrade — v9
📦 [[Release v9](https://github.com/BitBadges/bitbadgeschain/releases/tag/v9)](https://github.com/BitBadges/bitbadgeschain/releases/tag/v9) • 🆕 [[Latest Release](https://github.com/BitBadges/bitbadgeschain/releases/latest)](https://github.com/BitBadges/bitbadgeschain/releases/latest)

Upgrade Name: v9
Upgrade Block Height: [INSERT_BLOCK_HEIGHT]
Estimated Time: [INSERT_DATE_TIME]

Summary of update: [INSERT_SUMMARY]

🚨 Important Instructions for Node Operators
To successfully upgrade your node to v9 please follow these steps:

Download the New Binary
Download the latest v9 binary from the [[release page](https://github.com/BitBadges/bitbadgeschain/releases/tag/v9)](https://github.com/BitBadges/bitbadgeschain/releases/tag/v9).

Place Binary in the Correct Path
Move the new binary to:

<your_node_home>/cosmovisor/upgrades/v9/bin/
Ensure Correct Naming
Rename the binary file to:

bitbadgeschaind
This is required for nodes using the default Cosmovisor configuration.

Check Executable Permissions
Ensure the binary is executable:

chmod +x <your_node_home>/cosmovisor/upgrades/v9/bin/bitbadgeschaind
Verify Setup
Run the following to confirm the version and setup:

<your_node_home>/cosmovisor/upgrades/v9/bin/bitbadgeschaind version
Cosmovisor Will Auto-Switch
Cosmovisor will automatically switch to the new binary at the specified block height. If your node does not have the correct setup, it will halt and could be slashed.

📚 For full setup and operational details, see the [[Run a Node documentation](https://docs.bitbadges.io/for-developers/bitbadges-blockchain/run-a-node)](https://docs.bitbadges.io/for-developers/bitbadges-blockchain/run-a-node).

📝 Notable Changes in v9

-   [INSERT_CHANGE_1]
-   [INSERT_CHANGE_2]
-   [INSERT_CHANGE_3]
```

#### 3.2 Mainnet Proposal

Create a `mainnet-proposal.json` file in the `release-info/v9` directory:

**IMPORTANT: Use find-and-replace to update all version references:**

-   Replace all instances of `v9` with the current version number
-   Replace all instances of `[INSERT_BLOCK_HEIGHT]` with the actual mainnet block height

```json
{
    "messages": [
        {
            "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
            "authority": "bb10d07y265gmmuvt4z0w9aw880jnsr700jelmk2z",
            "plan": {
                "name": "v9",
                "height": "[INSERT_BLOCK_HEIGHT]",
                "info": "Upgrade to v9",
                "upgraded_client_state": null
            }
        }
    ],
    "expedited": true,
    "deposit": "1000000000000ustake",
    "title": "Upgrade to v9",
    "summary": "This proposal upgrades the chain to version v9."
}
```

#### 3.3 Testnet Proposal

Create a `testnet-proposal.json` file in the `release-info/v9` directory:

**IMPORTANT: Use find-and-replace to update all version references:**

-   Replace all instances of `v9` with the current version number
-   Replace all instances of `[INSERT_TESTNET_BLOCK_HEIGHT]` with the actual testnet block height

```json
{
    "messages": [
        {
            "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
            "authority": "bb10d07y265gmmuvt4z0w9aw880jnsr700jelmk2z",
            "plan": {
                "name": "v9",
                "height": "[INSERT_TESTNET_BLOCK_HEIGHT]",
                "info": "Upgrade to v9",
                "upgraded_client_state": null
            }
        }
    ],
    "expedited": true,
    "deposit": "1000000000000ustake",
    "title": "Upgrade to v9",
    "summary": "This proposal upgrades the chain to version v9."
}
```

#### 3.4 Discord Announcement

Create a `discord-announcement.md` file in the `release-info/v9` directory:

**NOTE: Discord announcement should use MAINNET information only (not testnet).**

**IMPORTANT: Use find-and-replace to update all version references:**

-   Replace all instances of `v9` with the current version number
-   Replace all instances of `[INSERT_BLOCK_HEIGHT]` with the actual mainnet block height
-   Replace all instances of `[INSERT_DATE_TIME]` with the actual mainnet date/time
-   Replace all instances of `[INSERT_PROPOSAL_ID]` with the actual mainnet proposal ID
-   Replace all instances of `[INSERT_SUMMARY]` with the actual summary

```markdown
Upgrade Name: v9
Upgrade Block Height: [INSERT_BLOCK_HEIGHT]
Estimated Time: [INSERT_DATE_TIME]
Voting Period: Now + 24 hours (https://explorer.bitbadges.io/BitBadges%20Mainnet/gov/[INSERT_PROPOSAL_ID])
Binaries / Instructions: https://github.com/BitBadges/bitbadgeschain/releases/tag/v9
Summary of Upgrade: [INSERT_SUMMARY]

@Validator
```

## Important Notes

-   Always increment the version number by 1 from the previous release
-   **Use find-and-replace to update ALL version references** (v8 → v9, v9 → v10, etc.)
-   Update all version references in the release notes template
-   Ensure the upgrade block height and estimated time are accurate
-   Include a comprehensive summary of changes in the release notes
-   Test the build process before creating the release
-   Verify that all binaries are properly generated for different platforms
-   **Double-check that no old version numbers remain in any generated files**
-   **Create version-specific subfolders** (e.g., `release-info/v9/`, `release-info/v10/`) to keep releases organized
