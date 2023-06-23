package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSafeAdd() {
	result, err := keeper.SafeAdd(sdk.NewUint(0), sdk.NewUint(1))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, sdk.NewUint(1))

	result, err = keeper.SafeAdd(sdk.NewUint(math.MaxUint64), sdk.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, sdk.NewUint(math.MaxUint64))

	_, err = keeper.SafeAdd(sdk.NewUint(math.MaxUint), sdk.NewUint(1))
	// suite.Require().EqualError(err, keeper.ErrOverflow.Error()) With Cosmos SDK Uint now, this error is not returned
}

func (suite *TestSuite) TestSafeSubtract() {
	result, err := keeper.SafeSubtract(sdk.NewUint(1), sdk.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, sdk.NewUint(1))

	result, err = keeper.SafeSubtract(sdk.NewUint(math.MaxUint64), sdk.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, sdk.NewUint(math.MaxUint64))

	_, err = keeper.SafeSubtract(sdk.NewUint(1), sdk.NewUint(math.MaxUint64))
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}

func (suite *TestSuite) TestUpdateAndGetBalancesForIds() {
	err := *new(error)
	balances := []*types.Balance{
		{
			Amount: sdk.NewUint(1),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(1),
				},
			},
		},
	}

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(1),
		},
	}, sdk.NewUint(10), balances)
	suite.Require().Nil(err, "Error updating balances: %s")

	fetchedBalances, err := keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(1),
		},
	}, balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	suite.Require().Equal(balances, []*types.Balance{
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(1),
				},
			},
		},
	})
	suite.Require().Equal(balances, fetchedBalances)

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(1),
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(2),
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: sdk.NewUint(0),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(2),
					End:   sdk.NewUint(2),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(1),
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(math.MaxUint64),
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: sdk.NewUint(0),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(2),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(1),
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(3),
			End:   sdk.NewUint(math.MaxUint64),
		},
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(2),
		},
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(1),
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: sdk.NewUint(0),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(2),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(1),
				},
			},
		},
	})

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
	}, sdk.NewUint(5), balances)

	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
		{
			Amount: sdk.NewUint(5),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(1),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(0),
				},
			},
		},
	}))

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(2),
			End:   sdk.NewUint(math.MaxUint64),
		},
	}, sdk.NewUint(5), balances)

	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
		{
			Amount: sdk.NewUint(5),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(0),
				},
			},
		},
	}))

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: sdk.NewUint(2),
			End:   sdk.NewUint(2),
		},
	}, sdk.NewUint(10), balances)

	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
		{
			Amount: sdk.NewUint(5),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(1),
				},
				{
					Start: sdk.NewUint(3),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount: sdk.NewUint(10),
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(0),
				},
				{
					Start: sdk.NewUint(2),
					End:   sdk.NewUint(2),
				},
			},
		},
	}))
}

func (suite *TestSuite) TestSubtractBalances() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(0),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End:   sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
		{
			Start: sdk.NewUint(135),
			End:   sdk.NewUint(200),
		},
		{
			Start: sdk.NewUint(235),
			End:   sdk.NewUint(300),
		},
		{
			Start: sdk.NewUint(335),
			End:   sdk.NewUint(400),
		},
	}

	err := *new(error)
	UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, badgeIdRanges, sdk.NewUint(1000))
	suite.Require().NoError(err)

	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: sdk.NewUint(0), End: sdk.NewUint(100),
		},
		{
			Start: sdk.NewUint(135),
			End:   sdk.NewUint(200),
		},
		{
			Start: sdk.NewUint(235),
			End:   sdk.NewUint(300),
		},
		{
			Start: sdk.NewUint(335),
			End:   sdk.NewUint(400),
		},
	})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
	}
	for _, badgeIdRangeToRemove := range badgeIdRangesToRemove {
		UserBalance.Balances, err = keeper.SubtractBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRangeToRemove}, sdk.NewUint(1))
		suite.Require().NoError(err)
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(998))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, sdk.NewUint(999))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})
}

func (suite *TestSuite) TestAddBalancesForIdRanges() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(0),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End:   sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdk.NewUint(1000))
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(34)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, sdk.NewUint(2000))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

	badgeIdRangesToAdd := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
	}

	for _, badgeIdRangeToAdd := range badgeIdRangesToAdd {
		UserBalance.Balances, _ = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRangeToAdd}, sdk.NewUint(1))
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().True(types.IdRangeEquals(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(0)}, {Start: sdk.NewUint(2), End: sdk.NewUint(34)}}))

	suite.Require().Equal(UserBalance.Balances[1].Amount, sdk.NewUint(1001))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

	suite.Require().Equal(UserBalance.Balances[2].Amount, sdk.NewUint(2002))
	suite.Require().Equal(UserBalance.Balances[2].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})
}

// func (suite *TestSuite) TestAddBalancesOverflow() {
// 	UserBalance := types.UserBalanceStore{}
// 	badgeIdRanges := []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(1),
// 		},
// 		{
// 			Start: sdk.NewUint(0),
// 			End:   sdk.NewUint(0),
// 		},
// 		{
// 			Start: sdk.NewUint(35),
// 			End: sdk.NewUint(35),
// 		},
// 		{
// 			Start: sdk.NewUint(2),
// 			End: sdk.NewUint(34),
// 		},
// 		{
// 			Start: sdk.NewUint(35),
// 			End: sdk.NewUint(100),
// 		},
// 	}

// 	err := *new(error)
// 	for _, badgeIdRange := range badgeIdRanges {
// 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdk.NewUint(1000))
// 		suite.Require().Nil(err, "error adding balance to approval")
// 	}

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(1000))
// 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(34)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

// 	suite.Require().Equal(UserBalance.Balances[1].Amount, sdk.NewUint(2000))
// 	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

// 	badgeIdRangesToAdd := []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(0),
// 			End: sdk.NewUint(1000),
// 		},
// 	}

// 	for _, badgeIdRange := range badgeIdRangesToAdd {
// 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdk.NewUint(math.MaxUint64))
// 		suite.Require().EqualError(err, keeper.ErrOverflow.Error())
// 	}
// }

func (suite *TestSuite) TestRemoveBalancesUnderflow() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(0),
			End:   sdk.NewUint(0),
		},
		{
			Start: sdk.NewUint(2),
			End:   sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{&badgeIdRange}, sdk.NewUint(1000))
		suite.Require().NoError(err)
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(100)}})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End:   sdk.NewUint(100),
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		UserBalance.Balances, err = keeper.SubtractBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdk.NewUint(math.MaxUint64))
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}

//Commented because it takes +- 25 seconds to run

// func (suite *TestSuite) TestBalancesFuzz() {
// 	for i := 0; i < 25; i++ {
// 		userBalance := types.UserBalanceStore{}
// 		balances := make([]uint64, sdk.NewUint(1000))

// 		adds := make([]*types.IdRange, 100)
// 		subs := make([]*types.IdRange, 100)
// 		for i := 0; i < 100; i++ { //10000 iterations
// 			//Get random start value
// 			start := sdk.NewUint(rand.Intn(500))
// 			//Get random end value
// 			end := sdk.NewUint(500 + rand.Intn(500))

// 			amount := sdk.NewUint(rand.Intn(100))
// 			err := *new(error)
// 			userBalance.Balances, err = keeper.AddBalancesForIdRanges(userBalance.Balances, []*types.IdRange{
// 				{
// 					Start: start,
// 					End:   end,
// 				},
// 			}, amount)
// 			suite.Require().Nil(err, "error adding balance to approval")

// 			adds = append(adds, &types.IdRange{
// 				Start: start,
// 				End:   end,
// 			})
// 			// println("adding", start, end, amount)

// 			for j := start; j <= end; j++ {
// 				balances[j] += amount
// 			}

// 			amount = sdk.NewUint(rand.Intn(100))
// 			start = sdk.NewUint(rand.Intn(1000))
// 			end = sdk.NewUint(500 + rand.Intn(500))
// 			userBalance.Balances, err = keeper.SubtractBalancesForIdRanges(userBalance.Balances, []*types.IdRange{
// 				{
// 					Start: start,
// 					End:   end,
// 				},
// 			}, amount)

// 			if err != nil {
// 				suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// 			} else {
// 				// if sdk.NewUint(256) < start && sdk.NewUint(256) < end {
// 				// 	println("removing", start, end, amount)
// 				// }
// 				subs = append(subs, &types.IdRange{
// 					Start: start,
// 					End:   end,
// 				})
// 				for j := start; j <= end; j++ {
// 					balances[j] -= amount
// 				}
// 			}

// 		}

// 		for i := 0; i < 1000; i++ {
// 			suite.Require().Equal(keeper.GetBalancesForIdRanges(
// 				[]*types.IdRange{
// 					{
// 						Start: sdk.NewUint(i),
// 						End:   sdk.NewUint(i),
// 					},
// 				},
// 				userBalance.Balances,
// 			)[0].Amount, balances[i], "balance mismatch at index %d", i)
// 		}
// 	}
// }
