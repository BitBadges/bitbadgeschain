package keeper_test

import (
	"math"

	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSetApprovals() {
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

	randomAcccountNum := uint64(30)
	err := *new(error)
	for _, subbadgeRange := range subbadgeRanges {
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 1000, aliceAccountNum, subbadgeRange)
		suite.Require().NoError(err)
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 1000, charlieAccountNum, subbadgeRange)
		suite.Require().NoError(err)
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 1000, randomAcccountNum, subbadgeRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 100}})


	suite.Require().Equal(userBalanceInfo.Approvals[1].Address, charlieAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[1].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[1].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 100}})

	subbadgeRangesToRemove := []types.IdRange{
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

	for _, subbadgeRange := range subbadgeRangesToRemove {
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, aliceAccountNum, &subbadgeRange)
		suite.Require().NoError(err)
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, charlieAccountNum, &subbadgeRange)
		suite.Require().NoError(err)
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, randomAcccountNum, &subbadgeRange)
		suite.Require().NoError(err)
	}
	
	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 2, End: 34}})

	userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, aliceAccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, charlieAccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 0, randomAcccountNum, &types.IdRange{Start: 2, End: 34})
	suite.Require().NoError(err)
	suite.Require().Equal(len(userBalanceInfo.Approvals), 0)
}

func (suite *TestSuite) TestRemoveApprovals() {
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
	for _, subbadgeRange := range subbadgeRanges {
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 1000, aliceAccountNum, subbadgeRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{
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

	for _, subbadgeRange := range subbadgeRangesToRemove {
		userBalanceInfo, err = keeper.RemoveBalanceFromApproval(userBalanceInfo, 1, aliceAccountNum, []*types.IdRange{subbadgeRange})
		suite.Require().NoError(err)
	}
	
	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(998))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 35, End: 0}})


	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].Balance, uint64(999))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].IdRanges, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})
}

func (suite *TestSuite) TestAddApprovals() {
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
		userBalanceInfo, err = keeper.AddBalanceToApproval(userBalanceInfo, 1000, aliceAccountNum,  []*types.IdRange{subbadgeRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	
	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].Balance, uint64(2000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].IdRanges, []*types.IdRange{{Start: 35, End: 0}})

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
		userBalanceInfo, _ = keeper.AddBalanceToApproval(userBalanceInfo, 1, aliceAccountNum,  []*types.IdRange{subbadgeRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}
	
	
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 34}})

	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].Balance, uint64(1001))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].IdRanges, []*types.IdRange{{Start: 1, End: 0}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[2].Balance, uint64(2002))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[2].IdRanges, []*types.IdRange{{Start: 35, End: 0}})
}

func (suite *TestSuite) TestAddApprovalsOverflow() {
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
		userBalanceInfo, err = keeper.AddBalanceToApproval(userBalanceInfo, 1000, aliceAccountNum,  []*types.IdRange{subbadgeRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	
	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 34}, {Start: 36, End: 100}})

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].Balance, uint64(2000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[1].IdRanges, []*types.IdRange{{Start: 35, End: 0}})

	subbadgeRangesToAdd := []*types.IdRange{
		{
			Start: 0,
			End:   1000,
		},
	}

	for _, subbadgeRange := range subbadgeRangesToAdd {
		userBalanceInfo, err = keeper.AddBalanceToApproval(userBalanceInfo, math.MaxUint64, aliceAccountNum,  []*types.IdRange{subbadgeRange})
		suite.Require().Nil(err, "we should just set to uint64 max and not overflow")
	}

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(math.MaxUint64))
}



func (suite *TestSuite) TestRemoveApprovalsUnderflow() {
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
		userBalanceInfo, err = keeper.SetApproval(userBalanceInfo, 1000, aliceAccountNum, subbadgeRange)
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalanceInfo.Approvals[0].Address, aliceAccountNum)
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance, uint64(1000))
	suite.Require().Equal(userBalanceInfo.Approvals[0].ApprovalAmounts[0].IdRanges, []*types.IdRange{{Start: 0, End: 100}})

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
		userBalanceInfo, err = keeper.RemoveBalanceFromApproval(userBalanceInfo, math.MaxUint64, aliceAccountNum,  []*types.IdRange{subbadgeRange})
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}