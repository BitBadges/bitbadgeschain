package cli

import "fmt"

const (
	docsBaseURL    = "https://docs.bitbadges.io"
	protoBaseURL   = "https://github.com/BitBadges/bitbadgeschain/blob/master/proto/tokenization"
	repoBaseURL    = "https://github.com/BitBadges/bitbadgeschain/blob/master"
)

// msgDocLinks maps CLI command names to their docs page paths and proto file.
var msgDocLinks = map[string]struct {
	docsPath  string
	protoFile string
}{
	"create-collection":        {"x-tokenization/messages/msg-create-collection", "tx.proto"},
	"universal-update-collection": {"x-tokenization/messages/msg-universal-update-collection", "tx.proto"},
	"update-collection":        {"x-tokenization/messages/msg-update-collection", "tx.proto"},
	"delete-collection":        {"x-tokenization/messages/msg-delete-collection", "tx.proto"},
	"transfer-tokens":          {"x-tokenization/messages/msg-transfer-tokens", "tx.proto"},
	"set-incoming-approval":    {"x-tokenization/messages/msg-set-incoming-approval", "tx.proto"},
	"set-outgoing-approval":    {"x-tokenization/messages/msg-set-outgoing-approval", "tx.proto"},
	"delete-incoming-approval": {"x-tokenization/messages/msg-delete-incoming-approval", "tx.proto"},
	"delete-outgoing-approval": {"x-tokenization/messages/msg-delete-outgoing-approval", "tx.proto"},
	"set-collection-approvals": {"x-tokenization/messages/msg-set-collection-approvals", "tx.proto"},
	"set-collection-metadata":  {"x-tokenization/messages/msg-set-collection-metadata", "tx.proto"},
	"set-token-metadata":       {"token-standard/messages/msgsettokenmetadata", "tx.proto"},
	"set-custom-data":          {"x-tokenization/messages/msg-set-custom-data", "tx.proto"},
	"set-manager":              {"x-tokenization/messages/msg-set-manager", "tx.proto"},
	"set-standards":            {"x-tokenization/messages/msg-set-standards", "tx.proto"},
	"set-valid-token-ids":      {"x-tokenization/messages/msg-set-valid-token-ids", "tx.proto"},
	"set-is-archived":          {"x-tokenization/messages/msg-set-is-archived", "tx.proto"},
	"purge-approvals":          {"x-tokenization/messages/msg-purge-approvals", "tx.proto"},
	"create-address-lists":     {"x-tokenization/messages/msg-create-address-lists", "tx.proto"},
	"create-dynamic-store":     {"x-tokenization/messages/msg-create-dynamic-store", "tx.proto"},
	"update-dynamic-store":     {"x-tokenization/messages/msg-update-dynamic-store", "tx.proto"},
	"delete-dynamic-store":     {"x-tokenization/messages/msg-delete-dynamic-store", "tx.proto"},
	"set-dynamic-store-value":  {"x-tokenization/messages/msg-set-dynamic-store-value", "tx.proto"},
	"cast-vote":                {"", "tx.proto"},
	"update-user-approved-transfers": {"x-tokenization/messages/msg-update-user-approvals", "tx.proto"},
	// Aliases for set-set* CLI command names
	"set-setcollectionapprovals": {"x-tokenization/messages/msg-set-collection-approvals", "tx.proto"},
	"set-setcollectionmetadata":  {"x-tokenization/messages/msg-set-collection-metadata", "tx.proto"},
	"set-setcustomdata":          {"x-tokenization/messages/msg-set-custom-data", "tx.proto"},
	"set-setisarchived":          {"x-tokenization/messages/msg-set-is-archived", "tx.proto"},
	"set-setstandards":           {"x-tokenization/messages/msg-set-standards", "tx.proto"},
	"set-settokenmetadata":       {"token-standard/messages/msgsettokenmetadata", "tx.proto"},
}

// queryDocLinks maps CLI query command names to their docs page paths.
var queryDocLinks = map[string]struct {
	docsPath  string
	protoFile string
}{
	"collection":               {"x-tokenization/queries/get-collection", "query.proto"},
	"collection-stats":         {"x-tokenization/queries/get-collection", "query.proto"},
	"balance":                  {"x-tokenization/queries/get-balance", "query.proto"},
	"balance-for-token":        {"x-tokenization/queries/get-balance", "query.proto"},
	"address-list":             {"x-tokenization/queries/get-address-list", "query.proto"},
	"approvals-trackers":       {"x-tokenization/queries/get-approval-tracker", "query.proto"},
	"num-used-for-merkle-challenge": {"x-tokenization/queries/get-challenge-tracker", "query.proto"},
	"num-used-for-eth-signature-challenge": {"x-tokenization/queries/get-eth-signature-tracker", "query.proto"},
	"dynamic-store":            {"x-tokenization/queries/get-dynamic-store", "query.proto"},
	"dynamic-store-value":      {"x-tokenization/queries/get-dynamic-store-value", "query.proto"},
	"wrappable-balances":       {"", "query.proto"},
}

// MsgHelpLinks returns help text with documentation links for a tx command.
func MsgHelpLinks(cmdName string) string {
	base := "Accepts JSON either inline or from a file path. If the argument is a valid file path, it will read the JSON from that file."

	links, ok := msgDocLinks[cmdName]
	if !ok {
		return base + schemaHelpFooter("tx.proto")
	}

	return base + schemaHelpFooter(links.protoFile) + docsLink(links.docsPath)
}

// QueryHelpLinks returns help text with documentation links for a query command.
func QueryHelpLinks(cmdName string) string {
	base := ""

	links, ok := queryDocLinks[cmdName]
	if !ok {
		return base + schemaHelpFooter("query.proto")
	}

	return base + schemaHelpFooter(links.protoFile) + docsLink(links.docsPath)
}

func schemaHelpFooter(protoFile string) string {
	return fmt.Sprintf(`

Schema & Documentation:
  Proto definition:  %s/%s
  Full OpenAPI spec: %s/docs/static/openapi.yml
  SDK CLI docs:      bitbadgeschaind cli docs messages (if SDK CLI is installed)`,
		protoBaseURL, protoFile, repoBaseURL)
}

func docsLink(docsPath string) string {
	if docsPath == "" {
		return ""
	}
	return fmt.Sprintf(`
  Documentation:     %s/%s`, docsBaseURL, docsPath)
}
