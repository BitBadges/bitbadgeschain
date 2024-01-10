package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgInstantiateContractCompat = "instantiateContractCompat"
)

func (msg MsgInstantiateContractCompat) Route() string {
	return RouterKey
}

func (msg MsgInstantiateContractCompat) Type() string {
	return TypeMsgInstantiateContractCompat
}

func (msg MsgInstantiateContractCompat) ValidateBasic() error {
	funds := sdk.Coins{}
	if msg.Funds != "0" {
		funds, _ = sdk.ParseCoinsNormalized(msg.Funds)
	}
	
	jsonMsgStr := "{}"
	bytesMsg := []byte(jsonMsgStr)

	oMsg := &wasmtypes.MsgInstantiateContract{
		Sender: msg.Sender,
		CodeID: msg.CodeId.Uint64(),
		Label: 	msg.Label,
		Funds: 	funds,
		Msg: 		bytesMsg,
	}

	if err := oMsg.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// Note ModuleCdc is Amino (see codec.go)
func (msg MsgInstantiateContractCompat) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

func (msg MsgInstantiateContractCompat) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
