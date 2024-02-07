package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	generateCmd.Flags().StringVarP(&lang, "lang", "l", "", "language to init project in (required)")
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "", "name of project (required)")
	generateCmd.Flags().StringVarP(&projectPath, "path", "o", "", "path to create project in (required)")

	generateCmd.MarkFlagRequired("lang")
	generateCmd.MarkFlagRequired("name")
	generateCmd.MarkFlagRequired("path")
}

type Language struct {
	Name         string
	CliEvoke     string
	CliEvokeInit string
	CliArgs      []string
	RequireChDir bool
	VersionFlag  []string
}

func (l Language) stringifyCliArgs() string {
	s := strings.Builder{}
	for _, a := range l.CliArgs {
		s.WriteString(a)
		s.WriteString(" ")
	}

	ts := strings.Trim(s.String(), " ")

	return ts
}

type Languages []Language

// var LANGUAGES = []string{"go", "python", "rust", "typescript", "zig"}
// TODO: Typescipt is going to be a bit more of a challenge to get working
// There is NPM, PNPM, Yarn, and Deno to consider
// I think I want to restrict to PNPM and Deno

// Actually, this brings up a good point; for languages that support a package
// manager, like PNPM, VENV, Poetry or the like, I need to change the parts
// of the code because it is based off a lot of assumptions and supporting tooling
// TODO: support a handful of package managers for each language where applicable

// to start, I want to use native built-ins for creating a brand new project
// but with a language like Zig or Rust, you also might want to create a new
// library instead of a binary

// oh, this is interesting, so even though we get the os.environ,
// the CliEvokeInit is not going to work because it just needs to be the command
// this points to the need to carefully think through this part
var LANGUAGES = Languages{
	Language{
		Name:         "go",
		CliEvoke:     "go",
		CliEvokeInit: "go",
		CliArgs:      []string{"mod", "init"},
		RequireChDir: true,
		VersionFlag:  []string{"version"},
	},
	Language{
		Name:         "python",
		CliEvoke:     "python3",
		CliEvokeInit: "python3",
		CliArgs:      []string{"-m", "venv"},
		RequireChDir: true,
		VersionFlag:  []string{"--version", "-V", "-VV"},
	},
	Language{
		Name:         "rust",
		CliEvoke:     "rustc",
		CliEvokeInit: "cargo",
		CliArgs:      []string{"new"},
		RequireChDir: false,
		VersionFlag:  []string{"--version"},
	},
	// TODO: this is going to be a bit more complex
	Language{
		Name:         "typescript",
		CliEvoke:     "tsc",
		CliEvokeInit: "npm",
		CliArgs:      []string{"init"},
		RequireChDir: true,
		VersionFlag:  []string{"--version", "-v"},
	},
	Language{
		Name:         "zig",
		CliEvoke:     "zig",
		CliEvokeInit: "zig",
		CliArgs:      []string{"init-exe"},
		RequireChDir: true,
		VersionFlag:  []string{"version"},
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate basic projects for a group of languages",
	Long: `Assumes you have the supporting tooling preinstalled within your environment, enabling generation of basic projects for the following languages:
- Go
- Python 3
- Rust
- TypeScript
- Zig`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: when parsing the args, should support
		// lowercase, uppercase, and titlecase for the
		// language name
		fmt.Println("proj-init generate called")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateLanguageAvailable()
		if err != nil {
			return err
		}

		err = createProject()
		if err != nil {
			return err
		}

		return nil
	},
}

// TODO: each of these have custom logic for creating a project
// swap from a FOR loop to Switch on the L.Name
func createProject() error {
	// check if the path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		// create the path
		err := os.MkdirAll(projectPath, 0755)
		if err != nil {
			return fmt.Errorf("could not create project: %v", err)
		}
	}
	// check if the project under the path exists
	if _, err := os.Stat(fmt.Sprintf("%s/%s", projectPath, projectName)); !os.IsNotExist(err) {
		return fmt.Errorf("project already exists")
	}
	// so this is also incorrect because of the way the
	// some language package managers create the project
	// with the name and then underlay the project files
	err := os.MkdirAll(fmt.Sprintf("%s/%s", projectPath, projectName), 0755)
	if err != nil {
		return fmt.Errorf("could not create project: %v", err)
	}

	// Assume that we cd into the project directory
	// before running the language's cli to create the project
	// this is incorrect because some languages don't need to be
	// in the project directory to create the project
	// err = os.Chdir(fmt.Sprintf("%s/%s", projectPath, projectName))
	// if err != nil {
	// 	return fmt.Errorf("could not change directory to project: %v", err)
	// }

	// I think this points to a need to tie our user input
	// to one of the languages in LANGUAGES for easier use

	// Even though I don't think this next is very good, it'll work
	var cmd *exec.Cmd

	// now we evoke the language's cli to create the project
	for _, l := range LANGUAGES {
		if lang == l.Name {
			if l.RequireChDir {
				err := os.Chdir(fmt.Sprintf("%s/%s", projectPath, projectName))
				if err != nil {
					return fmt.Errorf("could not change directory to project: %v", err)
				}
				if l.Name == "zig" {
					cmd = exec.Command(l.CliEvokeInit, l.stringifyCliArgs())
				} else {
					cmd = exec.Command(l.CliEvokeInit, l.stringifyCliArgs(), projectName)
				}

				// this highlights an bug
				// some projects with package managers don't need to be in the project directory
				// when creating the project
				// others do
				// TODO: identify the languages that need to be in the project directory
				// and those that don't
				// TODO: this is going to be hardcoded to prove it works under the happy path
				// cmd = exec.Command(l.Name, l.stringifyCliArgs(), projectName)
				// set the environment for the cmd
				cmd.Env = os.Environ()
				err = cmd.Run()
				if err != nil {
					return fmt.Errorf("issue running language, unable to generate project for %s", lang)
				}
			} else {
				cmd = exec.Command(l.CliEvokeInit, l.stringifyCliArgs(), projectName)
				// set the environment for the cmd
				cmd.Env = os.Environ()
				err = cmd.Run()
				if err != nil {
					return fmt.Errorf("issue running language, unable to generate project for %s", lang)
				}
			}

			// TODO: project path could include relative paths, so
			// we should clean the path first before printing to make it
			// more readable
			fmt.Printf("Created project %s at %s\n", projectName, projectPath)

			return nil
		}
	}

	return fmt.Errorf("language is not supported by proj-init")
}

func validateLanguageAvailable() error {
	// I'm not sure if I should also check the env to ensure the tooling is in place
	// or not including the paths; I'm going to assume the user has the tooling sorted
	// before using this tool

	// There are also no assumptions on the version of the language installed
	// First look at the language the user wants to use
	// Then check if the language is available on the system
	// based on one of the version flags
	for _, l := range LANGUAGES {
		if lang == l.Name {
			fmt.Println("running language version check")
			cmd := exec.Command(l.CliEvoke, l.VersionFlag[0])
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("issue running language check, cannot generate project for %s", lang)
			}
			fmt.Println("language is available")
			return nil
		}
	}

	return fmt.Errorf("language is not supported by proj-init")
}
