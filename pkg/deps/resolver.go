// Package deps provides artifact path discovery, cache operations, and reference resolution
// for SpecLedger dependencies.
//
// The package handles:
// - Auto-discovery of artifact_path from SpecLedger repositories
// - Manual artifact_path specification for non-SpecLedger repos
// - Cache operations for ~/.specledger/cache/
// - Reference resolution using alias:artifact syntax
package deps

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/specledger/specledger/pkg/cli/metadata"
)

// DetectArtifactPathFromSpecLedgerRepo reads specledger.yaml from a local repository
// and returns the artifact_path value.
//
// Parameters:
//   - repoPath: Local path to the cloned repository
//
// Returns:
//   - artifact_path value from specledger.yaml
//   - error if specledger.yaml not found or artifact_path is missing/empty
func DetectArtifactPathFromSpecLedgerRepo(repoPath string) (string, error) {
	// Check if specledger.yaml exists
	yamlPath := filepath.Join(repoPath, metadata.DefaultMetadataFile)
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		return "", fmt.Errorf("not a SpecLedger repository: %s not found", metadata.DefaultMetadataFile)
	}

	// Load metadata
	meta, err := metadata.Load(yamlPath)
	if err != nil {
		return "", fmt.Errorf("failed to read specledger.yaml: %w", err)
	}

	// Get artifact_path with default fallback
	artifactPath := meta.GetArtifactPath()
	if artifactPath == "" {
		return "", errors.New("artifact_path is empty in specledger.yaml")
	}

	return artifactPath, nil
}

// DetectArtifactPathFromRemote clones a repository (shallow clone for speed),
// reads its specledger.yaml, and returns the artifact_path value.
//
// This is useful for detecting artifact_path before adding a dependency.
//
// Parameters:
//   - repoURL: Git repository URL
//   - branch: Branch to clone (default "main")
//   - cacheDir: Temporary directory for cloning
//
// Returns:
//   - artifact_path value from specledger.yaml
//   - error if clone fails, not a SpecLedger repo, or artifact_path missing
func DetectArtifactPathFromRemote(repoURL, branch, cacheDir string) (string, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Shallow clone for speed
	args := []string{"clone", "--depth", "1", "--single-branch"}
	if branch != "" && branch != "main" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repoURL, cacheDir)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, string(output))
	}

	// Detect artifact_path from cloned repo
	artifactPath, err := DetectArtifactPathFromSpecLedgerRepo(cacheDir)
	if err != nil {
		// Clean up the clone if detection failed
		os.RemoveAll(cacheDir)
		return "", fmt.Errorf("artifact path detection failed: %w", err)
	}

	return artifactPath, nil
}
