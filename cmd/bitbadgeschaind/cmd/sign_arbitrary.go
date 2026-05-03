package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/cobra"
)

const (
	flagMessageFile = "message-file"
	flagOutputMode  = "output"

	algoSecp256k1 = "secp256k1"
	formatADR36   = "adr36"
)

// SignArbitraryCmd produces an ADR-36 signature over an arbitrary
// message using a key from the local keyring. Output is the JSON
// envelope the bitbadges-cli `auth verify` flow expects:
//
//	{ format, address, pubKey, signature, message, algo }
//
// The signed bytes are byte-equal to Keplr's serializeSignDoc(
// makeADR36AminoSignDoc(...)) and verify via @keplr-wallet/cosmos
// verifyADR36Amino.
//
// v1 supports secp256k1 keys only; eth_secp256k1 keys are rejected
// with a clear error pointing at a future --format eip191 flag.
func SignArbitraryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-arbitrary <key-name-or-address> [message]",
		Short: "Sign an arbitrary message via ADR-36 (offline; for CLI auth flows)",
		Long: `Sign an arbitrary message with a local keyring key, producing an
ADR-36-formatted signature suitable for posting to BitBadges
indexer's /api/v0/auth/verify endpoint.

This is a pure offline operation — no network calls, no chain state,
no transaction is broadcast. It exists to bridge the BitBadges CLI
auth flow ('bitbadges-cli auth verify --signature ...') with keys
held in the chain binary's keyring, without exposing the raw private
key.

Message is read from (in order): positional argument, --message-file,
or stdin when piped. Reading from a TTY without a message source
prints usage rather than hanging.

Examples:
  bitbadgeschaind sign-arbitrary mykey "auth challenge text"
  echo -n "auth challenge text" | bitbadgeschaind sign-arbitrary mykey
  bitbadgeschaind sign-arbitrary mykey --message-file challenge.txt`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			message, err := readMessage(cmd, args)
			if err != nil {
				return err
			}

			record, err := loadKey(clientCtx.Keyring, args[0])
			if err != nil {
				return fmt.Errorf("load key %q: %w", args[0], err)
			}

			pub, err := record.GetPubKey()
			if err != nil {
				return fmt.Errorf("get pubkey: %w", err)
			}
			if pub.Type() != algoSecp256k1 {
				return fmt.Errorf(
					"key %q has algo %q; only %q is supported in v1 "+
						"(eth_secp256k1 / EIP-191 personal_sign support is a planned --format eip191 follow-up). "+
						"Create a compatible key with: bitbadgeschaind keys add <name> --key-type secp256k1",
					args[0], pub.Type(), algoSecp256k1,
				)
			}

			addr, err := record.GetAddress()
			if err != nil {
				return fmt.Errorf("derive address: %w", err)
			}
			bech := addr.String()

			signBytes, err := buildADR36SignBytes(bech, message)
			if err != nil {
				return fmt.Errorf("build sign-doc: %w", err)
			}

			sig, _, err := clientCtx.Keyring.SignByAddress(
				addr,
				signBytes,
				signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
			)
			if err != nil {
				return fmt.Errorf("sign: %w", err)
			}

			outputMode, _ := cmd.Flags().GetString(flagOutputMode)
			return emitOutput(cmd.OutOrStdout(), outputMode, map[string]any{
				"format":    formatADR36,
				"algo":      pub.Type(),
				"address":   bech,
				"pubKey":    base64.StdEncoding.EncodeToString(pub.Bytes()),
				"signature": base64.StdEncoding.EncodeToString(sig),
				"message":   message,
			})
		},
	}

	cmd.Flags().String(flagMessageFile, "", "Read message from file (mutually exclusive with positional message and stdin)")
	cmd.Flags().String(flagOutputMode, "json", "Output mode: json | raw (raw prints only the base64 signature)")

	flags.AddKeyringFlags(cmd.Flags())
	cmd.Flags().String(flags.FlagHome, "", "The application home directory")
	return cmd
}

// buildADR36SignBytes constructs the canonical ADR-36 StdSignDoc and
// returns its serialized form (sorted keys, HTML-escaped <>&). Bytes
// are byte-equal to Keplr's serializeSignDoc(makeADR36AminoSignDoc(...))
// — verified by spike against @keplr-wallet/cosmos.
func buildADR36SignBytes(signer, message string) ([]byte, error) {
	doc := map[string]any{
		"account_number": "0",
		"chain_id":       "",
		"fee":            map[string]any{"amount": []any{}, "gas": "0"},
		"memo":           "",
		"msgs": []any{
			map[string]any{
				"type": "sign/MsgSignData",
				"value": map[string]any{
					"signer": signer,
					"data":   base64.StdEncoding.EncodeToString([]byte(message)),
				},
			},
		},
		"sequence": "0",
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return sdk.MustSortJSON(raw), nil
}

// loadKey accepts either a key name or a bech32 address.
func loadKey(kr cryptokeyring.Keyring, ref string) (*cryptokeyring.Record, error) {
	if addr, err := sdk.AccAddressFromBech32(ref); err == nil {
		return kr.KeyByAddress(addr)
	}
	return kr.Key(ref)
}

// readMessage resolves the message from (in priority order) the
// positional argument, --message-file, then stdin (only when stdin is
// not a TTY). Multiple sources is an error.
func readMessage(cmd *cobra.Command, args []string) (string, error) {
	hasPositional := len(args) >= 2
	filePath, _ := cmd.Flags().GetString(flagMessageFile)

	switch {
	case hasPositional && filePath != "":
		return "", errors.New("specify either a positional message or --message-file, not both")
	case hasPositional:
		return args[1], nil
	case filePath != "":
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("read --message-file %s: %w", filePath, err)
		}
		return string(raw), nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("stat stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("no message provided: pass as positional arg, --message-file, or pipe via stdin")
	}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return string(raw), nil
}

func emitOutput(w io.Writer, mode string, payload map[string]any) error {
	switch mode {
	case "raw":
		_, err := fmt.Fprintln(w, payload["signature"])
		return err
	case "", "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	default:
		return fmt.Errorf("unknown --output mode %q (expected json|raw)", mode)
	}
}
