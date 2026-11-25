package poolmanager_test

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	poolmanagerKeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func (s *KeeperTestSuite) TestMsgServerValidatesTakerFeesAndShareAgreements() {
	s.SetupTest()
	admin := s.TestAccs[0].String()
	params := s.App.PoolManagerKeeper.GetParams(s.Ctx)
	params.TakerFeeParams.AdminAddresses = []string{admin}
	s.App.PoolManagerKeeper.SetParams(s.Ctx, params)

	msgServer := poolmanagerKeeper.NewMsgServerImpl(&s.App.PoolManagerKeeper)

	// Nil TakerFee must return error without panicking
	s.Require().NotPanics(func() {
		_, err := msgServer.SetDenomPairTakerFee(s.Ctx, &types.MsgSetDenomPairTakerFee{
			Sender: admin,
			DenomPairTakerFee: []types.DenomPairTakerFee{{
				TokenInDenom:  "foo",
				TokenOutDenom: "bar",
				TakerFee:      osmomath.Dec{},
			}},
		})
		s.Require().Error(err)
	})

	// TakerFee >= 1 or negative must be rejected at both MsgServer and keeper levels
	for _, invalidFee := range []osmomath.Dec{
		osmomath.OneDec(),
		osmomath.MustNewDecFromStr("1.5"),
		osmomath.MustNewDecFromStr("-0.01"),
	} {
		_, err := msgServer.SetDenomPairTakerFee(s.Ctx, &types.MsgSetDenomPairTakerFee{
			Sender: admin,
			DenomPairTakerFee: []types.DenomPairTakerFee{{
				TokenInDenom:  "foo",
				TokenOutDenom: "bar",
				TakerFee:      invalidFee,
			}},
		})
		s.Require().Error(err)
		s.Require().Error(s.App.PoolManagerKeeper.SenderValidationSetDenomPairTakerFee(s.Ctx, admin, "foo", "bar", invalidFee))
	}

	// Nil / out-of-range SkimPercent on MsgSetTakerFeeShareAgreementForDenom must return error without panicking
	govAddr := s.App.AccountKeeper.GetModuleAccount(s.Ctx, govtypes.ModuleName).GetAddress().String()
	for _, invalidSkim := range []osmomath.Dec{
		{},
		osmomath.MustNewDecFromStr("-0.1"),
		osmomath.MustNewDecFromStr("1.01"),
	} {
		s.Require().NotPanics(func() {
			_, err := msgServer.SetTakerFeeShareAgreementForDenom(s.Ctx, &types.MsgSetTakerFeeShareAgreementForDenom{
				Sender:      govAddr,
				Denom:       "foo",
				SkimPercent: invalidSkim,
				SkimAddress: admin,
			})
			s.Require().Error(err)
		})
		s.Require().Error(s.App.PoolManagerKeeper.SetTakerFeeShareAgreementForDenom(s.Ctx, types.TakerFeeShareAgreement{
			Denom:       "foo",
			SkimPercent: invalidSkim,
			SkimAddress: admin,
		}))
	}
}

func (s *KeeperTestSuite) TestRouterPropagatesOutOfGasPanic() {
	s.SetupTest()
	poolID := s.PrepareBalancerPool()

	lowGasCtx := s.Ctx.WithGasMeter(storetypes.NewGasMeter(10))
	routesIn := []types.SwapAmountInRoute{{PoolId: poolID, TokenOutDenom: "bar"}}
	routesOut := []types.SwapAmountOutRoute{{PoolId: poolID, TokenInDenom: "foo"}}

	s.Require().Panics(func() {
		_, _ = s.App.PoolManagerKeeper.MultihopEstimateOutGivenExactAmountIn(lowGasCtx, routesIn, sdk.NewCoin("foo", osmomath.NewInt(100)))
	})
	s.Require().Panics(func() {
		_, _ = s.App.PoolManagerKeeper.MultihopEstimateInGivenExactAmountOut(lowGasCtx, routesOut, sdk.NewCoin("bar", osmomath.NewInt(100)))
	})
	s.Require().Panics(func() {
		_, _ = s.App.PoolManagerKeeper.RouteExactAmountOut(lowGasCtx, s.TestAccs[0], routesOut, osmomath.NewInt(100000), sdk.NewCoin("bar", osmomath.NewInt(100)))
	})
}

func (s *KeeperTestSuite) TestMsgUpdateParams() {
	s.SetupTest()
	msgServer := poolmanagerKeeper.NewMsgServerImpl(&s.App.PoolManagerKeeper)
	govAddr := s.App.AccountKeeper.GetModuleAccount(s.Ctx, govtypes.ModuleName).GetAddress().String()
	initialParams := s.App.PoolManagerKeeper.GetParams(s.Ctx)

	// Nil message rejected
	_, err := msgServer.UpdateParams(s.Ctx, nil)
	s.Require().Error(err)

	// Non-gov authority rejected without mutating params
	updated := initialParams
	updated.TakerFeeParams.DefaultTakerFee = osmomath.MustNewDecFromStr("0.01")
	updated.TakerFeeParams.AdminAddresses = []string{s.TestAccs[0].String()}
	_, err = msgServer.UpdateParams(s.Ctx, &types.MsgUpdateParams{
		Authority: s.TestAccs[0].String(),
		Params:    updated,
	})
	s.Require().ErrorIs(err, types.ErrUnauthorizedGov)
	s.Require().Equal(initialParams, s.App.PoolManagerKeeper.GetParams(s.Ctx))

	// Invalid params rejected without mutating params
	invalid := updated
	invalid.TakerFeeParams.DefaultTakerFee = osmomath.OneDec()
	_, err = msgServer.UpdateParams(s.Ctx, &types.MsgUpdateParams{
		Authority: govAddr,
		Params:    invalid,
	})
	s.Require().Error(err)
	s.Require().Equal(initialParams, s.App.PoolManagerKeeper.GetParams(s.Ctx))

	// Valid gov authority persists updated params and works via MsgServiceRouter
	msg := &types.MsgUpdateParams{
		Authority: govAddr,
		Params:    updated,
	}
	handler := s.App.MsgServiceRouter().Handler(msg)
	s.Require().NotNil(handler)
	_, err = handler(s.Ctx, msg)
	s.Require().NoError(err)

	got := s.App.PoolManagerKeeper.GetParams(s.Ctx)
	s.Require().Equal(updated.TakerFeeParams.DefaultTakerFee, got.TakerFeeParams.DefaultTakerFee)
	s.Require().Equal(updated.TakerFeeParams.AdminAddresses, got.TakerFeeParams.AdminAddresses)
}
