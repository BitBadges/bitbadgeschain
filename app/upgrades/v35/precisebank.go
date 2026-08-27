package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// PreciseBankStoreKey is the retired module's store name.
//
// The module itself is unwired as of this upgrade — at 18 decimals the base and
// extended denoms are the same and there is no precision gap left to bridge —
// but its store is mounted through the upgrade so the balances inside it can be
// paid out before it goes away. It carries no keeper and no module account
// permissions any more, so this migration reads the raw store.
const PreciseBankStoreKey = "precisebank"

// The retired module's key layout, mirrored here rather than imported. Importing
// its types would pull in the module's own ConversionFactor(), which reads the
// EVM's *current* decimals — 18 by the time this runs — and would misread every
// value in the store.
var (
	// PreciseBankFractionalPrefix keys per-account fractional balances.
	PreciseBankFractionalPrefix = []byte{0x01}
	// PreciseBankRemainderKey holds the unowned fractional remainder.
	PreciseBankRemainderKey = []byte{0x02}
)

// PreciseBankMigrationResult reports what MigratePreciseBank paid out.
type PreciseBankMigrationResult struct {
	Accounts  int
	Credited  sdkmath.Int
	Remainder sdkmath.Int
	Reserve   sdkmath.Int
	// Surplus is what the reserve held beyond what the fractional balances and
	// the remainder claim. Non-zero means somebody sent coins to the reserve
	// address, which anyone could do.
	Surplus sdkmath.Int
}

// MigratePreciseBank pays out the fractional balances x/precisebank held and
// empties its reserve.
//
// At 9 decimals the EVM saw 18-decimal balances by combining an integer ubadge
// bank balance with a per-account fractional remainder of 0..10^9-1, held in
// x/precisebank's own store and backed 1:1 by ubadge sitting on the
// precisebank module account. Deleting that store as part of the upgrade —
// which is what "Deleted: precisebank" in the store upgrades did, before the
// handler ever runs — destroys the fractional part of every EVM account that
// ever received a non-integer-ubadge transfer. eth_getBalance would go from
// ubadge*10^9 + fractional to ubadge*10^9, silently, and the reserve backing
// the difference would strand on a module address whose module no longer
// exists while still counting toward supply.
//
// Runs after RedenominateBank, so balances are already in the new denom and a
// fractional unit is exactly one unit of it: 1 abadge. Which is the whole point
// of the redenomination — the precision x/precisebank was faking is now real.
//
// The relation sum(fractional) + remainder == reserve is what x/precisebank's
// genesis validation checks. It is *not* an invariant the module enforces at
// runtime: x/precisebank registers no invariant function, so nothing re-checks
// it between genesis and this upgrade. In particular the reserve address is a
// plain module address that was never added to the bank's blocked list, so
// anyone can send coins to it — and on mainnet somebody has: at the time this
// was written the reserve held 3 ubadge of dust that no fractional balance
// claims.
//
// So the check here is deliberately one-sided:
//
//   - reserve < owed  -> fail. The reserve genuinely cannot cover what the
//     fractional balances claim, real value is at risk, and paying out anyway
//     would mint from nothing. A halted upgrade is recoverable; a wrong payout
//     is not.
//   - reserve > owed  -> succeed. Pay out exactly what is owed and dispose of
//     the surplus, logging it loudly.
//
// A two-sided equality check here was a free, public, chain-halting DoS: one
// ubadge sent to the reserve before the upgrade height becomes 10^9 abadge of
// surplus after step 2, the equality fails, CustomUpgradeHandlerLogic returns
// an error, and the upgrade BeginBlocker panics every node at the upgrade
// height. Cost to the attacker: 10^-9 BADGE plus gas.
//
// The surplus is burned rather than sent to the community pool. It is unowned
// dust on an address whose module no longer exists and which nothing can spend
// from after this upgrade, so leaving it anywhere keeps supply inflated against
// coins nobody can reach. Burning is also exactly what the unowned remainder
// below already gets, so both kinds of unowned residue are disposed of the same
// way rather than by two different rules.
func MigratePreciseBank(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	bk bankkeeper.BaseKeeper,
) (PreciseBankMigrationResult, error) {
	res := PreciseBankMigrationResult{
		Credited:  sdkmath.ZeroInt(),
		Remainder: sdkmath.ZeroInt(),
		Reserve:   sdkmath.ZeroInt(),
		Surplus:   sdkmath.ZeroInt(),
	}

	if storeKey == nil {
		return res, fmt.Errorf("precisebank store key is not wired; refusing to run the upgrade blind to it")
	}

	store := ctx.KVStore(storeKey)

	// Collected before writing: the payouts below touch the bank store, and the
	// deletions touch this one, underneath an open iterator.
	type fractional struct {
		addr   sdk.AccAddress
		amount sdkmath.Int
		key    []byte
	}
	var holdings []fractional

	iter := storetypes.KVStorePrefixIterator(store, PreciseBankFractionalPrefix)
	for ; iter.Valid(); iter.Next() {
		key := append([]byte(nil), iter.Key()...)
		addr := sdk.AccAddress(key[len(PreciseBankFractionalPrefix):])

		var amount sdkmath.Int
		if err := amount.Unmarshal(iter.Value()); err != nil {
			iter.Close()
			return res, fmt.Errorf("unmarshalling fractional balance for %s: %w", addr, err)
		}
		holdings = append(holdings, fractional{addr: addr, amount: amount, key: key})
		res.Credited = res.Credited.Add(amount)
	}
	iter.Close()
	res.Accounts = len(holdings)

	if bz := store.Get(PreciseBankRemainderKey); bz != nil {
		var remainder sdkmath.Int
		if err := remainder.Unmarshal(bz); err != nil {
			return res, fmt.Errorf("unmarshalling precisebank remainder: %w", err)
		}
		res.Remainder = remainder
	}

	reserveAddr := authtypes.NewModuleAddress(PreciseBankStoreKey)
	res.Reserve = bk.GetBalance(ctx, reserveAddr, newDenom()).Amount

	owed := res.Credited.Add(res.Remainder)
	if res.Reserve.LT(owed) {
		return res, fmt.Errorf(
			"precisebank reserve holds %s %s but owes %s (%s fractional + %s remainder): "+
				"refusing to pay out from a reserve that does not back it",
			res.Reserve, newDenom(), owed, res.Credited, res.Remainder,
		)
	}
	res.Surplus = res.Reserve.Sub(owed)
	if res.Surplus.IsPositive() {
		ctx.Logger().Error(
			"v35: precisebank reserve holds more than the fractional balances claim; "+
				"burning the surplus rather than halting the upgrade",
			"reserve", res.Reserve.String(),
			"owed", owed.String(),
			"surplus", res.Surplus.String(),
			"denom", newDenom(),
		)
	}

	if res.Accounts == 0 && res.Remainder.IsZero() && res.Reserve.IsZero() {
		return res, nil
	}

	// Credit each holder. A fractional unit is one unit of the new denom, so
	// this is an addition, not a conversion.
	for _, h := range holdings {
		current := bk.GetBalance(ctx, h.addr, newDenom()).Amount
		credited := sdk.NewCoin(newDenom(), current.Add(h.amount))
		if err := bk.UncheckedSetBalance(ctx, h.addr, credited); err != nil {
			return res, fmt.Errorf("crediting fractional balance to %s: %w", h.addr, err)
		}
		store.Delete(h.key)
	}
	store.Delete(PreciseBankRemainderKey)

	// Empty the reserve. Its module no longer exists, so anything left there is
	// unreachable but still counted in supply.
	if err := bk.UncheckedSetBalance(ctx, reserveAddr, sdk.NewCoin(newDenom(), sdkmath.ZeroInt())); err != nil {
		return res, fmt.Errorf("draining the precisebank reserve: %w", err)
	}

	// Two kinds of unowned residue come out of the reserve and both are burned.
	//
	// The remainder was backed but owned by nobody: it is the rounding residue
	// x/precisebank carried so its reserve stayed whole. The surplus is
	// whatever anyone sent to the reserve address on top of that. Neither can
	// be claimed by any account after this upgrade, so leaving either one
	// anywhere inflates supply against coins nothing can reach.
	//
	// Burned through the same module account RedenominateBank borrows, since
	// precisebank has no burner permission any more — and no module account.
	unowned := res.Remainder.Add(res.Surplus)
	if unowned.IsPositive() {
		burnAddr := authtypes.NewModuleAddress(redenominationModule)
		staged := bk.GetBalance(ctx, burnAddr, newDenom()).Amount.Add(unowned)
		if err := bk.UncheckedSetBalance(ctx, burnAddr, sdk.NewCoin(newDenom(), staged)); err != nil {
			return res, fmt.Errorf("staging the precisebank remainder for burn: %w", err)
		}
		burn := sdk.NewCoins(sdk.NewCoin(newDenom(), unowned))
		if err := bk.BurnCoins(ctx, redenominationModule, burn); err != nil {
			return res, fmt.Errorf("burning the precisebank remainder: %w", err)
		}
	}

	ctx.Logger().Info(
		"v35: paid out precisebank fractional balances",
		"accounts", res.Accounts,
		"credited", res.Credited.String(),
		"remainder_burned", res.Remainder.String(),
		"surplus_burned", res.Surplus.String(),
		"reserve_drained", res.Reserve.String(),
	)

	return res, nil
}
