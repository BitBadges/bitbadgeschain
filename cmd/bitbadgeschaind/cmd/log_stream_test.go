package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// newServerCmd builds a command with a server context attached and both output
// streams captured, the shape the root command's PersistentPreRunE leaves
// behind.
func newServerCmd(t *testing.T, name string) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	var out, errOut bytes.Buffer
	cmd := &cobra.Command{Use: name}
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	serverCtx := server.NewDefaultContext()

	// Attach the context the way InterceptConfigsPreRunHandler does. Note that
	// server.SetCmdServerContext does NOT attach one — it only copies over an
	// already-attached *Context — so a fixture that relied on it would leave
	// GetServerContextFromCmd returning a fresh default whose logger writes to
	// the process's real stdout, and both assertions below would pass or fail
	// for reasons unrelated to the code under test.
	cmd.SetContext(context.WithValue(context.Background(), server.ServerContextKey, serverCtx))
	serverCtx.Viper.Set(flags.FlagLogLevel, "info")

	// Mirror InterceptConfigsPreRunHandler, which builds the logger on
	// cmd.OutOrStdout() (server/util.go). Without this the fixture would use
	// NewDefaultContext's logger, which writes to the process's real stdout and
	// would make both assertions below vacuous.
	logger, err := server.CreateSDKLogger(serverCtx, cmd.OutOrStdout())
	require.NoError(t, err)
	serverCtx.Logger = logger
	require.NoError(t, server.SetCmdServerContext(cmd, serverCtx))

	return cmd, &out, &errOut
}

// TestExportKeepsDiagnosticsOffStdout pins that `bitbadgeschaind export`
// produces a stdout stream containing nothing but the exported genesis.
//
// The SDK builds the server logger with `CreateSDKLogger(serverCtx,
// cmd.OutOrStdout())` (server/util.go), and ExportCmd writes the genesis JSON
// to `cmd.OutOrStdout()` too (server/export.go:62,100). Data and diagnostics
// therefore share one stream, and whether that breaks depends entirely on
// whether the app happens to log during export.
//
// On v33 it did not, so `bitbadgeschaind export > genesis.json` produced a
// valid file. On v34 it logs twice during app construction — a max_gas warning
// and a RecheckMempool context error — so the same command produced a file
// starting with two log lines, which jq rejects and any genesis parser rejects.
// The genesis itself was intact on the last line; the stream around it was not.
//
// Relying on the app staying silent is not a fix, because any future log line
// on a construction path silently breaks exports again. Diagnostics belong on
// stderr; stdout belongs to the data.
func TestExportKeepsDiagnosticsOffStdout(t *testing.T) {
	cmd, out, errOut := newServerCmd(t, "export")

	require.NoError(t, routeLogsAwayFromDataStdout(cmd))

	server.GetServerContextFromCmd(cmd).Logger.Info("app construction chatter")

	require.Empty(t, out.String(),
		"export's stdout must carry only the genesis; a log line here means "+
			"`bitbadgeschaind export > genesis.json` writes a file no parser accepts")
	require.Contains(t, errOut.String(), "app construction chatter",
		"the diagnostics must still be emitted, on stderr")
}

// TestStartKeepsLoggingToStdout pins the other half: this must not become a
// blanket "logs go to stderr" change.
//
// `start` is the command 22 validators run under systemd, and its stdout is not
// data. Moving its logs to stderr would silently break any operator whose unit
// file or supervisor redirects stdout to a log file, in the middle of an
// upgrade. Only commands whose stdout carries machine-readable output need the
// redirect.
func TestStartKeepsLoggingToStdout(t *testing.T) {
	cmd, out, errOut := newServerCmd(t, "start")

	require.NoError(t, routeLogsAwayFromDataStdout(cmd))

	server.GetServerContextFromCmd(cmd).Logger.Info("node is running")

	require.Contains(t, out.String(), "node is running",
		"start's logging destination must not change; operators redirect its stdout")
	require.Empty(t, errOut.String())
}
