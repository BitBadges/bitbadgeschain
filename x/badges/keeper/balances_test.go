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

func (suite *TestSuite) TestUpdateBalancesForIds() {
	// ranges []*types.IdRange, newAmount uint64, balanceObjects []*types.BalanceObject
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