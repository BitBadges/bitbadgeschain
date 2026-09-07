// evmtx signs and broadcasts one legacy (EIP-155) value transfer over
// JSON-RPC and waits for its receipt. It exists because the rehearsal
// container has no cast/foundry; it is built alongside the chain binaries so
// go-ethereum resolves from the same offline module cache.
//
//	evmtx --rpc http://127.0.0.1:8545 --key <hex privkey> --to 0x... --value-wei N
//
// Prints hash=, gasPrice=, gasUsed=, status= lines for the caller to parse.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func rpc(url, method string, params ...any) (json.RawMessage, error) {
	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Error != nil {
		return nil, fmt.Errorf("%s: %s", method, out.Error.Message)
	}
	return out.Result, nil
}

func rpcBig(url, method string, params ...any) (*big.Int, error) {
	raw, err := rpc(url, method, params...)
	if err != nil {
		return nil, err
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return hexutil.DecodeBig(s)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "evmtx:", err)
	os.Exit(1)
}

func main() {
	url := flag.String("rpc", "http://127.0.0.1:8545", "JSON-RPC endpoint")
	keyHex := flag.String("key", "", "sender private key (hex, with or without 0x)")
	to := flag.String("to", "", "recipient 0x address")
	valueStr := flag.String("value-wei", "", "value to send, in wei")
	flag.Parse()
	if *keyHex == "" || *to == "" || *valueStr == "" {
		fail(fmt.Errorf("--key, --to and --value-wei are required"))
	}
	value, ok := new(big.Int).SetString(*valueStr, 10)
	if !ok {
		fail(fmt.Errorf("bad --value-wei %q", *valueStr))
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(*keyHex, "0x"))
	if err != nil {
		fail(err)
	}
	from := crypto.PubkeyToAddress(key.PublicKey)

	chainID, err := rpcBig(*url, "eth_chainId")
	if err != nil {
		fail(err)
	}
	nonce, err := rpcBig(*url, "eth_getTransactionCount", from.Hex(), "pending")
	if err != nil {
		fail(err)
	}
	gasPrice, err := rpcBig(*url, "eth_gasPrice")
	if err != nil {
		fail(err)
	}
	if floor := big.NewInt(10_000_000_000); gasPrice.Cmp(floor) < 0 { // 10 gwei
		gasPrice = floor
	}
	// Round up to a whole ubadge per gas (1e9 wei on this 9-decimal chain) so
	// gasUsed*gasPrice is an integer number of ubadge and callers can check
	// bank balances exactly.
	ubadge := big.NewInt(1_000_000_000)
	if rem := new(big.Int).Mod(gasPrice, ubadge); rem.Sign() != 0 {
		gasPrice.Add(gasPrice, new(big.Int).Sub(ubadge, rem))
	}

	dest := common.HexToAddress(*to)
	tx := types.NewTx(&types.LegacyTx{Nonce: nonce.Uint64(), GasPrice: gasPrice, Gas: 21000, To: &dest, Value: value})
	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), key)
	if err != nil {
		fail(err)
	}
	raw, err := signed.MarshalBinary()
	if err != nil {
		fail(err)
	}
	hashRaw, err := rpc(*url, "eth_sendRawTransaction", hexutil.Encode(raw))
	if err != nil {
		fail(err)
	}
	var hash string
	_ = json.Unmarshal(hashRaw, &hash)
	fmt.Printf("hash=%s\nfrom=%s\nchainId=%s\ngasPrice=%s\n", hash, from.Hex(), chainID, gasPrice)

	for i := 0; i < 60; i++ {
		rcpt, err := rpc(*url, "eth_getTransactionReceipt", hash)
		if err == nil && string(rcpt) != "null" {
			var r struct {
				Status  string `json:"status"`
				GasUsed string `json:"gasUsed"`
				Block   string `json:"blockNumber"`
			}
			if err := json.Unmarshal(rcpt, &r); err != nil {
				fail(err)
			}
			gasUsed, _ := hexutil.DecodeBig(r.GasUsed)
			fmt.Printf("gasUsed=%s\nblock=%s\nstatus=%s\n", gasUsed, r.Block, r.Status)
			if r.Status != "0x1" {
				os.Exit(2)
			}
			return
		}
		time.Sleep(time.Second)
	}
	fail(fmt.Errorf("no receipt for %s after 60s", hash))
}
