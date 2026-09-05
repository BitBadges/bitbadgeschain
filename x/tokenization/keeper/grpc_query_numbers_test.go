package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testutilkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestQueryMalformedNumbersReturnInvalidArgument(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	queries := map[string]func(string) error{
		"collection": func(s string) error {
			_, err := k.GetCollection(ctx, &types.QueryGetCollectionRequest{CollectionId: s})
			return err
		},
		"balance": func(s string) error {
			_, err := k.GetBalance(ctx, &types.QueryGetBalanceRequest{CollectionId: s})
			return err
		},
		"token collection": func(s string) error {
			_, err := k.GetBalanceForToken(ctx, &types.QueryGetBalanceForTokenRequest{CollectionId: s})
			return err
		},
		"token id": func(s string) error {
			_, err := k.GetBalanceForToken(ctx, &types.QueryGetBalanceForTokenRequest{CollectionId: "1", TokenId: s})
			return err
		},
		"token time": func(s string) error {
			_, err := k.GetBalanceForToken(ctx, &types.QueryGetBalanceForTokenRequest{CollectionId: "1", TokenId: "1", Time: s})
			return err
		},
		"store": func(s string) error {
			_, err := k.GetDynamicStore(ctx, &types.QueryGetDynamicStoreRequest{StoreId: s})
			return err
		},
		"store value": func(s string) error {
			_, err := k.GetDynamicStoreValue(ctx, &types.QueryGetDynamicStoreValueRequest{StoreId: s})
			return err
		},
		"eth tracker": func(s string) error {
			_, err := k.GetETHSignatureTracker(ctx, &types.QueryGetETHSignatureTrackerRequest{CollectionId: s})
			return err
		},
		"approval tracker": func(s string) error {
			_, err := k.GetApprovalTracker(ctx, &types.QueryGetApprovalTrackerRequest{CollectionId: s})
			return err
		},
		"challenge tracker": func(s string) error {
			_, err := k.GetChallengeTracker(ctx, &types.QueryGetChallengeTrackerRequest{CollectionId: s, LeafIndex: "0"})
			return err
		},
		"challenge leaf": func(s string) error {
			_, err := k.GetChallengeTracker(ctx, &types.QueryGetChallengeTrackerRequest{CollectionId: "1", LeafIndex: s})
			return err
		},
		"vote": func(s string) error {
			_, err := k.GetVote(ctx, &types.QueryGetVoteRequest{CollectionId: s})
			return err
		},
		"votes": func(s string) error {
			_, err := k.GetVotes(ctx, &types.QueryGetVotesRequest{CollectionId: s})
			return err
		},
	}
	for name, query := range queries {
		for _, value := range []string{"invalid", "-1", strings.Repeat("9", 80)} {
			t.Run(name+"/"+value, func(t *testing.T) {
				require.NotPanics(t, func() { require.Equal(t, codes.InvalidArgument, status.Code(query(value))) })
			})
		}
	}
}
