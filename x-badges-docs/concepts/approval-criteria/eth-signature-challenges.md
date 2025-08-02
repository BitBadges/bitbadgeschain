# ETH Signature Challenges

ETH Signature Challenges are a type of approval criteria that require users to provide valid Ethereum signatures from a predetermined signer to complete transfers. The signer approves the address by signing a message that contains the address and a nonce. This feature allows for secure, on-chain verification of off-chain authorization without the complexity of Merkle trees.

## Overview

ETH Signature Challenges work by requiring users to provide Ethereum signatures that prove they have authorization from specific Ethereum addresses. Each signature can only be used once, preventing replay attacks and ensuring the security of the approval system.

## How It Works

### Signature Scheme

The signature scheme follows the pattern:

```
ETHSign(nonce + "-" + creatorAddress)
```

Where:

-   `nonce`: A unique identifier provided by the user
-   `creatorAddress`: The address of the collection creator
-   `-`: A literal dash character separating the two values

### Challenge Structure

Each ETH Signature Challenge contains:

-   `signer`: The Ethereum address that must sign the challenge
-   `challengeTrackerId`: Unique identifier for tracking used signatures
-   `uri`: Optional metadata URI
-   `customData`: Optional custom data

### Proof Structure

Users provide ETH Signature Proofs containing:

-   `nonce`: The nonce that was signed
-   `signature`: The Ethereum signature of the nonce

## Key Features

### One-Time Use Signatures

Each signature can only be used once per challenge tracker. This prevents:

-   Replay attacks
-   Double-spending of approvals
-   Unauthorized reuse of signatures

### Multiple Signers

You can require signatures from multiple Ethereum addresses in a single approval:

```json
{
    "ethSignatureChallenges": [
        {
            "signer": "0x1234567890123456789012345678901234567890",
            "challengeTrackerId": "challenge1"
        },
        {
            "signer": "0x0987654321098765432109876543210987654321",
            "challengeTrackerId": "challenge2"
        }
    ]
}
```

## Implementation Details

### Signature Verification

The system verifies signatures by:

1. Reconstructing the signed message: `nonce + "-" + creatorAddress`
2. Recovering the signer address from the signature
3. Comparing the recovered address with the expected `signer` address
4. Checking that the signature hasn't been used before

### Storage

Used signatures are tracked in the blockchain state using:

-   **Key**: `ETHSignatureTrackerKey` with challenge tracker ID
-   **Value**: Number of times the signature has been used (increment-only per tracker ID)

## Quick Reference

### Interface Definitions

```typescript
interface ETHSignatureChallenge {
    signer: string; // Ethereum address that must sign
    challengeTrackerId: string; // Unique ID for tracking used signatures
    uri?: string; // Optional metadata URI
    customData?: string; // Optional custom data
}

interface ETHSignatureProof {
    nonce: string; // The nonce that was signed
    signature: string; // Ethereum signature
}
```

## Error Handling

Common error scenarios:

-   **Invalid Signature**: Signature doesn't match the expected signer
-   **Already Used**: Signature has been used before
-   **Missing Proof**: Required ETH signature proof not provided
-   **Invalid Nonce**: Nonce format or content is invalid

The system provides clear error messages to help users understand and resolve issues.
