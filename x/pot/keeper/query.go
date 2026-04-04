package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// QueryParamsResponse is the response for the Params query.
type QueryParamsResponse struct {
	Params types.Params `json:"params"`
}

// QueryCredentialStatusRequest is the request for querying a validator's credential status.
type QueryCredentialStatusRequest struct {
	ValidatorAddress string `json:"validator_address"`
}

// QueryCredentialStatusResponse is the response for the CredentialStatus query.
type QueryCredentialStatusResponse struct {
	HasCredential      bool   `json:"has_credential"`
	CredentialBalance  uint64 `json:"credential_balance"`
	IsValidator        bool   `json:"is_validator"`
	IsBonded           bool   `json:"is_bonded"`
	IsJailed           bool   `json:"is_jailed"`
	IsComplianceJailed bool   `json:"is_compliance_jailed"`
}

// QueryParams returns the current module parameters.
// Once proto/gRPC is added, this will become a proper gRPC query handler.
func (k Keeper) QueryParams(ctx sdk.Context) QueryParamsResponse {
	return QueryParamsResponse{Params: k.GetParams(ctx)}
}

// QueryCredentialStatus returns the credential status for a given validator address.
// This is useful for dashboards and debugging.
func (k Keeper) QueryCredentialStatus(ctx sdk.Context, req QueryCredentialStatusRequest) QueryCredentialStatusResponse {
	params := k.GetParams(ctx)
	resp := QueryCredentialStatusResponse{}

	if !params.IsEnabled() {
		return resp
	}

	// The query accepts either a validator operator address or an account address.
	// We cannot call the old staking keeper directly — use the abstract interface.
	// For now, the query provides credential balance and compliance-jailed status.
	// Detailed staking info (IsBonded, IsJailed) is only available when the
	// ValidatorSetKeeper can look up by the provided address.

	// Check credential balance using the provided address directly.
	balance, err := k.tokenizationKeeper.GetCredentialBalance(
		ctx,
		params.CredentialCollectionId,
		params.CredentialTokenId,
		req.ValidatorAddress,
	)
	if err != nil {
		return resp
	}

	resp.CredentialBalance = balance
	resp.HasCredential = balance >= params.MinCredentialBalance
	resp.IsValidator = true // Assume valid if queried; full validation requires staking context.
	return resp
}
