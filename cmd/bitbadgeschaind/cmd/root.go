package cmd

import (
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log/v2"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	txsigning "github.com/cosmos/cosmos-sdk/x/tx/signing"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	evmhd "github.com/cosmos/evm/crypto/hd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/bitbadges/bitbadgeschain/app"
)

const (
	flagVersion  = "version"
	flagHelpJSON = "help-json"
)

// NewRootCmd creates a new root command for bitbadgeschaind. It is called once in the main function.
func NewRootCmd() *cobra.Command {
	var (
		autoCliOpts        autocli.AppOptions
		moduleBasicManager module.BasicManager
		clientCtx          client.Context
	)

	if err := depinject.Inject(
		depinject.Configs(app.AppConfig(),
			depinject.Supply(
				log.NewNopLogger(),
			),
			depinject.Provide(
				ProvideClientContext,
				ProvideEVMGetSigners,
			),
		),
		&autoCliOpts,
		&moduleBasicManager,
		&clientCtx,
	); err != nil {
		panic(err)
	}
	// Apply keyring options to the initial client context BEFORE creating root command
	// This ensures keyring options are available when keys.Commands() is called in initRootCmd
	// Use Cosmos EVM's keyring option directly
	clientCtx = clientCtx.WithKeyringOptions(evmhd.EthSecp256k1Option())

	var rootCmd *cobra.Command
	rootCmd = &cobra.Command{
		Use:           app.Name + "d",
		Short:         "Start bitbadgeschain node",
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if helpJSON, _ := cmd.Flags().GetBool(flagHelpJSON); helpJSON {
				return nil
			}
			return cobra.NoArgs(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Skip server/client config initialization for root-only flags
			// (`bb`, `bb --version`, `bb --help-json`).
			if cmd == rootCmd {
				return nil
			}

			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			// SDK-forwarded commands delegate verbatim to bitbadges-cli and do
			// not consume Cosmos client.toml or server config files. Bypassing
			// config loading here allows `bb dev --help`, `bb build --help`,
			// etc. to work cleanly with an empty, missing, or read-only --home.
			if isSDKForwarderCmd(cmd) {
				return nil
			}

			clientCtx = clientCtx.WithCmdContext(cmd.Context())
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			clientCtx, err = config.ReadFromClientConfig(clientCtx)
			if err != nil {
				return err
			}

			// Apply keyring options to support eth_secp256k1
			// This must be done after ReadFromClientConfig to ensure they're not overwritten
			// Use Cosmos EVM's keyring option directly
			clientCtx = clientCtx.WithKeyringOptions(evmhd.EthSecp256k1Option())

			if err := client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			if err := server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig); err != nil {
				return err
			}

			return routeLogsAwayFromDataStdout(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion, _ := cmd.Flags().GetBool(flagVersion); showVersion {
				fmt.Fprintln(cmd.OutOrStdout(), formatRootVersion())
				return nil
			}
			if helpJSON, _ := cmd.Flags().GetBool(flagHelpJSON); helpJSON {
				return execNodeCLIWithIO("--help-json", args, cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			}
			return cmd.Help()
		},
	}

	rootCmd.Flags().Bool(flagVersion, false, "Print version information")
	rootCmd.Flags().Bool(flagHelpJSON, false, "Output SDK CLI commands as structured JSON (forwards to bitbadges-cli --help-json)")

	// Since the IBC modules don't support dependency injection, we need to
	// manually register the modules on the client side.
	// This needs to be removed after IBC supports App Wiring.
	ibcModules := app.RegisterIBC(clientCtx.InterfaceRegistry)
	for name, mod := range ibcModules {
		moduleBasicManager[name] = module.CoreAppModuleBasicAdaptor(name, mod)
		autoCliOpts.Modules[name] = mod
	}

	gammModules := app.RegisterGamm(clientCtx.InterfaceRegistry)
	for name, mod := range gammModules {
		moduleBasicManager[name] = module.CoreAppModuleBasicAdaptor(name, mod)
		autoCliOpts.Modules[name] = mod
	}

	// Install `bb --help` Cobra groups before commands are added so the
	// initRootCmd / autocli / SDK-forwarder registrations can attach
	// each command to its group as it's wired up.
	registerHelpGroups(rootCmd)

	initRootCmd(rootCmd, clientCtx.TxConfig, moduleBasicManager)

	overwriteFlagDefaults(rootCmd, map[string]string{
		flags.FlagChainID:        strings.ReplaceAll(app.Name, "-", ""),
		flags.FlagKeyringBackend: "test",
	})

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	// Tag every chain-native top-level command with its GroupID for
	// the `bb --help` grouping. Runs after autocli enhancement so it
	// catches commands the enhancer adds (e.g. module tx/query
	// shortcuts). SDK CLI forwarders register their own GroupID
	// inline in sdk_forwarders.go and are skipped here.
	tagChainCommandGroups(rootCmd)

	return rootCmd
}

func overwriteFlagDefaults(c *cobra.Command, defaults map[string]string) {
	set := func(s *pflag.FlagSet, key, val string) {
		if f := s.Lookup(key); f != nil {
			f.DefValue = val
			_ = f.Value.Set(val)
		}
	}
	for key, val := range defaults {
		set(c.Flags(), key, val)
		set(c.PersistentFlags(), key, val)
	}
	for _, c := range c.Commands() {
		overwriteFlagDefaults(c, defaults)
	}
}

// ProvideEVMGetSigners provides the custom signer for EVM MsgEthereumTx
func ProvideEVMGetSigners() txsigning.CustomGetSigner {
	return evmtypes.MsgEthereumTxCustomGetSigner
}

func ProvideClientContext(
	appCodec codec.Codec,
	interfaceRegistry codectypes.InterfaceRegistry,
	txConfigOpts tx.ConfigOptions,
	legacyAmino *codec.LegacyAmino,
) client.Context {
	// Register legacy types for CLI backward compatibility.
	// These are needed so CLI queries (like gov proposals) can resolve old type URLs
	// from renamed/removed modules. Without this, queries fail with "unable to resolve type URL".
	// The same registration happens in app.go for the full node, but the CLI doesn't
	// initialize the full app for queries.
	//
	// IMPORTANT: We must register these BEFORE creating the codec, and then create
	// a new codec from the augmented interface registry. The appCodec parameter
	// was created by depinject before we had a chance to register our legacy types.
	app.RegisterLegacyWasmInterfaces(interfaceRegistry)
	app.RegisterLegacyBadgesInterfaces(interfaceRegistry)

	// Create a new codec that includes the legacy types we just registered
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Create client context with keyring options set BEFORE ReadFromClientConfig
	// This matches the Cosmos EVM reference implementation pattern
	// Reference: https://github.com/cosmos/evm/blob/f1f4c2aee76243f0ffe0dd444e05b9cb2ef9898b/evmd/cmd/evmd/cmd/root.go#L85
	clientCtx := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(interfaceRegistry).
		WithLegacyAmino(legacyAmino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper(app.Name).                           // env variable prefix
		WithKeyringOptions(evmhd.EthSecp256k1Option()) // Set keyring options BEFORE ReadFromClientConfig

	// Read the config again to overwrite the default values with the values from the config file
	// Note: ReadFromClientConfig should preserve keyring options
	clientCtx, _ = config.ReadFromClientConfig(clientCtx)

	txConfig, err := tx.NewTxConfigWithOptions(clientCtx.Codec, txConfigOpts)
	if err != nil {
		panic(err)
	}
	clientCtx = clientCtx.WithTxConfig(txConfig)

	return clientCtx
}

// routeLogsAwayFromDataStdout moves the server logger onto stderr for commands
// whose stdout carries data.
//
// InterceptConfigsPreRunHandler builds the logger with
// CreateSDKLogger(serverCtx, cmd.OutOrStdout()), and ExportCmd writes the
// exported genesis to cmd.OutOrStdout() as well, so the two share a stream.
// Whether that corrupts anything depends only on whether the app happens to log
// while the command runs — which makes it a latent bug that surfaces the first
// time someone adds a log line on a construction path.
//
// It surfaced in v34. v33 logged nothing during export, so
// `bitbadgeschaind export > genesis.json` produced a valid file; v34 logs a
// max_gas warning and a RecheckMempool context error while building the app, so
// the same command wrote a file beginning with two log lines that no genesis
// parser accepts. The genesis was intact on the last line, which is worse than
// an outright failure: the command still exits 0.
//
// Only these commands are redirected, deliberately. `start` is what validators
// run under systemd and its stdout is not data; moving its logs to stderr mid
// upgrade would break any operator whose unit file redirects stdout to a log
// file, to fix a problem `start` does not have.
func routeLogsAwayFromDataStdout(cmd *cobra.Command) error {
	if cmd.Name() != "export" {
		return nil
	}

	serverCtx := server.GetServerContextFromCmd(cmd)
	logger, err := server.CreateSDKLogger(serverCtx, cmd.ErrOrStderr())
	if err != nil {
		return err
	}
	serverCtx.Logger = logger

	return server.SetCmdServerContext(cmd, serverCtx)
}

// formatRootVersion builds the output for `bb --version` / `bitbadgeschaind --version`.
// It always prints the chain binary version and appends the SDK CLI version when
// bitbadges-cli is directly installed, without failing when bitbadges-cli is absent.
func formatRootVersion() string {
	info := version.NewInfo()
	ver := strings.TrimSpace(info.Version)
	if ver == "" {
		ver = "dev"
	}
	line := fmt.Sprintf("%sd %s", app.Name, ver)
	if commit := strings.TrimSpace(info.GitCommit); commit != "" {
		line += fmt.Sprintf(" (%s)", commit)
	}
	if sdkVer := lookupNodeCLIVersion(); sdkVer != "" {
		line += fmt.Sprintf("\nbitbadges-cli %s", sdkVer)
	}
	return line
}

// MatchRootHelpJSON reports whether args represent a root-level `--help-json`
// invocation (before any subcommand or `-h`/`--help`) and returns any extra
// flags/arguments to pass through to `bitbadges-cli --help-json`.
func MatchRootHelpJSON(args []string) ([]string, bool) {
	helpJSONIdx := -1
	for i, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-h" || arg == "--help" {
			return nil, false
		}
		if arg == "--"+flagHelpJSON {
			helpJSONIdx = i
			break
		}
		if !strings.HasPrefix(arg, "-") {
			// A subcommand precedes --help-json (e.g. `bb build --help-json`);
			// let Cobra route it to the subcommand forwarder.
			return nil, false
		}
	}
	if helpJSONIdx < 0 {
		return nil, false
	}

	extra := make([]string, 0, len(args)-1)
	for i, arg := range args {
		if i == helpJSONIdx {
			continue
		}
		if i > helpJSONIdx && arg == "--" {
			extra = append(extra, args[i+1:]...)
			break
		}
		if arg == "-h" || arg == "--help" {
			return nil, false
		}
		extra = append(extra, arg)
	}
	return extra, true
}

// ExecRootHelpJSON forwards a root-level `--help-json` invocation to bitbadges-cli.
func ExecRootHelpJSON(extraArgs []string) error {
	return execNodeCLI("--"+flagHelpJSON, extraArgs)
}
