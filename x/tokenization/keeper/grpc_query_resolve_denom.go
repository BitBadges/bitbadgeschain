package keeper

import (
	"context"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ResolveDenom resolves a wrapped (`badges:<collectionId>:<baseDenom>`) or alias
// (`badgeslp:<collectionId>:<baseDenom>`) denom to its collection ID and matching path metadata.
func (k Keeper) ResolveDenom(goCtx context.Context, req *types.QueryResolveDenomRequest) (*types.QueryResolveDenomResponse, error) {
	if req == nil || req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	denom := req.Denom

	isWrapper := strings.HasPrefix(denom, WrappedDenomPrefix)
	isAlias := strings.HasPrefix(denom, AliasDenomPrefix)
	if !isWrapper && !isAlias {
		return nil, status.Errorf(codes.InvalidArgument, "denom %q is not a tokenization wrapped (%s) or alias (%s) denom", denom, WrappedDenomPrefix, AliasDenomPrefix)
	}

	collectionIdUint64, err := ParseDenomCollectionId(denom)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid tokenization denom %q: %v", denom, err)
	}

	baseDenom, err := ParseDenomPath(denom)
	if err != nil || baseDenom == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid tokenization denom path %q", denom)
	}

	collectionId := sdkmath.NewUint(collectionIdUint64)
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "collection %s not found", collectionId.String())
	}

	if isAlias {
		for _, path := range collection.AliasPaths {
			if path != nil && path.Denom == baseDenom {
				return &types.QueryResolveDenomResponse{
					PathType:     "alias",
					CollectionId: collectionId,
					BaseDenom:    baseDenom,
					AliasPath:    path,
				}, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "alias path %q not found in collection %s", baseDenom, collectionId.String())
	}

	// Wrapper denom (`badges:`): check exact match first, then `{id}` placeholder paths.
	for _, path := range collection.CosmosCoinWrapperPaths {
		if path != nil && path.Denom == baseDenom {
			return &types.QueryResolveDenomResponse{
				PathType:     "wrapper",
				CollectionId: collectionId,
				BaseDenom:    baseDenom,
				WrapperPath:  path,
			}, nil
		}
	}

	for _, path := range collection.CosmosCoinWrapperPaths {
		if path == nil || !path.AllowOverrideWithAnyValidToken || !strings.Contains(path.Denom, "{id}") {
			continue
		}
		matchedId, ok := matchWrapperPlaceholderDenom(path.Denom, baseDenom, collection.ValidTokenIds)
		if ok {
			return &types.QueryResolveDenomResponse{
				PathType:       "wrapper",
				CollectionId:   collectionId,
				BaseDenom:      baseDenom,
				MatchedTokenId: matchedId,
				WrapperPath:    path,
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "wrapper path %q not found in collection %s", baseDenom, collectionId.String())
}

// matchWrapperPlaceholderDenom matches a concrete baseDenom against a wrapper path template
// containing "{id}" (e.g. "token-{id}" matching "token-42") and verifies the extracted token ID
// is a valid uint within validTokenIds.
func matchWrapperPlaceholderDenom(template string, baseDenom string, validTokenIds []*types.UintRange) (string, bool) {
	firstIdx := strings.Index(template, "{id}")
	lastIdx := strings.LastIndex(template, "{id}")
	if firstIdx < 0 {
		return "", false
	}

	prefix := template[:firstIdx]
	suffix := template[lastIdx+len("{id}"):]
	if !strings.HasPrefix(baseDenom, prefix) || !strings.HasSuffix(baseDenom, suffix) {
		return "", false
	}
	if len(baseDenom) < len(prefix)+len(suffix) {
		return "", false
	}

	// Extract the candidate token ID for the first placeholder occurrence.
	var candidate string
	if firstIdx == lastIdx {
		candidate = baseDenom[len(prefix) : len(baseDenom)-len(suffix)]
	} else {
		// Multiple {id} occurrences: take the maximal numeric prefix right after the first prefix.
		rest := baseDenom[len(prefix):]
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}
		candidate = rest[:i]
	}

	if candidate == "" {
		return "", false
	}

	tokenId, err := sdkmath.ParseUint(candidate)
	if err != nil || tokenId.IsZero() {
		return "", false
	}
	// Reject non-canonical representations like "01" that differ from tokenId.String().
	if tokenId.String() != candidate {
		return "", false
	}

	// Confirm substituting {id} reproduces baseDenom exactly.
	if strings.ReplaceAll(template, "{id}", candidate) != baseDenom {
		return "", false
	}

	// Verify tokenId falls within collection.ValidTokenIds.
	inRange, err := types.SearchUintRangesForUint(tokenId, validTokenIds)
	if err != nil || !inRange {
		return "", false
	}

	return candidate, true
}
