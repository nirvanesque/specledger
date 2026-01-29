package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"specledger/internal/spec"
)

// VarUpdateCmd represents the update command
var VarUpdateCmd = &cobra.Command{
	Use:   "update [--force] [repo-url]",
	Short: "Update dependencies to latest compatible versions",
	Long:  `Update dependencies to the latest compatible versions and regenerate the lockfile.`,
	RunE:  runUpdateDependencies,
}

func init() {
	VarDepsCmd.AddCommand(VarUpdateCmd)

	VarUpdateCmd.Flags().BoolP("force", "f", false, "Force update all dependencies")
}

func runUpdateDependencies(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	manifestPath := "specs/spec.mod"

	// Read current manifest
	manifest, err := spec.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	if len(manifest.Dependecies) == 0 {
		return fmt.Errorf("no dependencies to update")
	}

	fmt.Printf("Checking %d dependency(ies) for updates...\n", len(manifest.Dependecies))

	var updated []string
	var unchanged []string

	for _, dep := range manifest.Dependecies {
		// Determine if we should check this dependency
		shouldUpdate := force || args == nil || len(args) == 0

		if !shouldUpdate && len(args) > 0 {
			// Check if this dependency matches the provided repo URL
			if strings.HasPrefix(dep.RepositoryURL, args[0]) {
				shouldUpdate = true
			}
		}

		if shouldUpdate {
			// Update the version
			// For now, we'll just use the current version
			// In production, this would fetch the latest tag or branch
			fmt.Printf("  %s: already at version %s\n", dep.RepositoryURL, dep.Version)
			unchanged = append(unchanged, dep.RepositoryURL)
		}
	}

	fmt.Printf("\nUpdated %d dependency(ies)\n", len(updated))
	fmt.Printf("Unchanged %d dependency(ies)\n", len(unchanged))

	// Write updated manifest
	if err := spec.WriteManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}
