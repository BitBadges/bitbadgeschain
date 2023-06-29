package keeper_test

// import (
// sdkmath "cosmossdk.io/math"
// 	"math"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *TestSuite) TestSafeAdd() {
// 	result, err := keeper.SafeAdd(sdkmath.NewUint(0), sdkmath.NewUint(1))
// 	suite.Require().Nil(err, "Error adding: %s")
// 	suite.Require().Equal(result, sdkmath.NewUint(1))

// 	result, err = keeper.SafeAdd(sdkmath.NewUint(math.MaxUint64), sdkmath.NewUint(0))
// 	suite.Require().Nil(err, "Error adding: %s")
// 	suite.Require().Equal(result, sdkmath.NewUint(math.MaxUint64))

// 	_, err = keeper.SafeAdd(sdkmath.NewUint(math.MaxUint), sdkmath.NewUint(1))
// 	// suite.Require().EqualError(err, keeper.ErrOverflow.Error()) With Cosmos SDK Uint now, this error is not returned
// }

// func (suite *TestSuite) TestSafeSubtract() {
// 	result, err := keeper.SafeSubtract(sdkmath.NewUint(1), sdkmath.NewUint(0))
// 	suite.Require().Nil(err, "Error adding: %s")
// 	suite.Require().Equal(result, sdkmath.NewUint(1))

// 	result, err = keeper.SafeSubtract(sdkmath.NewUint(math.MaxUint64), sdkmath.NewUint(0))
// 	suite.Require().Nil(err, "Error adding: %s")
// 	suite.Require().Equal(result, sdkmath.NewUint(math.MaxUint64))

// 	_, err = keeper.SafeSubtract(sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))
// 	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// }

// func (suite *TestSuite) TestUpdateAndGetBalancesForIds() {
// 	err := *new(error)
// 	balances := []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(1),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	}

// 	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(1),
// 		},
// 	}, sdkmath.NewUint(10), balances)
// 	suite.Require().Nil(err, "Error updating balances: %s")

// 	fetchedBalances, err := keeper.GetBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(1),
// 		},
// 	}, balances)
// 	suite.Require().Nil(err, "Error fetching balances: %s")

// 	suite.Require().Equal(balances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Equal(balances, fetchedBalances)

// 	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 	}, balances)

// 	suite.Require().Equal(fetchedBalances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	})

// 	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(2),
// 		},
// 	}, balances)

// 	suite.Require().Equal(fetchedBalances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(0),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(2),
// 					End:   sdkmath.NewUint(2),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	})

// 	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(math.MaxUint64),
// 		},
// 	}, balances)

// 	suite.Require().Equal(fetchedBalances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(0),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(2),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	})

// 	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(3),
// 			End:   sdkmath.NewUint(math.MaxUint64),
// 		},
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(2),
// 		},
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(1),
// 		},
// 	}, balances)

// 	suite.Require().Equal(fetchedBalances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(0),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(2),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	})

// 	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 	}, sdkmath.NewUint(5), balances)

// 	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(5),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(0),
// 				},
// 			},
// 		},
// 	}))

// 	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(2),
// 			End:   sdkmath.NewUint(math.MaxUint64),
// 		},
// 	}, sdkmath.NewUint(5), balances)

// 	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(5),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(0),
// 				},
// 			},
// 		},
// 	}))

// 	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(2),
// 			End:   sdkmath.NewUint(2),
// 		},
// 	}, sdkmath.NewUint(10), balances)

// 	suite.Require().True(types.BalancesEqual(balances, []*types.Balance{
// 		{
// 			Amount: sdkmath.NewUint(5),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(1),
// 				},
// 				{
// 					Start: sdkmath.NewUint(3),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 		{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdkmath.NewUint(0),
// 					End:   sdkmath.NewUint(0),
// 				},
// 				{
// 					Start: sdkmath.NewUint(2),
// 					End:   sdkmath.NewUint(2),
// 				},
// 			},
// 		},
// 	}))
// }

// func (suite *TestSuite) TestSubtractBalances() {
// 	UserBalance := types.UserBalanceStore{}
// 	badgeIdRanges := []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(0),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(35),
// 		},
// 		{
// 			Start: sdkmath.NewUint(2),
// 			End:   sdkmath.NewUint(34),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 		{
// 			Start: sdkmath.NewUint(135),
// 			End:   sdkmath.NewUint(200),
// 		},
// 		{
// 			Start: sdkmath.NewUint(235),
// 			End:   sdkmath.NewUint(300),
// 		},
// 		{
// 			Start: sdkmath.NewUint(335),
// 			End:   sdkmath.NewUint(400),
// 		},
// 	}

// 	err := *new(error)
// 	UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, badgeIdRanges, sdkmath.NewUint(1000))
// 	suite.Require().NoError(err)

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(1000))
// 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(0), End: sdkmath.NewUint(100),
// 		},
// 		{
// 			Start: sdkmath.NewUint(135),
// 			End:   sdkmath.NewUint(200),
// 		},
// 		{
// 			Start: sdkmath.NewUint(235),
// 			End:   sdkmath.NewUint(300),
// 		},
// 		{
// 			Start: sdkmath.NewUint(335),
// 			End:   sdkmath.NewUint(400),
// 		},
// 	})

// 	badgeIdRangesToRemove := []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(35),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 	}
// 	for _, badgeIdRangeToRemove := range badgeIdRangesToRemove {
// 		UserBalance.Balances, err = keeper.SubtractBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRangeToRemove}, sdkmath.NewUint(1))
// 		suite.Require().NoError(err)
// 	}

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(998))
// 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(35), End: sdkmath.NewUint(35)}})

// 	suite.Require().Equal(UserBalance.Balances[1].Amount, sdkmath.NewUint(999))
// 	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}, {Start: sdkmath.NewUint(36), End: sdkmath.NewUint(100)}})
// }

// func (suite *TestSuite) TestAddBalancesForIdRanges() {
// 	UserBalance := types.UserBalanceStore{}
// 	badgeIdRanges := []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(0),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(35),
// 		},
// 		{
// 			Start: sdkmath.NewUint(2),
// 			End:   sdkmath.NewUint(34),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 	}

// 	err := *new(error)
// 	for _, badgeIdRange := range badgeIdRanges {
// 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdkmath.NewUint(1000))
// 		suite.Require().Nil(err, "error adding balance to approval")
// 	}

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(1000))
// 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(34)}, {Start: sdkmath.NewUint(36), End: sdkmath.NewUint(100)}})

// 	suite.Require().Equal(UserBalance.Balances[1].Amount, sdkmath.NewUint(2000))
// 	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(35), End: sdkmath.NewUint(35)}})

// 	badgeIdRangesToAdd := []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(35),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 	}

// 	for _, badgeIdRangeToAdd := range badgeIdRangesToAdd {
// 		UserBalance.Balances, _ = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRangeToAdd}, sdkmath.NewUint(1))
// 		suite.Require().Nil(err, "error adding balance to approval")
// 	}

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(1000))
// 	suite.Require().True(types.IdRangeEquals(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(0)}, {Start: sdkmath.NewUint(2), End: sdkmath.NewUint(34)}}))

// 	suite.Require().Equal(UserBalance.Balances[1].Amount, sdkmath.NewUint(1001))
// 	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}, {Start: sdkmath.NewUint(36), End: sdkmath.NewUint(100)}})

// 	suite.Require().Equal(UserBalance.Balances[2].Amount, sdkmath.NewUint(2002))
// 	suite.Require().Equal(UserBalance.Balances[2].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(35), End: sdkmath.NewUint(35)}})
// }

// // func (suite *TestSuite) TestAddBalancesOverflow() {
// // 	UserBalance := types.UserBalanceStore{}
// // 	badgeIdRanges := []*types.IdRange{
// // 		{
// // 			Start: sdkmath.NewUint(1),
// // 			End:   sdkmath.NewUint(1),
// // 		},
// // 		{
// // 			Start: sdkmath.NewUint(0),
// // 			End:   sdkmath.NewUint(0),
// // 		},
// // 		{
// // 			Start: sdkmath.NewUint(35),
// // 			End: sdkmath.NewUint(35),
// // 		},
// // 		{
// // 			Start: sdkmath.NewUint(2),
// // 			End: sdkmath.NewUint(34),
// // 		},
// // 		{
// // 			Start: sdkmath.NewUint(35),
// // 			End: sdkmath.NewUint(100),
// // 		},
// // 	}

// // 	err := *new(error)
// // 	for _, badgeIdRange := range badgeIdRanges {
// // 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdkmath.NewUint(1000))
// // 		suite.Require().Nil(err, "error adding balance to approval")
// // 	}

// // 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(1000))
// // 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(34)}, {Start: sdkmath.NewUint(36), End: sdkmath.NewUint(100)}})

// // 	suite.Require().Equal(UserBalance.Balances[1].Amount, sdkmath.NewUint(2000))
// // 	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(35), End: sdkmath.NewUint(35)}})

// // 	badgeIdRangesToAdd := []*types.IdRange{
// // 		{
// // 			Start: sdkmath.NewUint(0),
// // 			End: sdkmath.NewUint(1000),
// // 		},
// // 	}

// // 	for _, badgeIdRange := range badgeIdRangesToAdd {
// // 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdkmath.NewUint(math.MaxUint64))
// // 		suite.Require().EqualError(err, keeper.ErrOverflow.Error())
// // 	}
// // }

// func (suite *TestSuite) TestRemoveBalancesUnderflow() {
// 	UserBalance := types.UserBalanceStore{}
// 	badgeIdRanges := []types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(0),
// 			End:   sdkmath.NewUint(0),
// 		},
// 		{
// 			Start: sdkmath.NewUint(2),
// 			End:   sdkmath.NewUint(34),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 	}

// 	err := *new(error)
// 	for _, badgeIdRange := range badgeIdRanges {
// 		UserBalance.Balances, err = keeper.AddBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{&badgeIdRange}, sdkmath.NewUint(1000))
// 		suite.Require().NoError(err)
// 	}

// 	suite.Require().Equal(UserBalance.Balances[0].Amount, sdkmath.NewUint(1000))
// 	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(100)}})

// 	badgeIdRangesToRemove := []*types.IdRange{
// 		{
// 			Start: sdkmath.NewUint(1),
// 			End:   sdkmath.NewUint(1),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(35),
// 		},
// 		{
// 			Start: sdkmath.NewUint(35),
// 			End:   sdkmath.NewUint(100),
// 		},
// 	}

// 	for _, badgeIdRange := range badgeIdRangesToRemove {
// 		UserBalance.Balances, err = keeper.SubtractBalancesForIdRanges(UserBalance.Balances, []*types.IdRange{badgeIdRange}, sdkmath.NewUint(math.MaxUint64))
// 		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// 	}
// }

// //Commented because it takes +- 25 seconds to run

// // func (suite *TestSuite) TestBalancesFuzz() {
// // 	for i := 0; i < 25; i++ {
// // 		userBalance := types.UserBalanceStore{}
// // 		balances := make([]uint64, sdkmath.NewUint(1000))

// // 		adds := make([]*types.IdRange, 100)
// // 		subs := make([]*types.IdRange, 100)
// // 		for i := 0; i < 100; i++ { //10000 iterations
// // 			//Get random start value
// // 			start := sdkmath.NewUint(rand.Intn(500))
// // 			//Get random end value
// // 			end := sdkmath.NewUint(500 + rand.Intn(500))

// // 			amount := sdkmath.NewUint(rand.Intn(100))
// // 			err := *new(error)
// // 			userBalance.Balances, err = keeper.AddBalancesForIdRanges(userBalance.Balances, []*types.IdRange{
// // 				{
// // 					Start: start,
// // 					End:   end,
// // 				},
// // 			}, amount)
// // 			suite.Require().Nil(err, "error adding balance to approval")

// // 			adds = append(adds, &types.IdRange{
// // 				Start: start,
// // 				End:   end,
// // 			})
// // 			// println("adding", start, end, amount)

// // 			for j := start; j <= end; j++ {
// // 				balances[j] += amount
// // 			}

// // 			amount = sdkmath.NewUint(rand.Intn(100))
// // 			start = sdkmath.NewUint(rand.Intn(1000))
// // 			end = sdkmath.NewUint(500 + rand.Intn(500))
// // 			userBalance.Balances, err = keeper.SubtractBalancesForIdRanges(userBalance.Balances, []*types.IdRange{
// // 				{
// // 					Start: start,
// // 					End:   end,
// // 				},
// // 			}, amount)

// // 			if err != nil {
// // 				suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// // 			} else {
// // 				// if sdkmath.NewUint(256) < start && sdkmath.NewUint(256) < end {
// // 				// 	println("removing", start, end, amount)
// // 				// }
// // 				subs = append(subs, &types.IdRange{
// // 					Start: start,
// // 					End:   end,
// // 				})
// // 				for j := start; j <= end; j++ {
// // 					balances[j] -= amount
// // 				}
// // 			}

// // 		}

// // 		for i := 0; i < 1000; i++ {
// // 			suite.Require().Equal(keeper.GetBalancesForIdRanges(
// // 				[]*types.IdRange{
// // 					{
// // 						Start: sdkmath.NewUint(i),
// // 						End:   sdkmath.NewUint(i),
// // 					},
// // 				},
// // 				userBalance.Balances,
// // 			)[0].Amount, balances[i], "balance mismatch at index %d", i)
// // 		}
// // 	}
// // }
