package tokenization_test

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *ERC3643TokenizationTestSuite) TestTemplateInitializedSupplyCanTransfer() {
	contractABI, err := helpers.GetContractABIByType(helpers.ContractTypeERC3643Template)
	suite.Require().NoError(err)
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypeERC3643Template)
	suite.Require().NoError(err)
	args, err := contractABI.Pack("", "Transfer", "TX")
	suite.Require().NoError(err)
	address, result, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.AliceKey, append(bytecode, args...), suite.getChainID())
	suite.Require().NoError(err)
	suite.Require().Empty(result.VmError)
	call := func(method string, view bool, args ...interface{}) []interface{} {
		suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		ret, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, address, contractABI, method, args, suite.getChainID(), view)
		suite.Require().NoError(err)
		reason, _ := abi.UnpackRevert(res.Ret)
		suite.Require().Empty(res.VmError, "%s reverted: %s", method, reason)
		values, err := contractABI.Methods[method].Outputs.Unpack(ret)
		suite.Require().NoError(err)
		return values
	}
	call("initializeCollection", false, big.NewInt(100))
	collectionId := call("collectionId", true)[0].(*big.Int)
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, sdkmath.NewUintFromBigInt(collectionId))
	suite.Require().True(found)
	suite.Require().Len(collection.CollectionApprovals, 1)
	suite.Require().Equal("!Mint", collection.CollectionApprovals[0].FromListId)
	suite.Require().Equal(sdk.AccAddress(address.Bytes()).String(), collection.CollectionApprovals[0].InitiatedByListId)

	reject := func(method string, args ...interface{}) {
		suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		_, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, address, contractABI, method, args, suite.getChainID(), false)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(res.VmError, "%s must reject", method)
	}
	call("registerIdentity", false, suite.AliceEVM, true)
	reject("transfer", suite.BobEVM, big.NewInt(1))
	call("registerIdentity", false, suite.BobEVM, true)
	call("freezeAddress", false, suite.BobEVM)
	reject("transfer", suite.BobEVM, big.NewInt(1))
	call("unfreezeAddress", false, suite.BobEVM)
	suite.Require().Equal("100", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Require().Equal(true, call("canTransfer", true, suite.AliceEVM, suite.BobEVM, big.NewInt(1))[0])
	call("transfer", false, suite.BobEVM, big.NewInt(1))
	suite.Require().Equal("99", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Require().Equal("1", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())
	suite.Require().Equal("0", call("balanceOf", true, address)[0].(*big.Int).String())
	suite.Require().Equal("100", call("totalSupply", true)[0].(*big.Int).String())

	// Direct precompile calls cannot bypass the wrapper even for a verified holder.
	for _, from := range []string{suite.Alice.String(), suite.Bob.String(), "Mint"} {
		payload := fmt.Sprintf(`{"collectionId":"%s","transfers":[{"from":"%s","toAddresses":["%s"],"balances":[{"amount":"1","tokenIds":[{"start":"1","end":"1"}],"ownershipTimes":[{"start":"1","end":"18446744073709551615"}]}]}]}`, collectionId, from, suite.Bob.String())
		suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		_, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, common.HexToAddress("0x1001"), suite.Precompile.ABI, "transferTokens", []interface{}{payload}, suite.getChainID(), false)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(res.VmError, "direct precompile from %s must reject", from)
	}
	reject("initializeCollection", big.NewInt(100))
	suite.Require().Equal("99", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Require().Equal("1", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())

	// The caller cannot select another holder as the debit source.
	suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	_, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.BobKey, address, contractABI, "transfer", []interface{}{suite.AliceEVM, big.NewInt(2)}, suite.getChainID(), false)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(res.VmError)
	suite.Require().Equal("99", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	_, res, err = helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.BobKey, address, contractABI, "transfer", []interface{}{suite.AliceEVM, big.NewInt(1)}, suite.getChainID(), false)
	suite.Require().NoError(err)
	suite.Require().Empty(res.VmError)
	suite.Require().Equal("100", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
	suite.Require().Equal("0", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())
}
