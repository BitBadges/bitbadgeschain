//go:build test
// +build test

package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmcryptocodec "github.com/cosmos/evm/crypto/codec"
	evmhd "github.com/cosmos/evm/crypto/hd"

	"github.com/bitbadges/bitbadgeschain/app/params"
)

func init() {
	// SignArbitraryCmd derives bech32 addresses with the chain's
	// configured prefix; tests must set it before any address ops.
	params.SetAddressPrefixes()
}

// TestBuildADR36SignBytes_GoldenVector locks the canonical-JSON
// encoding to bytes precomputed via the verification spike (Go output
// piped into @keplr-wallet/cosmos's serializeSignDoc(makeADR36AminoSignDoc(...)),
// confirmed byte-equal and verifyADR36Amino-PASS). Any drift in Go's
// JSON marshaller or sdk.MustSortJSON ordering will fail this test.
func TestBuildADR36SignBytes_GoldenVector(t *testing.T) {
	const (
		signer  = "bb1mnyn7x24xj6vraxeeq56dfkxa009tvhg4wlgny"
		message = "BitBadges sign-arbitrary spike message — testing < > & escapes"
	)
	const expected = `{"account_number":"0","chain_id":"","fee":{"amount":[],"gas":"0"},"memo":"","msgs":[{"type":"sign/MsgSignData","value":{"data":"Qml0QmFkZ2VzIHNpZ24tYXJiaXRyYXJ5IHNwaWtlIG1lc3NhZ2Ug4oCUIHRlc3RpbmcgPCA+ICYgZXNjYXBlcw==","signer":"bb1mnyn7x24xj6vraxeeq56dfkxa009tvhg4wlgny"}}],"sequence":"0"}`

	got, err := buildADR36SignBytes(signer, message)
	if err != nil {
		t.Fatalf("buildADR36SignBytes: %v", err)
	}
	if string(got) != expected {
		t.Fatalf("sign bytes drifted from spike-verified golden vector\nwant: %s\ngot : %s", expected, string(got))
	}
}

// TestSignArbitrary_RoundTrip_Secp256k1 exercises the full command
// against an in-memory keyring with a true cosmos secp256k1 key, then
// verifies the produced signature locally via PubKey.VerifySignature.
func TestSignArbitrary_RoundTrip_Secp256k1(t *testing.T) {
	kr, addr, pub := newInMemoryKeyring(t, "tester", hd.Secp256k1)

	out := executeSignArbitrary(t, kr, []string{"tester", "hello world"}, "")

	var got struct {
		Format    string `json:"format"`
		Algo      string `json:"algo"`
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("parse output: %v\noutput: %s", err, out)
	}
	if got.Format != "adr36" || got.Algo != "secp256k1" {
		t.Fatalf("wrong format/algo: %+v", got)
	}
	if got.Address != addr.String() {
		t.Fatalf("wrong address: want %s got %s", addr.String(), got.Address)
	}

	signBytes, err := buildADR36SignBytes(addr.String(), "hello world")
	if err != nil {
		t.Fatalf("rebuild signBytes: %v", err)
	}
	sig, err := base64.StdEncoding.DecodeString(got.Signature)
	if err != nil {
		t.Fatalf("decode sig: %v", err)
	}
	if !pub.VerifySignature(signBytes, sig) {
		t.Fatalf("VerifySignature returned false — signing path produced an invalid signature")
	}
}

// TestSignArbitrary_RejectsEthSecp256k1 confirms v1 refuses
// eth_secp256k1 keys with a clear actionable error rather than
// silently producing a signature CosmosDriver will reject.
func TestSignArbitrary_RejectsEthSecp256k1(t *testing.T) {
	kr, _, _ := newInMemoryKeyring(t, "ethkey", evmhd.EthSecp256k1)

	cmd := SignArbitraryCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"ethkey", "hello"})

	cmd.SetContext(context.Background())
	clientCtx := client.Context{}.WithKeyring(kr)
	if err := client.SetCmdClientContext(cmd, clientCtx); err != nil {
		t.Fatalf("set cmd ctx: %v", err)
	}

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error rejecting eth_secp256k1, got success: %s", out.String())
	}
	if !strings.Contains(err.Error(), "eth_secp256k1") || !strings.Contains(err.Error(), "secp256k1") {
		t.Fatalf("error should mention both algos for clarity, got: %v", err)
	}
}

func TestReadMessage_Modes(t *testing.T) {
	t.Run("positional and file are mutually exclusive", func(t *testing.T) {
		cmd := SignArbitraryCmd()
		_ = cmd.Flags().Set(flagMessageFile, "/tmp/whatever")
		_, err := readMessage(cmd, []string{"keyname", "inline-msg"})
		if err == nil || !strings.Contains(err.Error(), "not both") {
			t.Fatalf("expected mutual-exclusion error, got: %v", err)
		}
	})
	t.Run("positional wins when alone", func(t *testing.T) {
		cmd := SignArbitraryCmd()
		got, err := readMessage(cmd, []string{"keyname", "the-message"})
		if err != nil || got != "the-message" {
			t.Fatalf("want the-message, got %q err %v", got, err)
		}
	})
	t.Run("file is read when provided", func(t *testing.T) {
		path := writeTempFile(t, "from-the-file")
		cmd := SignArbitraryCmd()
		_ = cmd.Flags().Set(flagMessageFile, path)
		got, err := readMessage(cmd, []string{"keyname"})
		if err != nil || got != "from-the-file" {
			t.Fatalf("want from-the-file, got %q err %v", got, err)
		}
	})
}

// ── helpers ─────────────────────────────────────────────────────────

func newInMemoryKeyring(t *testing.T, name string, algo keyring.SignatureAlgo) (keyring.Keyring, sdk.AccAddress, cryptotypes.PubKey) {
	t.Helper()
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	evmcryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kr := keyring.NewInMemory(cdc, evmhd.EthSecp256k1Option())

	rec, _, err := kr.NewMnemonic(
		name,
		keyring.English,
		sdk.GetConfig().GetFullBIP44Path(),
		"",
		algo,
	)
	if err != nil {
		t.Fatalf("NewMnemonic: %v", err)
	}
	addr, err := rec.GetAddress()
	if err != nil {
		t.Fatalf("GetAddress: %v", err)
	}
	pub, err := rec.GetPubKey()
	if err != nil {
		t.Fatalf("GetPubKey: %v", err)
	}
	return kr, addr, pub
}

// executeSignArbitrary runs the command end-to-end against the given
// keyring and returns stdout (the JSON envelope).
func executeSignArbitrary(t *testing.T, kr keyring.Keyring, args []string, output string) string {
	t.Helper()
	cmd := SignArbitraryCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs(args)
	if output != "" {
		_ = cmd.Flags().Set(flagOutputMode, output)
	}
	cmd.SetContext(context.Background())
	clientCtx := client.Context{}.WithKeyring(kr)
	if err := client.SetCmdClientContext(cmd, clientCtx); err != nil {
		t.Fatalf("set cmd ctx: %v", err)
	}
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v\noutput: %s", err, out.String())
	}
	return out.String()
}

func writeTempFile(t *testing.T, contents string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "msg.txt")
	if err := os.WriteFile(p, []byte(contents), 0o600); err != nil {
		t.Fatalf("write tmp file: %v", err)
	}
	return p
}
