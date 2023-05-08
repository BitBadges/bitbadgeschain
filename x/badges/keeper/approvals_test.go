package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSetApprovals() {
	userBalance := types.UserBalanceStore{}
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End: sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	
	err := *new(error)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(1000), alice, badgeIdRanges)
	suite.Require().NoError(err)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(1000), charlie, badgeIdRanges)
	suite.Require().NoError(err)
	
	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(100)}})

	suite.Require().Equal(userBalance.Approvals[1].Address, charlie)
	suite.Require().Equal(userBalance.Approvals[1].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[1].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(100)}})

	badgeIdRangesToRemove := []*types.IdRange{
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(0), alice, badgeIdRangesToRemove)
	suite.Require().NoError(err)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(0), charlie, badgeIdRangesToRemove)
	suite.Require().NoError(err)

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(34)}})

	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(0), alice, []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(34)}})
	suite.Require().NoError(err)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(0), charlie, []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(34)}})
	suite.Require().NoError(err)
	suite.Require().Equal(len(userBalance.Approvals), 0)
}

func (suite *TestSuite) TestRemoveApprovals() {
	userBalance := types.UserBalanceStore{}
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End: sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
		{
			Start: sdk.NewUint(135),
			End: sdk.NewUint(200),
		},
		{
			Start: sdk.NewUint(235),
			End: sdk.NewUint(300),
		},
		{
			Start: sdk.NewUint(335),
			End: sdk.NewUint(400),
		},
	}

	err := *new(error)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(1000), alice, badgeIdRanges)
	suite.Require().NoError(err)

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: sdk.NewUint(0), End: sdk.NewUint(100),
		},
		{
			Start: sdk.NewUint(135),
			End: sdk.NewUint(200),
		},
		{
			Start: sdk.NewUint(235),
			End: sdk.NewUint(300),
		},
		{
			Start: sdk.NewUint(335),
			End: sdk.NewUint(400),
		},
	})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance.Approvals, err = keeper.RemoveBalanceFromApproval(userBalance.Approvals, sdk.NewUint(1), alice, []*types.IdRange{badgeIdRange})
		suite.Require().NoError(err)
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(998))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Amount, sdk.NewUint(999))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})
}

func (suite *TestSuite) TestAddApprovals() {
	userBalance := types.UserBalanceStore{}
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End: sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		userBalance.Approvals, err = keeper.AddBalanceToApproval(userBalance.Approvals, sdk.NewUint(1000), alice, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(34)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Amount, sdk.NewUint(2000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance.Approvals, _ = keeper.AddBalanceToApproval(userBalance.Approvals, sdk.NewUint(1), alice, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().True(types.IdRangeEquals(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(0)}, {Start: sdk.NewUint(2), End: sdk.NewUint(34)}}))
	
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Amount, sdk.NewUint(1001))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[2].Amount, sdk.NewUint(2002))
	suite.Require().Equal(userBalance.Approvals[0].Balances[2].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})
}

func (suite *TestSuite) TestAddApprovalsOverflow() {
	userBalance := types.UserBalanceStore{}
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End: sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	err := *new(error)
	for _, badgeIdRange := range badgeIdRanges {
		userBalance.Approvals, err = keeper.AddBalanceToApproval(userBalance.Approvals, sdk.NewUint(1000), alice, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "error adding balance to approval")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(34)}, {Start: sdk.NewUint(36), End: sdk.NewUint(100)}})

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].Amount, sdk.NewUint(2000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[1].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(35), End: sdk.NewUint(35)}})

	badgeIdRangesToAdd := []*types.IdRange{
		{
			Start: sdk.NewUint(0),
			End: sdk.NewUint(1000),
		},
	}

	for _, badgeIdRange := range badgeIdRangesToAdd {
		userBalance.Approvals, err = keeper.AddBalanceToApproval(userBalance.Approvals, sdk.NewUint(math.MaxUint64), alice, []*types.IdRange{badgeIdRange})
		suite.Require().Nil(err, "we should just set to uint64 max and not overflow")
	}

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(math.MaxUint64))
}

func (suite *TestSuite) TestRemoveApprovalsUnderflow() {
	userBalance := types.UserBalanceStore{}
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
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(2),
			End: sdk.NewUint(34),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	err := *new(error)
	userBalance.Approvals, err = keeper.SetApproval(userBalance.Approvals, sdk.NewUint(1000), alice, badgeIdRanges)
	suite.Require().NoError(err)

	suite.Require().Equal(userBalance.Approvals[0].Address, alice)
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].Amount, sdk.NewUint(1000))
	suite.Require().Equal(userBalance.Approvals[0].Balances[0].BadgeIds, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(100)}})

	badgeIdRangesToRemove := []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(1),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(35),
		},
		{
			Start: sdk.NewUint(35),
			End: sdk.NewUint(100),
		},
	}

	for _, badgeIdRange := range badgeIdRangesToRemove {
		userBalance.Approvals, err = keeper.RemoveBalanceFromApproval(userBalance.Approvals, sdk.NewUint(math.MaxUint64), alice, []*types.IdRange{badgeIdRange})
		suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
	}
}
