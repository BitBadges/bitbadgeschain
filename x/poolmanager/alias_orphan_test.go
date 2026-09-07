package poolmanager_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (s *KeeperTestSuite) TestOrphanAliasBankCoinsRemainSpendable() {
	coin := sdk.NewInt64Coin("badgeslp:99999:orphan", 1000)
	s.FundAcc(s.TestAccs[0], sdk.NewCoins(coin))
	s.Require().NoError(s.App.SendmanagerKeeper.SendCoinWithAliasRouting(s.Ctx, s.TestAccs[0], s.TestAccs[1], &coin))
	s.Require().Equal(coin, s.App.BankKeeper.GetBalance(s.Ctx, s.TestAccs[1], coin.Denom))
}
