Moving forward, the plan for validator operations is going to be a mix of existing proof-of-stake rewards and a new proof-of-authority model.

**1. Existing Validator Rewards Algorithm**

This program applies until 8/12/2026 — one year from the initial start date of 8/12/2025.

Uptime is measured by sampling block signatures every 2,500 blocks from block 6,882,000 (~Nov 13, 2025) to the current block height. At each sample, we check whether each validator's signature appears in the block. If a validator is not found in the block signatures — whether due to being jailed, offline, or not in the active set — it counts as a miss. Uptime % is calculated as (signed blocks / total sampled blocks) \* 100.

**Eligible Validators:**

- **Mission Decentralization candidates**: Reserved for the first ~40 validators who participated in Mission Decentralization, an incentivized program.
- **Everyone else**: Any validator that has registered by block height 8,998,000. New validators after this height are not eligible for this program.

**Reward Tiers:**

- **Mission Decentralization candidates**: Base allocation of 200,000 BADGE, scaled by uptime %, with a minimum floor of 50,000 BADGE (honoring 100% uptime for period 1 of 8/12 - 11/12)
- **All other eligible validators**: Base allocation of 100,000 BADGE, scaled by uptime %, no minimum
- **Excluded validators** (BitBadges, ChainTools, WHEN MOON WHEN LAMBO, Scafire): 0 BADGE (team / already received)

**Formula:**

```
BADGE Awarded = Base Allocation × (Uptime % / 100)

where Base Allocation:
  Mission Decentralization: max(200,000 × uptime%, 50,000)
    (50,000 floor honoring 100% uptime for period 1 of 8/12 - 11/12)
  Other validators:        100,000 × uptime% (no minimum)
  Excluded:                0 (team / already received)
```

**Example:** A Mission Decentralization validator with 95.17% uptime receives 200,000 × 0.9517 = 190,340 BADGE. A non-MD validator with the same uptime receives 100,000 × 0.9517 = 95,170 BADGE. An MD validator with 10% uptime would receive the 50,000 BADGE floor.

**What happens on 8/12/2026:**
Upon reaching the end of this program, whatever BADGE has been earned based on the final uptime calculations will be awarded to each validator. Following this, all current delegations will fully shift to the proof of authority model. Note that current delegations are not reflective of awards. Most validators are delegated ~200K+ which will go away at this time.

**2. Proof of Authority Model**

The proof of authority model is effective immediately and will run in parallel with the existing validator rewards program until 8/12/2026, at which point it becomes the sole delegation model.

The remaining BADGE allocations — including the 50M community pool, other team delegations, and leftover awards — will be allocated to well-known, trusted validators via a proof of authority setup. This follows a know-your-validator setup, with preference given to institutions and well-known brands.

_How is this determined?_ Through governance proposals. The community will vote on which validators receive delegations under the proof of authority model.

_How do delegations work?_ The specifics are determined per proposal, but delegations will simply be delegations with no selling. Commissions are fine. These delegated tokens can only be delegated, never sold.

As a result, the effective circulating supply of BADGE is cut by >60M, since such delegated tokens will only ever be delegated.

**Other Notes:**

- **Disputes**: If a validator believes their uptime was miscalculated (e.g., node was up but missed signatures due to network issues), they can contact the team to review. Adjustments will be made on a case-by-case basis.
- **Relayer bonus**: An IBC relayer bonus has already been distributed. No further rewards will be issued for relayer status.
- **Malicious behavior**: Validators found to have engaged in malicious activity (double signing, etc.) may be disqualified from rewards entirely at the team's discretion.
