# GetETHSignatureTracker

Retrieves the number of times a given signature has been used for a specific ETH signature challenge tracker.

## Proto Definition

```protobuf
message QueryGetETHSignatureTrackerRequest {
  string collectionId = 1;
  string approvalLevel = 2; // "collection", "incoming", or "outgoing"
  string approverAddress = 3; // Leave blank if approvalLevel is "collection"
  string approvalId = 4;
  string challengeTrackerId = 5;
  string signature = 6;
}

message QueryGetETHSignatureTrackerResponse {
  string numUsed = 1; // Number of times this signature has been used
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-num-used-for-eth-signature-challenge [collectionId] [approvalLevel] [approverAddress] [approvalId] [challengeTrackerId] [signature]

# REST API
# Note for blank values, use "" so you may have // in the query
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_eth_signature_tracker/1/collection//approval-123/challenge-1/0x1234567890abcdef..."
```

### Response Example

```json
{
    "numUsed": "1"
}
```

## Notes

- Each signature can only be used once per challenge tracker
- If a signature has never been used, the response will be "0"
- The signature parameter should be the full Ethereum signature (0x-prefixed hex string) 