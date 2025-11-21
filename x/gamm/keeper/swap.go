package keeper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// swapExactAmountIn is an internal method for swapping an exact amount of tokens
// as input to a pool, using the provided spreadFactor. This is intended to allow
// different spread factors as determined by multi-hops, or when recovering from
// chain liveness failures.
// TODO: investigate if spreadFactor can be unexported
// https://github.com/osmosis-labs/osmosis/issues/3130
func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	pool poolmanagertypes.PoolI,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount osmomath.Int,
	spreadFactor osmomath.Dec,
	affiliates []poolmanagertypes.Affiliate,
) (tokenOutAmount osmomath.Int, err error) {
	if tokenIn.Denom == tokenOutDenom {
		return osmomath.Int{}, errors.New("cannot trade same denomination in and out")
	}
	poolSpreadFactor := pool.GetSpreadFactor(ctx)
	if spreadFactor.LT(poolSpreadFactor.QuoInt64(2)) {
		return osmomath.Int{}, fmt.Errorf("given spread factor (%s) must be greater than or equal to half of the pool's spread factor (%s)", spreadFactor, poolSpreadFactor)
	}
	tokensIn := sdk.Coins{tokenIn}

	defer func() {
		if r := recover(); r != nil {
			tokenOutAmount = osmomath.Int{}
			if isErr, d := osmoutils.IsOutOfGasError(r); isErr {
				err = fmt.Errorf("function swapExactAmountIn failed due to lack of gas: %v", d)
			} else {
				err = fmt.Errorf("function swapExactAmountIn failed due to internal reason: %v", r)
			}
		}
	}()

	cfmmPool, err := asCFMMPool(pool)
	if err != nil {
		return osmomath.Int{}, err
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	tokenOutCoin, err := cfmmPool.SwapOutAmtGivenIn(ctx, tokensIn, tokenOutDenom, spreadFactor)
	if err != nil {
		return osmomath.Int{}, err
	}

	tokenOutAmount = tokenOutCoin.Amount

	if !tokenOutAmount.IsPositive() {
		return osmomath.Int{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
	}

	if tokenOutAmount.LT(tokenOutMinAmount) {
		return osmomath.Int{}, errorsmod.Wrapf(types.ErrLimitMinAmount, "%s token is lesser than min amount", tokenOutDenom)
	}

	// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
	// Also emits swap event and updates related liquidity metrics
	// Affiliates are processed inside updatePoolForSwap
	err = k.updatePoolForSwap(ctx, pool, sender, tokenIn, tokenOutCoin, affiliates, tokenOutMinAmount)
	if err != nil {
		return osmomath.Int{}, err
	}

	// Calculate affiliate fees if provided
	totalFeeAmount := osmomath.ZeroInt()
	if len(affiliates) > 0 {
		var calcErr error
		totalFeeAmount, _, calcErr = k.calculateAffiliateFees(ctx, affiliates, tokenOutMinAmount, tokenOutDenom)
		if calcErr != nil {
			return osmomath.Int{}, calcErr
		}

		// Validate that fees don't exceed output
		if totalFeeAmount.IsPositive() {
			if tokenOutAmount.LT(totalFeeAmount) {
				return osmomath.Int{}, fmt.Errorf("affiliate fees exceed swap output: fees=%s, output=%s", totalFeeAmount.String(), tokenOutAmount.String())
			}
		}
	}

	// Return adjusted tokenOutAmount after affiliate fees
	if totalFeeAmount.IsPositive() {
		tokenOutAmount = tokenOutAmount.Sub(totalFeeAmount)
	}

	return tokenOutAmount, nil
}

// RouteExactAmountIn routes a swap through multiple pools using the poolmanager
// This method forwards to the poolmanager's RouteExactAmountIn for multi-hop swaps
func (k Keeper) RouteExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	routes []poolmanagertypes.SwapAmountInRoute,
	tokenIn sdk.Coin,
	tokenOutMinAmount osmomath.Int,
	affiliates []poolmanagertypes.Affiliate,
) (tokenOutAmount osmomath.Int, err error) {
	if k.poolManager == nil {
		return osmomath.Int{}, fmt.Errorf("pool manager not set")
	}
	return k.poolManager.RouteExactAmountIn(ctx, sender, routes, tokenIn, tokenOutMinAmount, affiliates)
}

// SwapExactAmountOut is a method for swapping to get an exact number of tokens out of a pool,
// using the provided spreadFactor.
// This is intended to allow different spread factors as determined by multi-hops,
// or when recovering from chain liveness failures.
func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	pool poolmanagertypes.PoolI,
	tokenInDenom string,
	tokenInMaxAmount osmomath.Int,
	tokenOut sdk.Coin,
	spreadFactor osmomath.Dec,
) (tokenInAmount osmomath.Int, err error) {
	if tokenInDenom == tokenOut.Denom {
		return osmomath.Int{}, errors.New("cannot trade same denomination in and out")
	}
	defer func() {
		if r := recover(); r != nil {
			tokenInAmount = osmomath.Int{}
			if isErr, d := osmoutils.IsOutOfGasError(r); isErr {
				err = fmt.Errorf("function swapExactAmountOut failed due to lack of gas: %v", d)
			} else {
				err = fmt.Errorf("function swapExactAmountOut failed due to internal reason: %v", r)
			}
		}
	}()

	liquidity, err := k.GetTotalPoolLiquidity(ctx, pool.GetId())
	if err != nil {
		return osmomath.Int{}, err
	}

	poolOutBal := liquidity.AmountOf(tokenOut.Denom)
	if tokenOut.Amount.GTE(poolOutBal) {
		return osmomath.Int{}, errorsmod.Wrapf(types.ErrTooManyTokensOut,
			"can't get more tokens out than there are tokens in the pool")
	}

	cfmmPool, err := asCFMMPool(pool)
	if err != nil {
		return osmomath.Int{}, err
	}

	tokenIn, err := cfmmPool.SwapInAmtGivenOut(ctx, sdk.Coins{tokenOut}, tokenInDenom, spreadFactor)
	if err != nil {
		return osmomath.Int{}, err
	}
	tokenInAmount = tokenIn.Amount

	if tokenInAmount.LTE(osmomath.ZeroInt()) {
		return osmomath.Int{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "token amount is zero or negative")
	}

	if tokenInAmount.GT(tokenInMaxAmount) {
		return osmomath.Int{}, errorsmod.Wrapf(types.ErrLimitMaxAmount, "Swap requires %s, which is greater than the amount %s", tokenIn, tokenInMaxAmount)
	}

	err = k.updatePoolForSwap(ctx, pool, sender, tokenIn, tokenOut, nil, osmomath.ZeroInt())
	if err != nil {
		return osmomath.Int{}, err
	}
	return tokenInAmount, nil
}

// CalcOutAmtGivenIn calculates the amount of tokenOut given tokenIn and the pool's current state.
// Returns error if the given pool is not a CFMM pool. Returns error on internal calculations.
func (k Keeper) CalcOutAmtGivenIn(
	ctx sdk.Context,
	poolI poolmanagertypes.PoolI,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	spreadFactor osmomath.Dec,
) (tokenOut sdk.Coin, err error) {
	cfmmPool, err := asCFMMPool(poolI)
	if err != nil {
		return sdk.Coin{}, err
	}
	return cfmmPool.CalcOutAmtGivenIn(ctx, sdk.NewCoins(tokenIn), tokenOutDenom, spreadFactor)
}

// CalcInAmtGivenOut calculates the amount of tokenIn given tokenOut and the pool's current state.
// Returns error if the given pool is not a CFMM pool. Returns error on internal calculations.
func (k Keeper) CalcInAmtGivenOut(
	ctx sdk.Context,
	poolI poolmanagertypes.PoolI,
	tokenOut sdk.Coin,
	tokenInDenom string,
	spreadFactor osmomath.Dec,
) (tokenIn sdk.Coin, err error) {
	cfmmPool, err := asCFMMPool(poolI)
	if err != nil {
		return sdk.Coin{}, err
	}
	return cfmmPool.CalcInAmtGivenOut(ctx, sdk.NewCoins(tokenOut), tokenInDenom, spreadFactor)
}

// calculateAffiliateFees calculates the total affiliate fees and validates affiliates
// Returns the total fee amount and individual fee amounts per affiliate
// Fees are calculated from tokenOutMinAmount (minimum expected output), not the actual swap output
func (k Keeper) calculateAffiliateFees(
	ctx sdk.Context,
	affiliates []poolmanagertypes.Affiliate,
	tokenOutMinAmount osmomath.Int,
	tokenOutDenom string,
) (totalFeeAmount osmomath.Int, affiliateFees []sdk.Coin, err error) {
	totalFeeAmount = osmomath.ZeroInt()
	affiliateFees = make([]sdk.Coin, 0, len(affiliates))

	if len(affiliates) == 0 {
		return totalFeeAmount, affiliateFees, nil
	}

	// Validate affiliates
	totalBasisPoints := uint64(0)
	for i, affiliate := range affiliates {
		// Validate address
		if affiliate.Address == "" {
			return osmomath.ZeroInt(), nil, fmt.Errorf("affiliate_%d.address is required", i)
		}

		_, err := sdk.AccAddressFromBech32(affiliate.Address)
		if err != nil {
			return osmomath.ZeroInt(), nil, fmt.Errorf("invalid affiliate_%d.address: %w", i, err)
		}

		// Validate basis_points_fee
		if affiliate.BasisPointsFee == "" {
			return osmomath.ZeroInt(), nil, fmt.Errorf("affiliate_%d.basis_points_fee is required", i)
		}

		basisPoints, err := strconv.ParseUint(affiliate.BasisPointsFee, 10, 64)
		if err != nil {
			return osmomath.ZeroInt(), nil, fmt.Errorf("invalid affiliate_%d.basis_points_fee: %w", i, err)
		}

		// Basis points should not exceed 10000 (100%)
		if basisPoints > 10000 {
			return osmomath.ZeroInt(), nil, fmt.Errorf("affiliate_%d.basis_points_fee cannot exceed 10000 (100%%)", i)
		}

		totalBasisPoints += basisPoints
	}

	// Total basis points should not exceed 10000 (100%)
	if totalBasisPoints > 10000 {
		return osmomath.ZeroInt(), nil, fmt.Errorf("total affiliate basis_points_fee cannot exceed 10000 (100%%), got %d", totalBasisPoints)
	}

	// Calculate fees for each affiliate
	for i, affiliate := range affiliates {
		basisPoints, err := strconv.ParseUint(affiliate.BasisPointsFee, 10, 64)
		if err != nil {
			// This should not happen as we validated earlier, but handle it just in case
			return osmomath.ZeroInt(), nil, fmt.Errorf("invalid affiliate_%d.basis_points_fee: %w", i, err)
		}

		// Calculate fee from minimum output amount: (tokenOutMinAmount * basisPoints) / 10000
		feeAmount := tokenOutMinAmount.Mul(osmomath.NewIntFromUint64(basisPoints)).Quo(osmomath.NewInt(10000))

		if feeAmount.IsPositive() {
			affiliateFee := sdk.NewCoin(tokenOutDenom, feeAmount)
			affiliateFees = append(affiliateFees, affiliateFee)
			totalFeeAmount = totalFeeAmount.Add(feeAmount)
		} else {
			// Add zero coin to maintain index alignment
			affiliateFees = append(affiliateFees, sdk.NewCoin(tokenOutDenom, osmomath.ZeroInt()))
		}
	}

	return totalFeeAmount, affiliateFees, nil
}

// updatePoolForSwap takes a pool, sender, and tokenIn, tokenOut amounts
// It then updates the pool's balances to the new reserve amounts, and
// sends the in tokens from the sender to the pool, and the out tokens from the pool to the sender.
// If affiliates are provided, fees are deducted from tokenOut before sending to sender,
// and fees are sent directly from the pool to affiliates.
func (k Keeper) updatePoolForSwap(
	ctx sdk.Context,
	pool poolmanagertypes.PoolI,
	sender sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	affiliates []poolmanagertypes.Affiliate,
	tokenOutMinAmount osmomath.Int,
) error {
	tokensIn := sdk.Coins{tokenIn}
	poolAddress := pool.GetAddress()

	err := k.setPool(ctx, pool)
	if err != nil {
		return err
	}

	// 1. Calculate affiliate fees (no sends)
	totalFeeAmount := osmomath.ZeroInt()
	var affiliateFees []sdk.Coin
	if len(affiliates) > 0 {
		var err error
		totalFeeAmount, affiliateFees, err = k.calculateAffiliateFees(ctx, affiliates, tokenOutMinAmount, tokenOut.Denom)
		if err != nil {
			return err
		}

		// Validate that fees don't exceed tokenOut
		if totalFeeAmount.IsPositive() {
			if tokenOut.Amount.LT(totalFeeAmount) {
				return fmt.Errorf("affiliate fees exceed swap output: fees=%s, output=%s", totalFeeAmount.String(), tokenOut.Amount.String())
			}
		}
	}

	// 2. Send tokenIn from sender to pool
	err = k.SendCoinsToPoolWithWrapping(ctx, sender, poolAddress, sdk.Coins{
		tokenIn,
	})
	if err != nil {
		return err
	}

	// 3. Send (tokenOut - fees) from pool to sender
	tokenOutToSender := tokenOut.Amount.Sub(totalFeeAmount)
	if tokenOutToSender.IsPositive() {
		tokenOutCoinToSender := sdk.NewCoin(tokenOut.Denom, tokenOutToSender)
		err = k.SendCoinsFromPoolWithUnwrapping(ctx, pool.GetAddress(), sender, sdk.NewCoins(tokenOutCoinToSender))
		if err != nil {
			return err
		}
	}

	// 4. Send fees from pool to affiliates
	if len(affiliates) > 0 && totalFeeAmount.IsPositive() {
		for i, affiliate := range affiliates {
			if affiliateFees[i].Amount.IsPositive() {
				affiliateAddr, err := sdk.AccAddressFromBech32(affiliate.Address)
				if err != nil {
					return fmt.Errorf("invalid affiliate_%d.address: %w", i, err)
				}

				// Send affiliate fee from pool to affiliate using wrapper functions
				err = k.SendCoinsFromPoolWithUnwrapping(ctx, pool.GetAddress(), affiliateAddr, sdk.NewCoins(affiliateFees[i]))
				if err != nil {
					return fmt.Errorf("failed to send affiliate fee to %s: %w", affiliate.Address, err)
				}

				// Emit event for affiliate fee distribution
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"affiliate_fee_distributed",
						sdk.NewAttribute("module", "gamm"),
						sdk.NewAttribute("affiliate_address", affiliate.Address),
						sdk.NewAttribute("basis_points_fee", affiliate.BasisPointsFee),
						sdk.NewAttribute("fee_amount", affiliateFees[i].String()),
					),
				)
			}
		}
	}

	// Calculate tokensOut for liquidity tracking (full amount, before fee deduction)
	tokensOut := sdk.Coins{tokenOut}

	// Emit swap event. Note that we emit these at the layer of each pool module rather than the poolmanager module
	// since poolmanager has many swap wrapper APIs that we would need to consider.
	// Search for references to this function to see where else it is used.
	// Each new pool module will have to emit this event separately
	k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	k.RecordTotalLiquidityDecrease(ctx, tokensOut)

	// Global pool invariant check: ensure pool has enough underlying assets for all recorded liquidity
	// This check happens after liquidity increments/decrements
	err = k.CheckPoolLiquidityInvariant(ctx, pool)
	if err != nil {
		return err
	}

	// Emit event for fallback path
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"pool_swap_success",
			sdk.NewAttribute("module", "gamm"),
			sdk.NewAttribute("sender", sender.String()),
			sdk.NewAttribute("pool_id", strconv.FormatUint(pool.GetId(), 10)),
			sdk.NewAttribute("token_in", tokensIn.String()),
			sdk.NewAttribute("token_out", tokensOut.String()),
			sdk.NewAttribute("denom_in", tokenIn.Denom),
			sdk.NewAttribute("denom_out", tokenOut.Denom),
		),
	)

	return nil
}
