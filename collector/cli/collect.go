package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// collectCmd is a backwards-compatibility alias for the pre-split CLI verb.
// `agenthound collect <kind> [args]` forwards to `agenthound scan --<kind> [args]`
// and prints a deprecation note. Hidden from --help (the canonical verb is `scan`).
var collectCmd = &cobra.Command{
	Use:    "collect",
	Short:  "Deprecated: use 'scan' instead",
	Hidden: true,
	Long: `'collect' was the pre-split CLI verb. It is preserved as a backwards-
compatibility alias that forwards to 'scan' and prints a deprecation note.

  agenthound collect config              → agenthound scan --config
  agenthound collect mcp   --url URL     → agenthound scan --mcp --url URL
  agenthound collect a2a   --target T    → agenthound scan --a2a --target T

Will be removed in a future release.`,
}

func newCollectAlias(kind string) *cobra.Command {
	return &cobra.Command{
		Use:                kind,
		Short:              fmt.Sprintf("Deprecated: use 'scan --%s' instead", kind),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintf(os.Stderr,
				"Note: 'collect %s' is deprecated; use 'scan --%s' instead.\n",
				kind, kind)
			newArgs := append([]string{"scan", "--" + kind}, args...)
			rootCmd.SetArgs(newArgs)
			return rootCmd.Execute()
		},
	}
}

func init() {
	collectCmd.AddCommand(newCollectAlias("mcp"))
	collectCmd.AddCommand(newCollectAlias("a2a"))
	collectCmd.AddCommand(newCollectAlias("config"))
	rootCmd.AddCommand(collectCmd)
}
