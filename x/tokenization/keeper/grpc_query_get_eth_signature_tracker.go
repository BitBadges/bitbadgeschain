package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries how many times a signature has been used for an ETH signature challenge
func (k Keeper) GetETHSignatureTracker(goCtx context.Context, req *types.QueryGetETHSignatureTrackerRequest) (*types.QueryGetETHSignatureTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId := sdkmath.NewUintFromString(req.CollectionId)

	// Construct the signature key using the same pattern as in challenges.go
	signatureKey := ConstructETHSignatureTrackerKey(collectionId, req.ApproverAddress, req.ApprovalLevel, req.ApprovalId, req.ChallengeTrackerId, req.Signature)

	numUsed, exists := k.GetETHSignatureTrackerFromStore(ctx, signatureKey)
	if !exists {
		numUsed = sdkmath.NewUint(0)
	}

	return &types.QueryGetETHSignatureTrackerResponse{
		NumUsed: numUsed.String(),
	}, nil
}
