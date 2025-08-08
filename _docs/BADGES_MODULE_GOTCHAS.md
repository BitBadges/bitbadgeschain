# BitBadges Module Development Gotchas and Common Issues

This document captures common pitfalls, gotchas, and important considerations when working with the BitBadges `x/badges` module.

## Table of Contents

- [Timeline System Gotchas](#timeline-system-gotchas)
- [Approval System Complexities](#approval-system-complexities)
- [State Management Issues](#state-management-issues)
- [Validation Gotchas](#validation-gotchas)
- [Permission System Traps](#permission-system-traps)
- [Gas Optimization Issues](#gas-optimization-issues)
- [Cross-Chain Considerations](#cross-chain-considerations)
- [Development Workflow Issues](#development-workflow-issues)

## Timeline System Gotchas

### 1. Timeline Immutability After Block Time

**Issue**: Once a timeline time has passed (current block time >= timeline time), that timeline entry becomes immutable and cannot be modified.

```go
// ❌ This will fail if current block time is > 1000000
managerTimeline := []ManagerTimeline{
    {
        TimelineTimes: []UintRange{{Start: NewUint(500000), End: NewUint(1000000)}},
        Manager: "bb1alice...",
    },
}
// Trying to update this timeline entry after block 1000000 will be rejected
```

**Solution**: Always check current block time before attempting timeline updates. Use future block times for scheduled changes.

```go
// ✅ Correct approach - use future block times
currentTime := ctx.BlockTime().UnixMilli()
futureTime := currentTime + (24 * 60 * 60 * 1000) // 24 hours from now

managerTimeline := []ManagerTimeline{
    {
        TimelineTimes: []UintRange{{Start: NewUint(uint64(futureTime)), End: MaxUint}},
        Manager: "bb1bob...",
    },
}
```

### 2. Overlapping Timeline Ranges

**Issue**: Timeline times within the same timeline array cannot overlap, but this is only validated at message processing time.

```go
// ❌ This will fail - overlapping timeline ranges
managerTimeline := []ManagerTimeline{
    {
        TimelineTimes: []UintRange{{Start: NewUint(1), End: NewUint(1000)}},
        Manager: "bb1alice...",
    },
    {
        TimelineTimes: []UintRange{{Start: NewUint(500), End: NewUint(1500)}}, // Overlaps!
        Manager: "bb1bob...",
    },
}
```

**Solution**: Ensure timeline ranges are sequential and non-overlapping.

```go
// ✅ Correct approach - sequential, non-overlapping ranges
managerTimeline := []ManagerTimeline{
    {
        TimelineTimes: []UintRange{{Start: NewUint(1), End: NewUint(1000)}},
        Manager: "bb1alice...",
    },
    {
        TimelineTimes: []UintRange{{Start: NewUint(1001), End: MaxUint}},
        Manager: "bb1bob...",
    },
}
```

### 3. Current Value Resolution

**Issue**: Timeline current values are resolved based on the latest timeline entry where `current_time >= timeline_start`, not the active range.

```go
// Timeline with gap
managerTimeline := []ManagerTimeline{
    {
        TimelineTimes: []UintRange{{Start: NewUint(1), End: NewUint(1000)}},
        Manager: "bb1alice...",
    },
    {
        TimelineTimes: []UintRange{{Start: NewUint(2000), End: MaxUint}},
        Manager: "bb1bob...",
    },
}
// At block 1500: Manager is still "bb1alice" (not "bb1bob")
// At block 2000: Manager becomes "bb1bob"
```

## Approval System Complexities

### 4. Three-Tier Approval Requirement

**Issue**: ALL three approval levels must approve a transfer: collection, outgoing, and incoming. Missing any one causes transfer failure.

```go
// ❌ Transfer will fail if ANY level rejects
// 1. Collection approval must allow the transfer
// 2. Sender's outgoing approval must allow the transfer  
// 3. Recipient's incoming approval must allow the transfer
```

**Solution**: Always verify all three approval levels are configured correctly before attempting transfers.

### 5. First-Match Permission Policy

**Issue**: The first matching permission rule determines the outcome. Order matters!

```go
// ❌ Dangerous permission ordering
permissions := []ActionPermission{
    {
        // This matches everything and denies all
        PermanentlyForbiddenTimes: []UintRange{{Start: NewUint(1), End: MaxUint}},
    },
    {
        // This will never be reached due to first-match policy
        PermanentlyPermittedTimes: []UintRange{{Start: NewUint(1), End: MaxUint}},
    },
}
```

**Solution**: Order permissions from most specific to most general.

### 6. Approval Version Invalidation

**Issue**: Updating user approvals increments the approval version, invalidating all existing approval usage tracking.

```go
// ❌ This invalidates all existing approval trackers for the user
err := k.UpdateUserApprovals(ctx, msg)
// All previous approval usage is reset to zero
```

**Solution**: Be aware that approval updates reset usage counters. Plan approval updates carefully.

## State Management Issues

### 7. Default Balance Inheritance

**Issue**: Users automatically inherit collection default balances, which can lead to unexpected behavior if not properly configured.

```go
// ❌ Dangerous default - gives everyone unlimited tokens
defaultBalances := UserBalanceStore{
    Balances: []Balance{
        {
            Amount: []UintRange{{Start: NewUint(1), End: MaxUint}},
            BadgeIds: []UintRange{{Start: NewUint(1), End: NewUint(1000)}},
            OwnershipTimes: []UintRange{{Start: NewUint(1), End: MaxUint}},
        },
    },
}
```

**Solution**: Set conservative defaults (usually zero amounts) and explicitly grant balances.

```go
// ✅ Safe default - no badges by default
defaultBalances := UserBalanceStore{
    Balances: []Balance{
        {
            Amount: []UintRange{{Start: NewUint(0), End: NewUint(0)}},
            BadgeIds: []UintRange{{Start: NewUint(1), End: NewUint(1000)}},
            OwnershipTimes: []UintRange{{Start: NewUint(1), End: MaxUint}},
        },
    },
}
```

### 8. Lazy State Initialization

**Issue**: User balance stores are only created when explicitly modified, leading to "not found" vs "default values" confusion.

```go
// When querying balances:
userBalance, found := k.GetUserBalanceStoreFromStore(ctx, collectionId, address)
if !found {
    // User inherits collection defaults - this is normal!
    userBalance = collection.DefaultBalances
}
```

**Solution**: Always check for default inheritance when user balance store is not found.

## Validation Gotchas

### 9. UintRange Validation Rules

**Issue**: Multiple validation rules apply to UintRange arrays that are not immediately obvious.

```go
// ❌ Common UintRange mistakes
ranges := []UintRange{
    {Start: NewUint(10), End: NewUint(1)},   // Start > End - INVALID
    {Start: NewUint(1), End: NewUint(5)},    
    {Start: NewUint(3), End: NewUint(8)},    // Overlaps with previous - INVALID
    {Start: NewUint(0), End: NewUint(0)},    // Zero amounts in balances - INVALID
}
```

**Solution**: Ensure ranges are:
- Start ≤ End
- Non-overlapping within the same array
- Non-zero for balance amounts

### 10. Address Validation Edge Cases

**Issue**: Special address handling and validation edge cases.

```go
// ❌ Common address mistakes
addresses := []string{
    "Mint",      // Reserved address - valid in some contexts
    "Total",     // Reserved address - valid in some contexts  
    "*",         // Wildcard - valid in address lists
    "",          // Empty string - INVALID
    "invalid",   // Invalid format - INVALID
}
```

**Solution**: Use proper address validation and understand reserved address contexts.

## Permission System Traps

### 11. Permanent Permission Decisions

**Issue**: Once a time range is marked as permanently permitted or forbidden, it cannot be changed.

```go
// ❌ This permission cannot be reverted later
permission := ActionPermission{
    PermanentlyForbiddenTimes: []UintRange{{Start: NewUint(1), End: MaxUint}},
}
// This permanently forbids the action for all time - cannot be undone!
```

**Solution**: Use time-limited permissions for reversible decisions.

```go
// ✅ Time-limited permission that can be changed later
permission := ActionPermission{
    PermanentlyForbiddenTimes: []UintRange{{Start: NewUint(1), End: NewUint(1000000)}},
    // After block 1000000, this permission no longer applies
}
```

### 12. Permission Inheritance and Defaults

**Issue**: When no explicit permission matches, the system uses default behavior which may not be what you expect.

```go
// If no permission rule matches:
// - Some operations default to ALLOWED
// - Some operations default to FORBIDDEN
// - Behavior varies by operation type
```

**Solution**: Always provide explicit permission rules rather than relying on defaults.

## Gas Optimization Issues

### 13. Large Range Iterations

**Issue**: Operations on large ranges can consume excessive gas.

```go
// ❌ This could be expensive for large ranges
badgeIds := []UintRange{{Start: NewUint(1), End: NewUint(1000000)}}
// Processing 1M token IDs individually would be prohibitively expensive
```

**Solution**: The module uses range-based algorithms, but be aware of practical limits.

### 14. Excessive Timeline Entries

**Issue**: Too many timeline entries can make queries and updates expensive.

```go
// ❌ Excessive timeline granularity
var managerTimeline []ManagerTimeline
for i := 0; i < 1000; i++ {
    managerTimeline = append(managerTimeline, ManagerTimeline{
        TimelineTimes: []UintRange{{Start: NewUint(uint64(i)), End: NewUint(uint64(i))}},
        Manager: fmt.Sprintf("bb1manager%d...", i),
    })
}
```

**Solution**: Use reasonable timeline granularity based on actual business needs.

## Cross-Chain Considerations

### 15. EIP712 Schema Completeness

**Issue**: EIP712 schemas must include ALL possible fields, even optional ones, or Ethereum signatures will fail.

```go
// ❌ Missing fields in EIP712 schema
schema := `{
    "type": "badges/CreateCollection",
    "value": {
        "creator": "",
        "balancesType": ""
        // Missing other fields - signatures will fail!
    }
}`
```

**Solution**: Always include all message fields in EIP712 schemas with appropriate default values.

### 16. Multi-Chain Address Format Differences

**Issue**: Different blockchains have different address formats, but the module stores addresses as strings.

```go
// Valid addresses from different chains:
addresses := []string{
    "bb1abc123...",           // BitBadges
    "cosmos1xyz789...",       // Cosmos
    "0x1234567890abcdef...",  // Ethereum
    "bc1qw508d6qejxtdg4y...", // Bitcoin
}
```

**Solution**: Implement proper address validation for each supported chain type.

## Development Workflow Issues

### 17. Proto Generation Order Dependencies

**Issue**: The order of proto generation and API cleanup matters.

```bash
# ❌ Wrong order - can cause build failures
ignite generate proto-go --yes
go build ./cmd/bitbadgeschaind  # May fail with versioned API imports
rm -rf api/badges/v*

# ✅ Correct order
ignite generate proto-go --yes
rm -rf api/badges/v*
go build ./cmd/bitbadgeschaind
```

**Solution**: Always clean up API versions immediately after proto generation.

### 18. Skip Proto Flag Requirement

**Issue**: Must use `--skip-proto` flag with Ignite commands due to manual proto corrections.

```bash
# ❌ This will overwrite manual corrections
ignite chain serve

# ✅ Always use --skip-proto
ignite chain serve --skip-proto
```

**Solution**: Always use `--skip-proto` flag with Ignite CLI commands.

### 19. Test Data Preparation Complexity

**Issue**: Creating valid test data requires understanding all validation rules and dependencies.

```go
// ❌ Incomplete test collection - will fail validation
collection := &types.BadgeCollection{
    CollectionId: NewUint(1),
    // Missing required fields like manager timeline, permissions, etc.
}
```

**Solution**: Use helper functions or copy from working examples for test data preparation.

### 20. Message Handler Testing Order

**Issue**: Message handlers have dependencies on keeper state that must be set up correctly.

```go
// ❌ Testing transfer without proper setup
func TestTransfer(t *testing.T) {
    k, ctx := setupKeeper(t)
    
    // Missing: Create collection, set up approvals, mint initial balances
    msg := &types.MsgTransferBadges{...}
    _, err := k.TransferBadges(ctx, msg) // Will fail - collection doesn't exist
}
```

**Solution**: Set up all required state (collections, approvals, balances) before testing operations.

## Common Error Messages and Solutions

### "Collection not found"
- Ensure collection ID exists
- Check if using correct collection ID format (string representation of uint)

### "Timeline validation failed"
- Check for overlapping timeline ranges
- Verify timeline times are in correct order
- Ensure no gaps in required timelines

### "Approval not satisfied"
- Verify all three approval levels (collection, outgoing, incoming)
- Check approval criteria (Merkle proofs, must-own tokens, etc.)
- Ensure approval versions are current

### "Permission denied"
- Check first-match permission policy
- Verify current time is in permitted ranges
- Ensure manager has required permissions

### "Invalid range"
- Verify start ≤ end for all ranges
- Check for overlapping ranges in arrays
- Ensure non-zero amounts in balance ranges

This document will be continuously updated as new gotchas and common issues are discovered during development and deployment.