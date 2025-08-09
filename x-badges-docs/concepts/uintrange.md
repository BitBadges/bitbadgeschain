# UintRanges

The `UintRange` is the fundamental data structure used throughout the badges module to represent inclusive ranges of unsigned integers efficiently. This type enables powerful range-based operations and is primarily used for token IDs, time ranges, and amounts.

## Proto Definition

```protobuf
message UintRange {
  string start = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
  string end = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Patterns

UintRanges are used to represent:

-   **Token ID ranges**: `[1-100]` represents token IDs 1 through 100 (inclusive)
-   **Time ranges**: `[1640995200000-1672531200000]` represents a year in UNIX milliseconds
-   **Amount ranges**: `[1-5]` represents quantities from 1 to 5
-   **Ownership time ranges**: When tokens are valid for ownership

## Restrictions & Valid Values

Unless otherwise specified, we only allow numbers in the ranges to be from **1 to Go Max UInt64**:

-   **Valid range**: 1 to 18446744073709551615 (Go's `math.MaxUint64`)
-   **Zero and negative values**: Not allowed
-   **Values greater than maximum**: Not allowed

## Validation Rules

-   `start` must be ≤ `end`
-   Ranges in the same array cannot overlap
-   Zero amounts are not allowed in balance ranges
-   All values must be within the valid range (1 to MaxUint64)

## Special Cases

### Full Range

To represent a complete range covering all possible values:

```protobuf
// Full range from 1 to maximum
{
  start: "1",
  end: "18446744073709551615"
}
```

### Single Value

To represent a single value, use the same value for start and end:

```protobuf
// Single token ID 5
{
  start: "5",
  end: "5"
}
```

### Range Inversion

Inverting a range results in all values from 1 to 18446744073709551615 that are **not** in the current range. This is useful for exclusion logic.

## Examples

### Token ID Examples

```typescript
// Token IDs 1-10 (inclusive)
const badgeRange: UintRange[] = [{ start: '1', end: '10' }];

// Multiple non-overlapping ranges
const multipleBadges: UintRange[] = [
    { start: '1', end: '10' },
    { start: '20', end: '50' },
];
```

### Go Code Examples

```go
// Token IDs 1-10
badgeIdRange := UintRange{Start: NewUint(1), End: NewUint(10)}

// Unlimited amount
unlimitedAmount := UintRange{Start: NewUint(1), End: MaxUint}

// Single token ID
singleBadge := UintRange{Start: NewUint(5), End: NewUint(5)}
```

## Efficiency Benefits

-   **Compact representation**: Ranges avoid storing individual values
-   **Range operations**: Efficient intersection, union, and containment checks
-   **Gas optimization**: Reduces transaction size and computational costs
-   **Scalability**: Handles large ranges without performance degradation
