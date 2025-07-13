# Action Permission

ActionPermissions are the simplest (no criteria). Just denotes what times the action is executable or not.

<pre class="language-json"><code class="lang-json"><strong>"collectionPermissions": {
</strong>    "canDeleteCollection": [...],
    ...
}
</code></pre>

```json
"userPermissions": {
    ...
    "canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [...],
    "canUpdateAutoApproveSelfInitiatedIncomingTransfers": [...],
    "canUpdateAutoApproveAllIncomingTransfers": [...],
}
```

```typescript
export interface ActionPermission<T extends NumberType> {
  permanentlyPermittedTimes: UintRange<T>[];
  permanentlyForbiddenTimes: UintRange<T>[];
}
```

**Examples**

Below, this forbids the action from ever being executed.

```json
"canDeleteCollection": [
  {
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
  }
]
```
