package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// reBadgeIDString can be 3 ~ 60 characters long and support letters, followed by either
	// a letter, a number or a slash ('/') or a colon (':') or ('-').
	// reBadgeIDString = `[a-zA-Z][a-zA-Z0-9/:-]{2,60}`
	// reBadgeID       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reBadgeIDString))

	// URI must be a valid URI. Method <= 35 characters long. Path <= 1000 characters long.
	reUriString = `\w{0,35}:(\/?\/?)[^\s]{0,1000}`
	reUri       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reUriString))

	// Metadata hash must be <= 32 characters long.
	reMetadataString = `(\/?\/?)[^\s]{0,32}`
	reMetadata       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reMetadataString))
)

// ValidateURI returns whether the uri is valid
func ValidateURI(uri string) error {
	if !reUri.MatchString(uri) {
		return sdkerrors.Wrapf(ErrInvalidBadgeURI, "invalid uri: %s", uri)
	}
	return nil
}

func ValidateAddress(address string) error {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

func ValidateMetadata(hash string) error {
	if !reMetadata.MatchString(hash) {
		return sdkerrors.Wrapf(ErrInvalidBadgeHash, "invalid metadata hash: %s", hash)
	}
	return nil
}
