package keeper

import (
	"strconv"
	"strings"
)

const (
	BalanceKeyDelimiter = "-"
)

type BalanceKeyDetails struct {
	badge_id    uint64
	subasset_id uint64
	account_num uint64
}

func GetBalanceKey(accountNumber uint64, id uint64, subasset_id uint64) string {
	badge_id_str := strconv.FormatUint(id, 10)
	subasset_id_str := strconv.FormatUint(subasset_id, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return account_num_str + BalanceKeyDelimiter + badge_id_str + BalanceKeyDelimiter + subasset_id_str
}

func GetDetailsFromBalanceKey(id string) BalanceKeyDetails {
	result := strings.Split(id, BalanceKeyDelimiter)
	account_num, _ := strconv.ParseUint(result[0], 10, 64)
	badge_id, _ := strconv.ParseUint(result[1], 10, 64)
	subasset_id, _ := strconv.ParseUint(result[2], 10, 64)

	return BalanceKeyDetails{
		account_num: account_num,
		badge_id:    badge_id,
		subasset_id: subasset_id,
	}
}

func GetManagerRequestKey(badgeId uint64, accountNumber uint64) string {
	badge_id_str := strconv.FormatUint(badgeId, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return badge_id_str + BalanceKeyDelimiter + account_num_str + BalanceKeyDelimiter
}
