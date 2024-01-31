package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	generateCmd.Flags().StringVarP(&lang, "lang", "l", "", "language to init project in")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "", "name of project")
	generateCmd.Flags().StringVarP(&projectPath, "path", "o", "", "path to create project in")

	generateCmd.MarkFlagRequired("lang")
	generateCmd.MarkFlagRequired("name")
	generateCmd.MarkFlagRequired("path")
}

var LANGUAGES = []string{"go", "python", "rust", "typescript", "zig"}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate basic projects for a group of languages",
	Long: `Assumes you have the supporting tooling preinstalled within your environment, enabling generation of basic projects for the following languages:
- Go
- Python
- Rust
- TypeScript
- Zig`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: when parsing the args, should support
		// lowercase, uppercase, and titlecase for the
		// language name
		fmt.Println("proj-init generate called")
	},
}
