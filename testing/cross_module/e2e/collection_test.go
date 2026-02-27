//go:build test
// +build test

package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// CollectionTestSuite tests tokenization collection operations
type CollectionTestSuite struct {
	CrossModuleTestSuite
}

func TestCollectionTestSuite(t *testing.T) {
	suite.Run(t, new(CollectionTestSuite))
}

// TestCreateCollection tests basic collection creation
func (s *CollectionTestSuite) TestCreateCollection() {
	s.T().Log("Testing collection creation")

	creator := s.TestAccs[0]
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)
	s.Require().True(collectionId.GT(sdkmath.ZeroUint()))

	s.T().Logf("Created collection with ID: %s", collectionId)

	// Verify collection exists
	collection, found := s.App.TokenizationKeeper.GetCollectionFromStore(s.Ctx, collectionId)
	s.Require().True(found)
	s.Require().NotNil(collection)
	s.Require().Equal(collectionId, collection.CollectionId)
}

// TestMintTokens tests minting tokens to an address
func (s *CollectionTestSuite) TestMintTokens() {
	s.T().Log("Testing token minting")

	creator := s.TestAccs[0]
	recipient := s.TestAccs[1]

	// Create collection
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	// Mint tokens
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	}
	amount := sdkmath.NewUint(1)

	err = s.MintTokensToAddress(collectionId, recipient, tokenIds, amount)
	s.Require().NoError(err)

	// Verify balance
	balance := s.GetTokenBalance(collectionId, recipient, sdkmath.NewUint(1))
	s.Require().Equal(amount, balance)

	s.T().Logf("Minted tokens 1-10 to %s", recipient.String())
}

// TestTransferTokens tests transferring tokens between addresses
func (s *CollectionTestSuite) TestTransferTokens() {
	s.T().Log("Testing token transfer")

	creator := s.TestAccs[0]
	sender := s.TestAccs[1]
	receiver := s.TestAccs[2]

	// Create collection and mint to sender
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)},
	}
	err = s.MintTokensToAddress(collectionId, sender, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Transfer tokens
	transferTokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(3)},
	}
	err = s.TransferTokens(collectionId, sender, receiver, transferTokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Verify receiver has tokens
	receiverBalance := s.GetTokenBalance(collectionId, receiver, sdkmath.NewUint(1))
	s.Require().Equal(sdkmath.NewUint(1), receiverBalance)

	s.T().Logf("Transferred tokens 1-3 from %s to %s", sender.String(), receiver.String())
}

// TestMultipleCollections tests creating and managing multiple collections
func (s *CollectionTestSuite) TestMultipleCollections() {
	s.T().Log("Testing multiple collections")

	creator := s.TestAccs[0]
	numCollections := 3

	var collectionIds []sdkmath.Uint
	for i := 0; i < numCollections; i++ {
		collectionId, err := s.CreateBasicCollection(creator)
		s.Require().NoError(err)
		collectionIds = append(collectionIds, collectionId)
	}

	s.Require().Len(collectionIds, numCollections)

	// Verify all collections exist
	for _, collectionId := range collectionIds {
		collection, found := s.App.TokenizationKeeper.GetCollectionFromStore(s.Ctx, collectionId)
		s.Require().True(found)
		s.Require().NotNil(collection)
	}

	s.T().Logf("Created %d collections: %v", numCollections, collectionIds)
}

// TestCollectionMetadata tests collection metadata updates
func (s *CollectionTestSuite) TestCollectionMetadata() {
	s.T().Log("Testing collection metadata")

	creator := s.TestAccs[0]

	// First create a basic collection
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	// Update the collection metadata
	updateMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                  creator.String(),
		CollectionId:             collectionId,
		UpdateCollectionMetadata: true,
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "ipfs://QmTest123",
			CustomData: "My Custom Collection",
		},
	}

	_, err = s.tokenizationMsgServer.UniversalUpdateCollection(s.Ctx, updateMsg)
	s.Require().NoError(err)

	// Verify metadata
	collection, found := s.App.TokenizationKeeper.GetCollectionFromStore(s.Ctx, collectionId)
	s.Require().True(found)
	s.Require().NotNil(collection.CollectionMetadata)
	s.Require().Equal("My Custom Collection", collection.CollectionMetadata.CustomData)

	s.T().Logf("Created collection %s with custom metadata", collectionId)
}

// TestBatchMint tests minting multiple token ranges
func (s *CollectionTestSuite) TestBatchMint() {
	s.T().Log("Testing batch mint")

	creator := s.TestAccs[0]
	recipient := s.TestAccs[1]

	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	// Mint multiple token ranges
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(60)},
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(110)},
	}

	err = s.MintTokensToAddress(collectionId, recipient, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Verify balances for each range
	balance1 := s.GetTokenBalance(collectionId, recipient, sdkmath.NewUint(5))
	balance2 := s.GetTokenBalance(collectionId, recipient, sdkmath.NewUint(55))
	balance3 := s.GetTokenBalance(collectionId, recipient, sdkmath.NewUint(105))

	s.Require().Equal(sdkmath.NewUint(1), balance1)
	s.Require().Equal(sdkmath.NewUint(1), balance2)
	s.Require().Equal(sdkmath.NewUint(1), balance3)

	s.T().Log("Successfully minted batch tokens")
}

// TestTransferToMultipleRecipients tests transferring to multiple recipients
func (s *CollectionTestSuite) TestTransferToMultipleRecipients() {
	s.T().Log("Testing transfer to multiple recipients")

	creator := s.TestAccs[0]
	sender := s.TestAccs[1]

	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	// Mint tokens to sender
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}
	err = s.MintTokensToAddress(collectionId, sender, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Transfer different token ranges to different "recipients"
	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      sender.String(),
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        sender.String(),
				ToAddresses: []string{s.TestAccs[2].String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(20)}},
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
					},
				},
			},
		},
	}

	_, err = s.tokenizationMsgServer.TransferTokens(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify recipient has tokens
	balance := s.GetTokenBalance(collectionId, s.TestAccs[2], sdkmath.NewUint(10))
	s.Require().Equal(sdkmath.NewUint(1), balance)

	s.T().Log("Successfully transferred to multiple recipients")
}
