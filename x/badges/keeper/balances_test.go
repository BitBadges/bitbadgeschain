package keeper_test

import (
	"math"

	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
	balanceObjects := []*types.BalanceObject{
		{
			Balance: 1,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	}

	balanceObjects = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	}, 10, balanceObjects)

	gottenBalanceObjects := keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	}, balanceObjects)

	suite.Require().Equal(balanceObjects, []*types.BalanceObject{
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})
	suite.Require().Equal(balanceObjects, gottenBalanceObjects)

	gottenBalanceObjects = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
	}, balanceObjects)

	suite.Require().Equal(gottenBalanceObjects, []*types.BalanceObject{
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 1,
					End:   0,
				},
			},
		},
	})

	gottenBalanceObjects = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   2,
		},
	}, balanceObjects)

	suite.Require().Equal(gottenBalanceObjects, []*types.BalanceObject{
		{
			Balance: 0,
			IdRanges: []*types.IdRange{
				{
					Start: 2,
					End:   0,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 1,
					End:   0,
				},
			},
		},
	})

	gottenBalanceObjects = keeper.GetBalancesForIdRanges([]*types.IdRange{
		{
			Start: 0,
			End:   math.MaxUint64,
		},
	}, balanceObjects)

	suite.Require().Equal(gottenBalanceObjects, []*types.BalanceObject{
		{
			Balance: 0,
			IdRanges: []*types.IdRange{
				{
					Start: 2,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})

	gottenBalanceObjects = keeper.GetBalancesForIdRanges([]*types.IdRange{
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
		
	}, balanceObjects)

	suite.Require().Equal(gottenBalanceObjects, []*types.BalanceObject{
		{
			Balance: 0,
			IdRanges: []*types.IdRange{
				{
					Start: 2,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})

	balanceObjects = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 1,
			End:   1,
		},
	}, 5, balanceObjects)

	suite.Require().Equal(balanceObjects, []*types.BalanceObject{
		{
			Balance: 5,
			IdRanges: []*types.IdRange{
				{
					Start: 1,
					End:   0,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
			},
		},
	})
	
	balanceObjects = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 2,
			End:   math.MaxUint64,
		},
	}, 5, balanceObjects)

	suite.Require().Equal(balanceObjects, []*types.BalanceObject{
		{
			Balance: 5,
			IdRanges: []*types.IdRange{
				{
					Start: 1,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
			},
		},
	})

	balanceObjects = keeper.UpdateBalancesForIdRanges([]*types.IdRange{
		{
			Start: 2,
			End:   0,
		},
	}, 10, balanceObjects)

	suite.Require().Equal(balanceObjects, []*types.BalanceObject{
		{
			Balance: 5,
			IdRanges: []*types.IdRange{
				{
					Start: 1,
					End:   0,
				},
				{
					Start: 3,
					End:   math.MaxUint64,
				},
			},
		},
		{
			Balance: 10,
			IdRanges: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
				{
					Start: 2,
					End:   0,
				},
			},
		},
	})
}

func (suite *TestSuite) TestSubtractBalances() {
	userBalanceInfo := types.UserBalanceInfo{}
	subbadgeRanges := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
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
	userBalanceInfo, err = keeper.AddBalancesForIdRanges(userBalanceInfo, subbadgeRanges, 1000)
	suite.Require().NoError(err)

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{
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

	subbadgeRangesToRemove := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 35,
			End:   0,
		},
		{
			Start: 35,
			End:   100,
		},
	}
	for _, subbadgeRangeToRemove := range subbadgeRangesToRemove {
		userBalanceInfo, err = keeper.SubtractBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRangeToRemove}, 1)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(998))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{{Start: 35, End: 0}})

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].Balance, uint64(999))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].IdRanges, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})
}

func (suite *TestSuite) TestAddBalancesForIdRanges() {
	userBalanceInfo := types.UserBalanceInfo{}
	subbadgeRanges := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
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
	for _, subbadgeRange := range subbadgeRanges {
		userBalanceInfo, err = keeper.AddBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRange}, 1000)
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].Balance, uint64(2000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].IdRanges, []*types.IdRange{{Start: 35, End: 0}})

	subbadgeRangesToAdd := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 35,
			End:   0,
		},
		{
			Start: 35,
			End:   100,
		},
	}

	for _, subbadgeRangeToAdd := range subbadgeRangesToAdd {
		userBalanceInfo, _ = keeper.AddBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRangeToAdd}, 1)
		suite.Require().Nil(err, "error adding balance to approval")
	}
	
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 34}})

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].Balance, uint64(1001))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].IdRanges, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[2].Balance, uint64(2002))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[2].IdRanges, []*types.IdRange{{Start: 35, End: 0}})
}

func (suite *TestSuite) TestAddBalancesOverflow() {
	userBalanceInfo := types.UserBalanceInfo{}
	subbadgeRanges := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 0,
			End:   0,
		},
		{
			Start: 35,
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
	for _, subbadgeRange := range subbadgeRanges {
		userBalanceInfo, err = keeper.AddBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRange}, 1000)
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].Balance, uint64(2000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[1].IdRanges, []*types.IdRange{{Start: 35, End: 0}})

	subbadgeRangesToAdd := []*types.IdRange{
		{
			Start: 0,
			End:   1000,
		},
	}

	for _, subbadgeRange := range subbadgeRangesToAdd {
		userBalanceInfo, err = keeper.AddBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRange}, math.MaxUint64)
		suite.Require().EqualError(err, keeper.ErrOverflow.Error())
	}
}



func (suite *TestSuite) TestRemoveBalancesUnderflow() {
	userBalanceInfo := types.UserBalanceInfo{}
	subbadgeRanges := []types.IdRange{
		{
			Start: 1,
			End:   0,
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
	for _, subbadgeRange := range subbadgeRanges {
		userBalanceInfo, err = keeper.AddBalancesForIdRanges(userBalanceInfo, []*types.IdRange{&subbadgeRange}, 1000, )
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 100}})

	subbadgeRangesToRemove := []*types.IdRange{
		{
			Start: 1,
			End:   0,
		},
		{
			Start: 35,
			End:   0,
		},
		{
			Start: 35,
			End:   100,
		},
	}
	
	for _, subbadgeRange := range subbadgeRangesToRemove {
		userBalanceInfo, err = keeper.SubtractBalancesForIdRanges(userBalanceInfo, []*types.IdRange{subbadgeRange}, math.MaxUint64, )
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}