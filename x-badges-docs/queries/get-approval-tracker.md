# GetApprovalTracker

Retrieves tracking information for approval usage.

## Proto Definition

```protobuf
message QueryGetApprovalTrackerRequest {
  string amountTrackerId = 1; 
  string approvalLevel = 2; // "collection", "incoming", or "outgoing"
  string approverAddress = 3; // Leave blank if approvalLevel is "collection"
  string trackerType = 4; // "overall", "to", "from", "initiatedBy"
  string collectionId = 5;
  string approvedAddress = 6; // Leave blank if trackerType is "overall"
  string approvalId = 7;
}

message QueryGetApprovalTrackerResponse {
  ApprovalTracker tracker = 1;
}

message ApprovalTracker {
  string numTransfers = 1; // Number of transfers that have been processed
  repeated Balance amounts = 2; // Cumulative balances associated with processed transfers
  string lastUpdatedAt = 3; // Last updated at time (UNIX millisecond timestamp)
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-approval-tracker [collectionId] [approvalLevel] [approverAddress] [approvalId] [amountTrackerId] [trackerType] [approvedAddress]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_approvals_tracker/1/outgoing/bb1.../approval-1/tracker-1/overall/"
```

### Response Example
```json
{
  "tracker": {
    "numTransfers": "5",
    "amounts": [
      {
        "amount": "100",
        "badgeIds": [{"start": "1", "end": "10"}],
        "ownershipTimes": [{"start": "1672531200000", "end": "18446744073709551615"}]
      }
    ],
    "lastUpdatedAt": "1672531200000"
  }
}
```