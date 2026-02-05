package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"specledger/pkg/cli/prerequisites"
)

var (
	doctorJSONOutput bool
)

// VarDoctorCmd represents the doctor command
var VarDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check installation status of required and optional tools",
	Long: `Check the installation status of all tools required by SpecLedger.

This command verifies that:
- Core tools (mise, bd, perles) are installed and accessible
- Framework tools (specify, openspec) are installed (optional)

Use --json flag for machine-readable output suitable for CI/CD pipelines.`,
	Example: `  sl doctor           # Human-readable output
  sl doctor --json    # JSON output for CI/CD`,
	RunE: runDoctor,
}

func init() {
	VarDoctorCmd.Flags().BoolVar(&doctorJSONOutput, "json", false, "Output results in JSON format")
}

// DoctorOutput represents the JSON output structure for doctor command
type DoctorOutput struct {
	Status              string              `json:"status"`
	Tools               []DoctorToolStatus  `json:"tools"`
	Missing             []string            `json:"missing,omitempty"`
	InstallInstructions string              `json:"install_instructions,omitempty"`
}

// DoctorToolStatus represents a tool's status in JSON output
type DoctorToolStatus struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version,omitempty"`
	Path      string `json:"path,omitempty"`
	Category  string `json:"category"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	check := prerequisites.CheckPrerequisites()

	if doctorJSONOutput {
		return outputDoctorJSON(check)
	}

	return outputDoctorHuman(check)
}

func outputDoctorJSON(check prerequisites.PrerequisiteCheck) error {
	output := DoctorOutput{
		Status: "pass",
		Tools:  []DoctorToolStatus{},
	}

	// Add all tools to output
	for _, result := range check.CoreResults {
		output.Tools = append(output.Tools, DoctorToolStatus{
			Name:      result.Tool.Name,
			Installed: result.Installed,
			Version:   result.Version,
			Path:      result.Path,
			Category:  string(result.Tool.Category),
		})
	}

	for _, result := range check.FrameworkResults {
		output.Tools = append(output.Tools, DoctorToolStatus{
			Name:      result.Tool.Name,
			Installed: result.Installed,
			Version:   result.Version,
			Path:      result.Path,
			Category:  string(result.Tool.Category),
		})
	}

	// Set status and missing tools
	if !check.AllCoreInstalled {
		output.Status = "fail"
		output.Missing = []string{}
		for _, tool := range check.MissingCore {
			output.Missing = append(output.Missing, tool.Name)
		}
		output.InstallInstructions = check.Instructions
	}

	// Marshal and print JSON
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}

func outputDoctorHuman(check prerequisites.PrerequisiteCheck) error {
	fmt.Println()
	fmt.Println("SpecLedger Tool Status")
	fmt.Println("======================")
	fmt.Println()

	// Core tools section
	fmt.Println("Core Tools (Required):")
	for _, result := range check.CoreResults {
		fmt.Printf("  %s\n", prerequisites.FormatToolStatus(result))
	}
	fmt.Println()

	// Framework tools section
	fmt.Println("Framework Tools (Optional):")
	for _, result := range check.FrameworkResults {
		fmt.Printf("  %s\n", prerequisites.FormatToolStatus(result))
	}
	fmt.Println()

	// Overall status
	if check.AllCoreInstalled {
		fmt.Println("✅ All core tools are installed!")
		fmt.Println()

		// Check if any framework tools are installed
		anyFrameworkInstalled := false
		for _, result := range check.FrameworkResults {
			if result.Installed {
				anyFrameworkInstalled = true
				break
			}
		}

		if !anyFrameworkInstalled {
			fmt.Println("ℹ️  No SDD framework tools detected.")
			fmt.Println("   To use Spec Kit or OpenSpec, install via mise:")
			fmt.Println()
			for _, tool := range prerequisites.FrameworkTools {
				fmt.Printf("   %s\n", tool.InstallCmd)
			}
			fmt.Println()
		}

		return nil
	}

	// Missing tools
	fmt.Println("❌ Missing required tools:")
	fmt.Println()
	fmt.Println(check.Instructions)

	return fmt.Errorf("missing required tools")
}
