# MsgTransferBadges

Executes badge transfers between addresses.

## Proto Definition

```protobuf
message MsgTransferBadges {
  string creator = 1; // Address initiating the transfer
  string collectionId = 2; // Collection containing badges to transfer
  repeated Transfer transfers = 3; // Transfer operations (must pass approvals)
}

message MsgTransferBadgesResponse {}

message Transfer {
  // The address of the sender of the transfer.
  string from = 1;
  // The addresses of the recipients of the transfer.
  repeated string toAddresses = 2;
  // The balances to be transferred.
  repeated Balance balances = 3;
  // If defined, we will use the predeterminedBalances from the specified approval to calculate the balances at execution time.
  // We will override the balances field with the precalculated balances. Only applicable for approvals with predeterminedBalances set.
  ApprovalIdentifierDetails precalculateBalancesFromApproval = 4;
  // The Merkle proofs / solutions for all Merkle challenges required for the transfer.
  repeated MerkleProof merkleProofs = 5;
  // The ETH signature proofs / solutions for all ETH signature challenges required for the transfer.
  repeated ETHSignatureProof ethSignatureProofs = 6;
  // The memo for the transfer.
  string memo = 7;
  // The prioritized approvals for the transfer. By default, we scan linearly through the approvals and use the first match.
  // This field can be used to prioritize specific approvals and scan through them first.
  repeated ApprovalIdentifierDetails prioritizedApprovals = 8;
  // Whether to only check prioritized approvals for the transfer.
  // If true, we will only check the prioritized approvals and fail if none of them match (i.e. do not check any non-prioritized approvals).
  // If false, we will check the prioritized approvals first and then scan through the rest of the approvals.
  bool onlyCheckPrioritizedCollectionApprovals = 9;
  // Whether to only check prioritized approvals for the transfer.
  // If true, we will only check the prioritized approvals and fail if none of them match (i.e. do not check any non-prioritized approvals).
  // If false, we will check the prioritized approvals first and then scan through the rest of the approvals.
  bool onlyCheckPrioritizedIncomingApprovals = 10;
  // Whether to only check prioritized approvals for the transfer.
  // If true, we will only check the prioritized approvals and fail if none of them match (i.e. do not check any non-prioritized approvals).
  // If false, we will check the prioritized approvals first and then scan through the rest of the approvals.
  bool onlyCheckPrioritizedOutgoingApprovals = 11;
  // The options for precalculating the balances.
  PrecalculationOptions precalculationOptions = 12;
  // Affiliate address for the transfer.
  string affiliateAddress = 13;
  // The number of times to attempt approval validation. If 0 / not specified, we default to only one.
  string numAttempts = 14;
}

message PrecalculationOptions {
  // The timestamp to override with when calculating the balances.
  string overrideTimestamp = 1;
  // The badgeIdsOverride to use for the transfer.
  repeated UintRange badgeIdsOverride = 2;
}
```

## Auto-Scan vs Prioritized Approvals

The transfer approval system operates in two modes to balance efficiency and precision:

### Auto-Scan Mode (Default)

By default, the system automatically scans through available approvals to find a match for the transfer. This mode:

-   **Works with**: Approvals using [Empty Approval Criteria](../examples/empty-approval-criteria.md) (no side effects)
-   **Behavior**: Automatically finds and uses the first matching approval
-   **Use case**: Simple transfers without custom logic or side effects
-   **No versioning required**: The system handles approval selection automatically

### Prioritized Approvals (Required for Side Effects)

**CRITICAL REQUIREMENT**: Any transfer with side effects or custom approval criteria MUST always be prioritized with proper versioning set. No exceptions.

#### Race Condition Protection

The versioning control ensures that before submitting, the user knows the exact approval they are using:

```typescript
"prioritizedApprovals": [
    {
        "approvalId": "abc123",
        "approvalLevel": "collection",
        "approverAddress": "",
        "version": "2" // Must specify exact version
    }
]
```

#### Example: Coin Transfer Approval

```typescript
// MUST be prioritized - has coin transfer side effects
"prioritizedApprovals": [
    {
        "approvalId": "reward-approval",
        "approvalLevel": "collection",
        "approverAddress": "",
        "version": "1"
    }
],
"onlyCheckPrioritizedCollectionApprovals": true
```

#### Example: Auto-Scan Safe Transfer

```typescript
// Can use auto-scan - no side effects
"prioritizedApprovals": [], // Empty - will auto-scan

// Only will succeed if it finds an approval has empty approval criteria with no custom logic
```

### Control Flags

-   `onlyCheckPrioritizedCollectionApprovals`: If true, only check prioritized approvals
-   `onlyCheckPrioritizedIncomingApprovals`: If true, only check prioritized incoming approvals
-   `onlyCheckPrioritizedOutgoingApprovals`: If true, only check prioritized outgoing approvals

**Setting these to `true` is recommended when using prioritized approvals to ensure deterministic behavior.**

### Related Documentation

-   [Empty Approval Criteria](../examples/empty-approval-criteria.md) - Template for auto-scan compatible approvals
-   [Approval Criteria](../concepts/approval-criteria/README.md) - Understanding approval complexity
-   [Coin Transfers](../concepts/approval-criteria/usdbadge-transfers.md) - Side effect examples

## Transfer Validation Process

Each transfer undergoes a systematic validation process to ensure security and proper authorization:

### Validation Steps

```
PRE. CALCULATE BALANCES (if needed)
  └── If precalculateBalancesFromApproval is specified, we will use the predeterminedBalances from the specified approval to pre-calculate the balances at execution time.

1. BALANCE CHECK
   └── Verify sender has sufficient badge balances for the transfer including ownership times
   └── FAIL if insufficient balances

2. COLLECTION APPROVAL CHECK
   └── Scan collection-level approvals (prioritized first, then auto-scan) to find a match for the entire transfer
   └── If match found:
       ├── Check approval criteria (merkle proofs, amounts, timing, etc.) and constraints
       ├── Check if it overrides sender approvals (overridesFromOutgoingApprovals)
       ├── Check if it overrides recipient approvals (overridesToIncomingApprovals)
       └── PROCEED with override flags set
   └── Else:
        └── Continue scanning
   └── If some attempted transfer balances have no valid collection approval: FAIL

3. SENDER APPROVAL CHECK (if not overridden)
   └── Check sender's outgoing approvals for this transfer
   └── Verify approval criteria and constraints
   └── FAIL if no valid outgoing approval found

4. RECIPIENT APPROVAL CHECK (if not overridden)
   └── Check each recipient's incoming approvals
   └── Verify approval criteria and constraints
   └── FAIL if any recipient lacks valid incoming approval

5. EXECUTE TRANSFER
   └── Update badge balances
   └── Execute any approved side effects
   └── Emit transfer events
   └── SUCCESS
```

### Override Behavior

Collection approvals can override user-level approvals:

-   **`overridesFromOutgoingApprovals: true`** - Forcefully skips sender approval check
-   **`overridesToIncomingApprovals: true`** - Forcefully skips recipient approval checks

This allows collection managers to enable transfers that would otherwise be blocked by user settings.

### Failure Points

Transfers fail at the first validation step that doesn't pass:

1. **Insufficient Balances** - Sender doesn't own the badges
2. **No Collection Approval** - No valid collection-level approval found
3. **Blocked by Sender** - Sender's outgoing approvals reject the transfer
4. **Blocked by Recipient** - Recipient's incoming approvals reject the transfer

### ETH Signature Proofs

ETH Signature Proofs are required when transfers use [ETH Signature Challenges](../concepts/approval-criteria/eth-signature-challenges.md). Each proof contains:

- **`nonce`**: The unique identifier that was signed
- **`signature`**: The Ethereum signature of the message `nonce + "-" + creatorAddress`

**Important**: Each signature can only be used once per challenge tracker. The system tracks used signatures to prevent replay attacks.

### Related Documentation

-   [Transferability / Approvals](../concepts/transferability-approvals.md) - Approval system overview
-   [Collection Approvals](../concepts/approval-criteria/README.md) - Collection-level controls
-   [User Approvals](../examples/building-user-approvals.md) - User-level settings
-   [ETH Signature Challenges](../concepts/approval-criteria/eth-signature-challenges.md) - Ethereum signature requirements

## Collection ID Auto-Lookup

If you specify `collectionId` as `"0"`, it will automatically lookup the latest collection ID created. This can be used if you are creating a collection and do not know the official collection ID yet but want to perform a multi-message transaction.

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges transfer-badges '[tx-json]' --from sender-key
```

### JSON Example

```json
{
    "creator": "bb1initiator123...",
    "collectionId": "1",
    "transfers": [
        {
            "from": "bb1sender123...",
            "toAddresses": ["bb1recipient123..."],
            // Balances to transfer (can be left blank if you are using precalculateBalancesFromApproval)
            "balances": [
                {
                    "amount": "10",
                    "ownershipTimes": [
                        { "start": "1", "end": "18446744073709551615" }
                    ],
                    "badgeIds": [{ "start": "1", "end": "5" }]
                }
            ],
            // Specific approval to calculate balances dynamically for (from the approvalCriteria.predeterminedBalances)
            "precalculateBalancesFromApproval": {
                "approvalId": "",
                "approvalLevel": "",
                "approverAddress": "",
                "version": "0"
            },
            // Additional options dependent on what is allowed (e.g. allow timestamp override, badge ID override, etc.)
            "precalculationOptions": {
                "overrideTimestamp": "0",
                "badgeIdsOverride": []
            },
            // Supply all merkle proofs for any merkle challenges that need to be satisfied
            "merkleProofs": [],
            // Supply all ETH signature proofs for any ETH signature challenges that need to be satisfied
            "ethSignatureProofs": [
                {
                    "nonce": "unique-nonce-001",
                    "signature": "0x..."
                }
            ],
            // Memo for the transfer (can be left blank)
            "memo": "",

            // Any approval IDs that you want to prioritize for this transfer
            // Note: All approvals with side effects must be prioritized with proper versioning
            "prioritizedApprovals": [
                {
                    "approvalId": "abc123",
                    "approvalLevel": "collection",
                    "approverAddress": "", // blank for collection, otherwise the address of the approver
                    "version": "0"
                }
            ],

            // If specified, we will stop checking after the prioritized approvals list.
            // If false, we will check prioritized first, but then continue to check the rest of the approvals in auto-scan mode
            "onlyCheckPrioritizedCollectionApprovals": false,
            "onlyCheckPrioritizedIncomingApprovals": false,
            "onlyCheckPrioritizedOutgoingApprovals": false,

            // Add your address if you want to claim part of the protocol fee
            "affiliateAddress": "",
            // Number of times to attempt this transfer (default is 1, 0 is empty and also defaults to 1)
            // Use this if you want to try this transfer multiple times
            "numAttempts": "1"
        }
    ]
}
```
