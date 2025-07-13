# Params

Retrieves the current module parameters.

## Proto Definition

```protobuf
message QueryParamsRequest {}

message QueryParamsResponse {
  Params params = 1;
}

message Params {
  // Array of allowed denominations for fee payments and escrow operations
  repeated string allowed_denoms = 1;
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges params

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/params"
```

### Response Example
```json
{
  "params": {
    "allowedDenoms": ["ubadge", "ibc/1234567890"]
  }
}
```
