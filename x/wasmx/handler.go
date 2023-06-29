package wasmx

import (
	"fmt"
	"runtime/debug"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	log "github.com/xlab/suplog"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/keeper"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {

	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (res *sdk.Result, err error) {
		defer Recover(&err) // nolint:all

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgExecuteContractCompat:
			res, err := msgServer.ExecuteContractCompat(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(badgestypes.ErrUnknownRequest,
				fmt.Sprintf("Unrecognized wasmx Msg type: %T", msg))
		}
	}
}

func Recover(err *error) { // nolint:all
	if r := recover(); r != nil {
		*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r) // nolint:all

		if e, ok := r.(error); ok {
			log.WithError(e).Errorln("wasmx msg handler panicked with an error")
			log.Debugln(string(debug.Stack()))
		} else {
			log.Errorln(r)
		}
	}
}
