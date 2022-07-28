package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) SetApproval(goCtx context.Context, msg *types.MsgSetApproval) (*types.MsgSetApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	account := k.accountKeeper.GetAccount(ctx, address)
	if account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", account)
	}
	account_num := account.GetAccountNumber()

	signer_account := k.accountKeeper.GetAccount(ctx, signer)
	if signer_account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", signer)
	}
	signer_account_num := signer_account.GetAccountNumber()

	balance_id := GetFullSubassetID(signer_account_num, msg.BadgeId, msg.SubbadgeId)

	err = k.Keeper.SetApproval(ctx, balance_id, msg.Amount, account_num)
	if err != nil {
		return nil, err
	}

	_ = ctx

	return &types.MsgSetApprovalResponse{
		Message: "Success!",
	}, nil
}
