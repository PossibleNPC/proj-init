package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of proj-init",
	Long:  `All software has versions. This is proj-init's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proj-init v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
