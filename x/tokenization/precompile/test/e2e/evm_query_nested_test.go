package tokenization_test

import (
	"encoding/hex"
	"math"
	"math/big"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	types "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func (s *EVMQueryChallengesE2ESuite) TestE2E_QueryTransferNestedInvariant() {
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypePrecompileTransfer)
	s.Require().NoError(err)
	contractABI, err := helpers.GetContractABIByType(helpers.ContractTypePrecompileTransfer)
	s.Require().NoError(err)
	contract, _, err := helpers.DeployContract(s.Ctx, s.App.EVMKeeper, s.DeployerKey, bytecode, s.ChainID)
	s.Require().NoError(err)
	contractAddress := sdk.AccAddress(contract.Bytes()).String()
	ranges := func() []*types.UintRange {
		return []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	}
	server := tokenizationkeeper.NewMsgServerImpl(s.TokenizationKeeper)
	created, err := server.UniversalUpdateCollection(s.Ctx, &types.MsgUniversalUpdateCollection{
		Creator: s.Alice.String(), CollectionId: sdkmath.ZeroUint(),
		UpdateValidTokenIds: true, ValidTokenIds: getValidIdRanges(),
		CollectionPermissions:     &types.CollectionPermissions{},
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{{
			ApprovalId: "transfer", FromListId: "AllWithoutMint", ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
			TransferTimes: ranges(), TokenIds: ranges(), OwnershipTimes: ranges(), Version: sdkmath.ZeroUint(),
			ApprovalCriteria: &types.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
		}, {
			ApprovalId: "mint", FromListId: "Mint", ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
			TransferTimes: ranges(), TokenIds: ranges(), OwnershipTimes: ranges(), Version: sdkmath.ZeroUint(),
			ApprovalCriteria: &types.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
		}},
	})
	s.Require().NoError(err)
	_, err = server.TransferTokens(s.Ctx, &types.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: created.CollectionId,
		Transfers: []*types.Transfer{{From: "Mint", ToAddresses: []string{contractAddress},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: getValidIdRanges(), OwnershipTimes: ranges()}},
		}},
	})
	s.Require().NoError(err)
	args := []interface{}{created.CollectionId.BigInt(), common.BytesToAddress(s.Bob), big.NewInt(1), big.NewInt(1), big.NewInt(1)}
	calldata, err := contractABI.Pack("testTransfer", args...)
	s.Require().NoError(err)
	collection, found := s.TokenizationKeeper.GetCollectionFromStore(s.Ctx, created.CollectionId)
	s.Require().True(found)
	before, _, err := s.TokenizationKeeper.GetBalanceOrApplyDefault(s.Ctx, collection, contractAddress)
	s.Require().NoError(err)
	recipientBefore, _, err := s.TokenizationKeeper.GetBalanceOrApplyDefault(s.Ctx, collection, s.Bob.String())
	s.Require().NoError(err)
	query := func() error {
		ctx := s.Ctx.WithGasMeter(storetypes.NewGasMeter(20_000_000))
		_, err := s.TokenizationKeeper.ExecuteEVMQuery(ctx, contract.Hex(), calldata, 10_000_000)
		s.Require().Positive(ctx.GasMeter().GasConsumed())
		return err
	}
	s.Require().NoError(query(), "the contract transfer must reach the precompile before adding the invariant")

	check, err := contractABI.Pack("testGetCollection", created.CollectionId.BigInt())
	s.Require().NoError(err)
	collection.Invariants = &types.CollectionInvariants{EvmQueryChallenges: []*types.EVMQueryChallenge{{
		ContractAddress: contract.Hex(), Calldata: hex.EncodeToString(check), GasLimit: sdkmath.NewUint(500_000),
	}}}
	s.Require().NoError(s.TokenizationKeeper.SetCollectionInStore(s.Ctx, collection, false))
	s.Require().Error(query(), "a query cannot start another query through a contract transfer invariant")
	after, _, err := s.TokenizationKeeper.GetBalanceOrApplyDefault(s.Ctx, collection, contractAddress)
	s.Require().NoError(err)
	s.Require().Equal(before.Balances, after.Balances, "query transfers must not persist balances")
	recipientAfter, _, err := s.TokenizationKeeper.GetBalanceOrApplyDefault(s.Ctx, collection, s.Bob.String())
	s.Require().NoError(err)
	s.Require().Equal(recipientBefore.Balances, recipientAfter.Balances)

	txCtx := s.Ctx.WithGasMeter(storetypes.NewGasMeter(20_000_000))
	_, response, err := helpers.CallContractMethod(txCtx, s.App.EVMKeeper, s.DeployerKey, contract, contractABI, "testTransfer", args, s.ChainID, false)
	s.Require().NoError(err, "the same invariant remains valid in an ordinary contract transaction")
	s.Require().Empty(response.VmError, "return data: %x", response.Ret)
	committed, _, err := s.TokenizationKeeper.GetBalanceOrApplyDefault(s.Ctx, collection, contractAddress)
	s.Require().NoError(err)
	s.Require().NotEqual(before.Balances, committed.Balances, "the ordinary transaction must commit its transfer")
}
