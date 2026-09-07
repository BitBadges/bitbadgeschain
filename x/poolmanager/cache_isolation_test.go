package poolmanager_test

import (
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func (s *KeeperTestSuite) TestDiscardedFeeShareWriteDoesNotEscapeContext() {
	k := &s.App.PoolManagerKeeper
	original := types.TakerFeeShareAgreement{Denom: "uatom", SkimPercent: osmomath.MustNewDecFromStr("0.01"), SkimAddress: "original"}
	s.Require().NoError(k.SetTakerFeeShareAgreementForDenom(s.Ctx, original))
	cached, _ := s.Ctx.CacheContext()
	changed := original
	changed.SkimAddress = "discarded"
	s.Require().NoError(k.SetTakerFeeShareAgreementForDenom(cached, changed))
	actual, found := k.GetTakerFeeShareAgreementFromDenomUNSAFE(s.Ctx, "uatom")
	s.Require().True(found)
	s.Require().Equal(original, actual)
}
