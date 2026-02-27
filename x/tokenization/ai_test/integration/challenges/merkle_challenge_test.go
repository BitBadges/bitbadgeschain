package challenges_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type MerkleChallengeTestSuite struct {
	testutil.AITestSuite
}

func TestMerkleChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(MerkleChallengeTestSuite))
}

func (suite *MerkleChallengeTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// hashLeafBytes hashes leaf data using SHA256 and returns the raw bytes
func hashLeafBytes(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// hashConcatBytes hashes two byte slices concatenated together
func hashConcatBytes(left, right []byte) []byte {
	combined := append(left, right...)
	h := sha256.Sum256(combined)
	return h[:]
}

// createSimpleMerkleTree creates a 2-leaf merkle tree and returns hex-encoded root and proof for leaf1
// Note: leaf values are raw strings, aunts in proof are hex-encoded hashes
func createSimpleMerkleTree(leaf1, leaf2 string) (root string, proof []*types.MerklePathItem) {
	h1 := hashLeafBytes([]byte(leaf1))
	h2 := hashLeafBytes([]byte(leaf2))

	// Root = hash(h1 || h2)
	rootBytes := hashConcatBytes(h1, h2)
	root = hex.EncodeToString(rootBytes)

	// Proof for leaf1: h2 is the aunt, on right side (hex-encoded)
	proof = []*types.MerklePathItem{
		{Aunt: hex.EncodeToString(h2), OnRight: true},
	}

	return root, proof
}

// createFourLeafMerkleTree creates a 4-leaf merkle tree
func createFourLeafMerkleTree(leaf1, leaf2, leaf3, leaf4 string) (root string, proofForLeaf1 []*types.MerklePathItem) {
	h1 := hashLeafBytes([]byte(leaf1))
	h2 := hashLeafBytes([]byte(leaf2))
	h3 := hashLeafBytes([]byte(leaf3))
	h4 := hashLeafBytes([]byte(leaf4))

	// Level 1: hash pairs
	leftBytes := hashConcatBytes(h1, h2)
	rightBytes := hashConcatBytes(h3, h4)

	// Root = hash(left || right)
	rootBytes := hashConcatBytes(leftBytes, rightBytes)
	root = hex.EncodeToString(rootBytes)

	// Proof for leaf1: h2 is first aunt (on right), then right subtree hash is second aunt (on right)
	proofForLeaf1 = []*types.MerklePathItem{
		{Aunt: hex.EncodeToString(h2), OnRight: true},
		{Aunt: hex.EncodeToString(rightBytes), OnRight: true},
	}

	return root, proofForLeaf1
}

// TestMerkleChallenge_ValidProofPasses tests that a valid merkle proof allows transfer
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_ValidProofPasses() {
	// Create merkle tree with two allowed addresses
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, proof := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_challenge",
	}

	// Create collection first without any approvals
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Mint tokens first
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Now add the merkle challenge approval after minting
	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Get current collection and add our approval
	collection := suite.GetCollection(collectionId)
	newApprovals := append(collection.CollectionApprovals, approval)

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       newApprovals,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "adding merkle approval should succeed")

	// Transfer with valid merkle proof
	// Note: Approvals with merkle challenges must be explicitly prioritized
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with valid merkle proof should succeed")
}

// TestMerkleChallenge_InvalidProofRejected tests that an invalid merkle proof rejects transfer
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_InvalidProofRejected() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, _ := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_challenge_invalid",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Create an invalid proof with wrong aunt hash
	invalidProof := []*types.MerklePathItem{
		{Aunt: "invalid_hash_1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", OnRight: true},
	}

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: invalidProof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with invalid merkle proof should fail")
}

// TestMerkleChallenge_ExpectedProofLengthEnforced_TooShort tests that proof shorter than expected is rejected
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_ExpectedProofLengthEnforced_TooShort() {
	// Create a 4-leaf tree which requires 2 proof items
	leaf1 := "leaf_1"
	leaf2 := "leaf_2"
	leaf3 := "leaf_3"
	leaf4 := "leaf_4"
	root, fullProof := createFourLeafMerkleTree(leaf1, leaf2, leaf3, leaf4)

	// Require proof length of 2
	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(2),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_proof_length_short",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Provide truncated proof (only 1 item instead of 2)
	truncatedProof := fullProof[:1]

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: truncatedProof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with too short proof should fail")
}

// TestMerkleChallenge_ExpectedProofLengthEnforced_TooLong tests that proof longer than expected is rejected
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_ExpectedProofLengthEnforced_TooLong() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, validProof := createSimpleMerkleTree(leaf1, leaf2)

	// Require proof length of 1
	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_proof_length_long",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Add extra proof item
	extendedProof := append(validProof, &types.MerklePathItem{
		Aunt:    "extra_hash_1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		OnRight: true,
	})

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: extendedProof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with too long proof should fail")
}

// TestMerkleChallenge_UseCreatorAddressAsLeaf tests that useCreatorAddressAsLeaf hashes creator address as leaf
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_UseCreatorAddressAsLeaf() {
	// Create merkle tree with actual addresses as leaves
	// When useCreatorAddressAsLeaf=true, the system uses the initiator's address as the leaf
	leaf1 := suite.Alice
	leaf2 := suite.Bob
	root, proofForAlice := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: true, // Use creator address as leaf
		MaxUsesPerLeaf:          sdkmath.NewUint(10),
		ChallengeTrackerId:      "test_creator_as_leaf",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer initiated by Alice - should use Alice's address as leaf automatically
	// The proof should be for Alice's address
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  "", // Empty because useCreatorAddressAsLeaf=true overrides this
						Aunts: proofForAlice,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with useCreatorAddressAsLeaf should succeed when initiator is in whitelist tree")
}

// TestMerkleChallenge_MaxUsesPerLeafTracking tests that maxUsesPerLeaf is tracked correctly
// Note: MaxUsesPerLeaf > 1 requires UseCreatorAddressAsLeaf = true
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_MaxUsesPerLeafTracking() {
	// Create merkle tree with actual addresses as leaves
	// When useCreatorAddressAsLeaf=true, the system uses the initiator's address as the leaf
	leaf1 := suite.Alice
	leaf2 := suite.Bob
	root, proofForAlice := createSimpleMerkleTree(leaf1, leaf2)

	// Allow leaf to be used twice (requires UseCreatorAddressAsLeaf = true)
	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: true, // Required when MaxUsesPerLeaf > 1
		MaxUsesPerLeaf:          sdkmath.NewUint(2),
		ChallengeTrackerId:      "test_max_uses",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(30, 1)})

	// First use - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  "", // Empty because useCreatorAddressAsLeaf=true
						Aunts: proofForAlice,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first use of leaf should succeed")

	// Second use - should succeed (maxUsesPerLeaf = 2)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  "", // Empty because useCreatorAddressAsLeaf=true
						Aunts: proofForAlice,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second use of leaf should succeed")
}

// TestMerkleChallenge_LeafExceedsMaxUses tests that using a leaf more than maxUsesPerLeaf is rejected
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_LeafExceedsMaxUses() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, proof := createSimpleMerkleTree(leaf1, leaf2)

	// Allow leaf to be used only once
	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_exceed_max_uses",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(30, 1)})

	// First use - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first use of leaf should succeed")

	// Second use - should fail (maxUsesPerLeaf = 1)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second use of leaf should fail when maxUsesPerLeaf=1")
}

// TestMerkleChallenge_ChallengeTrackerIdScoping tests that challengeTrackerId scopes leaf tracking
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_ChallengeTrackerIdScoping() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, proof := createSimpleMerkleTree(leaf1, leaf2)

	// Create first challenge with tracker ID "tracker_1"
	merkleChallenge1 := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "tracker_1",
	}

	// Create second challenge with different tracker ID "tracker_2"
	merkleChallenge2 := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "tracker_2",
	}

	// Create two approvals with different challenge tracker IDs
	approval1 := testutil.GenerateCollectionApproval("merkle_approval_1", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge1},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("merkle_approval_2", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(30, 1)})

	// Use leaf with first approval (tracker_1)
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval_1"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first use with tracker_1 should succeed")

	// Use same leaf with second approval (tracker_2) - should succeed because different tracker
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval_2"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "use with different tracker_2 should succeed")
}

// TestMerkleChallenge_WrongLeafRejected tests that providing a wrong leaf value is rejected
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_WrongLeafRejected() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, proof := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_wrong_leaf",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Provide wrong leaf but correct proof structure
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  "wrong_leaf_not_in_tree",
						Aunts: proof,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with wrong leaf should fail")
}

// TestMerkleChallenge_EmptyProofRejected tests that an empty proof is rejected when proof is required
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_EmptyProofRejected() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, _ := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_empty_proof",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Provide empty aunts
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1,
						Aunts: []*types.MerklePathItem{}, // Empty
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with empty proof should fail when proof is required")
}

// TestMerkleChallenge_NoMerkleProofProvidedRejected tests that transfer fails when no merkle proof provided
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_NoMerkleProofProvidedRejected() {
	leaf1 := "allowed_address_1"
	leaf2 := "allowed_address_2"
	root, _ := createSimpleMerkleTree(leaf1, leaf2)

	merkleChallenge := &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "test_no_proof",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// No merkle proofs provided at all
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:         suite.Alice,
				ToAddresses:  []string{suite.Bob},
				Balances:     []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{}, // Empty slice
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer without merkle proof should fail when challenge exists")
}

// TestMerkleChallenge_MultipleChallengesAllMustPass tests that all merkle challenges must pass
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_MultipleChallengesAllMustPass() {
	// Create two separate merkle trees
	leaf1a := "tree1_leaf_1"
	leaf1b := "tree1_leaf_2"
	root1, proof1 := createSimpleMerkleTree(leaf1a, leaf1b)

	leaf2a := "tree2_leaf_1"
	leaf2b := "tree2_leaf_2"
	root2, proof2 := createSimpleMerkleTree(leaf2a, leaf2b)

	merkleChallenge1 := &types.MerkleChallenge{
		Root:                    root1,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "challenge_1",
	}

	merkleChallenge2 := &types.MerkleChallenge{
		Root:                    root2,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "challenge_2",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge1, merkleChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer with both valid proofs - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1a,
						Aunts: proof1,
					},
					{
						Leaf:  leaf2a,
						Aunts: proof2,
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with all valid proofs should succeed")
}

// TestMerkleChallenge_MultipleChallengesOneInvalidFails tests that if one challenge fails, transfer fails
func (suite *MerkleChallengeTestSuite) TestMerkleChallenge_MultipleChallengesOneInvalidFails() {
	leaf1a := "tree1_leaf_1"
	leaf1b := "tree1_leaf_2"
	root1, proof1 := createSimpleMerkleTree(leaf1a, leaf1b)

	leaf2a := "tree2_leaf_1"
	leaf2b := "tree2_leaf_2"
	root2, _ := createSimpleMerkleTree(leaf2a, leaf2b)

	merkleChallenge1 := &types.MerkleChallenge{
		Root:                    root1,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "challenge_1",
	}

	merkleChallenge2 := &types.MerkleChallenge{
		Root:                    root2,
		ExpectedProofLength:     sdkmath.NewUint(1),
		UseCreatorAddressAsLeaf: false,
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      "challenge_2",
	}

	approval := testutil.GenerateCollectionApproval("merkle_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MerkleChallenges:               []*types.MerkleChallenge{merkleChallenge1, merkleChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer with first valid proof but invalid second proof
	invalidProof := []*types.MerklePathItem{
		{Aunt: "invalid_hash_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", OnRight: true},
	}

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  leaf1a,
						Aunts: proof1,
					},
					{
						Leaf:  leaf2a,
						Aunts: invalidProof, // Invalid proof for second challenge
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("merkle_approval"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer should fail when one challenge has invalid proof")
}
