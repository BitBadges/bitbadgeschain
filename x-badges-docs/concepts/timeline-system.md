# Timeline System

BitBadges uses timeline-based fields to allow dynamic, time-dependent values for various attributes. This feature enables automatic updates to field values based on the current time, without requiring additional blockchain transactions.

## Key Concepts

### Time Representation

-   Times are represented as UNIX time (milliseconds since the epoch)
-   Epoch: Midnight at the beginning of January 1, 1970, UTC
-   Time fields use UintRange format with valid values from 1 to 18446744073709551615 (Go MaxUint64)

### Value Assignment

-   Values are assigned to specific time ranges
-   No overlapping time ranges are allowed for a single field
-   Current value determined by finding the timeline entry that includes current time

### Default Behavior

-   If no value is set for the current time, the field assumes an empty/null/default value

## Structure

Timeline-based fields extend the `TimelineItem` interface:

```typescript
export interface TimelineItem<T extends NumberType> {
    timelineTimes: UintRange<T>[];
}
```

## Proto Definition Examples

```protobuf
message ManagerTimeline {
  repeated UintRange timelineTimes = 1;
  string manager = 2;
}

message CollectionMetadataTimeline {
  repeated UintRange timelineTimes = 1;
  CollectionMetadata collectionMetadata = 2;
}

message BadgeMetadataTimeline {
  repeated UintRange timelineTimes = 1;
  repeated UintRange badgeIds = 2;
  BadgeMetadata badgeMetadata = 3;
}
```

## Timeline Fields in Collections

The collection interface includes the following timeline-based fields:

-   `managerTimeline: ManagerTimeline<T>[]`
-   `collectionMetadataTimeline: CollectionMetadataTimeline<T>[]`
-   `badgeMetadataTimeline: BadgeMetadataTimeline<T>[]`
-   `offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline<T>[]`
-   `customDataTimeline: CustomDataTimeline<T>[]`
-   `standardsTimeline: StandardsTimeline<T>[]`
-   `isArchivedTimeline: IsArchivedTimeline<T>[]`

## Usage Example

### Collection Metadata Timeline

```json
{
    "collectionMetadataTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "1680307199000" }],
            "collectionMetadata": {
                "uri": "ipfs://abc123",
                "customData": ""
            }
        },
        {
            "timelineTimes": [
                { "start": "1680307200000", "end": "18446744073709551615" }
            ],
            "collectionMetadata": {
                "uri": "ipfs://xyz456",
                "customData": ""
            }
        }
    ]
}
```

In this example:

-   From time 1 to March 31, 2023, the collection metadata URI is 'ipfs://abc123'
-   From April 1, 2023 onwards, the collection metadata URI is 'ipfs://xyz456'
-   The change happens automatically without additional transactions

### Manager Timeline

```json
{
    "managerTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "1672531199000" }],
            "manager": "bb1alice..."
        },
        {
            "timelineTimes": [
                { "start": "1672531200000", "end": "18446744073709551615" }
            ],
            "manager": "bb1bob..."
        }
    ]
}
```

This setup transfers management from Alice to Bob on January 1, 2023.

## Key Principles

1. **Automatic Updates**: Values change automatically based on current time
2. **No Overlaps**: Timeline times within the same array should not overlap
3. **Future Scheduling**: Can schedule changes for future times
4. **Gas Efficiency**: No additional transactions needed for scheduled changes

## Practical Applications

Timeline-based fields enable:

-   **Scheduled ownership transfers** without manual intervention
-   **Automatic metadata updates** for evolving collections
-   **Time-based permission changes** for governance
-   **Seasonal content updates** for dynamic collections
-   **Archive scheduling** for temporary collections
