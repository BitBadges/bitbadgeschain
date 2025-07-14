package v9

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgesmoduletypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

const (
	UpgradeName = "v9"
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
		faucetAddr, err := sdk.AccAddressFromBech32("bb1kx9532ujful8vgg2dht6k544ax4k9qzsp0sany")
		if err != nil {
			return nil, err
		}

		currBalance := bankKeeper.GetBalance(sdk.UnwrapSDKContext(ctx), faucetAddr, "ustake")
		if err != nil {
			return nil, err
		}

		ustakeCoins := sdk.NewCoins(currBalance)
		err = bankKeeper.SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(ctx), faucetAddr, "transfer", ustakeCoins)
		if err != nil {
			return nil, err
		}

		err = bankKeeper.BurnCoins(sdk.UnwrapSDKContext(ctx), "transfer", ustakeCoins)
		if err != nil {
			return nil, err
		}

		// module account permissions
		moduleNames := []string{
			authtypes.FeeCollectorName,
			distrtypes.ModuleName,
			minttypes.ModuleName,
			stakingtypes.BondedPoolName,
			stakingtypes.NotBondedPoolName,
			govtypes.ModuleName,
			badgesmoduletypes.ModuleName,
			ibctransfertypes.ModuleName,
			ibcfeetypes.ModuleName,
			icatypes.ModuleName,
			packetforwardtypes.ModuleName,
		}

		moduleAddresses := []sdk.AccAddress{}
		for _, moduleName := range moduleNames {
			moduleAddresses = append(moduleAddresses, authtypes.NewModuleAddress(moduleName))
		}

		// Migrate everything from "ustake" to "ubadge"
		for _, balance := range bankKeeper.GetAccountsBalances(sdk.UnwrapSDKContext(ctx)) {
			address := balance.Address
			coins := balance.Coins
			for _, coin := range coins {
				if coin.Denom == "ustake" {
					accAddr, err := sdk.AccAddressFromBech32(address)
					if err != nil {
						return nil, err
					}

					coins := sdk.NewCoins(coin)
					ubadgeCoins := sdk.NewCoins(sdk.NewCoin("ubadge", coin.Amount))

					isModuleAccount := false
					moduleName := ""
					i := 0
					for _, moduleAddr := range moduleAddresses {
						if moduleAddr.Equals(accAddr) {
							isModuleAccount = true
							moduleName = moduleNames[i]
						}
						i++
					}

					if isModuleAccount {
						err = bankKeeper.SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(ctx), accAddr, "transfer", coins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.BurnCoins(sdk.UnwrapSDKContext(ctx), "transfer", coins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.MintCoins(sdk.UnwrapSDKContext(ctx), "transfer", ubadgeCoins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.SendCoinsFromModuleToModule(sdk.UnwrapSDKContext(ctx), "transfer", moduleName, ubadgeCoins)
						if err != nil {
							return nil, err
						}

					} else {

						err = bankKeeper.SendCoinsFromAccountToModule(sdk.UnwrapSDKContext(ctx), accAddr, "transfer", coins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.BurnCoins(sdk.UnwrapSDKContext(ctx), "transfer", coins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.MintCoins(sdk.UnwrapSDKContext(ctx), "transfer", ubadgeCoins)
						if err != nil {
							return nil, err
						}

						err = bankKeeper.SendCoinsFromModuleToAccount(sdk.UnwrapSDKContext(ctx), "transfer", accAddr, ubadgeCoins)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}

		// Iterate through all rewards and convert them to ubadge
		distrKeeper.IterateValidatorCurrentRewards(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorCurrentRewards) bool {
			ustakeAmount := sdkmath.LegacyNewDec(0)
			newRewards := sdk.NewDecCoins()

			// Calculate total ustake and build new rewards
			for _, reward := range rewards.Rewards {
				if reward.Denom == "ustake" {
					ustakeAmount = ustakeAmount.Add(reward.Amount)
				} else if reward.Denom == "ubadge" {
					newRewards = newRewards.Add(reward)
				} else {
					// Keep other denoms as they are
					newRewards = newRewards.Add(reward)
				}
			}

			// Add the converted ustake amount to ubadge
			if ustakeAmount.GT(sdkmath.LegacyNewDec(0)) {
				newRewards = newRewards.Add(sdk.NewDecCoinFromDec("ubadge", ustakeAmount))
			}

			distrKeeper.SetValidatorCurrentRewards(sdk.UnwrapSDKContext(ctx), valAddr, distrtypes.ValidatorCurrentRewards{
				Rewards: newRewards,
				Period:  rewards.Period,
			})
			return false
		})

		distrKeeper.IterateValidatorOutstandingRewards(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) bool {
			ustakeAmount := sdkmath.LegacyNewDec(0)
			newRewards := sdk.NewDecCoins()

			// Calculate total ustake and build new rewards
			for _, reward := range rewards.Rewards {
				if reward.Denom == "ustake" {
					ustakeAmount = ustakeAmount.Add(reward.Amount)
				} else if reward.Denom == "ubadge" {
					newRewards = newRewards.Add(reward)
				} else {
					// Keep other denoms as they are
					newRewards = newRewards.Add(reward)
				}
			}

			// Add the converted ustake amount to ubadge
			if ustakeAmount.GT(sdkmath.LegacyNewDec(0)) {
				newRewards = newRewards.Add(sdk.NewDecCoinFromDec("ubadge", ustakeAmount))
			}

			distrKeeper.SetValidatorOutstandingRewards(sdk.UnwrapSDKContext(ctx), valAddr, distrtypes.ValidatorOutstandingRewards{
				Rewards: newRewards,
			})
			return false
		})

		distrKeeper.IterateValidatorAccumulatedCommissions(sdk.UnwrapSDKContext(ctx), func(valAddr sdk.ValAddress, commission distrtypes.ValidatorAccumulatedCommission) bool {
			ustakeAmount := sdkmath.LegacyNewDec(0)
			newCommission := sdk.NewDecCoins()

			// Calculate total ustake and build new commission
			for _, reward := range commission.Commission {
				if reward.Denom == "ustake" {
					ustakeAmount = ustakeAmount.Add(reward.Amount)
				} else if reward.Denom == "ubadge" {
					newCommission = newCommission.Add(reward)
				} else {
					// Keep other denoms as they are
					newCommission = newCommission.Add(reward)
				}
			}

			// Add the converted ustake amount to ubadge
			if ustakeAmount.GT(sdkmath.LegacyNewDec(0)) {
				newCommission = newCommission.Add(sdk.NewDecCoinFromDec("ubadge", ustakeAmount))
			}

			distrKeeper.SetValidatorAccumulatedCommission(sdk.UnwrapSDKContext(ctx), valAddr, distrtypes.ValidatorAccumulatedCommission{
				Commission: newCommission,
			})
			return false
		})

		// Set staking parameters to ustake
		currStakingParams, err := stakingKeeper.GetParams(sdk.UnwrapSDKContext(ctx))
		if err != nil {
			return nil, err
		}

		currStakingParams.BondDenom = "ubadge"
		stakingKeeper.SetParams(sdk.UnwrapSDKContext(ctx), currStakingParams)

		// Set gov params to ustake
		currParams, err := govKeeper.Params.Get(sdk.UnwrapSDKContext(ctx))
		if err != nil {
			return nil, err
		}

		currParams.MinDeposit[0].Denom = "ubadge"
		currParams.ExpeditedMinDeposit[0].Denom = "ubadge"
		govKeeper.Params.Set(sdk.UnwrapSDKContext(ctx), currParams)

		// Set mint params to ustake
		currMintParams, err := mintKeeper.Params.Get(sdk.UnwrapSDKContext(ctx))
		if err != nil {
			return nil, err
		}

		currMintParams.MintDenom = "ubadge"
		mintKeeper.Params.Set(sdk.UnwrapSDKContext(ctx), currMintParams)

		// Burn all ustake from the community pool
		feePool, err := distrKeeper.FeePool.Get(sdk.UnwrapSDKContext(ctx))
		if err != nil {
			return nil, err
		}

		ustakeDecCoins := sdk.NewDecCoins()
		for _, coin := range feePool.CommunityPool {
			if coin.Denom == "ustake" {
				ustakeDecCoins = ustakeDecCoins.Add(coin)
			}
		}

		if !ustakeDecCoins.IsZero() {
			// Remove ustake from community pool
			feePool.CommunityPool = feePool.CommunityPool.Sub(ustakeDecCoins)
			feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinFromDec("ubadge", ustakeDecCoins[0].Amount))
			err = distrKeeper.FeePool.Set(sdk.UnwrapSDKContext(ctx), feePool)
			if err != nil {
				return nil, err
			}
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
