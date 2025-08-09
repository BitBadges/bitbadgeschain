# Merkle Challenges

## Overview

Merkle challenges provide cryptographic proof-based approval mechanisms using SHA256 Merkle trees. They enable secure, gas-efficient whitelisting and claim code systems without storing large address lists on-chain.

**Key Benefits**:

-   **Gas Efficiency**: Distribute gas costs among users instead of collection creators
-   **Security**: Cryptographic proof verification prevents unauthorized access
-   **Flexibility**: Support both whitelist trees and claim code systems
-   **Scalability**: Handle large user bases without on-chain storage

## Interface Definition

```typescript
export interface MerkleChallenge<T extends NumberType> {
    root: string; // SHA256 Merkle tree root hash
    expectedProofLength: T; // Required proof length (security)
    useCreatorAddressAsLeaf: boolean; // Use initiator address as leaf?
    maxUsesPerLeaf: T; // Maximum uses per leaf
    uri: string; // Metadata URI
    customData: string; // Custom data field
    challengeTrackerId: string; // Unique tracker identifier
    leafSigner: string; // Optional leaf signature authority
}
```

## Basic Example

```json
{
    "merkleChallenges": [
        {
            "root": "758691e922381c4327646a86e44dddf8a2e060f9f5559022638cc7fa94c55b77",
            "expectedProofLength": "1",
            "useCreatorAddressAsLeaf": false,
            "maxUsesPerLeaf": "1",
            "uri": "ipfs://Qmbbe75FaJyTHn7W5q8EaePEZ9M3J5Rj3KGNfApSfJtYyD",
            "customData": "",
            "challengeTrackerId": "uniqueId",
            "leafSigner": "0x"
        }
    ]
}
```

## Challenge Types

### 1. Claim Code Challenges

Create a Merkle tree of secret claim codes that users must provide to claim tokens.

**Use Case**: Private claim codes, invitation systems, promotional campaigns

**Process**:

1. Generate secret claim codes
2. Build Merkle tree from hashed codes
3. Distribute codes privately to users with leaf signatures
4. Users provide code + Merkle proof in transfer

### 2. Whitelist Challenges

Create a Merkle tree of user addresses for gas-efficient whitelisting.

**Use Case**: Large whitelists, community access, gas cost distribution

**Process**:

1. Collect user addresses
2. Build Merkle tree from hashed addresses
3. Users provide their address + Merkle proof
4. System verifies address is in whitelist / valid proof

**Gas Cost Distribution**: Instead of the collection creator paying gas to store N addresses on-chain, each user pays their own gas for proof verification.

## Understanding useCreatorAddressAsLeaf

The `useCreatorAddressAsLeaf` field determines how the system handles the leaf value in Merkle proofs:

### Whitelist Trees (`useCreatorAddressAsLeaf: true`)

**Purpose**: Verify that the transaction initiator is in the whitelist.

**How It Works**:

1. **Automatic Override**: The system expects the provided leaf to be the initiator's BitBadges address ("bb1...")
2. **Address Verification**: Checks if the initiator's address exists in the Merkle tree
3. **No Manual Leaf**: Users don't need to provide their address as the leaf - the system handles it

**Recommended Configuration**:

-   Set `initiatedByList` to "All" (whitelist tree handles the restriction)
-   Set `useCreatorAddressAsLeaf: true`
-   Build Merkle tree from BitBadges addresses as leaves ["bb1...", "bb2...", "bb3..."]

### Claim Code Trees (`useCreatorAddressAsLeaf: false`)

**Purpose**: Verify that the user possesses a valid claim code.

**How It Works**:

1. **Manual Leaf**: User must provide the actual claim code as the leaf
2. **Code Verification**: System verifies the provided code exists in the Merkle tree
3. **User Responsibility**: Users must know and provide their claim code

**Recommended Configuration**:

-   Set `useCreatorAddressAsLeaf: false`
-   Build Merkle tree from claim codes as leaves ["secret1", "secret2", "secret3"]
-   Post root hash on-chain as challenge
-   Distribute codes privately to users with leaf signatures

## Security Features

### Expected Proof Length

**Critical Security Feature**: All proofs must have the same length to prevent preimage and second preimage attacks.

```typescript
// All proofs must match this length
expectedProofLength: '2'; // 2-level proof required
```

**Design Requirement**: Your Merkle tree must be constructed so all leaves are at the same depth.

### Max Uses Per Leaf

Control how many times each leaf can be used:

| Setting         | Behavior          | Use Case             |
| --------------- | ----------------- | -------------------- |
| `"0"` or `null` | Unlimited uses    | Public claim codes   |
| `"1"`           | One-time use      | Single-use codes     |
| `"5"`           | Five uses maximum | Limited distribution |

**Critical Security Requirement**: For claim code challenges (`useCreatorAddressAsLeaf: false`), `maxUsesPerLeaf` must be `"1"` to prevent replay attacks.

### Replay Attack Protection

**⚠️ CRITICAL SECURITY RISK**: Non-address trees (claim codes) are vulnerable to front-running attacks.

**The Problem**:

1. User submits transaction with valid Merkle proof
2. Proof becomes visible in mempool (public blockchain)
3. Malicious actor sees the proof and front-runs the transaction
4. Original user's transaction fails, attacker gets the token

**Why This Happens**:

-   Merkle proofs for claim codes are reusable until consumed
-   Once in mempool, proofs are publicly visible
-   No built-in protection against proof reuse

**The Solution**: Leaf signatures provide cryptographic protection against this attack.

## Challenge Tracking

### Tracker System

Uses increment-only, immutable trackers to prevent double-spending:

```typescript
{
    collectionId: T;
    approvalId: string;
    approvalLevel: 'collection' | 'incoming' | 'outgoing';
    approverAddress: string; // blank if collection-level
    challengeTrackerId: string;
    leafIndex: T; // Leftmost base layer leaf index = 0, rightmost = numLeaves - 1
}
```

Note the fact we use leaf indices to track usage and not leaf values.

### Tracker Examples

```
1-collection- -approvalId-uniqueID-0  → USED 1 TIME
1-collection- -approvalId-uniqueID-1  → UNUSED
1-collection- -approvalId-uniqueID-2  → USED 3 TIMES
```

**Important**: Trackers are scoped to specific approvals and cannot be shared between different approval configurations.

### Tracker Management

-   **Increment-Only**: Once used, the number of uses cannot be decremented
-   **Immutable**: Tracker state cannot be modified
-   **Best Practice**: Use unique `challengeTrackerId` for fresh tracking of new approvals

## Leaf Signatures

### Protection Against Front-Running

Leaf signatures provide cryptographic protection against front-running attacks on claim code challenges.

**How It Works**:

```typescript
// Signature scheme
signature = ETHSign(leaf + '-' + bitbadgesAddressOfInitiator);
```

**Security Mechanism**:

1. **Address Binding**: Each proof is cryptographically tied to a specific BitBadges address
2. **Replay Prevention**: Even if proof is intercepted, it cannot be used by other addresses
3. **Mempool Safety**: Intercepted proofs in mempool are useless to attackers

### Implementation

```typescript
// Only Ethereum addresses supported currently
leafSigner: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6';
```

**Critical Benefits**:

-   **Front-Running Protection**: Prevents attackers from stealing tokens via mempool interception
-   **Address-Specific**: Each proof is cryptographically bound to the intended recipient
-   **Mempool Safety**: Makes intercepted proofs useless to malicious actors
-   **Required for Claim Codes**: Strongly recommended for all non-address tree challenges

**⚠️ IMPORTANT**: For claim code challenges, leaf signatures are not just recommended—they are essential for security against front-running attacks.

## Merkle Tree Construction

### Standard Configuration

```typescript
import { SHA256 } from 'crypto-js';
import MerkleTree from 'merkletreejs';

// For claim codes
const codes = ['secret1', 'secret2', 'secret3'];
const hashedCodes = codes.map((x) => SHA256(x).toString());

// For whitelists
const addresses = ['bb1...', 'bb1...', 'bb1...'];
const hashedAddresses = addresses.map((x) => SHA256(x));

// Tree options (tested configuration)
const treeOptions = {
    fillDefaultHash:
        '0000000000000000000000000000000000000000000000000000000000000000',
};

// Build tree
const tree = new MerkleTree(hashedCodes, SHA256, treeOptions);
const root = tree.getRoot().toString('hex');
const expectedProofLength = tree.getLayerCount() - 1;
```

### Critical Requirements

1. **Same Layer**: All leaves must be at the same depth
2. **Consistent Proof Length**: All proofs must have identical length
3. **Test Thoroughly**: Verify all paths work before deployment
4. **Use Tested Options**: Stick to the `fillDefaultHash` configuration

## Transfer Integration

### Providing Proofs

Include Merkle proofs in [MsgTransferTokens](../../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md):

```typescript
const txCosmosMsg: MsgTransferTokens<bigint> = {
    creator: chain.bitbadgesAddress,
    collectionId: collectionId,
    transfers: [
        {
            // ... other fields
            merkleProofs: [
                {
                    aunts: proofObj.map((proof) => ({
                        aunt: proof.data.toString('hex'),
                        onRight: proof.position === 'right',
                    })),
                    leaf: isWhitelist ? '' : passwordCodeToSubmit,
                    leafSignature: leafSignature, // if applicable
                },
            ],
        },
    ],
};
```

### Proof Generation

```typescript
// Generate proof for user submission
const passwordCodeToSubmit = 'secretCode123';
const leaf = isWhitelist
    ? SHA256(chain.bitbadgesAddress).toString()
    : SHA256(passwordCodeToSubmit).toString();

const proofObj = tree.getProof(leaf, whitelistIndex);
const isValidProof = proofObj && proofObj.length === tree.getLayerCount() - 1;

// Create signature if needed
const leafSignature = signLeaf(leaf + '-' + chain.bitbadgesAddress);
```

## Comparison with ETH Signature Challenges

Merkle challenges and ETH signature challenges are very similar. The main difference is that Merkle challenges must also check that the signed message was pre-committed to in the tree, whereas ETH signature challenges only need to check that the signature is valid and not used before.

For more information, see [ETH Signature Challenges](eth-signature-challenges.md).

## Best Practices

### Design Considerations

1. **Tree Structure**: Ensure all leaves at same depth
2. **Proof Length**: Test all proof lengths are identical
3. **Tracker Management**: Use unique IDs for fresh tracking
4. **Security**: **MANDATORY** - Enable leaf signatures for claim codes to prevent front-running
5. **Testing**: Verify all paths work before mainnet

### Performance Optimization

1. **Small Lists**: For <100 users, consider regular address lists
2. **Gas Distribution**: Merkle trees excel with large user bases
3. **Proof Verification**: On-chain verification is gas-efficient
4. **Storage**: No on-chain storage of large lists required
