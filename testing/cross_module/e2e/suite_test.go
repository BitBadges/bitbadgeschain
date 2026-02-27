//go:build test
// +build test

package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// CrossModuleTestSuite is the base suite for cross-module E2E tests
type CrossModuleTestSuite struct {
	apptesting.KeeperTestHelper

	gammMsgServer         gammtypes.MsgServer
	tokenizationMsgServer tokenizationtypes.MsgServer
}

func TestCrossModuleTestSuite(t *testing.T) {
	suite.Run(t, new(CrossModuleTestSuite))
}

func (s *CrossModuleTestSuite) SetupTest() {
	s.Setup()
	s.gammMsgServer = gammkeeper.NewMsgServerImpl(&s.App.GammKeeper)
	s.tokenizationMsgServer = tokenizationkeeper.NewMsgServerImpl(s.App.TokenizationKeeper)
}

// Helper to create a pool with specified coins
func (s *CrossModuleTestSuite) CreatePoolWithCoins(coins ...sdk.Coin) uint64 {
	var poolAssets []balancer.PoolAsset
	for _, coin := range coins {
		poolAsset := balancer.PoolAsset{
			Weight: sdkmath.NewInt(1),
			Token:  coin,
		}
		poolAssets = append(poolAssets, poolAsset)
	}

	// Fund the test account
	fundCoins := sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(10000000000)))
	for _, coin := range coins {
		fundCoins = fundCoins.Add(coin)
	}
	s.FundAcc(s.TestAccs[0], fundCoins)

	msg := balancer.NewMsgCreateBalancerPool(s.TestAccs[0], balancer.PoolParams{
		SwapFee: sdkmath.LegacyZeroDec(),
		ExitFee: sdkmath.LegacyZeroDec(),
	}, poolAssets)

	poolId, err := s.App.PoolManagerKeeper.CreatePool(s.Ctx, msg)
	s.Require().NoError(err)
	return poolId
}

// Helper to perform a swap
func (s *CrossModuleTestSuite) SwapTokens(sender sdk.AccAddress, poolId uint64, tokenIn sdk.Coin, tokenOutDenom string) (sdk.Coin, error) {
	s.FundAcc(sender, sdk.NewCoins(tokenIn))

	msg := &gammtypes.MsgSwapExactAmountIn{
		Sender: sender.String(),
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{PoolId: poolId, TokenOutDenom: tokenOutDenom},
		},
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdkmath.NewInt(1),
	}

	res, err := s.gammMsgServer.SwapExactAmountIn(s.Ctx, msg)
	if err != nil {
		return sdk.Coin{}, err
	}

	return sdk.NewCoin(tokenOutDenom, res.TokenOutAmount), nil
}

// Helper to create a basic tokenization collection using MsgCreateCollection
func (s *CrossModuleTestSuite) CreateBasicCollection(creator sdk.AccAddress) (sdkmath.Uint, error) {
	msg := &tokenizationtypes.MsgCreateCollection{
		Creator: creator.String(),
		DefaultBalances: &tokenizationtypes.UserBalanceStore{
			Balances:                                   []*tokenizationtypes.Balance{},
			AutoApproveSelfInitiatedOutgoingTransfers:  true,
			AutoApproveSelfInitiatedIncomingTransfers:  true,
			AutoApproveAllIncomingTransfers:            true,
		},
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		},
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
		Manager:               creator.String(),
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "ipfs://test-collection",
			CustomData: "",
		},
		TokenMetadata:       []*tokenizationtypes.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}

	res, err := s.tokenizationMsgServer.CreateCollection(s.Ctx, msg)
	if err != nil {
		return sdkmath.ZeroUint(), err
	}

	// Set up mint approval for minting tokens
	s.setupMintApproval(res.CollectionId, creator)

	return res.CollectionId, nil
}

// setupMintApproval sets up a collection approval that allows minting from Mint address
func (s *CrossModuleTestSuite) setupMintApproval(collectionId sdkmath.Uint, creator sdk.AccAddress) {
	mintApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "mint-approval",
		FromListId:        "Mint",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
		TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
		OwnershipTimes:    []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}

	transferApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "transfer-approval",
		FromListId:        "!Mint",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
		TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
		OwnershipTimes:    []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
	}

	updateMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   creator.String(),
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{mintApproval, transferApproval},
	}
	_, err := s.tokenizationMsgServer.UniversalUpdateCollection(s.Ctx, updateMsg)
	s.Require().NoError(err, "failed to set up mint approval")
}

// Helper to mint tokens to an address (via transfer from Mint address)
func (s *CrossModuleTestSuite) MintTokensToAddress(collectionId sdkmath.Uint, recipient sdk.AccAddress, tokenIds []*tokenizationtypes.UintRange, amount sdkmath.Uint) error {
	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      recipient.String(),
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{recipient.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount:         amount,
						TokenIds:       tokenIds,
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
					},
				},
			},
		},
	}

	_, err := s.tokenizationMsgServer.TransferTokens(s.Ctx, msg)
	return err
}

// Helper to transfer tokens between addresses
func (s *CrossModuleTestSuite) TransferTokens(collectionId sdkmath.Uint, from, to sdk.AccAddress, tokenIds []*tokenizationtypes.UintRange, amount sdkmath.Uint) error {
	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      from.String(),
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        from.String(),
				ToAddresses: []string{to.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount:         amount,
						TokenIds:       tokenIds,
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
					},
				},
			},
		},
	}

	_, err := s.tokenizationMsgServer.TransferTokens(s.Ctx, msg)
	return err
}

// GetTokenBalance gets the token balance for an address
func (s *CrossModuleTestSuite) GetTokenBalance(collectionId sdkmath.Uint, address sdk.AccAddress, tokenId sdkmath.Uint) sdkmath.Uint {
	balanceKey := tokenizationkeeper.ConstructBalanceKey(address.String(), collectionId)
	balance, found := s.App.TokenizationKeeper.GetUserBalanceFromStore(s.Ctx, balanceKey)
	if !found {
		return sdkmath.ZeroUint()
	}

	for _, bal := range balance.Balances {
		for _, tokenIdRange := range bal.TokenIds {
			if tokenIdRange.Start.LTE(tokenId) && tokenIdRange.End.GTE(tokenId) {
				return bal.Amount
			}
		}
	}
	return sdkmath.ZeroUint()
}

// GetCoinBalance gets the coin balance for an address
func (s *CrossModuleTestSuite) GetCoinBalance(address sdk.AccAddress, denom string) sdk.Coin {
	return s.App.BankKeeper.GetBalance(s.Ctx, address, denom)
}
