package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ThreatFox CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ThreatFox CLI\nVersion:    %s\nCommit:     %s\nBuilt At:   %s\n", 
			rootCmd.Version, commit, buildDate)
	},
}

func init() {
	AddCommandFunc(versionCmd)
}
