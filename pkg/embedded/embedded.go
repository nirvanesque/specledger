package embedded

import (
	"embed"
)

//go:embed templates/.beads
//go:embed templates/.claude
//go:embed templates/AGENTS.md
//go:embed templates/mise.toml
//go:embed templates/specledger
//go:embed templates/.gitattributes
var TemplatesFS embed.FS
