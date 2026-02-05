package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/metadata"
)

// VarDepsCmd represents the deps command
var VarDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage specification dependencies",
	Long: `Manage external specification dependencies for your project.

Dependencies are stored in specledger/specledger.yaml and cached locally for offline use.

Examples:
  sl deps list                           # List all dependencies
  sl deps add git@github.com:org/spec    # Add a dependency
  sl deps remove git@github.com:org/spec # Remove a dependency`,
}

// VarAddCmd represents the add command
var VarAddCmd = &cobra.Command{
	Use:     "add <repo-url> [branch] [spec-path]",
	Short:   "Add a dependency",
	Long:    `Add an external specification dependency to your project. The dependency will be tracked in specledger.yaml and cached locally for offline use.`,
	Example: `  sl deps add git@github.com:org/api-spec
  sl deps add git@github.com:org/api-spec v1.0 specs/api.md
  sl deps add git@github.com:org/api-spec main spec.md --alias api`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAddDependency,
}

// VarDepsListCmd represents the list command
var VarDepsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all dependencies",
	Long:    `List all declared dependencies from specledger.yaml, showing their repository, version, and resolved status.`,
	Example: `  sl deps list`,
	RunE:    runListDependencies,
}

// VarRemoveCmd represents the remove command
var VarRemoveCmd = &cobra.Command{
	Use:     "remove <repo-url>",
	Short:   "Remove a dependency",
	Long:    `Remove a dependency from specledger.yaml. The local cache will be kept for future use.`,
	Example: `  sl deps remove git@github.com:org/api-spec`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRemoveDependency,
}

// VarResolveCmd represents the resolve command
var VarResolveCmd = &cobra.Command{
	Use:     "resolve",
	Short:   "Download and cache dependencies",
	Long:    `Download all dependencies from specledger.yaml and cache them locally at ~/.specledger/cache/.`,
	Example: `  sl deps resolve`,
	RunE:    runResolveDependencies,
}

// VarDepsUpdateCmd represents the update command
var VarDepsUpdateCmd = &cobra.Command{
	Use:     "update [repo-url]",
	Short:   "Update dependencies to latest versions",
	Long:    `Update dependencies to their latest versions. If no URL is given, updates all dependencies.`,
	Example: `  sl deps update                    # Update all
  sl deps update git@github.com:org/spec # Update one`,
	RunE:    runUpdateDependencies,
}

func init() {
	VarDepsCmd.AddCommand(VarAddCmd, VarDepsListCmd, VarResolveCmd, VarDepsUpdateCmd, VarRemoveCmd)

	VarAddCmd.Flags().StringP("alias", "a", "", "Optional alias for the dependency")
	VarResolveCmd.Flags().BoolP("no-cache", "n", false, "Ignore cached specifications")
}

func runAddDependency(cmd *cobra.Command, args []string) error {
	// Get current directory or find project root
	projectDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Load existing metadata
	meta, err := metadata.LoadFromProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// Extract flags
	alias, _ := cmd.Flags().GetString("alias")

	// Parse arguments
	repoURL := args[0]
	branch := "main" // default
	specPath := "spec.md"

	if len(args) >= 2 {
		branch = args[1]
	}
	if len(args) >= 3 {
		specPath = args[2]
	}

	// Validate URL
	if !isValidGitURL(repoURL) {
		return fmt.Errorf("invalid repository URL: %s", repoURL)
	}

	// Create dependency
	dep := metadata.Dependency{
		URL:    repoURL,
		Branch: branch,
		Path:   specPath,
		Alias:  alias,
	}

	// Check for duplicates
	for _, existing := range meta.Dependencies {
		if existing.URL == repoURL {
			return fmt.Errorf("dependency already exists: %s", repoURL)
		}
		if alias != "" && existing.Alias == alias {
			return fmt.Errorf("alias already exists: %s", alias)
		}
	}

	// Add dependency
	meta.Dependencies = append(meta.Dependencies, dep)

	// Save metadata
	if err := metadata.SaveToProject(meta, projectDir); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Printf("Added dependency: %s\n", repoURL)
	if alias != "" {
		fmt.Printf("  Alias: %s\n", alias)
	}
	fmt.Printf("\nRun 'sl deps resolve' to download and cache this dependency.\n")

	return nil
}

func runListDependencies(cmd *cobra.Command, args []string) error {
	projectDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	meta, err := metadata.LoadFromProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if len(meta.Dependencies) == 0 {
		fmt.Println("No dependencies declared.")
		fmt.Println("\nAdd dependencies with:")
		fmt.Println("  sl deps add git@github.com:org/spec")
		return nil
	}

	fmt.Printf("Dependencies (%d):\n", len(meta.Dependencies))
	fmt.Println()

	for i, dep := range meta.Dependencies {
		fmt.Printf("%d. %s\n", i+1, dep.URL)
		if dep.Branch != "" && dep.Branch != "main" {
			fmt.Printf("   Branch: %s\n", dep.Branch)
		}
		if dep.Path != "" && dep.Path != "spec.md" {
			fmt.Printf("   Path: %s\n", dep.Path)
		}
		if dep.Alias != "" {
			fmt.Printf("   Alias: %s\n", dep.Alias)
		}
		if dep.ResolvedCommit != "" {
			fmt.Printf("   Resolved: %s\n", dep.ResolvedCommit[:8])
		} else {
			fmt.Printf("   Status: not resolved (run 'sl deps resolve')\n")
		}
		fmt.Println()
	}

	return nil
}

func runRemoveDependency(cmd *cobra.Command, args []string) error {
	projectDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	meta, err := metadata.LoadFromProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	target := args[0]

	// Find and remove dependency
	removed := false
	removedIndex := -1

	for i, dep := range meta.Dependencies {
		if dep.URL == target || dep.Alias == target {
			removedIndex = i
			removed = true
			break
		}
	}

	if !removed {
		return fmt.Errorf("dependency not found: %s", target)
	}

	// Remove from slice
	meta.Dependencies = append(meta.Dependencies[:removedIndex], meta.Dependencies[removedIndex+1:]...)

	// Save metadata
	if err := metadata.SaveToProject(meta, projectDir); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Printf("Removed dependency: %s\n", target)

	return nil
}

func runResolveDependencies(cmd *cobra.Command, args []string) error {
	projectDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	meta, err := metadata.LoadFromProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if len(meta.Dependencies) == 0 {
		fmt.Println("No dependencies to resolve.")
		return nil
	}

	// TODO: Implement actual dependency resolution
	// For now, just show what would be resolved
	fmt.Printf("Would resolve %d dependencies:\n", len(meta.Dependencies))
	for _, dep := range meta.Dependencies {
		fmt.Printf("  - %s", dep.URL)
		if dep.Alias != "" {
			fmt.Printf(" (alias: %s)", dep.Alias)
		}
		fmt.Println()
	}

	fmt.Println("\nDependency resolution not yet implemented.")
	fmt.Println("Dependencies are tracked in specledger/specledger.yaml")

	return nil
}

func runUpdateDependencies(cmd *cobra.Command, args []string) error {
	projectDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	meta, err := metadata.LoadFromProject(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	if len(meta.Dependencies) == 0 {
		return fmt.Errorf("no dependencies to update")
	}

	fmt.Printf("Checking %d dependency(ies) for updates...\n", len(meta.Dependencies))

	for _, dep := range meta.Dependencies {
		// TODO: Implement actual update checking
		if dep.ResolvedCommit != "" {
			fmt.Printf("  %s: at %s\n", dep.URL, dep.ResolvedCommit[:8])
		} else {
			fmt.Printf("  %s: not resolved yet\n", dep.URL)
		}
	}

	fmt.Println("\nDependency updates not yet implemented.")

	return nil
}

func isValidGitURL(s string) bool {
	// Simple check for common Git URLs
	return len(s) > 0 && (strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@"))
}

func findProjectRoot() (string, error) {
	// Start from current directory and work up
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check current directory
	if metadata.HasYAMLMetadata(dir) {
		return dir, nil
	}

	// Check parent directories
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return "", fmt.Errorf("not in a SpecLedger project (no specledger/specledger.yaml found)")
		}
		dir = parent

		if metadata.HasYAMLMetadata(dir) {
			return dir, nil
		}
	}
}
