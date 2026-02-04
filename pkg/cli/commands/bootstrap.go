package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/config"
	"specledger/pkg/cli/dependencies"
	"specledger/pkg/cli/logger"
	"specledger/pkg/cli/tui"
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
	Short: "Start interactive TUI for project bootstrap",
	Long: `Bootstrap a new SpecLedger project with all necessary infrastructure:
- Claude Code skills and commands
- Beads issue tracker
- SpecKit templates and scripts
- Tool configuration (mise)`,

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
	l.Info("Starting interactive bootstrap...")

	// Check dependencies
	r := dependencies.New()
	hasGum := r.CheckGum()
	hasMise := r.CheckMise()

	l.Debug(fmt.Sprintf("Dependencies: gum=%v, mise=%v", hasGum, hasMise))

	// If gum is missing, prompt user
	if !hasGum {
		shouldInstall, err := r.PromptForInstall(r.GetGumDep())
		if err != nil {
			return fmt.Errorf("failed to check for gum: %w", err)
		}
		if shouldInstall {
			l.Info(fmt.Sprintf("Installing %s...", r.GetGumDep().Name))
			if _, err := r.Install(r.GetGumDep()); err != nil {
				return fmt.Errorf("failed to install gum: %w", err)
			}
			l.Info("Gum installed successfully")
		}
	}

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

	l.Info(fmt.Sprintf("Selected: project=%s, dir=%s, code=%s, playbook=%s, shell=%s",
		projectName, projectDir, shortCode, playbook, shell))

	l.Info(fmt.Sprintf("Project: %s (short code: %s, playbook: %s, shell: %s)",
		projectName, shortCode, playbook, shell))

	// Create project path
	projectPath := filepath.Join(projectDir, projectName)
	l.Debug(fmt.Sprintf("Creating project at: %s", projectPath))

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
	l.Info(fmt.Sprintf("Created directory: %s", projectPath))

	// TODO: Copy template files
	// TODO: Initialize git repo
	// TODO: Configure Beads
	// TODO: Install tools

	// Success message
	fmt.Printf("\n✓ Project created: %s\n", projectPath)
	fmt.Printf("✓ Beads prefix: %s\n", shortCode)
	fmt.Printf("✓ Playbook: %s\n", playbook)
	fmt.Printf("✓ Agent Shell: %s\n", shell)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Printf("  claude\n")

	return nil
}

// runBootstrapNonInteractive runs bootstrap without TUI
func runBootstrapNonInteractive(cmd *cobra.Command, l *logger.Logger, cfg *config.Config) error {
	l.Info("Running bootstrap in non-interactive mode...")

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
	l.Debug(fmt.Sprintf("Creating project at: %s", projectPath))

	// Check if directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return ErrProjectExists(projectName)
	}

	// Create directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return ErrPermissionDenied(projectPath)
	}
	l.Info(fmt.Sprintf("Created directory: %s", projectPath))

	// TODO: Copy template files
	// TODO: Initialize git repo
	// TODO: Configure Beads
	// TODO: Install tools

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
