# GetChallengeTracker

Retrieves the number of times a given leaf has been used for a specific challenge tracker.

## Proto Definition

```protobuf
message QueryGetChallengeTrackerRequest {
  string collectionId = 1;
  string approvalLevel = 2; // "collection", "incoming", or "outgoing"
  string approverAddress = 3; // Leave blank if approvalLevel is "collection"
  string challengeTrackerId = 4;
  string leafIndex = 5;
  string approvalId = 6;
}

message QueryGetChallengeTrackerResponse {
  string numUsed = 1; // Number of times this leaf has been used
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-challenge-tracker [collectionId] [approvalLevel] [approverAddress] [approvalId] [challengeTrackerId] [leafIndex]

# REST API
# Note for blank values, use "" so you may have // in the query
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_challenge_tracker/1/collection//approval-123/challenge-1/42"
```

### Response Example

```json
{
    "numUsed": "1"
}
```
