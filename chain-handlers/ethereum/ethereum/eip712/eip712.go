package eip712

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func WrapTxToTypedData(
	chainID uint64,
	data []byte,
	chain string,
) (apitypes.TypedData, error) {
	messagePayload, err := CreateEIP712MessagePayload(data, chain)
	message := messagePayload.Message
	if err != nil {
		return apitypes.TypedData{}, err
	}

	types, err := CreateEIP712Types(messagePayload)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	domain := CreateEIP712Domain(chainID)

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: txField,
		Domain:      domain,
		Message:     message,
	}

	return typedData, nil
}
