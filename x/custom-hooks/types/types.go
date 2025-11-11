package types

import (
	"encoding/json"
)

// HookData represents the data that can be executed via IBC hooks
type HookData struct {
	SwapAndAction *SwapAndAction `json:"swap_and_action,omitempty"`
}

// SwapAndAction represents the Skip-style format for swap and action
type SwapAndAction struct {
	UserSwap         *UserSwap       `json:"user_swap,omitempty"`
	MinAsset         *MinAsset       `json:"min_asset,omitempty"`
	TimeoutTimestamp *uint64         `json:"timeout_timestamp,omitempty"`
	PostSwapAction   *PostSwapAction `json:"post_swap_action,omitempty"`
}

// UserSwap contains the swap operation
type UserSwap struct {
	SwapExactAssetIn *SwapExactAssetIn `json:"swap_exact_asset_in,omitempty"`
}

// SwapExactAssetIn contains the swap operations
type SwapExactAssetIn struct {
	SwapVenueName string      `json:"swap_venue_name,omitempty"`
	Operations    []Operation `json:"operations"`
}

// Operation represents a single swap operation
type Operation struct {
	Pool     string `json:"pool"` // Pool ID as string
	DenomIn  string `json:"denom_in"`
	DenomOut string `json:"denom_out"`
}

// MinAsset represents the minimum asset output
type MinAsset struct {
	Native *NativeAsset `json:"native,omitempty"`
}

// NativeAsset represents a native asset
type NativeAsset struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// PostSwapAction contains actions to execute after swap
type PostSwapAction struct {
	IBCTransfer *IBCTransferInfo `json:"ibc_transfer,omitempty"`
	Transfer    *TransferInfo    `json:"transfer,omitempty"`
}

// TransferInfo contains local transfer information
type TransferInfo struct {
	ToAddress string `json:"to_address"`
}

// IBCTransferInfo contains IBC transfer information
type IBCTransferInfo struct {
	IBCInfo *IBCInfo `json:"ibc_info,omitempty"`
}

// IBCInfo contains IBC transfer details
type IBCInfo struct {
	SourceChannel  string `json:"source_channel"`
	Receiver       string `json:"receiver"`
	Memo           string `json:"memo,omitempty"`
	RecoverAddress string `json:"recover_address,omitempty"`
}

// ParseHookDataFromMemo parses hook data from IBC memo
func ParseHookDataFromMemo(memo string) (*HookData, error) {
	if len(memo) == 0 {
		return nil, nil
	}

	var memoObj map[string]interface{}
	if err := json.Unmarshal([]byte(memo), &memoObj); err != nil {
		return nil, err
	}

	// Check for swap_and_action key
	if swapAndActionRaw, ok := memoObj["swap_and_action"]; ok {
		swapAndActionBytes, err := json.Marshal(swapAndActionRaw)
		if err != nil {
			return nil, err
		}

		var swapAndAction SwapAndAction
		if err := json.Unmarshal(swapAndActionBytes, &swapAndAction); err != nil {
			return nil, err
		}

		return &HookData{
			SwapAndAction: &swapAndAction,
		}, nil
	}

	return nil, nil
}
