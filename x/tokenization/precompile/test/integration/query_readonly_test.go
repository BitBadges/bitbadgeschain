package tokenization_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// A view method of the precompile must not change state when it is called
// from inside a transaction (a plain call to the precompile, or a contract
// that uses CALL rather than STATICCALL).
type QueryReadOnlyTestSuite struct {
	EVMKeeperIntegrationTestSuite
}

func TestQueryReadOnlyTestSuite(t *testing.T) {
	suite.Run(t, new(QueryReadOnlyTestSuite))
}

// defaultIncomingApprovalVersion returns the version a first access to the
// address's default incoming approval would receive, without persisting it.
func (suite *QueryReadOnlyTestSuite) defaultIncomingApprovalVersion(collection *tokenizationtypes.TokenCollection, address string) sdkmath.Uint {
	cacheCtx, _ := suite.Ctx.CacheContext()
	balance, appliedDefault, err := suite.TokenizationKeeper.GetBalanceOrApplyDefault(cacheCtx, collection, address)
	suite.Require().NoError(err)
	suite.Require().True(appliedDefault, "address must still have no stored balance")
	suite.Require().Len(balance.IncomingApprovals, 1)
	return balance.IncomingApprovals[0].Version
}

func (suite *QueryReadOnlyTestSuite) TestViewMethodsInsideTxLeaveApprovalVersionsUntouched() {
	// A default incoming approval whose version gets initialised on the
	// first access to an address's balance.
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, suite.CollectionId)
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
	suite.Require().NoError(suite.TokenizationKeeper.SetCollectionInStore(suite.Ctx, collection, false))

	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)

	queries := []struct {
		method string
		json   func(address sdk.AccAddress) string
	}{
		{
			method: "getBalanceAmount",
			json: func(address sdk.AccAddress) string {
				return fmt.Sprintf(`{"collectionId":"%s","address":"%s","tokenId":"1","ownershipTime":"1"}`, suite.CollectionId, address)
			},
		},
		{
			method: "getBalance",
			json: func(address sdk.AccAddress) string {
				return fmt.Sprintf(`{"collectionId":"%s","address":"%s"}`, suite.CollectionId, address)
			},
		},
	}

	for _, q := range queries {
		suite.Run(q.method, func() {
			_, _, address := helpers.CreateEVMAccount()
			suite.Require().True(suite.defaultIncomingApprovalVersion(collection, address.String()).IsZero(),
				"precondition: the version has never been initialised")

			method := suite.Precompile.ABI.Methods[q.method]
			input, err := helpers.PackMethodWithJSON(&method, q.json(address))
			suite.Require().NoError(err)

			// A transaction, not an eth_call: whatever the query writes would be committed.
			tx, err := helpers.BuildEVMTransaction(suite.AliceKey, &precompileAddr, input, big.NewInt(0), 1_000_000, big.NewInt(0), suite.getNonce(suite.AliceEVM), chainID)
			suite.Require().NoError(err)
			response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)
			suite.Require().NoError(err)
			suite.Require().Empty(response.VmError, "%s must succeed inside a transaction", q.method)

			suite.True(suite.defaultIncomingApprovalVersion(collection, address.String()).IsZero(),
				"%s inside a transaction must not initialise the address's default approval version", q.method)
		})
	}
}
