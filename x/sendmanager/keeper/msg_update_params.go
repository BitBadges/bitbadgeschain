package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"

    "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Type assert addressCodec
	addressCodec, ok := k.addressCodec.(interface {
		StringToBytes(string) ([]byte, error)
		BytesToString([]byte) (string, error)
	})
	if !ok {
		return nil, errorsmod.Wrap(errorsmod.New("invalid", 1, "address codec not available"), "invalid address codec")
	}

	authority, err := addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	if err := req.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
