package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitbadges/bitbadgeschain/lifecycle"
	"io"
	"os"
)

func main() { os.Exit(run()) }
func run() (code int) {
	defer func() {
		if v := recover(); v != nil {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"version": 1, "passed": false, "error": fmt.Sprint(v)})
			code = 2
		}
	}()
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"protocolVersion": 1, "name": "bitbadges-lifecycle", "execution": "module", "chainCommit": lifecycle.ChainCommit})
		return 0
	}
	if len(os.Args) == 2 && os.Args[1] == "--schema" {
		fmt.Println(lifecycle.Schema)
		return 0
	}
	var err error
	var data []byte
	if len(os.Args) != 1 {
		err = fmt.Errorf("usage: bitbadges-lifecycle [--version|--schema]; read scenario from stdin")
	} else {
		data, err = io.ReadAll(io.LimitReader(os.Stdin, (8<<20)+1))
	}
	if err == nil {
		var result *lifecycle.Result
		result, err = lifecycle.RunJSON(data)
		if err == nil {
			if e := json.NewEncoder(os.Stdout).Encode(result); e != nil {
				return 2
			}
			if result.Passed {
				return 0
			}
			return 1
		}
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"version": 1, "passed": false, "error": err.Error()})
	return 2
}
