package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(rootCmd.Version)
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
