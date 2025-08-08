# MsgDeleteCollection

Deletes a collection.

## Authorization

Collection deletion can only be performed by the **current manager** of the collection and requires the `canDeleteCollection` permission to be enabled at the current time in the collection's permissions.

## Proto Definition

```protobuf
message MsgDeleteCollection {
  string creator = 1; // Address requesting deletion (must be manager)
  string collectionId = 2; // ID of collection to delete
}

message MsgDeleteCollectionResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges delete-collection [collection-id] --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "collectionId": "1"
}
```
