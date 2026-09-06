package tokenization_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// The precompile's view methods must not write, whichever context they run
// in: here the plain SDK context the unit helpers provide.
type QueryReadOnlyTestSuite struct {
	suite.Suite
	TestSuite *helpers.TestSuite
}

func TestQueryReadOnlyTestSuite(t *testing.T) {
	suite.Run(t, new(QueryReadOnlyTestSuite))
}

func (suite *QueryReadOnlyTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
}

func (suite *QueryReadOnlyTestSuite) TestGetBalanceAmountLeavesApprovalVersionsUntouched() {
	ts := suite.TestSuite
	collectionId, err := ts.CreateTestCollection(ts.Alice.String())
	suite.Require().NoError(err)

	collection, found := ts.Keeper.GetCollectionFromStore(ts.Ctx, collectionId)
	suite.Require().True(found)
	full := []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	collection.DefaultBalances = &tokenizationtypes.UserBalanceStore{
		IncomingApprovals: []*tokenizationtypes.UserIncomingApproval{{
			ApprovalId:        "default-in",
			FromListId:        "All",
			InitiatedByListId: "All",
			TransferTimes:     full,
			TokenIds:          full,
			OwnershipTimes:    full,
			ApprovalCriteria:  &tokenizationtypes.IncomingApprovalCriteria{},
		}},
	}
	suite.Require().NoError(ts.Keeper.SetCollectionInStore(ts.Ctx, collection, false))

	versionOnFirstAccess := func() sdkmath.Uint {
		cacheCtx, _ := ts.Ctx.CacheContext()
		balance, appliedDefault, err := ts.Keeper.GetBalanceOrApplyDefault(cacheCtx, collection, ts.Bob.String())
		suite.Require().NoError(err)
		suite.Require().True(appliedDefault)
		return balance.IncomingApprovals[0].Version
	}
	suite.Require().True(versionOnFirstAccess().IsZero(), "precondition")

	for _, q := range []struct{ method, json string }{
		{"getBalanceAmount", fmt.Sprintf(`{"collectionId":"%s","address":"%s","tokenId":"1","ownershipTime":"1"}`, collectionId, ts.Bob)},
		{"getBalance", fmt.Sprintf(`{"collectionId":"%s","address":"%s"}`, collectionId, ts.Bob)},
	} {
		method := ts.Precompile.ABI.Methods[q.method]
		input, err := helpers.PackMethodWithJSON(&method, q.json)
		suite.Require().NoError(err)
		_, err = ts.CallPrecompile(ts.AliceEVM, input)
		suite.Require().NoError(err)
		suite.True(versionOnFirstAccess().IsZero(), "%s must not initialise the default approval version", q.method)
	}
}
