# Address Lists

Address lists define collections of addresses for use in approval configurations. They support both static lists stored on-chain and dynamic reserved patterns for common access control scenarios.

## Usage in Approval Configurations

Address lists are referenced by ID in three approval contexts. IDs can either be reserved, shorthand IDs or user-created lists via `MsgCreateAddressLists`.

### Collection Approvals

```protobuf
message CollectionApproval {
  string fromListId = 1;        // Who can send badges
  string toListId = 2;          // Who can receive badges
  string initiatedByListId = 3; // Who can initiate transfers
  // ... other fields
}
```

### User Outgoing Approvals

```protobuf
message UserOutgoingApproval {
  string toListId = 1;          // Who user can send to
  string initiatedByListId = 2; // Who can initiate on user's behalf
  // ... other fields
}
```

### User Incoming Approvals

```protobuf
message UserIncomingApproval {
  string fromListId = 1;        // Who can send to user
  string initiatedByListId = 2; // Who can initiate transfers to user
  // ... other fields
}
```

### Usage Examples

#### Universal Access

```json
{
    "fromListId": "AllWithoutMint", // Everyone except Mint
    "toListId": "All", // Everyone including Mint
    "initiatedByListId": "All" // Anyone can initiate
}
```

#### Restricted Access

```json
{
    "fromListId": "vipMembers", // Only VIP members can send
    "toListId": "!banned", // Everyone except banned users
    "initiatedByListId": "AllWithoutMint"
}
```

#### Quick Address Lists

```json
{
    "fromListId": "bb1alice...:bb1bob...:bb1charlie...", // Direct addresses
    "toListId": "AllWithoutMint:bb1blocked...", // Everyone except these
    "initiatedByListId": "All"
}
```

## Proto Definition

```protobuf
message AddressList {
  string listId = 1;           // Unique identifier
  repeated string addresses = 2; // List of addresses
  bool whitelist = 3;          // true = whitelist, false = blacklist
  string uri = 4;              // Metadata URI
  string customData = 5;       // Custom data
  string createdBy = 6;        // Creator address
}
```

## Reserved Address List IDs

BitBadges provides built-in reserved list IDs that are dynamically generated without storage overhead:

### Core Reserved Lists

#### "Mint"

-   **Purpose**: Contains only the "Mint" address
-   **Logic**: Whitelist (addresses: ["Mint"], whitelist: true)
-   **Use case**: Minting operations and initial badge distribution

#### "All" and "AllWithMint"

-   **Purpose**: Represents all addresses including Mint
-   **Logic**: Blacklist with empty addresses list (addresses: [], whitelist: false)
-   **Use case**: Universal access, public collections

#### "None"

-   **Purpose**: Represents no addresses
-   **Logic**: Whitelist with empty addresses list (addresses: [], whitelist: true)
-   **Use case**: Blocking all access, disabled transfers

### Dynamic Patterns

#### AllWithout Pattern

-   **Format**: `"AllWithout<addresses>"` where addresses are colon-separated
-   **Example**: `"AllWithoutMint"`, `"AllWithoutMint:bb1user123"`
-   **Logic**: Blacklist containing the specified addresses (addresses: ["Mint", "bb1user123"], whitelist: false)
-   **Use case**: Allow everyone except specific addresses

#### Colon-Separated Addresses

-   **Format**: `"address1:address2:address3"`
-   **Logic**: Whitelist containing the specified addresses (addresses: ["bb1user123", "bb1user234", "bb1user345"], whitelist: true)
-   **Use case**: Quick address lists without creating stored lists

#### Inversion Patterns

-   **Format**: `"!listId"` or `"!(listId)"`
-   **Effect**: Inverts the whitelist/blacklist behavior of the referenced list
-   **Example**: `"!5"` inverts list ID 5's behavior

## Whitelist vs Blacklist Logic

Address lists use a boolean `whitelist` field to determine inclusion/exclusion behavior:

### Whitelist Logic (`whitelist: true`)

```javascript
function isAddressIncluded(address, addressList) {
    const found = addressList.addresses.includes(address);
    return addressList.whitelist ? found : !found;
}
```

-   **Listed addresses**: Explicitly included
-   **Unlisted addresses**: Explicitly excluded
-   **Use case**: "Only these addresses are allowed"

### Blacklist Logic (`whitelist: false`)

-   **Listed addresses**: Explicitly excluded
-   **Unlisted addresses**: Explicitly included
-   **Use case**: "All addresses except these are allowed"

### Inversion Effect

When a list ID has the `"!"` prefix, the final whitelist boolean is inverted:

-   Whitelist becomes blacklist behavior
-   Blacklist becomes whitelist behavior

## Usage in Approval Configurations

Address lists are referenced by ID in three approval contexts:

### Collection Approvals

```protobuf
message CollectionApproval {
  string fromListId = 1;        // Who can send badges
  string toListId = 2;          // Who can receive badges
  string initiatedByListId = 3; // Who can initiate transfers
  // ... other fields
}
```

### User Outgoing Approvals

```protobuf
message UserOutgoingApproval {
  string toListId = 1;          // Who user can send to
  string initiatedByListId = 2; // Who can initiate on user's behalf
  // ... other fields
}
```

### User Incoming Approvals

```protobuf
message UserIncomingApproval {
  string fromListId = 1;        // Who can send to user
  string initiatedByListId = 2; // Who can initiate transfers to user
  // ... other fields
}
```

## Address List Creation

### User-Created Lists

Created through `MsgCreateAddressLists` with the requirements below. Once created, this list is immutable and cannot be modified.

#### ID Validation Rules

-   Must be alphanumeric characters only
-   Cannot be empty or reserved keywords
-   Cannot contain `:` or `!` characters
-   Cannot be valid addresses themselves
-   Cannot conflict with reserved IDs

#### Address Validation

-   All addresses must be valid Bech32 format
-   No duplicate addresses allowed
-   Special addresses like "Mint" are permitted

### Example User-Created List

```json
{
    "listId": "vipMembers",
    "addresses": ["bb1alice...", "bb1bob...", "bb1charlie..."],
    "whitelist": true,
    "uri": "https://api.example.com/vip-list",
    "customData": "VIP members with exclusive access",
    "createdBy": "bb1manager..."
}
```

## Performance Characteristics

Tip: Use reserved IDs for common patterns. Use user-created lists for large lists used multiple times.

### Reserved Lists

-   **Storage**: Zero on-chain storage
-   **ID Length**: Dependent on the pattern used
-   **Resolution**: Dynamic generation at runtime
-   **Gas efficiency**: Minimal overhead for common patterns

### User-Created Lists

-   **Storage**: On-chain storage per list
-   **ID Length**: Reusable short ID that can be used to reference complex lists
-   **Resolution**: Direct store lookup
-   **Gas cost**: Proportional to validation complexity

## Address Lists vs Dynamic Stores

Both address lists and dynamic stores can control who is approved for transfers, but they serve different purposes:

### Address Lists

-   **Purpose**: Immutable shorthand references for collections of addresses
-   **Mutability**: Cannot be modified after creation - addresses are fixed
-   **Storage**: Direct list of addresses stored on-chain
-   **Use case**: Static whitelists/blacklists that don't change over time

### Dynamic Stores

-   **Purpose**: Mutable on-chain approval management with CRUD operations
-   **Mutability**: Can be updated dynamically using CRUD messages (Create, Update, Delete, Set)
-   **Storage**: Boolean values per address (true/false approval status)
-   **Use case**: Dynamic approval systems that need real-time updates (typically contract logic)
