package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-proto/anyutil"
	"github.com/cosmos/cosmos-sdk/codec"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	msgv1 "cosmossdk.io/api/cosmos/msg/v1"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ReplaceCreatorField(cdc codec.Codec, msg *codectypes.Any, creatorAddress string) (*codectypes.Any, error) {
	msgV2, err := anyutil.Unpack(&anypb.Any{
		TypeUrl: msg.TypeUrl,
		Value:   msg.Value,
	}, cdc.InterfaceRegistry().SigningContext().FileResolver(), nil)
	if err != nil {
		return nil, err
	}

	descriptor := msgV2.ProtoReflect().Descriptor()
	signersFields := proto.GetExtension(descriptor.Options(), msgv1.E_Signer).([]string)
	if len(signersFields) == 0 {
		return nil, fmt.Errorf("no cosmos.msg.v1.signer option found for message %s; use DefineCustomGetSigners to specify a custom getter", descriptor.FullName())
	}

	for _, signersField := range signersFields {
		signerFieldDescriptor := msgV2.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(signersField))

		msgV2.ProtoReflect().Set(signerFieldDescriptor, protoreflect.ValueOf(creatorAddress))
	}

	//Little hacky but we encode to bytes then encode to a cosmos proto Any and unpack to sdk.Msg
	protoMsg := msgV2.ProtoReflect().Interface()
	msgBytes, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified message: %v", err)
	}

	// Create an Any protobuf message
	anyMsg := &codectypes.Any{
		TypeUrl: msg.TypeUrl, // Use the original TypeUrl
		Value:   msgBytes,
	}

	return anyMsg, nil
}

func (k Keeper) ExecuteGenericMsgs(ctx sdk.Context, msgs []*codectypes.Any, creatorAddress string) error {
	for _, msg := range msgs {
		msg, err := ReplaceCreatorField(k.cdc, msg, creatorAddress)
		if err != nil {
			return err
		}

		var sudoedMsg sdk.Msg
		err = k.cdc.UnpackAny(msg, &sudoedMsg)
		if err != nil {
			return fmt.Errorf("failed to unpack Any: %v", err)
		}

		// make sure this account can send it
		signers, _, err := k.cdc.GetMsgV1Signers(sudoedMsg)
		if err != nil {
			return err
		}

		for _, acct := range signers {
			if sdk.AccAddress(acct).String() != creatorAddress {
				return fmt.Errorf("invalid signer: %s, expected: %s", sdk.AccAddress(acct).String(), creatorAddress)
			}
		}

		// check if the message implements the HasValidateBasic interface
		if m, ok := sudoedMsg.(sdk.HasValidateBasic); ok {
			if err := m.ValidateBasic(); err != nil {
				return errors.Wrapf(err, "invalid sudo-ed message: %s", err)
			}
		}

		handler := k.msgRouter.Handler(sudoedMsg)
		if handler == nil {
			return fmt.Errorf("message handler not found for %T", sudoedMsg)
		}

		//TODO: Do something with the responses?
		_, err = handler(ctx, sudoedMsg)
		if err != nil {
			return err
		}
	}

	return nil
}
