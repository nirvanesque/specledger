# Research: Project Template & Coding Agent Selection

**Feature**: 593-init-project-templates
**Date**: 2026-02-20
**Status**: Complete

This document consolidates research findings for implementing template and agent selection in `sl new`.

---

## Prior Work

### Feature 011-streamline-onboarding (TUI Implementation)

**Location**: `pkg/cli/tui/sl_new.go`

**Key Insights**:
- Implemented 7-step TUI flow using Bubble Tea framework (Bubbles, Lipgloss)
- Pattern: `Model` struct with `step` field, `Update()` for events, `View()` for rendering
- Text input steps: project name, directory, short code
- List selection steps: playbook, constitution, agent preference with `›` cursor + `◉`/`○` buttons
- Confirmation review step: displays all selections before creating project
- Navigation: Enter (next), Ctrl+C (quit), ↑/↓ (selection)
- Styling: Lipgloss colors (gold #13 primary, green success, red error)

**Reusable Components**:
- `textInput` from Bubbles (`github.com/charmbracelet/bubbles/textinput`)
- Radio selection pattern (selection index + arrow key navigation)
- Step-based state machine (stepProjectName, stepDirectory, etc.)
- Answers map (`map[string]string`) for collecting user input across steps

**Integration Point**: Add two new steps (stepTemplate, stepAgent) in existing TUI flow.

### Feature 005-embedded-templates (Template System)

**Location**: `pkg/cli/playbooks/`, `pkg/embedded/templates/`

**Key Insights**:
- Well-architected playbook system with `PlaybookSource` interface
- `EmbeddedSource` loads templates from embedded filesystem (`embed.FS`)
- Manifest-driven: YAML manifest (`manifest.yaml`) defines template metadata
- Copy strategy: Pattern matching with glob patterns, skip existing by default
- Templates compiled into binary via Go's `embed` package (`//go:embed all:templates`)

**Current Manifest Structure**:
```yaml
version: "1.0"
playbooks:
  - name: specledger
    description: "SpecLedger playbook..."
    version: "1.0.0"
    path: "specledger"
    patterns: ["**"]
    structure: [".claude/", ".gitattributes", "AGENTS.md", "mise.toml", ".specledger/"]
```

**Extension Path**: Add template definitions to manifest with directory structures per template type.

**Files**:
- `pkg/cli/playbooks/template.go` - Core interfaces (PlaybookSource, Playbook)
- `pkg/cli/playbooks/embedded.go` - EmbeddedSource implementation
- `pkg/cli/playbooks/copy.go` - File copying logic with pattern matching
- `pkg/embedded/templates/manifest.yaml` - Current manifest

### Feature 004-thin-wrapper-redesign (Metadata System)

**Location**: `pkg/cli/metadata/`

**Current Schema** (`schema.go`):
```go
type ProjectMetadata struct {
    Version      string           `yaml:"version"`       // Currently "1.0.0"
    Project      ProjectInfo      `yaml:"project"`
    Playbook     PlaybookInfo     `yaml:"playbook"`
    TaskTracker  TaskTrackerInfo  `yaml:"task_tracker,omitempty"`
    Dependencies []Dependency     `yaml:"dependencies,omitempty"`
}

type ProjectInfo struct {
    Name      string    `yaml:"name"`
    ShortCode string    `yaml:"short_code"`
    Created   time.Time `yaml:"created"`
    Modified  time.Time `yaml:"modified"`
    Version   string    `yaml:"version"`  // Project version
}
```

**YAML Library**: `gopkg.in/yaml.v3` (v3.0.1)
- Standard marshaling/unmarshaling via `yaml.Marshal()` and `yaml.Unmarshal()`
- No custom `MarshalYAML`/`UnmarshalYAML` implementations needed
- Validation via `Validate()` method after unmarshal, before marshal
- Auto-update of `Modified` timestamp before saving

**Extension Required**: Add `ID`, `Template`, `Agent` fields to `ProjectInfo` struct, bump version to "1.1.0".

---

## Research Findings

### Decision 1: UUID Generation for Project ID

**Decision**: Use `github.com/google/uuid` package

**Rationale**:
- Industry standard Go UUID library maintained by Google
- Based on RFC 9562 and DCE 1.1 standards
- Uses `crypto/rand` for cryptographically secure randomness
- Built-in YAML marshaling via `MarshalText()`/`UnmarshalText()` interfaces
- Zero additional code needed for YAML integration

**Implementation**:
```go
import "github.com/google/uuid"

// Generate UUID v4
projectID := uuid.New()  // Simple version, panics on error (acceptable for new projects)

// Or explicit error handling
projectID, err := uuid.NewRandom()
if err != nil {
    return fmt.Errorf("failed to generate UUID: %w", err)
}
```

**YAML Format**:
```yaml
project:
    id: 550e8400-e29b-41d4-a716-446655440000
    name: my-project
```

**Schema Extension**:
```go
type ProjectInfo struct {
    ID        uuid.UUID `yaml:"id"`                   // NEW
    Name      string    `yaml:"name"`
    ShortCode string    `yaml:"short_code"`
    Template  string    `yaml:"template,omitempty"`   // NEW
    Agent     string    `yaml:"agent,omitempty"`      // NEW
    Created   time.Time `yaml:"created"`
    Modified  time.Time `yaml:"modified"`
    Version   string    `yaml:"version"`
}
```

**Alternatives Considered**:
- `github.com/gofrs/uuid`: Older, less actively maintained
- `github.com/satori/go.uuid`: Deprecated, points to gofrs/uuid
- Manual UUID generation: Reinventing the wheel, potential security issues

**Installation**: `go get github.com/google/uuid@v1.6.0`

### Decision 2: Template Structures Based on Industry Best Practices

#### Template 1: General Purpose (Go CLI/Library)

**Directory Structure**:
```
project-name/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── services/
│   └── models/
├── pkg/
├── configs/
├── scripts/
├── tests/
│   ├── integration/
│   └── unit/
├── .gitignore
├── go.mod
├── Makefile
├── README.md
└── LICENSE
```

**Rationale**: Based on `golang-standards/project-layout` (most referenced Go structure), follows Go compiler conventions for `internal/` (private) and `pkg/` (public).

**Key Files**:
- `Makefile`: Automate testing, linting, building (standard in Go community)
- `cmd/`: Small main functions, imports from internal/pkg
- `internal/`: Private application code (enforced by Go compiler)

#### Template 2: Full-Stack Application (Go + TypeScript/React)

**Directory Structure**:
```
project-name/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── models/
│   │   ├── middleware/
│   │   └── database/
│   ├── go.mod
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── hooks/
│   │   └── App.tsx
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── .eslintrc.js
├── docker-compose.yml
└── README.md
```

**Rationale**: Separates frontend/backend cleanly. Uses Vite (modern build tool, 40x faster than CRA). Go backend with Gin framework (most popular). Includes docker-compose for local development with services.

**Key Technologies**:
- Backend: Go + Gin (web framework) + GORM (ORM)
- Frontend: React 18+ + TypeScript + Vite + TanStack Query
- Docker Compose: Multi-service orchestration

#### Template 3: Batch Data Processing

**Directory Structure**:
```
project-name/
├── workflows/        # Temporal workflows (Go)
│   ├── workflow.go
│   └── activity.go
├── cmd/
│   ├── worker/main.go
│   └── starter/main.go
├── internal/
│   ├── extractors/
│   ├── transformers/
│   └── loaders/
├── config/
├── tests/
├── docker-compose.yml
└── README.md
```

**Rationale**: Uses Temporal for Go-based batch processing (better type safety than Airflow, native Go support). Separates extract/transform/load logic. Temporal provides durability, retry, and state management.

**Key Technologies**:
- Temporal Go SDK: Workflow orchestration
- Temporal CLI: Development and debugging
- PostgreSQL: Temporal's persistence

#### Template 4: Real-Time Workflow (Temporal)

**Directory Structure**:
```
project-name/
├── cmd/
│   ├── worker/main.go
│   └── starter/main.go
├── internal/
│   ├── workflows/
│   ├── activities/
│   ├── models/
│   └── config/
├── tests/
├── docker-compose.yml
└── README.md
```

**Rationale**: Dedicated Temporal workflow structure. Separates workflows (orchestration) from activities (execution). Includes worker and starter applications.

**Key Technologies**:
- Temporal Go SDK: Workflow framework
- Docker Compose: Temporal server + PostgreSQL + Elasticsearch

#### Template 5: ML Image Processing

**Directory Structure**:
```
project-name/
├── src/
│   ├── data/
│   │   ├── dataset.py
│   │   ├── preprocessing.py
│   │   └── augmentation.py
│   ├── models/
│   │   └── cnn_model.py
│   ├── training/
│   │   ├── train.py
│   │   └── evaluate.py
│   ├── inference/
│   │   └── predict.py
│   └── utils/
├── data/
│   ├── raw/
│   ├── processed/
│   └── interim/
├── models/checkpoints/
├── notebooks/
├── configs/
├── tests/
├── pyproject.toml
├── requirements.txt
└── README.md
```

**Rationale**: Based on Cookiecutter Data Science structure (industry standard). Separates data loading, model definition, training, and inference. PyTorch-centric (most popular in 2026).

**Key Technologies**:
- PyTorch + PyTorch Lightning: Deep learning framework
- TorchVision: Pre-trained models
- Weights & Biases: Experiment tracking
- OpenCV: Image processing

#### Template 6: Real-Time Data Pipeline (Kafka)

**Directory Structure**:
```
project-name/
├── cmd/
│   ├── producer/main.go
│   ├── consumer/main.go
│   └── processor/main.go
├── internal/
│   ├── kafka/
│   │   ├── producer.go
│   │   ├── consumer.go
│   │   └── config.go
│   ├── handlers/
│   ├── models/
│   └── processors/
├── configs/
├── deployments/
│   └── docker-compose.yml
├── tests/
├── go.mod
└── README.md
```

**Rationale**: Uses segmentio/kafka-go (pure Go, best performance). Separates producer, consumer, and stream processor. Includes Kafka cluster setup via docker-compose.

**Key Technologies**:
- segmentio/kafka-go: Kafka client library
- Goka: Stream processing framework
- Docker Compose: Kafka + Zookeeper + Schema Registry

#### Template 7: AI Chatbot (Multi-Platform)

**Directory Structure**:
```
project-name/
├── src/
│   ├── agents/
│   ├── tools/
│   ├── prompts/
│   ├── chains/
│   ├── memory/
│   ├── middleware/
│   ├── integrations/
│   │   ├── slack/
│   │   ├── discord/
│   │   ├── telegram/
│   │   └── web/
│   ├── vectorstore/
│   └── utils/
├── data/
│   └── documents/
├── tests/
├── configs/
├── pyproject.toml
├── requirements.txt
├── langgraph.json
└── README.md
```

**Rationale**: Based on LangChain/LangGraph patterns (most popular AI framework). Unified architecture with platform adapters. Supports RAG with vector databases.

**Key Technologies**:
- LangChain + LangGraph: AI framework
- OpenAI/Anthropic: LLM providers
- Pinecone/Weaviate: Vector database
- Platform SDKs: Slack, Discord, Telegram

### Decision 3: Agent Configuration Formats

**Decision**: Create separate configuration directories for each agent

**Agent Options**:

1. **Claude Code** (Default):
   - Directory: `.claude/`
   - Files: `commands/`, `skills/`, `settings.json`
   - Session capture: Configured in `settings.json`

2. **OpenCode**:
   - Directory: `.opencode/`
   - Files: `commands/`, `skills/`, `opencode.json`
   - Backward compatible: Can read `.claude/` directories

3. **None**:
   - No agent-specific directories
   - Only creates `AGENTS.md` for context

**Claude Code Settings Structure**:
```json
{
  "saveTranscripts": true,
  "transcriptsDirectory": "~/.claude/sessions",
  "hooks": {
    "PostToolUse": [{
      "matcher": "Bash",
      "hooks": [{"type": "command", "command": "sl session capture"}]
    }]
  }
}
```

**Rationale**: Separates agent concerns, allows future agent additions, maintains backward compatibility.

### Decision 4: Metadata Schema Version 1.1.0

**Changes from 1.0.0**:
- Add `project.id` (uuid.UUID): Unique project identifier
- Add `project.template` (string, omitempty): Selected template ID
- Add `project.agent` (string, omitempty): Selected agent ID

**Migration Strategy**:
- Old projects (v1.0.0) without UUID: Auto-generate UUID on first load
- Update version to 1.1.0 automatically
- Backward compatible: Old fields remain unchanged, new fields optional with `omitempty`

**Example v1.1.0 YAML**:
```yaml
version: 1.1.0
project:
    id: 550e8400-e29b-41d4-a716-446655440000
    name: my-project
    short_code: mp
    template: full-stack
    agent: claude-code
    created: 2026-02-20T10:30:00Z
    modified: 2026-02-20T10:30:00Z
    version: 0.1.0
```

---

## Implementation Checklist

### Phase 0: Foundation (Complete)
- ✅ Research TUI patterns from feature 011
- ✅ Research template system from feature 005
- ✅ Research metadata system from feature 004
- ✅ Research UUID generation libraries
- ✅ Research industry best practices for 7 template types

### Phase 1: Data Model & Metadata (Next)
- [ ] Add `github.com/google/uuid` dependency
- [ ] Extend `ProjectInfo` struct with `ID`, `Template`, `Agent` fields
- [ ] Create `TemplateDefinition` struct in `pkg/models/`
- [ ] Create `AgentConfig` struct in `pkg/models/`
- [ ] Update `metadata.Validate()` to validate new fields
- [ ] Update metadata version to "1.1.0"
- [ ] Implement migration from v1.0.0 to v1.1.0

### Phase 2: Template System
- [ ] Define 7 template directories in `pkg/embedded/templates/`:
  - `general-purpose/` (copy current specledger playbook)
  - `full-stack/` (backend/ + frontend/ dirs)
  - `batch-data/` (workflows/, cmd/worker, cmd/starter)
  - `realtime-workflow/` (workflows/, activities/)
  - `ml-image/` (src/data, src/models, src/training)
  - `realtime-data/` (cmd/producer, cmd/consumer, internal/kafka)
  - `ai-chatbot/` (src/agents, src/tools, src/integrations)
- [ ] Extend `manifest.yaml` with template definitions
- [ ] Update `EmbeddedSource` to load template list from manifest
- [ ] Create README files for each template
- [ ] Create starter files for each template type

### Phase 3: Agent Configuration
- [ ] Create `supportedAgents` slice in `pkg/models/agent.go`
- [ ] Create OpenCode template structure in templates with:
  - `.opencode/commands/` (port from `.claude/commands/`)
  - `.opencode/skills/` (port from `.claude/skills/`)
  - `opencode.json` with schema reference
  - `AGENTS.md` (shared with Claude Code)
- [ ] Implement `.claude/settings.json` generation with project UUID

### Phase 4: TUI Updates
- [ ] Add `stepTemplate` constant
- [ ] Add `stepAgent` constant (after playbook step)
- [ ] Add template fields to `Model`: `templates []TemplateDefinition`, `selectedTemplateIndex int`
- [ ] Add agent fields to `Model`: `agents []AgentConfig`, `selectedAgentIndex int`
- [ ] Implement `Update()` handler for `stepTemplate` (arrow keys + Enter)
- [ ] Implement `View()` renderer for `stepTemplate`
- [ ] Implement `Update()` handler for `stepAgent`
- [ ] Implement `View()` renderer for `stepAgent`
- [ ] Update confirmation review to display template and agent selections

### Phase 5: Bootstrap Integration
- [ ] Update `bootstrap.go` to read template/agent from answers map
- [ ] Implement template directory copying based on selection
- [ ] Implement agent config directory creation based on selection
- [ ] Generate UUID and write to `specledger.yaml`
- [ ] Record template and agent in `specledger.yaml`
- [ ] Create `.claude/settings.json` with session capture hooks when Claude Code selected

### Phase 6: CLI Flags
- [ ] Add `--template <id>` flag to `sl new`
- [ ] Add `--agent <id>` flag to `sl new`
- [ ] Add `--list-templates` flag to `sl new`
- [ ] Implement non-interactive flow (skip TUI when flags provided)
- [ ] Validate flag values against supported options
- [ ] Implement TTY detection and require flags in non-interactive mode

### Phase 7: Testing
- [ ] Create `pkg/cli/tui/sl_new_test.go` with unit tests
- [ ] Create `tests/integration/bootstrap_tui_test.go` with integration tests
- [ ] Test all 7 templates create correct structures
- [ ] Test all 3 agents create correct config directories
- [ ] Test backward compatibility (general-purpose + claude-code = current behavior)
- [ ] Test non-interactive mode with flags
- [ ] Test TTY detection and error messages

---

## References

**Dependencies**:
- `github.com/google/uuid v1.6.0` - UUID generation (NEW)
- `gopkg.in/yaml.v3` - YAML marshaling (already in use)
- `github.com/charmbracelet/bubbletea v1.3.10` - TUI framework (already in use)
- `github.com/charmbracelet/bubbles v0.21.1` - TUI components (already in use)
- `github.com/charmbracelet/lipgloss v1.1.0` - TUI styling (already in use)

**External Resources**:
- [google/uuid Package Documentation](https://pkg.go.dev/github.com/google/uuid)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) - Go project structure
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Temporal Go SDK](https://github.com/temporalio/sdk-go)
- [LangChain Documentation](https://python.langchain.com/)

**Internal References**:
- Feature 011-streamline-onboarding: TUI implementation patterns
- Feature 005-embedded-templates: Template embedding and copying
- Feature 004-thin-wrapper-redesign: Metadata schema and validation
- `pkg/cli/tui/sl_new.go` - Current TUI implementation
- `pkg/cli/metadata/schema.go` - Current metadata schema
- `pkg/cli/playbooks/template.go` - Playbook system interfaces
