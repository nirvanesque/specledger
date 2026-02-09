package embedded

import (
	"embed"
)

//go:embed skills
var SkillsFS embed.FS

// TemplatesFS provides template file system access for SpecLedger playbooks.
// It includes the templates/ directory structure and some SpecLedger-specific
// configuration files for initialization.
//
//go:embed templates templates/specledger/.claude templates/specledger/.specledger
var TemplatesFS embed.FS
