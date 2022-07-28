package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	creator_account := k.accountKeeper.GetAccount(ctx, creator)
	if creator_account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", creator)
	}
	creator_account_num := creator_account.GetAccountNumber()

	balance_id := GetFullSubassetID(
		creator_account_num,
		msg.BadgeId,
		msg.SubbadgeId,
	)


	err = k.Keeper.HandlePendingTransfer(ctx, msg.Accept, balance_id, msg.PendingId)
	if err != nil {
		return nil, err
	}

	_ = ctx
	return &types.MsgHandlePendingTransferResponse{
		Message: "Success!",
	}, nil
}
