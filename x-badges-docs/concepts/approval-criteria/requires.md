# Requires

Additional address relationship restrictions for transfer approval.

## Interface

```typescript
interface ApprovalCriteria<T extends NumberType> {
    requireToEqualsInitiatedBy?: boolean;
    requireToDoesNotEqualInitiatedBy?: boolean;
    requireFromEqualsInitiatedBy?: boolean;
    requireFromDoesNotEqualInitiatedBy?: boolean;
}
```

## How It Works

Enforce additional checks on address relationships:

-   **`requireToEqualsInitiatedBy`**: Recipient must equal initiator
-   **`requireToDoesNotEqualInitiatedBy`**: Recipient must not equal initiator
-   **`requireFromEqualsInitiatedBy`**: Sender must equal initiator
-   **`requireFromDoesNotEqualInitiatedBy`**: Sender must not equal initiator

## Constraints

All checks are bounded by the respective address lists (`toList`, `fromList`, `initiatedByList`).
