# GetBalance

Retrieves balances for a specific address in a collection.

## Proto Definition

```protobuf
message QueryGetBalanceRequest {
  string collectionId = 1; // Collection ID to query
  string address = 2; // Address to get balances for
}

message QueryGetBalanceResponse {
  UserBalanceStore balance = 1;
}

message UserBalanceStore {
  repeated Balance balances = 1; // List of balances associated with this user
  repeated UserOutgoingApproval outgoingApprovals = 2; // Approved outgoing transfers
  repeated UserIncomingApproval incomingApprovals = 3; // Approved incoming transfers
  bool autoApproveSelfInitiatedOutgoingTransfers = 4; // Auto-approve self-initiated outgoing transfers
  bool autoApproveSelfInitiatedIncomingTransfers = 5; // Auto-approve self-initiated incoming transfers
  bool autoApproveAllIncomingTransfers = 6; // Auto-approve all incoming transfers
  UserPermissions userPermissions = 7; // Permissions for this user's actions
}

// See all the proto definitions [here](https://github.com/BitBadges/bitbadgeschain/tree/master/proto/badges)
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-balance [collection-id] [address]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_balance/1/bb1..."
```

### Response Example

```json
{
    "balance": {
        "balances": [
            {
                "amount": "1",
                "tokenIds": [{ "start": "1", "end": "1" }],
                "ownershipTimes": [
                    { "start": "1672531200000", "end": "18446744073709551615" }
                ]
            }
        ],
        "outgoingApprovals": [
            // ...
        ],
        "incomingApprovals": [
            // ...
        ],
        "autoApproveSelfInitiatedOutgoingTransfers": true,
        "autoApproveSelfInitiatedIncomingTransfers": true,
        "autoApproveAllIncomingTransfers": true,
        "userPermissions": {
            // ...
        }
    }
}
```
