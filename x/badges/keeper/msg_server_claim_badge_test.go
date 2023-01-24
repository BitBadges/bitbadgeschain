package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/merkle"
)

func (suite *TestSuite) TestSendAllToClaimsAndClaim() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	rootHash, merkleProofs := merkle.ProofsFromByteSlices([][]byte{[]byte(alice), []byte(bob), []byte(charlie)})

	

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	claimToAdd := types.Claim{
		Balance: &types.Balance{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
		AmountPerClaim: 1,
		Data:       rootHash,
		Type: 	 	types.ClaimType_AccountNum,
		Uri: "",
		TimeRange: &types.IdRange{
			Start: 0,
			End:   math.MaxUint64,
		},
	}

	err = CreateBadges(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	},
		[]*types.Transfers{},
		[]*types.Claim{
			&claimToAdd,
		})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)

	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.MaxSupplys)

	claim, err := GetClaim(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(claimToAdd, claim)

	err = ClaimBadge(suite, wctx, alice, 0, 0,  []byte(alice), (*types.Proof)(merkleProofs[0]), "", &types.IdRange{
		Start: 0,
		End:   math.MaxUint64,
	})
	suite.Require().Nil(err, "Error claiming badge")

	aliceBalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(1), aliceBalance.Balances[0].Balance)
	suite.Require().Equal([]*types.IdRange{{Start: 0, End: 0}}, aliceBalance.Balances[0].BadgeIds)

	claim, err = GetClaim(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(uint64(9), claim.Balance.Balance)
	// suite.Require().Equal([]*types.IdRange{{Start: 0, End: 0}}, aliceBalance.Balances[0].BadgeIds)
}


func (suite *TestSuite) TestSendAllToClaimsAccountTypeInvalid() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	rootHash, merkleProofs := merkle.ProofsFromByteSlices([][]byte{[]byte("121241234"), []byte(bob), []byte(charlie)})

	

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	claimToAdd := types.Claim{
		Balance: &types.Balance{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
		AmountPerClaim: 1,
		Data:       rootHash,
		Type: 	 	types.ClaimType_AccountNum,
	}

	err = CreateBadges(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	},
		[]*types.Transfers{},
		[]*types.Claim{
			&claimToAdd,
		})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)

	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.MaxSupplys)

	claim, err := GetClaim(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(claimToAdd, claim)

	err = ClaimBadge(suite, wctx, alice, 0, 0, []byte("121241234"), (*types.Proof)(merkleProofs[0]), "", &types.IdRange{
		Start: 0,
		End:   math.MaxUint64,
	})
	suite.Require().EqualError(err, keeper.ErrClaimDataInvalid.Error())
}



func (suite *TestSuite) TestSendAllToClaimsAccountTypeCodes() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	rootHash, merkleProofs := merkle.ProofsFromByteSlices([][]byte{[]byte("121241234"), []byte(bob), []byte(charlie)})

	

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	claimToAdd := types.Claim{
		Balance: &types.Balance{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
		AmountPerClaim: 1,
		Data:       rootHash,
		Type: 	 	types.ClaimType_Code,
	}

	err = CreateBadges(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	},
		[]*types.Transfers{},
		[]*types.Claim{
			&claimToAdd,
		})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)

	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.MaxSupplys)

	claim, err := GetClaim(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(claimToAdd, claim)

	err = ClaimBadge(suite, wctx, alice, 0, 0, []byte("121241234"), (*types.Proof)(merkleProofs[0]), "", &types.IdRange{
		Start: 0,
		End:   math.MaxUint64,
	})
	suite.Require().Nil(err, "Error claiming badge")

	aliceBalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(1), aliceBalance.Balances[0].Balance)
	suite.Require().Equal([]*types.IdRange{{Start: 0, End: 0}}, aliceBalance.Balances[0].BadgeIds)

	claim, err = GetClaim(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(uint64(9), claim.Balance.Balance)
	// suite.Require().Equal([]*types.IdRange{{Start: 0, End: 0}}, aliceBalance.Balances[0].BadgeIds)
}
