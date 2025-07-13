# MsgTransferBadges

Executes badge transfers between addresses.

## Collection ID Auto-Lookup

If you specify `collectionId` as `"0"`, it will automatically lookup the latest collection ID created. This can be used if you are creating a collection and do not know the official collection ID yet but want to perform a multi-message transaction.

## Proto Definition

```protobuf
message MsgTransferBadges {
  string creator = 1; // Address initiating the transfer
  string collectionId = 2; // Collection containing badges to transfer
  repeated Transfer transfers = 3; // Transfer operations (must pass approvals)
}

message MsgTransferBadgesResponse {}
```

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
      "balances": [
        {
          "amount": "10",
          "ownershipTimes": [{"start": "1", "end": "18446744073709551615"}],
          "badgeIds": [{"start": "1", "end": "5"}]
        }
      ],
      "precalculateBalancesFromApproval": {
        "approvalId": "",
        "approvalLevel": "",
        "approverAddress": "",
        "version": "0"
      },
      "merkleProofs": [],
      "memo": "",
      "prioritizedApprovals": [],
      "onlyCheckPrioritizedCollectionApprovals": false,
      "onlyCheckPrioritizedIncomingApprovals": false,
      "onlyCheckPrioritizedOutgoingApprovals": false,
      "precalculationOptions": {
        "overrideTimestamp": "0",
        "badgeIdsOverride": []
      },
      "affiliateAddress": "",
      "numAttempts": "1"
    }
  ]
}
```