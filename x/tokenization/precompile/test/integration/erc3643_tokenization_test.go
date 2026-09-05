package tokenization_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ERC3643TokenizationTestSuite deploys the shipped contracts/ERC3643Tokenization.sol
// (compiled artifact in contracts/test) and drives its transfer/balanceOf/
// totalSupply through the EVM, so the sample's precompile calls are checked
// against the real ABI rather than only compiled.
type ERC3643TokenizationTestSuite struct {
	EVMKeeperIntegrationTestSuite
}

func TestERC3643TokenizationTestSuite(t *testing.T) {
	suite.Run(t, new(ERC3643TokenizationTestSuite))
}

func (suite *ERC3643TokenizationTestSuite) TestSampleContractTransfersThroughPrecompile() {
	contractABI, err := helpers.GetContractABIByType(helpers.ContractTypeERC3643Tokenization)
	suite.Require().NoError(err)
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypeERC3643Tokenization)
	suite.Require().NoError(err)

	chainID := suite.getChainID()

	// Constructor: (uint256 collectionId)
	ctorArgs, err := contractABI.Pack("", suite.CollectionId.BigInt())
	suite.Require().NoError(err)
	contractAddr, deployRes, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.AliceKey, append(bytecode, ctorArgs...), chainID)
	suite.Require().NoError(err)
	suite.Require().Empty(deployRes.VmError, "deployment must succeed")

	// The contract transfers as itself, so give it the whole supply of token
	// ID 1. The suite seeds balances directly rather than minting, so the
	// "Total" store that getTotalSupply reads is seeded the same way.
	const contractSupply = uint64(100)
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, suite.CollectionId)
	suite.Require().True(found)
	full := []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	contractAcc := sdk.AccAddress(contractAddr.Bytes())
	for _, holder := range []string{contractAcc.String(), "Total"} {
		suite.Require().NoError(suite.TokenizationKeeper.SetBalanceForAddress(suite.Ctx, collection, holder, &tokenizationtypes.UserBalanceStore{
			Balances: []*tokenizationtypes.Balance{{
				Amount:         sdkmath.NewUint(contractSupply),
				TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: full,
			}},
			AutoApproveSelfInitiatedOutgoingTransfers: true,
			AutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveAllIncomingTransfers:           true,
		}))
	}

	// The suite runs every transaction on one context whose gas meter is
	// never reset, and a precompile call starts by charging that meter's
	// running total against the gas it was forwarded. Give each call a
	// fresh meter, as the ante handler does for a real transaction.
	call := func(method string, args []interface{}, isView bool) ([]byte, *evmtypes.MsgEthereumTxResponse, error) {
		suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		ret, res, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, contractAddr, contractABI, method, args, chainID, isView)
		if err == nil && res != nil && res.VmError != "" {
			reason, _ := abi.UnpackRevert(res.Ret)
			suite.T().Logf("%s reverted: %s (reason %q, gas used %d)", method, res.VmError, reason, res.GasUsed)
		}
		return ret, res, err
	}
	uintResult := func(method string, args []interface{}) *big.Int {
		ret, res, err := call(method, args, true)
		suite.Require().NoError(err)
		suite.Require().Empty(res.VmError, "%s must not revert", method)
		out, err := contractABI.Methods[method].Outputs.Unpack(ret)
		suite.Require().NoError(err)
		return out[0].(*big.Int)
	}
	balanceOf := func(account common.Address) *big.Int {
		return uintResult("balanceOf", []interface{}{account})
	}
	totalSupply := func() *big.Int {
		return uintResult("totalSupply", nil)
	}

	suite.Require().Equal(big.NewInt(int64(contractSupply)).String(), balanceOf(contractAddr).String(), "contract balance before transfer")
	bobBefore := balanceOf(suite.BobEVM)
	supplyBefore := totalSupply()
	suite.Require().Equal(big.NewInt(int64(contractSupply)).String(), supplyBefore.String(), "total supply of token ID 1 must be visible through the sample")

	const amount = int64(7)
	ret, res, err := call("transfer", []interface{}{suite.BobEVM, big.NewInt(amount)}, false)
	suite.Require().NoError(err)
	suite.Require().Empty(res.VmError, "transfer must not revert")
	out, err := contractABI.Methods["transfer"].Outputs.Unpack(ret)
	suite.Require().NoError(err)
	suite.Require().True(out[0].(bool))

	transferEvent, ok := contractABI.Events["Transfer"]
	suite.Require().True(ok)
	sawTransfer := false
	for _, log := range res.Logs {
		if len(log.Topics) > 0 && log.Topics[0] == transferEvent.ID.Hex() {
			sawTransfer = true
		}
	}
	suite.True(sawTransfer, "the sample must emit Transfer")

	suite.Equal(big.NewInt(int64(contractSupply)-amount).String(), balanceOf(contractAddr).String(), "contract balance after transfer")
	suite.Equal(new(big.Int).Add(bobBefore, big.NewInt(amount)).String(), balanceOf(suite.BobEVM).String(), "recipient balance after transfer")
	suite.Equal(supplyBefore.String(), totalSupply().String(), "a transfer must not change total supply")
}
