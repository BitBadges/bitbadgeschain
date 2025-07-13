# MsgCreateAddressLists

Creates reusable address lists by ID for gas optimizations.

## Important Notes

1. **Create Only**: There are no update, edit, or delete functions for address lists. Once created, they are immutable.

2. **Optional Efficiency Tool**: This is completely optional and serves as a reusable shorthand ID to avoid repetition of long reserved address list IDs. The primary purpose is gas efficiency.

3. **Minimal Metadata**: Typically, `uri` and `customData` are left blank as these fields are not supported on the BitBadges site and are different from off-chain lists you may see elsewhere.

## Proto Definition

```protobuf
message MsgCreateAddressLists {
  string creator = 1; // Address creating the address lists
  repeated AddressList addressLists = 2; // Lists to create in single transaction
}

message MsgCreateAddressListsResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges create-address-lists '[tx-json]' --from creator-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "addressLists": [
        {
            "listId": "",
            "addresses": ["bb1...", "bb1..."],
            "whitelist": true,
            "createdBy": "", // Leave blank - auto-generated
            "aliasAddress": "", // Leave blank - auto-generated
            "uri": "",
            "customData": ""
        }
    ]
}
```
