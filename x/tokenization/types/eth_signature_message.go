package types

import (
	"strconv"
	"strings"
)

func encodeETHSignatureFields(fields ...string) string {
	var out strings.Builder
	for _, field := range fields {
		out.WriteString(strconv.Itoa(len(field)))
		out.WriteByte(':')
		out.WriteString(field)
	}
	return out.String()
}

func ETHSignatureChallengeMessage(chainID, nonce, initiator, collectionID, approver, level, approvalID, challengeID string) string {
	return "BitBadges ETH Signature Challenge v2\n" + encodeETHSignatureFields(chainID, nonce, initiator, collectionID, approver, level, approvalID, challengeID)
}

func ETHSignatureTrackerScope(collectionID, approver, level, approvalID, challengeID, nonce string) string {
	return "BitBadges ETH Signature Tracker v2\n" + encodeETHSignatureFields(collectionID, approver, level, approvalID, challengeID, nonce)
}
