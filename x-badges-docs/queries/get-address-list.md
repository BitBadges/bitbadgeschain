# GetAddressList

Retrieves information about a specific address list.

## Proto Definition

```protobuf
message QueryGetAddressListRequest {
  string listId = 1; // ID of address list to retrieve
}

message QueryGetAddressListResponse {
  AddressList list = 1;
}

message AddressList {
  string listId = 1; // Unique identifier for the address list
  repeated string addresses = 2; // List of addresses included in the list
  bool whitelist = 3; // Whether list includes (true) or excludes (false) specified addresses
  string uri = 4; // URI providing metadata, if applicable
  string customData = 5; // Custom arbitrary data or additional information
  string createdBy = 6; // The user or entity who created the address list
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-address-list [id]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_address_list/1"
```

### Response Example

```json
{
    "list": {
        "listId": "1",
        "addresses": ["bb1...", "bb1..."],
        "whitelist": true,
        "uri": "",
        "customData": "",
        "createdBy": "bb1..."
    }
}
```
