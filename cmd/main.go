package main

import (
	"os"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/commands"
	"specledger/pkg/cli/logger"
)

var (
	logLevel string
	version  bool
)

var rootCmd = &cobra.Command{
	Use:   "sl",
	Short: "SpecLedger - Unified CLI for bootstrap and dependency management",
	Long: `SpecLedger is a unified CLI tool that provides both project bootstrap (with
interactive TUI) and specification dependency management. Use 'sl' for all
operations.`,
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		// Default run shows help
		cmd.Help()
	},
}

func init() {
	// Setup logging
	log := logger.New(logger.Debug)
	log.Debug("CLI initialized")

	// Add flags
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "debug", "Set log level (debug, info, warn, error)")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print version")

	// Set up command groups
	rootCmd.AddGroup(&cobra.Group{
		ID:    "bootstrap",
		Title: "Bootstrap",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "deps",
		Title: "Dependencies",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "refs",
		Title: "References",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "graph",
		Title: "Graph",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "vendor",
		Title: "Vendor",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "tools",
		Title: "Tools",
	})

	// Add subcommands
	rootCmd.AddCommand(commands.VarBootstrapCmd)
	rootCmd.AddCommand(commands.VarDepsCmd)
	rootCmd.AddCommand(commands.VarRefsCmd)
	rootCmd.AddCommand(commands.VarGraphCmd)
	rootCmd.AddCommand(commands.VarVendorCmd)
	rootCmd.AddCommand(commands.VarConflictCmd)
	rootCmd.AddCommand(commands.VarUpdateCmd)

	// Setup specledger alias for backward compatibility
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
