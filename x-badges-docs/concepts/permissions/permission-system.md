# Overview

First, read [Permissions](broken-reference) for an overview.

Note: The [Approved Transfers](../balances-transfers/transferability-approvals.md) and [Permissions ](broken-reference)are the most powerful features of the interface, but they can also be the most confusing. Please ask for help if needed.

```json
"collectionPermissions": {
    ...
}
```

```json
"userPermissions": {
    ...
}
```

### Collection Permissions

#### Manager

The collectionPermissions only apply to the current manager of the collection. In other words, the manager is the only one who is able to execute permissions. If there is no manager for a collection, no permissions can be executed.

The current manager is determined by the **managerTimeline.** Transferring the manager is facilitated via the **canUpdateManager** permission.

```json
"managerTimeline": [
  {
    "manager": "bb1kfr2xajdvs46h0ttqadu50nhu8x4v0tc2wajnh",
    "timelineTimes": [
      {
        "start": 1,
        "end": "18446744073709551615"
      }
    ]
  }
]
```

### **User Permissions**

Besides the collection permissions, there are also userPermissions that can be set. Typically, these will remain empty / unset, so that the user can always have full control over their approvals. If empty, they are permitted by default (but not frozen).

However, setting user permissions can be leveraged in some cases for specific purposes.

* Locking that a specific badge can never be transferred out of the account
* Locking that a specific approval is always set and uneditable so that two mutually distrusting parties can use the address as an escrow

**Defaults**

We also give the option for the collection to define default user permissions. These will be used as the starting values when the balance is initially created in storage. This can be used in tandem with the other defaults. The default permissions are also not typically used, but again can be used in certain situations. For example, by default, approve all incoming transfers and lock the permission so all transfers always have incoming approvals and can never be disapproved.

### **Permitted and Forbidden Times**

Permissions allow you to define permitted or forbidden times to be able to execute a permission.

<figure><img src="../../../.gitbook/assets/image (2) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1).png" alt=""><figcaption></figcaption></figure>

**States**

There are three states that a permission can be in at any one time:

1. **Forbidden + Permanently Frozen (permanentlyForbiddenTimes):** This permission is forbidden and will always remain forbidden.
   1. If a permission is explicitly allowed via the **permanentlyPermittedTimes, it will ALWAYS be allowed** during those permanentlyPermittedTimes (can't change it).
2. **Permitted + Not Frozen (Unhandled):** This permission is currently permitted but can be changed to one of the other two states.
   1. If not explicitly permitted or forbidden - NEUTRAL (not defined or unhandled), **permissions are ALLOWED by default** but can later be set to be permanently allowed or disallowed. There is no "forbidden currently but updatable" state.
3. **Permitted + Permanently Frozen (permanentlyPermittedTimes):** This permission is forbidden and will always remain permitted
   1. If a permission is explicitly forbidden via the **permanentlyForbiddenTimes, it will ALWAYS be disallowed** during those permanentlyForbiddenTimes.

There is no forbidden + not frozen state because theoretically, it could be updated to permitted at any time and executed (thus making it permitted).

**Examples**

This means the permission is permanently forbidden and frozen.

```typescriptreact
permanentlyPermittedTimes: []
permanentlyForbiddenTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
```

This means it is permanently allowed and frozen.

```typescriptreact
permanentlyForbiddenTimes: []
permanentlyPermittedTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
```

This means it is allowed currently but neutral and can be changed to be always permitted or always forbidden in the future.

```
permanentlyForbiddenTimes: []
permanentlyPermittedTimes: []
```

### First Match Policy

All permissions are a linear array where each element may have some criteria as well as **permanentlyForbiddenTimes** or **permanentlyPermittedTimes.** It can be interpreted as if the criteria matches, the permission is permitted or forbidden according to the defined times, respectively.

We do not allow times to be in both the permanentlyPermittedTimes and permanentlyForbiddenTimes array simultaneously.

**Unlike approvals, we only allow taking the first match in the case criteria satisfies multiple elements in the permissions array.** All subsequent matches are ignored. This makes it so that for any time and for each criteria combination, there is a deterministic permission state (permitted, forbidden, or neutral). This means you have to carefully design your permissions because order and overlaps matter.

Ex: If we have the following permission definitions in an array \[elem1, elem2]:

1. ```
   timelineTimes: [{ start: 1, end: 10 }]

   permanentlyPermittedTimes: []
   permanentlyForbiddenTimes: [{ start: 1, end: 10 }]
   ```
2. ```
   timelineTimes: [{ start: 1, end: 100 }]

   permanentlyForbiddenTimes: []
   permanentlyPermittedTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
   ```

In this case, the timeline times 1-10 will be forbidden ONLY from times 1-10 because we take the first element that matches for that specific criteria (which is permanentlyPermittedTimes: \[], permanentlyForbiddenTimes: \[1 to 10]).

Times 11-100 would be permanently permitted since the first match for those times is the second element.

Similar to approved transfers, even though we allow range logic to be specified, we first expand everything maintaining order to their singular values (one value, no ranges) before checking for our first match.

### Satisfying Criteria

For permissions, all criteria must be satisfied for it to be a match. If you satisfy N-1 criteria, it is not a match.

For example, lets say you had a permission with badge IDs and ownership times:

```
badgeIds: [{ start: 1, end: 10 }]
ownershipTimes: [{ start 1, end: 10 }]

permanentlyForbiddenTimes: []
permanentlyPermittedTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
```

This would result in the manager being able to create more of badges IDs 1-10 which can be owned from times 1-10.

However, this permission **does not** specify whether they can create more of badge ID 1 at time 11 or badge ID 11 at time 1. These combinations are considered unhandled or not defined by the permission definition above.

**Common Misunderstanding**

A common misunderstanding is that if the permission below is appended after the above one, this would forbid badges 11+ from ever being created. However, creating badge IDs 11+ at times 11+ would still be unhandled and **allowed by default**.

```
badgeIds: [{ start: 11, end: Max }]
ownershipTimes: [{ start 1, end: 10 }]

permanentlyForbiddenTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
permanentlyPermittedTimes: []
```

To permanently forbid all badgeIds, you must brute force ALL other combinations such as

```
badgeIds: [{ start: 11, end: Max }]
ownershipTimes: [{ start: 1, end: Max }] // 1-10 never gets matched to bc of first match
//can also do start: 11

permanentlyForbiddenTimes: [{ start: 1, end: GO_MAX_UINT_64 }]
permanentlyPermittedTimes: []
```

**Brute-Forcing**

A common pattern you will see is to brute force all possible combinations. For example, in the above example we brute forced all possible combinations for badge IDs 11+. No subsequent element specifying a badge ID 11+ will ever get matched to.

To brute force a specific criteria (such as IDs 11+), you specify it, then for all other N - 1 criteria, you set them equal to ALL values. All values in the case of UintRanges is 1 - max Uint64. For address lists / IDs, this is all possible addresses / IDs.

The following brute forces badge IDs 1-10.

<pre class="language-json"><code class="lang-json"><strong>{
</strong>  "fromListId": "All",
  "toListId": "All",
  "initiatedByListId": "All",
  "badgeIds":  [{ "start": "1", "end": "10" }],
  "transferTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "ownershipTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "approvalId": "All", //forbids approval "xyz" from being updated

  "permanentlyPermittedTimes": [],
  "permanentlyForbiddenTimes": [{ "start": "1", "end": "18446744073709551615" }]
}
</code></pre>

### **Permission Categories**

There are five categories of permissions, each with different criteria that must be matched with. If you get confused with the different time types, refer to [Different Time Types](../different-time-fields.md) for examples and explanations.

{% content-ref url="action-permission.md" %}
[action-permission.md](action-permission.md)
{% endcontent-ref %}

{% content-ref url="timed-update-permission.md" %}
[timed-update-permission.md](timed-update-permission.md)
{% endcontent-ref %}

{% content-ref url="timed-update-with-badge-ids-permission.md" %}
[timed-update-with-badge-ids-permission.md](timed-update-with-badge-ids-permission.md)
{% endcontent-ref %}

{% content-ref url="balances-action-permission.md" %}
[balances-action-permission.md](balances-action-permission.md)
{% endcontent-ref %}

{% content-ref url="update-approval-permission.md" %}
[update-approval-permission.md](update-approval-permission.md)
{% endcontent-ref %}

### **Examples**

See [Example Msgs](../../core-concepts/broken-reference/) for further examples. Or, see the page for each permission category.

```json
"collectionPermissions": {
    "canArchiveCollection": [],
    "canCreateMoreBadges": [
      {
        "badgeIds": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "ownershipTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "permanentlyPermittedTimes": [],
        "permanentlyForbiddenTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
      }
    ],
    "canDeleteCollection": [],
    "canUpdateBadgeMetadata": [],
    "canUpdateCollectionApprovals": [
      {
        "fromListId": "AllWithMint",
        "toListId": "AllWithMint",
        "initiatedByListId": "AllWithMint",
        "timelineTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "transferTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "badgeIds": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "ownershipTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ],
        "approvalId": "All",
        "permanentlyPermittedTimes": [],
        "permanentlyForbiddenTimes": [
          {
            "start": "1",
            "end": "18446744073709551615"
          }
        ]
      }
    ],
    "canUpdateCollectionMetadata": [],
    "canUpdateContractAddress": [],
    "canUpdateCustomData": [],
    "canUpdateManager": [],
    "canUpdateOffChainBalancesMetadata": [],
    "canUpdateStandards": []
  }
```
