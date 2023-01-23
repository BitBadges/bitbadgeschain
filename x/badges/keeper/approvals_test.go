package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSetApprovals() {
	userBalance := types.UserBalance{}
	badgeIdRanges := []*types.IdRange{
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

	randomAcccountNum := uint64(30)
	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		userBalance, err = keeper.SetApproval(userBalance, 1000, aliceAccountNum, badgeIdRange)
		suite.Require().NoError(err)
		userBalance, err = keeper.SetApproval(userBalance, 1000, charlieAccountNum, badgeIdRange)
		suite.Require().NoError(err)
		userBalance, err = keeper.SetApproval(userBalance, 1000, randomAcccountNum, badgeIdRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 100}})

	suite.Require().Equal(userBalance.Approvals[1].Address, charlieAccountNum)
	suite.Require().Equal(userBalance.Approvals[1].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[1].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 100}})

	badgeIdRangesToRemove := []types.IdRange{
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
			Start: 35,
			End:   100,
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance, err = keeper.SetApproval(userBalance, 0, aliceAccountNum, &badgeIdRange)
		suite.Require().NoError(err)
		userBalance, err = keeper.SetApproval(userBalance, 0, charlieAccountNum, &badgeIdRange)
		suite.Require().NoError(err)
		userBalance, err = keeper.SetApproval(userBalance, 0, randomAcccountNum, &badgeIdRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 2, End: 34}})

	userBalance, err = keeper.SetApproval(userBalance, 0, aliceAccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	userBalance, err = keeper.SetApproval(userBalance, 0, charlieAccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	userBalance, err = keeper.SetApproval(userBalance, 0, randomAcccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	suite.Require().Equal(len(userBalance.Approvals), 0)
}

func (suite *TestSuite) TestRemoveApprovals() {
	userBalance := types.UserBalance{}
	badgeIdRanges := []*types.IdRange{
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
	for _, badgeIdRange := range badgeIdRanges {
		userBalance, err = keeper.SetApproval(userBalance, 1000, aliceAccountNum, badgeIdRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{
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

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance, err = keeper.RemoveBalanceFromApproval(userBalance, 1, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(998))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 35, End: 0}})

	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Balance, uint64(999))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})
}

func (suite *TestSuite) TestAddApprovals() {
	userBalance := types.UserBalance{}
	badgeIdRanges := []*types.IdRange{
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
	for _, badgeIdRange := range badgeIdRanges {
		userBalance, err = keeper.AddBalanceToApproval(userBalance, 1000, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Balance, uint64(2000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: 35, End: 0}})

	badgeIdRangesToRemove := []*types.IdRange{
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

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance, _ = keeper.AddBalanceToApproval(userBalance, 1, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 34}})

	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Balance, uint64(1001))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[2].Balance, uint64(2002))
	suite.Require().Equal(userBalance.Approvals[0].Balances[2].BadgeIds, []*types.IdRange{{Start: 35, End: 0}})
}

func (suite *TestSuite) TestAddApprovalsOverflow() {
	userBalance := types.UserBalance{}
	badgeIdRanges := []*types.IdRange{
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
	for _, badgeIdRange := range badgeIdRanges {
		userBalance, err = keeper.AddBalanceToApproval(userBalance, 1000, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Balance, uint64(2000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: 35, End: 0}})

	badgeIdRangesToAdd := []*types.IdRange{
		{
			Start: 0,
			End:   1000,
		},
	}

	for _, badgeIdRange := range badgeIdRangesToAdd {
		userBalance, err = keeper.AddBalanceToApproval(userBalance, math.MaxUint64, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "we should just set to uint64 max and not overflow")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(math.MaxUint64))
}

func (suite *TestSuite) TestRemoveApprovalsUnderflow() {
	userBalance := types.UserBalance{}
	badgeIdRanges := []*types.IdRange{
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
	for _, badgeIdRange := range badgeIdRanges {
		userBalance, err = keeper.SetApproval(userBalance, 1000, aliceAccountNum, badgeIdRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Balance, uint64(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: 0, End: 100}})

	badgeIdRangesToRemove := []*types.IdRange{
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

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance, err = keeper.RemoveBalanceFromApproval(userBalance, math.MaxUint64, aliceAccountNum, []*types.IdRange{badgeIdRange})
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}
