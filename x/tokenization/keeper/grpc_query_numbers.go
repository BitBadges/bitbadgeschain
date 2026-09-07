package keeper

import (
	sdkmath "cosmossdk.io/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseQueryUint(value, field string) (sdkmath.Uint, error) {
	parsed, err := sdkmath.ParseUint(value)
	if err != nil {
		return sdkmath.Uint{}, status.Error(codes.InvalidArgument, "invalid "+field)
	}
	return parsed, nil
}
