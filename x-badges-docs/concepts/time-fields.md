# Different Time Fields

BitBadges uses various time-related fields to manage permissions, timelines, transfers, and ownership. Understanding these fields is crucial for effectively managing collections and tokens.

## Time Representation

All times in BitBadges are represented as UNIX time, which is the number of milliseconds elapsed since the epoch (midnight at the beginning of January 1, 1970, UTC).

Time fields use UintRange format with valid values from 1 to 18446744073709551615 (Go MaxUint64). For complete details on range formatting and restrictions, see the [UintRange concept](uintrange.md).

## Time Field Types

### 1. permanentlyPermittedTimes

-   **Purpose**: Defines the times when a permission will always be executable (permanent)
-   **Usage**: Setting allowed periods for specific actions

### 2. permanentlyForbiddenTimes

-   **Purpose**: Defines the times when a permission will always be forbidden (permanent)
-   **Usage**: Setting restricted periods for specific actions

### 3. timelineTimes

-   **Purpose**: Specifies when a field is scheduled to have a specific value in a timeline-based field
-   **Usage**: Scheduling changes to collection or token properties over time

### 4. transferTimes

-   **Purpose**: Defines when a transfer transaction can occur (i.e. when an approval is valid)
-   **Usage**: Setting periods when tokens can be transferred between addresses

### 5. ownershipTimes

-   **Purpose**: Specifies the times that a user owns a token
-   **Usage**: Defining the duration of token ownership for users

## Important Note

The `timelineTimes` in permissions correspond to the updatability of the timeline, while `timelineTimes` in the actual timeline represent the actual times for the values.

## Examples

### Example 1: Presidential Election Badges

Scenario: Users participate in a US presidential election by casting votes through token transfers.

-   T1: Conclusion of voting
-   T2: Start of presidential term
-   T3: End of presidential term

Setup:

-   `transferTimes`: [{ start: T1, end: T2 }] (President badge can be transferred after voting concludes)
-   `ownershipTimes`: [{ start: T2, end: T3 }] (Defines the presidential term)

### Example 2: Managing Collection Archival

Scenario: A collection can be optionally archived by the manager from T1 to T2, but is non-archivable at all other times.

Before archiving:

```
Permission:
permanentlyPermittedTimes: [{ start: T1, end: T2 }]
permanentlyForbiddenTimes: [everything but T1 to T2]
timelineTimes: [{ start: 1, end: MAX_TIME }]

Archived Timeline:
isArchived: false for [{ start: 1, end: MAX_TIME }]
```

After archiving for all times:

```
Permission: (unchanged)
permanentlyPermittedTimes: [{ start: T1, end: T2 }]
permanentlyForbiddenTimes: [everything but T1 to T2]
timelineTimes: [{ start: 1, end: MAX_TIME }]

Archived Timeline:
isArchived: true for [{ start: 1, end: MAX_TIME }]
```

## Best Practices

1. **Clear Timelines**: Always define clear and non-overlapping time ranges for each field to avoid confusion and conflicts
2. **Permission Management**: Carefully consider the implications of setting `permanentlyPermittedTimes` and `permanentlyForbiddenTimes`, as these can significantly impact the flexibility of your collection
3. **Timeline Planning**: When using `timelineTimes`, plan your collection's lifecycle in advance to minimize the need for frequent updates
4. **Transfer Windows**: Use `transferTimes` to create specific windows for token transfers, which can be useful for time-limited events or phased distributions
5. **Ownership Tracking**: Leverage `ownershipTimes` to create tokens with time-bound ownership, useful for temporary privileges or rotating responsibilities
6. **Permission Locking**: Be cautious when permanently locking permissions, as this action is irreversible and may limit future flexibility
7. **Time Synchronization**: Ensure all systems interacting with your BitBadges collection are properly time-synchronized to avoid discrepancies in time-based operations
