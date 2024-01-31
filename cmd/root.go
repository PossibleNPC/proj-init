package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	lang        string
	purpose     string
	username    string
	projectName string
	license     string
	projectPath string

	rootCmd = &cobra.Command{
		Use:   "proj-init",
		Short: "A CLI tool to init different kinds of personal projects",
		Long: `A CLI tool to init different kinds of personal projects for:
    - Go
    - Python
    - Rust
    - TypeScript
    - Zig`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("proj-init called")
		},
	}
)

func Execute() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}

func init() {
}
