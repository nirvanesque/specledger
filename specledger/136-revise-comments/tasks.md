# Tasks Index: Revise Comments CLI Command

Beads Issue Graph Index into the tasks and phases for this feature implementation.
This index does **not contain tasks directly**—those are fully managed through Beads CLI.

## Feature Tracking

* **Beads Epic ID**: `sl-0e7`
* **User Stories Source**: `specledger/136-revise-comments/spec.md`
* **Research Inputs**: `specledger/136-revise-comments/research.md`
* **Planning Details**: `specledger/136-revise-comments/plan.md`
* **Data Model**: `specledger/136-revise-comments/data-model.md`
* **Contract Definitions**: `specledger/136-revise-comments/contracts/`

## Beads Query Hints

```bash
# Find all open tasks for this feature
bd list --label spec:136-revise-comments --status open -n 10

# Find ready tasks to implement
bd ready --label spec:136-revise-comments -n 5

# See dependency tree
bd dep tree --reverse sl-0e7

# View by user story
bd list --label story:US1 --label spec:136-revise-comments
bd list --label story:US2 --label spec:136-revise-comments

# View by phase
bd list --type feature --label spec:136-revise-comments

# View by component
bd list --label component:revise --label spec:136-revise-comments
bd list --label component:cli --label spec:136-revise-comments
bd list --label component:tui --label spec:136-revise-comments
```

## Tasks and Phases Structure

```
sl-0e7 (epic) Revise Comments CLI Command
├── sl-ovq (feature) Phase 1: Setup
│   ├── sl-sa7 (task) Add charmbracelet/huh dependency
│   ├── sl-nd2 (task) Create pkg/cli/revise/ package with types
│   └── sl-2xb (task) Register sl revise command in main.go
│
├── sl-4xh (feature) Phase 2: Foundational [blocked by: sl-ovq]
│   ├── sl-5ao (task) PostgREST client with auth auto-retry
│   ├── sl-ci7 (task) Git helpers in pkg/deps/git.go (go-git + exec hybrid)
│   └── sl-tmq (task) Extend AgentLauncher with LaunchWithPrompt
│
├── sl-krp (feature) US1: Branch Selection & Fetching [blocked by: sl-4xh] 🎯 MVP
│   ├── sl-si4 (task) Branch detection and confirmation flow
│   └── sl-vo6 (task) Comment fetching via PostgREST query chain
│
├── sl-mkh (feature) US2: Artifact Multi-Select [blocked by: sl-krp]
│   └── sl-gkq (task) Group comments by artifact, build multi-select
│
├── sl-t8s (feature) US3: Comment Processing Loop [blocked by: sl-mkh]
│   └── sl-ccy (task) Comment display and process/skip/quit loop
│
├── sl-ni8 (feature) US4: Prompt Generation & Editor [blocked by: sl-t8s]
│   ├── sl-m98 (task) Embedded revision prompt template
│   ├── sl-8z4 (task) Token estimation and prompt rendering
│   └── sl-tv9 (task) Editor launch and confirm/re-edit flow
│
├── sl-c1h (feature) US5: Agent Launch [blocked by: sl-ni8]
│   └── sl-dc1 (task) Agent launch and post-agent state detection
│
├── sl-a7d (feature) US6: Commit/Push/Resolution [blocked by: sl-c1h]
│   ├── sl-ssr (task) File staging multi-select and commit/push
│   └── sl-x1o (task) Comment resolution multi-select and API
│
├── sl-yjc (feature) US7: Branch Checkout + Stash [blocked by: sl-4xh] ⚡ parallel
│   └── sl-ah5 (task) Stash detection, checkout, remote tracking
│
├── sl-7c7 (feature) US8: Automation Mode [blocked by: sl-ni8] ⚡ parallel
│   ├── sl-d8x (task) Fixture file parsing and comment matching
│   └── sl-cmz (task) Wire --auto and --dry-run flags
│
├── sl-7k3 (feature) US9: Summary Flag [blocked by: sl-4xh] ⚡ parallel
│   └── sl-6pk (task) --summary flag with compact output
│
└── sl-ndj (feature) Polish [blocked by: sl-a7d]
    ├── sl-0r5 (task) Edge case handling across all flows
    └── sl-2re (task) Lightweight pure-function tests
```

## Convention Summary

| Type    | Description                  | Labels                                 |
| ------- | ---------------------------- | -------------------------------------- |
| epic    | Full feature epic            | `spec:136-revise-comments`             |
| feature | Implementation phase / story | `phase:<name>`, `story:<US#>`          |
| task    | Implementation task          | `component:<x>`, `requirement:<FR-id>` |

## Dependencies & Execution Order

### Critical Path (sequential)

```
Setup → Foundational → US1 → US2 → US3 → US4 → US5 → US6 → Polish
```

### Parallel Opportunities

After Foundational completes, these can run in parallel with the critical path:
- **US7** (Branch Checkout + Stash) — enhances US1, independent implementation
- **US9** (Summary Flag) — only needs the PostgREST client

After US4 (Prompt Generation) completes:
- **US8** (Automation Mode) — needs prompt rendering but not agent launch

### Within Each Phase

- Tasks within a phase can generally run in parallel (different files)
- Go-git operations in `pkg/deps/git.go` — use go-git for most ops, exec for stash and remote tracking checkout (go-git gaps documented in sl-ci7)

## MVP Strategy

### MVP (US1 only) — Minimum Viable

```bash
bd list --label story:US1 --label spec:136-revise-comments
```

Delivers: `sl revise` fetches and displays comments for the current branch. Proves auth, API chain, and branch detection work.

### MVP+ (US1-US4) — Interactive Flow

```bash
bd list --label phase:us1 --label spec:136-revise-comments
bd list --label phase:us2 --label spec:136-revise-comments
bd list --label phase:us3 --label spec:136-revise-comments
bd list --label phase:us4 --label spec:136-revise-comments
```

Delivers: Full interactive flow from branch selection through prompt generation. User can manually copy the prompt to their agent.

### Full Feature (US1-US9)

All user stories implemented. Agent launch, commit/push, resolution, automation, and summary flag.

## Status Tracking

Status is tracked only in Beads:

* **Open** → default
* **In Progress** → task being worked on
* **Blocked** → dependency unresolved
* **Closed** → complete

```bash
bd ready --label spec:136-revise-comments -n 5
bd blocked --label spec:136-revise-comments
bd stats
```

---

> This file is intentionally light and index-only. Implementation data lives in Beads. Update this file only to point humans and agents to canonical query paths and feature references.
