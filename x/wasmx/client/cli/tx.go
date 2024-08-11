package cli

import (
	"fmt"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"bitbadgeschain/x/wasmx/types"

	cliflags "github.com/cosmos/cosmos-sdk/client/flags"
)

const (
	FlagContractGasLimit      = "contract-gas-limit"
	FlagContractGasPrice      = "contract-gas-price"
	FlagPinContract           = "pin-contract"
	FlagContractAddress       = "contract-address"
	FlagContractAdmin         = "contract-admin"
	FlagMigrationAllowed      = "migration-allowed"
	FlagCodeId                = "code-id"
	FlagBatchUploadProposal   = "batch-upload-proposal"
	FlagContractFiles         = "contract-files"
	FlagContractAddresses     = "contract-addresses"
	flagAmount                = "amount"
	FlagContractCallerAddress = "contract-caller-address"
	FlagContractExecMsg       = "contract-exec-msg"
)

// NewTxCmd returns a root CLI command handler for certain modules/wasmx transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Wasmx transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		ExecuteContractCompatCmd(),
	)
	return txCmd
}

// ExecuteContractCompatCmd will instantiate a contract from previously uploaded code.
func ExecuteContractCompatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "execute-compat [contract_addr_bech32] [json_encoded_send_args] --amount [coins,optional]",
		Short:   "Execute a command on a wasm contract",
		Aliases: []string{"run-compat", "call-compat", "exec-compat", "ex-compat", "e-compat"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := parseExecuteCompatArgs(args[0], args[1], clientCtx.GetFromAddress(), cmd.Flags())
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
		SilenceUsage: true,
	}

	cmd.Flags().String(flagAmount, "", "Coins to send to the contract along with command")
	cliflags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseExecuteCompatArgs(
	contractAddr string, execMsg string, sender sdk.AccAddress, fs *pflag.FlagSet,
) (types.MsgExecuteContractCompat, error) {
	fundsStr, err := fs.GetString(flagAmount)
	if err != nil {
		return types.MsgExecuteContractCompat{}, fmt.Errorf("amount: %w", err)
	}

	return types.MsgExecuteContractCompat{
		Sender:   sender.String(),
		Contract: contractAddr,
		Funds:    fundsStr,
		Msg:      execMsg,
	}, nil
}

func getWasmFile(contractSrc string) ([]byte, error) {
	wasm, err := os.ReadFile(contractSrc)
	if err != nil {
		return nil, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return nil, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return nil, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}
	return wasm, nil
}
