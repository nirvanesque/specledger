package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/metadata"
)

// VarMigrateCmd represents the migrate command
var VarMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate .mod files to YAML format",
	Long: `Migrate existing specledger.mod files to the new specledger.yaml format.

This command converts legacy .mod files to the new YAML metadata format.
The original .mod file is preserved for backup purposes.

Migration rules:
- Project name and short code are extracted from .mod comments
- Framework choice defaults to 'none' (you can edit the YAML later)
- Dependencies are preserved (if any were declared)
- Original .mod file is kept as specledger.spec.mod.backup

Examples:
  sl migrate           # Migrate in current project directory
  sl migrate --dry-run # Preview changes without writing`,
	RunE: runMigrate,
}

var migrateDryRun bool

func init() {
	VarMigrateCmd.Flags().BoolVarP(&migrateDryRun, "dry-run", "d", false, "Preview changes without writing files")
}

func runMigrate(cmd *cobra.Command, args []string) error {
	// Get current directory (or specified directory)
	projectDir := "."
	if len(args) > 0 {
		projectDir = args[0]
	}

	// Resolve to absolute path
	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if .mod file exists
	if !metadata.HasLegacyModFile(absDir) {
		return fmt.Errorf("no specledger.mod file found in %s", absDir)
	}

	// Check if YAML already exists
	if metadata.HasYAMLMetadata(absDir) {
		return fmt.Errorf("specledger.yaml already exists. Use --force to overwrite (not implemented yet)")
	}

	fmt.Printf("Migrating specledger.mod to specledger.yaml in %s\n", absDir)

	// Parse the .mod file to show preview
	modPath := filepath.Join(absDir, "specledger", "specledger.mod")
	modData, err := metadata.ParseModFile(modPath)
	if err != nil {
		return fmt.Errorf("failed to parse .mod file: %w", err)
	}

	// Show what would be done
	if migrateDryRun {
		fmt.Println("\nDry run - would create:")
		fmt.Printf("  Project: %s (short code: %s)\n", modData.ProjectName, modData.ShortCode)
		fmt.Printf("  Framework: none (default for migrated projects)\n")
		fmt.Printf("  Dependencies: 0\n")
		fmt.Println("\nOriginal .mod file would be backed up as: specledger.spec.mod.backup")
		return nil
	}

	// Perform migration
	meta, err := metadata.MigrateModToYAML(absDir)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	// Backup the .mod file
	backupPath := filepath.Join(absDir, "specledger", "specledger.spec.mod.backup")
	if err := os.Rename(modPath, backupPath); err != nil {
		fmt.Printf("Warning: failed to backup .mod file: %v\n", err)
	} else {
		fmt.Printf("  Backup: %s\n", backupPath)
	}

	fmt.Printf("\n✓ Migration complete!\n")
	fmt.Printf("  Project: %s (short code: %s)\n", meta.Project.Name, meta.Project.ShortCode)
	fmt.Printf("  Framework: %s\n", meta.Framework.Choice)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the generated specledger.yaml")
	fmt.Println("  2. Edit framework choice if desired (none/speckit/openspec/both)")
	fmt.Println("  3. Optionally remove the backup file: rm specledger.spec.mod.backup")

	return nil
}
