package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "v0.1-local"
	commit    = "manual"
	buildDate = "2025-04-16"
)

var rootCmd = &cobra.Command{
	Use:   "threatfox",
	Short: "ThreatFox CLI - collect and convert threat intelligence into LLM datasets",
	Long: fmt.Sprintf(`ThreatFox CLI - collect and convert threat intelligence into LLM datasets

Version: %s
Commit:  %s
Built:   %s
`, version, commit, buildDate),
}

// Execute triggers the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// AddCommandFunc is called from init()s in subcommands
func AddCommandFunc(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
