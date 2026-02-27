# Quint Formal Specifications

This directory contains [Quint](https://github.com/informalsystems/quint) formal specifications for BitBadges Chain's approval system.

## What is Quint?

Quint is a modern specification language designed for distributed systems and blockchain protocols. It allows us to:

- **Model state machines**: Define the valid states and transitions of our modules
- **Express invariants**: Formally specify properties that must always hold
- **Run simulations**: Find potential bugs through randomized exploration
- **Verify properties**: Use model checking to prove invariants hold

## Specifications

### Approval System (`tokenization/`)

| Spec | Description | Key Invariants |
|------|-------------|----------------|
| `approval_hierarchy.qnt` | Three-layer approval system (collection → outgoing → incoming) with override semantics | Collection approval required, override flags work correctly |
| `amount_limits.qnt` | Approval amount tracking - ensures users can't exceed approved transfer limits | Used amount never exceeds max, non-negative tracking |
| `replay_protection.qnt` | Version-based replay protection for approvals | Consumed version never exceeds approval version |

## Prerequisites

Install Quint CLI:

```bash
npm install -g @informalsystems/quint
```

Verify installation:

```bash
quint --version
```

## Running Specs Locally

### Type-check all specs

```bash
make quint-check
```

Or manually:

```bash
quint typecheck specs/quint/tokenization/*.qnt
```

### Run simulations

Find invariant violations through randomized execution:

```bash
make quint-run
```

Or manually:

```bash
quint run specs/quint/tokenization/approval_hierarchy.qnt \
  --invariant=inv_all \
  --max-steps=50
```

### Full verification (requires JDK 17+)

Use model checking to exhaustively verify invariants:

```bash
make quint-verify
```

## Security Properties Verified

1. **Approval Hierarchy**: Collection-level approval is always required; user-level approvals can be overridden by collection-level flags
2. **Amount Limits**: Transfer amounts tracked correctly; can't exceed approved limits
3. **Replay Protection**: Version increments prevent reuse of consumed approvals
4. **Balance Conservation**: Total supply never changes through transfers
5. **No Negative Balances**: Addresses can't go below zero

## Resources

- [Quint Documentation](https://quint-lang.org/)
- [Quint GitHub](https://github.com/informalsystems/quint)
- [Quint Tutorials](https://quint-lang.org/docs/tutorials)
