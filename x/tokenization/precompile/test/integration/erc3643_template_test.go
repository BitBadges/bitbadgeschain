package tokenization_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *ERC3643TokenizationTestSuite) TestTemplateInitializationAndRegistry() {
	suite.checkTemplateInitialization(big.NewInt(100))
}

func (suite *ERC3643TokenizationTestSuite) TestTemplateZeroSupply() {
	suite.checkTemplateInitialization(big.NewInt(0))
}

func (suite *ERC3643TokenizationTestSuite) checkTemplateInitialization(supply *big.Int) {
	contractABI, err := helpers.GetContractABIByType(helpers.ContractTypeERC3643Template)
	suite.Require().NoError(err)
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypeERC3643Template)
	suite.Require().NoError(err)
	args, err := contractABI.Pack("", "Example", "EX")
	suite.Require().NoError(err)
	address, result, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.AliceKey, append(bytecode, args...), suite.getChainID())
	suite.Require().NoError(err)
	suite.Require().Empty(result.VmError)
	call := func(method string, view bool, args ...interface{}) []interface{} {
		suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		ret, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, address, contractABI, method, args, suite.getChainID(), view)
		suite.Require().NoError(err)
		suite.Require().Empty(res.VmError, "%s: %x", method, res.Ret)
		values, err := contractABI.Methods[method].Outputs.Unpack(ret)
		suite.Require().NoError(err)
		return values
	}
	call("registerIdentity", false, suite.BobEVM, true)
	suite.Equal(true, call("isKYCVerified", true, suite.BobEVM)[0])
	suite.Equal(true, call("isAccredited", true, suite.BobEVM)[0])
	call("removeIdentity", false, suite.BobEVM)
	suite.Equal(false, call("isKYCVerified", true, suite.BobEVM)[0])
	suite.Equal(false, call("isAccredited", true, suite.BobEVM)[0])
	call("freezeAddress", false, suite.BobEVM)
	suite.Equal(true, call("isFrozen", true, suite.BobEVM)[0])
	call("initializeCollection", false, supply)
	suite.Equal(supply.String(), call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Equal("0", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())
	suite.Equal(supply.String(), call("totalSupply", true)[0].(*big.Int).String())
	collectionID := sdkmath.NewUintFromBigInt(call("collectionId", true)[0].(*big.Int))
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, collectionID)
	suite.Require().True(found)
	suite.Require().Len(collection.CollectionApprovals, 1)
	suite.Equal("wrapper-transfers", collection.CollectionApprovals[0].ApprovalId)
	suite.Equal("!Mint", collection.CollectionApprovals[0].FromListId, "initial issuance permission must be removed")
	suite.Equal(sdk.AccAddress(address.Bytes()).String(), collection.CollectionApprovals[0].InitiatedByListId)
	suite.Empty(collection.DefaultBalances.Balances)
	balance, _, err := suite.TokenizationKeeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice.String())
	suite.Require().NoError(err)
	suite.Empty(balance.IncomingApprovals)
	suite.False(balance.AutoApproveAllIncomingTransfers)
}
