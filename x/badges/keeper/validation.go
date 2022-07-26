package keeper

import (
	"fmt"
	"regexp"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// reBadgeIDString can be 3 ~ 100 characters long and support letters, followed by either
	// a letter, a number or a slash ('/') or a colon (':') or ('-').
	reBadgeIDString = `[a-zA-Z][a-zA-Z0-9/:-]{2,100}`
	reBadgeID       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reBadgeIDString))
)

// ValidateBadgeID returns whether the Badge id is valid
func ValidateBadgeID(id string) error {
	if !reBadgeID.MatchString(id) {
		return sdkerrors.Wrapf(ErrInvalidBadgeID, "invalid badge id: %s", id)
	}
	return nil
}