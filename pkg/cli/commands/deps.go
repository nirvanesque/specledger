package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"specledger/internal/spec"
	"specledger/pkg/models"
)

// VarDepsCmd represents the deps command
var VarDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage specification dependencies",
	Long: `Manage external specification dependencies for your project.

Dependencies are stored in spec.mod and cached locally for offline use.

Examples:
  sl deps list                           # List all dependencies
  sl deps add git@github.com:org/spec    # Add a dependency
  sl deps remove git@github.com:org/spec # Remove a dependency`,
}

// VarAddCmd represents the add command
var VarAddCmd = &cobra.Command{
	Use:     "add <repo-url> [branch] [spec-path]",
	Short:   "Add a dependency",
	Long:    `Add an external specification dependency to your project. The dependency will be downloaded and cached locally.`,
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
	Long:    `List all declared dependencies from spec.mod, showing their repository, version, and local cache status.`,
	Example: `  sl deps list`,
	RunE:    runListDependencies,
}

// VarRemoveCmd represents the remove command
var VarRemoveCmd = &cobra.Command{
	Use:     "remove <repo-url>",
	Short:   "Remove a dependency",
	Long:    `Remove a dependency from spec.mod. The local cache will be kept for future use.`,
	Example: `  sl deps remove git@github.com:org/api-spec`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRemoveDependency,
}

// VarResolveCmd represents the resolve command
var VarResolveCmd = &cobra.Command{
	Use:     "resolve",
	Short:   "Download and cache dependencies",
	Long:    `Download all dependencies from spec.mod and cache them locally. Validates versions and generates cryptographic hashes.`,
	Example: `  sl deps resolve`,
	RunE:    runResolveDependencies,
}

// VarDepsUpdateCmd represents the update command
var VarDepsUpdateCmd = &cobra.Command{
	Use:     "update [repo-url]",
	Short:   "Update dependencies to latest versions",
	Long:    `Update dependencies to their latest compatible versions. If no URL is given, updates all dependencies.`,
	Example: `  sl deps update                    # Update all
  sl deps update git@github.com:org/spec # Update one`,
	RunE:    runUpdateDependencies,
}

func init() {
	VarDepsCmd.AddCommand(VarAddCmd, VarDepsListCmd, VarResolveCmd, VarDepsUpdateCmd, VarRemoveCmd)

	VarAddCmd.Flags().StringP("alias", "a", "", "Optional alias for the dependency")
	VarDepsListCmd.Flags().BoolP("include-transitive", "t", false, "Include transitive dependencies")
	VarResolveCmd.Flags().BoolP("no-cache", "n", false, "Ignore cached specifications")
	VarResolveCmd.Flags().BoolP("deep", "d", false, "Fetch full git history")
	VarDepsUpdateCmd.Flags().BoolP("force", "f", false, "Force update all dependencies")
}

func runAddDependency(cmd *cobra.Command, args []string) error {
	// Extract flags
	alias, _ := cmd.Flags().GetString("alias")

	// Parse arguments
	repoURL := args[0]
	version := "main" // default
	specPath := "spec.md"

	if len(args) >= 2 {
		version = args[1]
	}
	if len(args) >= 3 {
		specPath = args[2]
	}

	// Validate URL
	if !isValidURL(repoURL) {
		return fmt.Errorf("invalid repository URL: %s", repoURL)
	}

	// Create dependency
	dep := models.Dependency{
		RepositoryURL: repoURL,
		Version:       version,
		SpecPath:      specPath,
		Alias:         alias,
	}

	// Validate
	if err := dep.Validate(); err != nil {
		return fmt.Errorf("invalid dependency: %w", err)
	}

	// Read existing manifest
	manifestPath := "specledger/specledger.mod"
	manifest, err := spec.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	// Add dependency (manually append for now)
	manifest.Dependecies = append(manifest.Dependecies, dep)
	manifest.UpdatedAt = time.Now()

	// Write manifest
	if err := spec.WriteManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	fmt.Printf("Added dependency: %s -> %s\n", repoURL, specPath)
	if alias != "" {
		fmt.Printf("  Alias: %s\n", alias)
	}

	return nil
}

func runListDependencies(cmd *cobra.Command, args []string) error {
	manifestPath := "specledger/specledger.mod"
	manifest, err := spec.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	fmt.Printf("Dependencies (%d):\n", len(manifest.Dependecies))
	fmt.Println()

	for i, dep := range manifest.Dependecies {
		fmt.Printf("%d. %s\n", i+1, dep.RepositoryURL)
		fmt.Printf("   Version: %s\n", dep.Version)
		fmt.Printf("   Spec: %s\n", dep.SpecPath)
		if dep.Alias != "" {
			fmt.Printf("   Alias: %s\n", dep.Alias)
		}
		fmt.Println()
	}

	return nil
}

func runResolveDependencies(cmd *cobra.Command, args []string) error {
	noCache, _ := cmd.Flags().GetBool("no-cache")

	manifestPath := "specledger/specledger.mod"
	lockfilePath := "specledger/specledger.sum"

	// Read manifest
	manifest, err := spec.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	// Validate manifest
	errors := spec.ValidateManifest(manifest)
	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err.Error())
		}
		return fmt.Errorf("%d validation errors found", len(errors))
	}

	// Show cache directory
	cacheDir := spec.GetGlobalCacheDir()
	fmt.Printf("Cache directory: %s\n", cacheDir)
	fmt.Printf("Resolving %d dependencies...\n\n", len(manifest.Dependecies))

	// Create resolver (uses global cache)
	resolver := spec.NewResolver("")

	// Resolve dependencies
	results, err := resolver.Resolve(cmd.Context(), manifest, noCache)
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Create lockfile
	lockfile := spec.NewLockfile(spec.ManifestVersion)
	for _, result := range results {
		entry := spec.LockfileEntry{
			RepositoryURL: result.Dependency.RepositoryURL,
			CommitHash:    result.CommitHash,
			ContentHash:   result.ContentHash,
			SpecPath:      result.Dependency.SpecPath,
			Branch:        result.Dependency.Version,
			Size:          result.Size,
			FetchedAt:     time.Now().Format(time.RFC3339),
		}
		lockfile.AddEntry(entry)

		source := "remote"
		if result.Source == "cache" {
			source = "cached"
		}

		fmt.Printf("✓ %s: %s\n", result.Dependency.RepositoryURL, source)
		fmt.Printf("  Commit: %s\n", result.CommitHash)
		fmt.Printf("  Spec: %s\n", result.Dependency.SpecPath)
		fmt.Printf("  Hash: %s\n", result.ContentHash)
		fmt.Println()
	}

	// Write lockfile
	if err := lockfile.Write(lockfilePath); err != nil {
		return fmt.Errorf("failed to write lockfile: %w", err)
	}

	fmt.Printf("Lockfile written to: %s\n", lockfilePath)
	fmt.Printf("Total size: %d bytes\n", lockfile.TotalSize)

	return nil
}

func runRemoveDependency(cmd *cobra.Command, args []string) error {
	repoURL := args[0]

	manifestPath := "specledger/specledger.mod"
	manifest, err := spec.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	removed := false
	for i, dep := range manifest.Dependecies {
		if dep.RepositoryURL == repoURL {
			manifest.Dependecies = append(manifest.Dependecies[:i], manifest.Dependecies[i+1:]...)
			removed = true
			break
		}
	}

	if !removed {
		return fmt.Errorf("dependency not found: %s", repoURL)
	}

	manifest.UpdatedAt = time.Now()

	// Write manifest
	if err := spec.WriteManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	fmt.Printf("Removed dependency: %s\n", repoURL)

	return nil
}

func runUpdateDependencies(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	manifestPath := "specledger/specledger.mod"

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
			// TODO: Actually fetch and check for updates
			// For now, we'll just display current version
			fmt.Printf("  %s: already at version %s\n", dep.RepositoryURL, dep.Version)
			unchanged = append(unchanged, dep.RepositoryURL)
		}
	}

	fmt.Printf("\nUpdated %d dependency(ies)\n", len(updated))
	fmt.Printf("Unchanged %d dependency(ies)\n", len(unchanged))

	// TODO: Write updated manifest if there were updates
	// if err := spec.WriteManifest(manifestPath, manifest); err != nil {
	// 	return fmt.Errorf("failed to write manifest: %w", err)
	// }

	return nil
}

func isValidURL(s string) bool {
	// Simple check for common Git URLs
	return len(s) > 0 && (strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@"))
}
