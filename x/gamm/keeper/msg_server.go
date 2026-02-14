package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

type msgServer struct {
	keeper *Keeper
}

func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func NewBalancerMsgServerImpl(keeper *Keeper) balancer.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var (
	_ types.MsgServer    = msgServer{}
	_ balancer.MsgServer = msgServer{}
)

// CreateBalancerPool is a create balancer pool message.
func (server msgServer) CreateBalancerPool(goCtx context.Context, msg *balancer.MsgCreateBalancerPool) (*balancer.MsgCreateBalancerPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	poolId, err := server.CreatePool(goCtx, msg)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "create_balancer_pool"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &balancer.MsgCreateBalancerPoolResponse{PoolID: poolId}, nil
}

// CreatePool attempts to create a pool returning the newly created pool ID or an error upon failure.
// The pool creation fee is used to fund the community pool.
// It will create a dedicated module account for the pool and sends the initial liquidity to the created module account.
func (server msgServer) CreatePool(goCtx context.Context, msg poolmanagertypes.CreatePoolMsg) (poolId uint64, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolId, err = server.keeper.poolManager.CreatePool(ctx, msg)
	if err != nil {
		return 0, err
	}

	return poolId, nil
}

// JoinPool routes `JoinPoolNoSwap` where we do an abstract calculation on needed lp liquidity coins to get the designated
// amount of shares for the pool. (This is done by taking the number of shares we want and then using the total number of shares
// to get the ratio of the pool it accounts for. Using this ratio, we iterate over all pool assets to get the number of tokens we need
// to get the specified number of shares).
// Using the number of tokens needed to actually join the pool, we do a basic sanity check on whether the token does not exceed
// `TokenInMaxs`. Then we hit the actual implementation of `JoinPool` defined by each pool model.
// `JoinPool` takes in the tokensIn calculated above as the parameter rather than using the number of shares provided in the msg.
// This can result in negotiable difference between the number of shares provided within the msg
// and the actual number of share amount resulted from joining pool.
// Internal logic flow for each pool model is as follows:
// Balancer: TokensInMaxs provided as the argument must either contain no tokens or containing all assets in the pool.
// * For the case of a not containing tokens, we simply perform calculation of sharesOut and needed amount of tokens for joining the pool
func (server msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	neededLp, sharesOut, err := server.keeper.JoinPoolNoSwap(ctx, sender, msg.PoolId, msg.ShareOutAmount, msg.TokenInMaxs)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "join_pool"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgJoinPoolResponse{
		ShareOutAmount: sharesOut,
		TokenIn:        neededLp,
	}, nil
}

func (server msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	exitCoins, err := server.keeper.ExitPool(ctx, sender, msg.PoolId, msg.ShareInAmount, msg.TokenOutMins)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "exit_pool"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgExitPoolResponse{
		TokenOut: exitCoins,
	}, nil
}

func (server msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Convert gamm Affiliate types to poolmanager Affiliate types
	var poolmanagerAffiliates []poolmanagertypes.Affiliate
	if len(msg.Affiliates) > 0 {
		poolmanagerAffiliates = make([]poolmanagertypes.Affiliate, len(msg.Affiliates))
		for i, affiliate := range msg.Affiliates {
			poolmanagerAffiliates[i] = poolmanagertypes.Affiliate{
				BasisPointsFee: affiliate.BasisPointsFee,
				Address:        affiliate.Address,
			}
		}
	}

	tokenOutAmount, err := server.keeper.poolManager.RouteExactAmountIn(ctx, sender, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount, poolmanagerAffiliates)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "swap_exact_amount_in"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSwapExactAmountInResponse{TokenOutAmount: tokenOutAmount}, nil
}

func (server msgServer) SwapExactAmountInWithIBCTransfer(goCtx context.Context, msg *types.MsgSwapExactAmountInWithIBCTransfer) (*types.MsgSwapExactAmountInWithIBCTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Get the token out denom from the last route
	if len(msg.Routes) == 0 {
		return nil, fmt.Errorf("routes cannot be empty")
	}
	lastRoute := msg.Routes[len(msg.Routes)-1]
	tokenOutDenom := lastRoute.TokenOutDenom

	// Convert gamm Affiliate types to poolmanager Affiliate types
	var poolmanagerAffiliates []poolmanagertypes.Affiliate
	if len(msg.Affiliates) > 0 {
		poolmanagerAffiliates = make([]poolmanagertypes.Affiliate, len(msg.Affiliates))
		for i, affiliate := range msg.Affiliates {
			poolmanagerAffiliates[i] = poolmanagertypes.Affiliate{
				BasisPointsFee: affiliate.BasisPointsFee,
				Address:        affiliate.Address,
			}
		}
	}

	// Perform the swap first (affiliates are processed inside updatePoolForSwap)
	tokenOutAmount, err := server.keeper.poolManager.RouteExactAmountIn(ctx, sender, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount, poolmanagerAffiliates)
	if err != nil {
		return nil, err
	}
	// Create the token out coin
	tokenOut := sdk.NewCoin(tokenOutDenom, tokenOutAmount)

	// Execute IBC transfer using the custom hooks keeper pattern
	// Since there's no intermediate sender, we use msg.Sender directly
	if err := server.keeper.ExecuteIBCTransfer(ctx, sender, &msg.IbcTransferInfo, tokenOut); err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "swap_exact_amount_in_with_ibc_transfer"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSwapExactAmountInWithIBCTransferResponse{TokenOutAmount: tokenOutAmount}, nil
}

func (server msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokenInAmount, err := server.keeper.poolManager.RouteExactAmountOut(ctx, sender, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "swap_exact_amount_out"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSwapExactAmountOutResponse{TokenInAmount: tokenInAmount}, nil
}

// JoinSwapExactAmountIn is an LP transaction, that will LP all of the provided tokensIn coins.
// * For the case of a single token, we simply perform single asset join (balancer notation: pAo, pool shares amount out,
// given single asset in).
// For more details on the calculation of the number of shares look at the CalcJoinPoolShares function for the appropriate pool style
func (server msgServer) JoinSwapExternAmountIn(goCtx context.Context, msg *types.MsgJoinSwapExternAmountIn) (*types.MsgJoinSwapExternAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokensIn := sdk.Coins{msg.TokenIn}
	shareOutAmount, err := server.keeper.JoinSwapExactAmountIn(ctx, sender, msg.PoolId, tokensIn, msg.ShareOutMinAmount)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "join_swap_extern_amount_in"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgJoinSwapExternAmountInResponse{ShareOutAmount: shareOutAmount}, nil
}

func (server msgServer) JoinSwapShareAmountOut(goCtx context.Context, msg *types.MsgJoinSwapShareAmountOut) (*types.MsgJoinSwapShareAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokenInAmount, err := server.keeper.JoinSwapShareAmountOut(ctx, sender, msg.PoolId, msg.TokenInDenom, msg.ShareOutAmount, msg.TokenInMaxAmount)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "join_swap_share_amount_out"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgJoinSwapShareAmountOutResponse{TokenInAmount: tokenInAmount}, nil
}

func (server msgServer) ExitSwapExternAmountOut(goCtx context.Context, msg *types.MsgExitSwapExternAmountOut) (*types.MsgExitSwapExternAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	shareInAmount, err := server.keeper.ExitSwapExactAmountOut(ctx, sender, msg.PoolId, msg.TokenOut, msg.ShareInMaxAmount)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "exit_swap_extern_amount_out"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgExitSwapExternAmountOutResponse{ShareInAmount: shareInAmount}, nil
}

func (server msgServer) ExitSwapShareAmountIn(goCtx context.Context, msg *types.MsgExitSwapShareAmountIn) (*types.MsgExitSwapShareAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokenOutAmount, err := server.keeper.ExitSwapShareAmountIn(ctx, sender, msg.PoolId, msg.TokenOutDenom, msg.ShareInAmount, msg.TokenOutMinAmount)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "gamm"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute("msg_type", "exit_swap_share_amount_in"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgExitSwapShareAmountInResponse{TokenOutAmount: tokenOutAmount}, nil
}
