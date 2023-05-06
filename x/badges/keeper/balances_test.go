package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSafeAdd() {
	result, err := keeper.SafeAdd(uint64(0), uint64(1))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, uint64(1))

	result, err = keeper.SafeAdd(uint64(math.MaxUint64), uint64(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, uint64(math.MaxUint64))

	_, err = keeper.SafeAdd(uint64(math.MaxUint64), uint64(1))
	suite.Require().EqualError(err, keeper.ErrOverflow.Error())
}

func (suite *TestSuite) TestSafeSubtract() {
	result, err := keeper.SafeSubtract(uint64(1), uint64(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, uint64(1))

	result, err = keeper.SafeSubtract(uint64(math.MaxUint64), uint64(0))
	suite.Require().Nil(err, "Error adding: %s")
	suite.Require().Equal(result, uint64(math.MaxUint64))

	_, err = keeper.SafeSubtract(uint64(1), uint64(math.MaxUint64))
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}

func (suite *TestSuite) TestUpdateAndGetBalancesForIds() {
	err := *new(error)
	balances := []*types.Balance{
		{
			Amount: 1,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	}

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	}, 10, balances)
	suite.Require().Nil(err, "Error updating balances: %s")

	fetchedBalances, err := keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	}, balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	suite.Require().Equal(balances, []*types.Balance{
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})
	suite.Require().Equal(balances, fetchedBalances)

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   1,
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   2,
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: 0,
			BadgeIds: []*types.IdRange{
				{
					Start: 2,
					End:   2,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   1,
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   math.MaxUint64,
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: 0,
			BadgeIds: []*types.IdRange{
				{
					Start: 2,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})

	fetchedBalances, err = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 3,
			End:   math.MaxUint64,
		},
		{
			Start: 0,
			End:   2,
		},
		{
			Start: 0,
			End:   1,
		},
	}, balances)

	suite.Require().Equal(fetchedBalances, []*types.Balance{
		{
			Amount: 0,
			BadgeIds: []*types.IdRange{
				{
					Start: 2,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
	}, 5, balances)

	suite.Require().Equal(balances, []*types.Balance{
		{
			Amount: 5,
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   1,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
			},
		},
	})

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 2,
			End:   math.MaxUint64,
		},
	}, 5, balances)

	suite.Require().Equal(balances, []*types.Balance{
		{
			Amount: 5,
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
			},
		},
	})

	balances, err = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 2,
			End:   2,
		},
	}, 10, balances)

	suite.Require().Equal(balances, []*types.Balance{
		{
			Amount: 5,
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   1,
				},
				{
					Start: 3,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Amount: 10,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
				{
					Start: 2,
					End:   2,
				},
			},
		},
	})
}

func (suite *TestSuite) TestSubtractBalances() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 2,
			End:   34,
		},
		{
			Start: 35,
			End:   100,
		},
		{
			Start: 135,
			End:   200,
		},
		{
			Start: 235,
			End:   300,
		},
		{
			Start: 335,
			End:   400,
		},
	}

	err := *new(error)
	UserBalance, err = keeper.AddBalancesForIdRanges(UserBalance, badgeIdRanges, 1000)
	suite.Require().NoError(err)

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: 0, End: 100,
		},
		{
			Start: 135,
			End:   200,
		},
		{
			Start: 235,
			End:   300,
		},
		{
			Start: 335,
			End:   400,
		},
	})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 35,
			End:   100,
		},
	}
	for _, badgeIdRangeToRemove := range badgeIdRangesToRemove {
		UserBalance, err = keeper.SubtractBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRangeToRemove}, 1)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(998))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: 35, End: 35}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, uint64(999))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: 1, End: 1}, {Start: 36, End: 100}})
}

func (suite *TestSuite) TestAddBalancesForIdRanges() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 2,
			End:   34,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		UserBalance, err = keeper.AddBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRange}, 1000)
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, uint64(2000))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: 35, End: 35}})

	badgeIdRangesToAdd := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	for _, badgeIdRangeToAdd := range badgeIdRangesToAdd {
		UserBalance, _ = keeper.AddBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRangeToAdd}, 1)
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 34}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, uint64(1001))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: 1, End: 1}, {Start: 36, End: 100}})

	suite.Require().Equal(UserBalance.Balances[2].Amount, uint64(2002))
	suite.Require().Equal(UserBalance.Balances[2].BadgeIds, []*types.IdRange{{Start: 35, End: 35}})
}

func (suite *TestSuite) TestAddBalancesOverflow() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 2,
			End:   34,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		UserBalance, err = keeper.AddBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRange}, 1000)
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(UserBalance.Balances[1].Amount, uint64(2000))
	suite.Require().Equal(UserBalance.Balances[1].BadgeIds, []*types.IdRange{{Start: 35, End: 35}})

	badgeIdRangesToAdd := []*types.IdRange{
		{
			Start: 0,
			End:   1000,
		},
	}

	for _, badgeIdRange := range badgeIdRangesToAdd {
		UserBalance, err = keeper.AddBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRange}, math.MaxUint64)
		suite.Require().EqualError(err, keeper.ErrOverflow.Error())
	}
}

func (suite *TestSuite) TestRemoveBalancesUnderflow() {
	UserBalance := types.UserBalanceStore{}
	badgeIdRanges := []types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 2,
			End:   34,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		UserBalance, err = keeper.AddBalancesForIdRanges(UserBalance, []*types.IdRange{&badgeIdRange}, 1000)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(UserBalance.Balances[0].Amount, uint64(1000))
	suite.Require().Equal(UserBalance.Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 100}})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 35,
			End:   35,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		UserBalance, err = keeper.SubtractBalancesForIdRanges(UserBalance, []*types.IdRange{badgeIdRange}, math.MaxUint64)
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}

//Commented because it takes +- 25 seconds to run

// func (suite *TestSuite) TestBalancesFuzz() {
// 	for i := 0; i < 25; i++ {
// 		userBalance := types.UserBalanceStore{}
// 		balances := make([]uint64, 1000)

// 		adds := make([]*types.IdRange, 100)
// 		subs := make([]*types.IdRange, 100)
// 		for i := 0; i < 100; i++ { //10000 iterations
// 			//Get random start value
// 			start := uint64(rand.Intn(500))
// 			//Get random end value
// 			end := uint64(500 + rand.Intn(500))

// 			amount := uint64(rand.Intn(100))
// 			err := *new(error)
// 			userBalance, err = keeper.AddBalancesForIdRanges(userBalance, []*types.IdRange{
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

// 			amount = uint64(rand.Intn(100))
// 			start = uint64(rand.Intn(1000))
// 			end = uint64(500 + rand.Intn(500))
// 			userBalance, err = keeper.SubtractBalancesForIdRanges(userBalance, []*types.IdRange{
// 				{
// 					Start: start,
// 					End:   end,
// 				},
// 			}, amount)

// 			if err != nil {
// 				suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// 			} else {
// 				// if uint64(256) < start && uint64(256) < end {
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
// 						Start: uint64(i),
// 						End:   uint64(i),
// 					},
// 				},
// 				userBalance.Balances,
// 			)[0].Amount, balances[i], "balance mismatch at index %d", i)
// 		}
// 	}
// }
