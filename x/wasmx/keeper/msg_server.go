package keeper

import (
	"context"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"
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

func (m msgServer) StoreCodeCompat(goCtx context.Context, msg *types.MsgStoreCodeCompat) (*types.MsgStoreCodeCompatResponse, error) {
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(&m.wasmKeeper)

	oMsg := &wasmtypes.MsgStoreCode{
		Sender:   msg.Sender,
		WASMByteCode: hexutil.MustDecode(msg.HexWasmByteCode),
	}

	res, err := wasmMsgServer.StoreCode(goCtx, oMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgStoreCodeCompatResponse{
		CodeId: sdk.NewUint(res.CodeID),
		Checksum: res.Checksum,
	}, nil
}
