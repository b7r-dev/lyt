package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "lyt",
	Short: "lyt - static site generator",
	Long: `lyt - yaml, markdown, templates, zero runtime JS.

A minimal static site generator: YAML for structure, Markdown for content,
Templ for components, pure HTML output.

Run from any project directory containing content/ and templates/.
Output defaults to ./dist in the current directory.

Commands:
  build       Build the static site
  serve       Start development server with live reload
  init        Initialize a new project
  validate    Validate content against schema

Use "lyt <command> --help" for more information about a command.
Use "lyt help agent" for AI agent-specific guidance.`,
	Version: "1.1.1",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.SetVersionTemplate("lyt {{printf \"%s\" .Version}}\n")
}

func initConfig() {
	// Config resolved at content-load time
}
