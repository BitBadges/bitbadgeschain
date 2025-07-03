package v7

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const (
	UpgradeName = "v8"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	govKeeper govkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Run migrations
		if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
			return nil, err
		}

		// For ustake merging

		// faucetAddr, err := sdk.AccAddressFromBech32("bb1kx9532ujful8vgg2dht6k544ax4k9qzsp0sany")
		// if err != nil {
		// 	return nil, err
		// }

		// currBalance := bankKeeper.GetBalance(sdk.UnwrapSDKContext(ctx), faucetAddr, "ustake")
		// if err != nil {
		// 	return nil, err
		// }

		// ustakeCoins := sdk.NewCoins(currBalance)
		// err = bankKeeper.SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(ctx), faucetAddr, "transfer", ustakeCoins)
		// if err != nil {
		// 	return nil, err
		// }

		// err = bankKeeper.BurnCoins(sdk.UnwrapSDKContext(ctx), "transfer", ustakeCoins)
		// if err != nil {
		// 	return nil, err
		// }

		// // Migrate everything from "ustake" to "ubadge"
		// // 1. Set up conversion for ustake -> ubadge
		// //    Note that currently delegated ustake will be converted to ubadge at 1:1 (in delegator shares)
		// for _, balance := range bankKeeper.GetAccountsBalances(sdk.UnwrapSDKContext(ctx)) {
		// 	address := balance.Address
		// 	coins := balance.Coins
		// 	for _, coin := range coins {
		// 		if coin.Denom == "ustake" {
		// 			accAddr, err := sdk.AccAddressFromBech32(address)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			coins := sdk.NewCoins(coin)
		// 			err = bankKeeper.SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(ctx), accAddr, "transfer", coins)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			err = bankKeeper.BurnCoins(sdk.UnwrapSDKContext(ctx), "transfer", coins)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			ubadgeCoins := sdk.NewCoins(sdk.NewCoin("ubadge", coin.Amount))
		// 			err = bankKeeper.MintCoins(sdk.UnwrapSDKContext(ctx), "transfer", ubadgeCoins)
		// 			if err != nil {
		// 				return nil, err
		// 			}
		// 		}
		// 	}

		// }

		// // Iterate through all rewards and convert them to ubadge
		// distrKeeper.IterateValidatorCurrentRewards(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorCurrentRewards) bool {
		// 	newRewards := sdk.DecCoins{}
		// 	for _, reward := range rewards.Rewards {
		// 		reward.Denom = "ubadge"
		// 		newRewards = newRewards.Add(reward)
		// 	}

		// 	distrKeeper.SetValidatorCurrentRewards(sdk.UnwrapSDKContext(ctx), valAddr, distrtypes.ValidatorCurrentRewards{
		// 		Rewards: newRewards,
		// 		Period:  rewards.Period,
		// 	})
		// 	return false
		// })

		// distrKeeper.IterateValidatorOutstandingRewards(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) bool {
		// 	newRewards := sdk.DecCoins{}
		// 	for _, reward := range rewards.Rewards {
		// 		reward.Denom = "ubadge"
		// 		newRewards = newRewards.Add(reward)
		// 	}

		// 	distrKeeper.SetValidatorOutstandingRewards(sdk.UnwrapSDKContext(ctx), valAddr, distrtypes.ValidatorOutstandingRewards{
		// 		Rewards: newRewards,
		// 	})
		// 	return false
		// })

		// distrKeeper.IterateValidatorAccumulatedCommissions(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, commission distrtypes.ValidatorAccumulatedCommission) bool {
		// 	newCommission := distrtypes.ValidatorAccumulatedCommission{
		// 		Commission: sdk.DecCoins{},
		// 	}

		// 	for _, reward := range commission.Commission {
		// 		reward.Denom = "ubadge"
		// 		newCommission.Commission = newCommission.Commission.Add(reward)
		// 	}

		// 	distrKeeper.SetValidatorAccumulatedCommission(sdk.UnwrapSDKContext(ctx), valAddr, newCommission)
		// 	return false
		// })

		// // Set staking parameters to ustake
		// currStakingParams, err := stakingKeeper.GetParams(sdk.UnwrapSDKContext(ctx))
		// if err != nil {
		// 	return nil, err
		// }

		// currStakingParams.BondDenom = "ubadge"
		// stakingKeeper.SetParams(sdk.UnwrapSDKContext(ctx), currStakingParams)

		// // Set gov params to ustake
		// currParams, err := govKeeper.Params.Get(sdk.UnwrapSDKContext(ctx))
		// if err != nil {
		// 	return nil, err
		// }

		// currParams.MinDeposit[0].Denom = "ubadge"
		// currParams.ExpeditedMinDeposit[0].Denom = "ubadge"
		// govKeeper.Params.Set(sdk.UnwrapSDKContext(ctx), currParams)

		// // Set mint params to ustake
		// currMintParams, err := mintKeeper.Params.Get(sdk.UnwrapSDKContext(ctx))
		// if err != nil {
		// 	return nil, err
		// }

		// currMintParams.MintDenom = "ubadge"
		// mintKeeper.Params.Set(sdk.UnwrapSDKContext(ctx), currMintParams)

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
