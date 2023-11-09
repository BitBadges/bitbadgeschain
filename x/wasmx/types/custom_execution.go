package types

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type InjectiveExecMsg struct {
	ExecutionData ExecutionData `json:"injective_exec"`
}

type ExecutionData struct {
	Origin string      `json:"origin"`
	Name   string      `json:"name"`
	Args   interface{} `json:"args"`
}

func NewInjectiveExecMsg(origin sdk.AccAddress, data string) (*InjectiveExecMsg, error) {
	var e ExecutionData
	if err := json.Unmarshal([]byte(data), &e); err != nil {
		return nil, sdkerrors.Wrap(err, data)
	}

	if e.Origin == "" && origin.Empty() {
		return nil, sdkerrors.Wrap(badgestypes.ErrInvalidAddress, "origin address is empty")
	}

	// override e.Origin for safety
	e.Origin = origin.String()

	return &InjectiveExecMsg{
		ExecutionData: e,
	}, nil
}
