package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	TypeMsgStoreCodeCompat = "storeCodeCompat"
)

func (msg MsgStoreCodeCompat) Route() string {
	return RouterKey
}

func (msg MsgStoreCodeCompat) Type() string {
	return TypeMsgStoreCodeCompat
}

func (msg MsgStoreCodeCompat) ValidateBasic() error {
	oMsg := &wasmtypes.MsgStoreCode{
		Sender:   msg.Sender,
		WASMByteCode: hexutil.MustDecode(msg.HexWasmByteCode),
		InstantiatePermission: &wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
			Addresses: []string{},
		},
	}

	if err := oMsg.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// Note ModuleCdc is Amino (see codec.go)
func (msg MsgStoreCodeCompat) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

func (msg MsgStoreCodeCompat) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
