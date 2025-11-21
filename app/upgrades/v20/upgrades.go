package v20

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Uncomment when configuring rate limits:
	cosmosmath "cosmossdk.io/math"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

const (
	UpgradeName = "v20"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper badgeskeeper.Keeper, poolManagerKeeper poolmanagerkeeper.Keeper, rateLimitKeeper ibcratelimitkeeper.Keeper) error {
	// Run badges migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	// Update poolmanager default taker fee to 0.1% (0.001)
	if err := badgeskeeper.MigratePoolManagerTakerFee(sdk.UnwrapSDKContext(ctx), poolManagerKeeper); err != nil {
		return err
	}

	channelsToSet := []string{
		"channel-0",
		"channel-2",
		"channel-3",
	}

	// Denoms to set rate limits for
	denomsToSet := []struct {
		denom    string
		decimals int
		label    string
	}{
		{
			denom:    "ubadge",
			decimals: 9,
			label:    "BADGE",
		},
		{
			denom:    "ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349",
			decimals: 6,
			label:    "USDC",
		},
		{
			denom:    "ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701",
			decimals: 6,
			label:    "ATOM",
		},
		{
			denom:    "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
			decimals: 6,
			label:    "OSMO",
		},
	}

	// Set IBC rate limits - relaxed limits that allow normal activity but catch exploits
	// These limits are intentionally high to not interfere with legitimate swaps
	// but will trigger on large exploit attempts
	for _, channelId := range channelsToSet {
		for _, denomInfo := range denomsToSet {
			// Calculate amounts based on decimals
			// For 9 decimals: 1 token = 10^9 base units
			// For 6 decimals: 1 token = 10^6 base units
			var hourlyLimit, dailyLimit, perAddressHourlyLimit, perAddressDailyLimit cosmosmath.Int

			if denomInfo.decimals == 9 {
				// BADGE: 9 decimals
				hourlyLimit = cosmosmath.NewInt(10000000000000000)          // 10,000,000 tokens
				dailyLimit = cosmosmath.NewInt(30000000000000000)           // 30,000,000 tokens
				perAddressHourlyLimit = cosmosmath.NewInt(1000000000000000) // 1,000,000 tokens
				perAddressDailyLimit = cosmosmath.NewInt(10000000000000000) // 10,000,000 tokens
			} else {
				// USDC, ATOM, OSMO: 6 decimals
				hourlyLimit = cosmosmath.NewInt(10000000000000)          // 10,000,000 tokens
				dailyLimit = cosmosmath.NewInt(30000000000000)           // 30,000,000 tokens
				perAddressHourlyLimit = cosmosmath.NewInt(1000000000000) // 1,000,000 tokens
				perAddressDailyLimit = cosmosmath.NewInt(10000000000000) // 10,000,000 tokens
			}

			rateLimitConfig := ibcratelimittypes.RateLimitConfig{
				ChannelId: channelId,
				Denom:     denomInfo.denom,
				UniqueSenderLimits: []ibcratelimittypes.UniqueSenderLimit{
					{
						MaxUniqueSenders:  10000,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR,
						TimeframeDuration: 1,
					},
					{
						MaxUniqueSenders:  100,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
						TimeframeDuration: 1,
					},
				},
				SupplyShiftLimits: []ibcratelimittypes.TimeframeLimit{
					{
						// Hourly limit: 10 million tokens per hour
						// Catches fast exploit attempts while allowing normal high-volume activity
						MaxAmount:         hourlyLimit,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR,
						TimeframeDuration: 1, // 1 hour
					},
					{
						// Daily limit: 30 million tokens per day
						// Catches large-scale exploits or slow drains while allowing normal daily activity
						MaxAmount:         dailyLimit,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY,
						TimeframeDuration: 1, // 1 day
					},
				},
				// Per-address limits to catch individual bad actors
				AddressLimits: []ibcratelimittypes.AddressLimit{
					{
						// Per address: 1 million tokens per hour
						// Prevents a single address from draining too quickly
						MaxAmount:         perAddressHourlyLimit,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR,
						TimeframeDuration: 1, // 1 hour
					},
					{
						// Per address: 10 million tokens per day
						// Allows high-volume legitimate users but catches exploiters
						MaxAmount:         perAddressDailyLimit,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY,
						TimeframeDuration: 1, // 1 day
					},
				},
			}
			if err := badgeskeeper.MigrateIBCRateLimit(sdk.UnwrapSDKContext(ctx), rateLimitKeeper, rateLimitConfig); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper badgeskeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, badgesKeeper, poolManagerKeeper, rateLimitKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
