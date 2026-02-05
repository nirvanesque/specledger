package commands

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/config"
	"specledger/pkg/cli/logger"
	"specledger/pkg/cli/tui"
	"specledger/pkg/embedded"
)

var (
	projectNameFlag string
	shortCodeFlag    string
	playbookFlag     string
	shellFlag        string
	demoDirFlag      string
	ciFlag           bool
)

// VarBootstrapCmd is the bootstrap command
var VarBootstrapCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new SpecLedger project",
	Long: `Create a new SpecLedger project with all necessary infrastructure:

Interactive mode:
  sl new

Non-interactive mode (for CI/CD):
  sl new --ci --project-name <name> --short-code <code> --project-dir <path>

The bootstrap creates:
- .claude/ directory with skills and commands
- .beads/ directory for issue tracking
- specledger/ directory for specifications
- specledger/specledger.mod file for project metadata`,

	// RunE is called when the command is executed
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create logger
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		l := logger.New(logger.Debug)

		// Check if non-interactive mode
		modeDetector := tui.NewModeDetector()
		if modeDetector.IsNonInteractive() || ciFlag {
			return runBootstrapNonInteractive(cmd, l, cfg)
		}

		// Interactive TUI mode
		return runBootstrapInteractive(l, cfg)
	},
}

// runBootstrapInteractive runs the bootstrap with Bubble Tea TUI
func runBootstrapInteractive(l *logger.Logger, cfg *config.Config) error {
	// Determine default project directory
	defaultDir := cfg.DefaultProjectDir
	if demoDirFlag != "" {
		defaultDir = demoDirFlag
	}

	// Run Bubble Tea TUI with default directory
	tuiProgram := tui.NewProgram(defaultDir)
	answers, err := tuiProgram.Run()
	if err != nil {
		return fmt.Errorf("TUI exited with error: %w", err)
	}

	// Check if user cancelled (Ctrl+C)
	if len(answers) == 0 || answers["project_name"] == "" {
		return fmt.Errorf("bootstrap cancelled by user")
	}

	projectName := answers["project_name"]
	projectDir := answers["project_dir"]
	shortCode := answers["short_code"]
	playbook := answers["playbook"]
	shell := answers["shell"]

	// Create project path
	projectPath := filepath.Join(projectDir, projectName)

	// Check if directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		shouldOverwrite, err := tui.ConfirmPrompt(fmt.Sprintf("Directory '%s' already exists. Overwrite? [y/N]: ", projectName))
		if err != nil {
			return fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !shouldOverwrite {
			return fmt.Errorf("bootstrap cancelled by user")
		}
	}

	// Create directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Copy template files
	if err := copyTemplates(projectPath, shortCode, projectName); err != nil {
		return fmt.Errorf("failed to copy templates: %w", err)
	}

	// Initialize git repo (but don't commit - user might bootstrap into existing repo)
	if err := initializeGitRepo(projectPath); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	// Success message
	fmt.Printf("\n✓ Project created: %s\n", projectPath)
	fmt.Printf("✓ Beads prefix: %s\n", shortCode)
	fmt.Printf("✓ Playbook: %s\n", playbook)
	fmt.Printf("✓ Agent Shell: %s\n", shell)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Printf("  mise install    # Install tools (bd)\n")
	fmt.Printf("  claude\n")

	return nil
}

// runBootstrapNonInteractive runs bootstrap without TUI
func runBootstrapNonInteractive(cmd *cobra.Command, l *logger.Logger, cfg *config.Config) error {
	// Validate required flags
	if projectNameFlag == "" {
		return fmt.Errorf("--project-name flag is required in non-interactive mode")
	}

	if shortCodeFlag == "" {
		return fmt.Errorf("--short-code flag is required in non-interactive mode")
	}

	projectName := projectNameFlag
	shortCode := strings.ToLower(shortCodeFlag)

	// Limit short code to 4 characters
	if len(shortCode) > 4 {
		shortCode = shortCode[:4]
	}

	// Get demo directory
	demoDir := cfg.DefaultProjectDir
	if demoDirFlag != "" {
		demoDir = demoDirFlag
	}

	projectPath := filepath.Join(demoDir, projectName)

	// Check if directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return ErrProjectExists(projectName)
	}

	// Create directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return ErrPermissionDenied(projectPath)
	}

	// Copy template files
	if err := copyTemplates(projectPath, shortCode, projectName); err != nil {
		return fmt.Errorf("failed to copy templates: %w", err)
	}

	// Initialize git repo (but don't commit - user might bootstrap into existing repo)
	if err := initializeGitRepo(projectPath); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	// Success message
	fmt.Printf("\n✓ Project created: %s\n", projectPath)
	fmt.Printf("✓ Beads prefix: %s\n", shortCode)
	if playbookFlag != "" {
		fmt.Printf("✓ Playbook: %s\n", playbookFlag)
	}
	if shellFlag != "" {
		fmt.Printf("✓ Agent Shell: %s\n", shellFlag)
	}
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Printf("  mise install    # Install tools (bd)\n")
	fmt.Printf("  claude\n")

	return nil
}

func init() {
	VarBootstrapCmd.PersistentFlags().StringVarP(&projectNameFlag, "project-name", "n", "", "Project name")
	VarBootstrapCmd.PersistentFlags().StringVarP(&shortCodeFlag, "short-code", "s", "", "Short code (2-4 letters)")
	VarBootstrapCmd.PersistentFlags().StringVarP(&playbookFlag, "playbook", "p", "", "Playbook type")
	VarBootstrapCmd.PersistentFlags().StringVarP(&shellFlag, "shell", "", "claude-code", "Agent shell")
	VarBootstrapCmd.PersistentFlags().StringVarP(&demoDirFlag, "project-dir", "d", "", "Project directory path")
	VarBootstrapCmd.PersistentFlags().BoolVarP(&ciFlag, "ci", "", false, "Force non-interactive mode (skip TUI)")
}

// copyTemplates copies SpecLedger template files to the new project using embedded templates
func copyTemplates(projectPath, shortCode, projectName string) error {
	// Files and directories to exclude from copying
	excludePaths := map[string]bool{
		"specledger/FORK.md":          true,
		"specledger/memory":           true,
		"specledger/scripts":          true,
		"spec-kit-version":            true,
		"specledger/spec-kit-version": true,
		"specledger/templates":        true,
		// Don't exclude specledger directory itself - we want it!
	}

	// Walk through the embedded filesystem
	err := fs.WalkDir(embedded.TemplatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path (remove "templates/" prefix)
		relPath := strings.TrimPrefix(path, "templates/")
		if relPath == "" || relPath == "." {
			return nil
		}

		// Check if this path should be excluded
		if excludePaths[relPath] {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		destPath := filepath.Join(projectPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read file from embedded FS
		data, err := embedded.TemplatesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// For .beads/config.yaml, replace the prefix
		if filepath.Base(path) == "config.yaml" && filepath.Dir(path) == "templates/.beads" {
			data = []byte(strings.ReplaceAll(string(data), "issue-prefix: \"sl\"", fmt.Sprintf("issue-prefix: \"%s\"", shortCode)))
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		// If we just copied mise.toml, run mise trust
		if filepath.Base(path) == "mise.toml" {
			cmd := exec.Command("mise", "trust")
			cmd.Dir = projectPath
			_ = cmd.Run() // Ignore errors if mise is not installed
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk embedded templates: %w", err)
	}

	// Create specledger.mod file as project artifact (empty manifest for now)
	specledgerMod := fmt.Sprintf(`# SpecLedger Dependency Manifest v1.0.0
# Generated by sl new on %s
# Project: %s
# Short Code: %s
#
# To add dependencies, use:
#   sl deps add git@github.com:org/spec main spec.md --alias alias

`, time.Now().Format("2006-01-02"), projectName, shortCode)

	return os.WriteFile(filepath.Join(projectPath, "specledger", "specledger.mod"), []byte(specledgerMod), 0644)
}

// initializeGitRepo initializes a git repository in the project directory
// Note: Only runs git init and git add, does NOT commit to support bootstrapping into existing repos
func initializeGitRepo(projectPath string) error {
	// Run git init
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init failed: %w\nOutput: %s", err, string(output))
	}

	// Run git add . to stage new files (ignore errors for existing repos)
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectPath
	_, _ = cmd.CombinedOutput() // Ignore errors - user might have custom .gitignore

	return nil
}
