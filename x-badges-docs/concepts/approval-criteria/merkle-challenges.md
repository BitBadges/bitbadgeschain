# Merkle Challenges

```typescript
export interface MerkleChallenge<T extends NumberType> {
    root: string;
    expectedProofLength: T;
    useCreatorAddressAsLeaf: boolean;
    maxUsesPerLeaf: T;
    uri: string;
    customData: string;
    challengeTrackerId: string;
    leafSigner: string;
}
```

<pre class="language-json"><code class="lang-json"><strong>"merkleChallenge": {
</strong>   "root": "758691e922381c4327646a86e44dddf8a2e060f9f5559022638cc7fa94c55b77",
   "expectedProofLength": "1",
   "useCreatorAddressAsLeaf": true,
   "maxOneUsePerLeaf": true,
   "uri": "ipfs://Qmbbe75FaJyTHn7W5q8EaePEZ9M3J5Rj3KGNfApSfJtYyD",
   "customData": "",
   "challengeTrackerId": "uniqueId",
   "leafSigner": "0x"
}
</code></pre>

Merkle challenges allow you to define a SHA256 Merkle tree, and to be approved for each transfer, the initiator of the transfer must provide a valid Merkle path for the tree when they transfer (via **merkleProofs** in [MsgTransferBadges](../../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md)).

For example, you can create a Merkle tree of claim codes. Then to be able to claim badges, each claimee must provide a valid unused Merkle path from the claim code to the **root**. You distribute the secret leaves / paths in any method you prefer.

Or, you can create an whitelist tree where the user's addresses are the leaves, and they must specify the valid Merkle path from their address to claim. This can be used to distribute gas costs among N users rather than the collection creator defining an address list with N users on-chain and paying all gas fees.

#### Expected Proof Length

The **expectedProofLength** defines the expected length for the Merkle proofs to be provided. This avoids preimage and second preimage attacks. **All proofs must be of the same length, which means you must design your trees accordingly. THIS IS CRITICAL.**

**Whitelist Trees**

Whitelist trees can be used to distribute gas costs among N users rather than the collection creator defining an expensive address list with N users on-chain and paying all gas fees. For small N, we recommend not using whitelist trees for user experience.

If defining a whitelist tree, note that the initiator must also be within the **initiatedByList** of the approval for it to make sense. Typically, **initiatedByList** will be set to "All" and then the whitelist tree restricts who can initiate.

To create a whitelist tree, you need to set **useCreatorAddressAsLeaf** to true. If **useCreatorAddressAsLeaf** is set to true, we will override the provided leaf of each Merkle proof with the BitBadges address of the initiator of the transfer transaction.

**Max Uses per Leaf**

For whitelist trees (**useCreatorAddressAsLeaf** is true), **maxUsesPerLeaf** can be set to any number. "0" or null means unlimited uses. "1" means max one use per leaf and so on. When **useCreatorAddressAsLeaf** is false, this must be set to "1" to avoid replay attacks. For example, ensure that a code / proof can only be used once because once used once, the blockchain is public and anyone then knows the secret code.

We track this in a challenge tracker, similar to the approvals trackers previously explained. We simply track if a leaf index (leftmost leaf of expected proof length layer (aka leaf layer) = index 0, ...) has been used and only allow it to be used **maxUsesPerLeaf** many times, if constrained.

The identifier for each challenge tracker consists of **challengeTrackerId** along with other identifying details seen below. The full ID contains the **approvalId,** so you know state will always be scoped to an approval and the tracker cannot be used by any other approval.

Like approval trackers, this is increment only and non-deletable. Thus, it is critical to not use a tracker with prior history if you intend for it to start tracking from scratch. This can be achieved by using an unused **challengeTrackerId**. If updating an approval with a challenge, please consider how the challenge tracker is working behind the scenes.

```typescript
{
  collectionId: T;
  approvalId: string;
  approvalLevel: "collection" | "incoming" | "outgoing";
  approverAddress?: string;
  challengeTrackerId: string;
  leafIndex: T;
}
```

**approvalLevel** corresponds to whether it is a collection-level approval, user incoming approval, or user outgoing approval. If it is user level, the **approverAddress** is the user setting the approval. **approverAddress** is blank for collection level.

Example:

`1-collection- -approvalId-uniqueID-0` -> USED 1 TIME

`1-collection- -approvalId-uniqueID-1` -> UNUSED

**Reserving Specific Leafs**

See Predetermined Balances below for reserving specific leaf indices for specific badges / ownership times.

**Leaf Signatures**

Leaf signatures are a protection against man-in-the-middle attacks. For code-based merkle challenges, there is always a risk that the code is intercepted while the transaction is in the mempool, and the malicious actor can try to claim the badge with the intercepted code before the user can.

If **leafSigner** is set, the leaf must be signed by the leaf signer. We currently only support leafSigner being an Ethereum address and signatures being ECDSA signatures.

The scheme we currently use is as follows:\
signature = ETHSign(leaf + "-" + bitbadgesAddressOfInitiator)

Then the user must provide the **leafSignature** in the **merkleProofs** field of the transfer transaction.

Note: The bitbadgesAddressOfInitiator is the converted BitBadges address (bb1...) of the initiator of the transfer transaction. This also helps to tie a specific code to a specific BitBadges address to prevent other users from using and intercepting the same code.

This is optional but strongly recommended for code-based merkle challenges.

#### **Creating a Merkle Tree**

We provide the **treeOptions** field in the SDK to let you define your own build options for the tree (see [Compatibility](../../../bitbadges-api/concepts/designing-for-compatibility.md) with the BitBadges API / Indexer). You may experiment with this, but please test all Merkle paths and claims work as intended first. The only tested build options so far are what you see below with the fillDefaultHash.

The important part is making sure all leaves are on the same layer and have the same proof length, or else, they will fail on-chain.

```typescript
import { SHA256 } from 'crypto-js';
import MerkleTree from 'merkletreejs';

const codes = [...]
const hashedCodes = codes.map(x => SHA256(x).toString());
const treeOptions = { fillDefaultHash: '0000000000000000000000000000000000000000000000000000000000000000' }
const codesTree = new MerkleTree(hashedCodes, SHA256, treeOptions);
const codesRoot = codesTree.getRoot().toString('hex');
const expectedMerkleProofLength = codesTree.getLayerCount() - 1;
```

For whitelists, replace with this code.

```typescript
addresses.push(...toAddresses.map((x) => convertToBitBadgesAddress(x)));

const addressesTree = new MerkleTree(
    addresses.map((x) => SHA256(x)),
    SHA256,
    treeOptions
);
const addressesRoot = addressesTree.getRoot().toString('hex');
```

A valid proof can then be created via where codeToSubmit is the code submitted by the user.

```typescript
const passwordCodeToSubmit = '....'
const leaf = isWhitelist ? SHA256(chain.bitbadgesAddress).toString() : SHA256(passwordCodeToSubmit).toString();
const proofObj = tree?.getProof(leaf, whitelistIndex !== undefined && whitelistIndex >= 0 ? whitelistIndex : undefined);
const isValidProof = proofObj && tree && proofObj.length === tree.getLayerCount() - 1;

const leafSignature = '...';


const codeProof = {
  aunts: proofObj ? proofObj.map((proof) => {
    return {
      aunt: proof.data.toString('hex'),
      onRight: proof.position === 'right'
    }
  }) : [],
  leaf: isWhitelist ? '' : passwordCodeToSubmit,
  leafSignature //if applicable
}

const txCosmosMsg: MsgTransferBadges<bigint> = {
  creator: chain.bitbadgesAddress,
  collectionId: collectionId,
  transfers: [{
    ...
    merkleProofs: requiresProof ? [codeProof] : [],
    ...
  }],
};
```
