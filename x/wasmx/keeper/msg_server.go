package keeper

import (
	"context"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbadgeschain/x/wasmx/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) ExecuteContractCompat(goCtx context.Context, msg *types.MsgExecuteContractCompat) (*types.MsgExecuteContractCompatResponse, error) {
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(&m.wasmKeeper)

	funds := sdk.Coins{}
	if msg.Funds != "0" {
		funds, _ = sdk.ParseCoinsNormalized(msg.Funds)
	}

	oMsg := &wasmtypes.MsgExecuteContract{
		Sender:   msg.Sender,
		Contract: msg.Contract,
		Msg:      []byte(msg.Msg),
		Funds:    funds,
	}

	res, err := wasmMsgServer.ExecuteContract(goCtx, oMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgExecuteContractCompatResponse{
		Data: res.Data,
	}, nil
}

func (m msgServer) InstantiateContractCompat(goCtx context.Context, msg *types.MsgInstantiateContractCompat) (*types.MsgInstantiateContractCompatResponse, error) {
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(&m.wasmKeeper)
	funds := sdk.Coins{}
	if msg.Funds != "0" {
		funds, _ = sdk.ParseCoinsNormalized(msg.Funds)
	}

	jsonMsgStr := "{}"
	bytesMsg := []byte(jsonMsgStr)

	oMsg := &wasmtypes.MsgInstantiateContract{
		Sender: msg.Sender,
		CodeID: msg.CodeId.Uint64(),
		Label:  msg.Label,
		Funds:  funds,
		Msg:    bytesMsg,
	}

	res, err := wasmMsgServer.InstantiateContract(goCtx, oMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgInstantiateContractCompatResponse{
		Address: res.Address,
		Data:    res.Data,
	}, nil
}
