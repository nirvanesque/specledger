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
	// Clone the repository using go-git
	cloneOpts := CloneOptions{
		URL:       repoURL,
		Branch:    branch,
		TargetDir: cacheDir,
		Shallow:   true,
	}

	_, _, err := Clone(cloneOpts)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
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
