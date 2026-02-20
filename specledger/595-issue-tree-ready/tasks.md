# Tasks Index: Issue Tree View and Ready Command

Issue tracking for this feature implementation.
This index provides navigation - all tasks are managed through `sl issue` CLI.

## Feature Tracking

* **Epic ID**: `SL-00d561`
* **User Stories Source**: `specledger/595-issue-tree-ready/spec.md`
* **Research Inputs**: `specledger/595-issue-tree-ready/research.md`
* **Planning Details**: `specledger/595-issue-tree-ready/plan.md`
* **Data Model**: `specledger/595-issue-tree-ready/data-model.md`
* **Validation Scenarios**: `specledger/595-issue-tree-ready/quickstart.md`

## Query Commands

```bash
# Find all tasks for this feature
sl issue list --label spec:595-issue-tree-ready

# Find ready tasks (not blocked)
sl issue ready --label spec:595-issue-tree-ready

# View dependency tree
sl issue show SL-00d561 --tree

# Filter by phase
sl issue list --label phase:foundational
sl issue list --label phase:us1
sl issue list --label phase:us2
sl issue list --label phase:us3
sl issue list --label phase:us4
sl issue list --label phase:polish

# Filter by story
sl issue list --label story:US1
sl issue list --label story:US2
sl issue list --label story:US3
sl issue list --label story:US4
```

## Phases and Tasks

### Phase: Foundational (Priority: 0)

**Feature**: SL-be8099 - Ready State Computation
**Purpose**: Core infrastructure that MUST complete before user stories

| ID | Title | Status |
|----|-------|--------|
| SL-b75157 | Add IsReady() method to Issue entity | open |
| SL-ab3620 | Add GetBlockers() method to Issue entity | open |
| SL-f06199 | Add ListReady() method to Store | open |

**Checkpoint**: Foundation ready → US1 and US2 can begin

---

### Phase: US1 - View Issue Dependencies as Tree (Priority: P1) 🎯 MVP

**Feature**: SL-939c5d
**Goal**: Display issues in hierarchical tree format with ASCII characters
**Independent Test**: Create issues with dependencies, verify tree display

| ID | Title | Dependencies |
|----|-------|--------------|
| SL-e05db0 | Create TreeRenderer in pkg/issues/tree.go | Foundational |
| SL-0f6a9c | Implement Render() and RenderForest() methods | SL-e05db0 |
| SL-26f469 | Add cycle detection and warning in tree output | SL-0f6a9c |
| SL-62f51c | Implement --tree flag handling in runIssueList | SL-26f469 |

**Checkpoint**: Tree view functional → Can visualize dependencies

---

### Phase: US2 - List Ready-to-Work Issues (Priority: P1) 🎯 MVP

**Feature**: SL-40f7df
**Goal**: Command to list only unblocked issues for quick task selection
**Independent Test**: Create blocked/unblocked issues, verify ready list

| ID | Title | Dependencies |
|----|-------|--------------|
| SL-f290ad | Add issueReadyCmd to CLI | Foundational |
| SL-ad1faa | Implement runIssueReady function | SL-f290ad, SL-f06199 |

**Checkpoint**: Ready command functional → Can identify unblocked work

---

### Phase: US3 - View Single Issue Dependency Tree (Priority: P2)

**Feature**: SL-044510
**Goal**: Show dependency context for specific issue (blocks/blocked_by)
**Independent Test**: Run sl issue show <id> --tree, verify context display

| ID | Title | Dependencies |
|----|-------|--------------|
| SL-eed51a | Enhance runIssueShow for --tree flag | SL-62f51c (tree renderer) |

**Checkpoint**: Single issue tree view functional

---

### Phase: US4 - Auto-Select Ready Tasks in Implement Workflow (Priority: P2)

**Feature**: SL-e5e420
**Goal**: /specledger.implement uses ready state for task selection
**Independent Test**: Run implement workflow, verify only unblocked tasks shown

| ID | Title | Dependencies |
|----|-------|--------------|
| SL-687a6e | Update specledger.implement.md to use sl issue ready | SL-ad1faa (ready command) |

**Checkpoint**: Implement workflow integrated with ready state

---

### Phase: Polish (Priority: 3)

**Feature**: SL-0dacde
**Purpose**: Documentation, validation, and cleanup

| ID | Title | Dependencies |
|----|-------|--------------|
| SL-e01a49 | Run quickstart.md validation tests | All user stories |
| SL-15bfe3 | Update CLAUDE.md with feature details | All user stories |

---

## Dependency Graph

```
Foundational
├── SL-b75157 (IsReady) ─────────────────────┐
├── SL-ab3620 (GetBlockers)                  │
└── SL-f06199 (ListReady) ◄──────────────────┤
    │                                        │
    ▼                                        ▼
US1: Tree View                         US2: Ready Command
├── SL-e05db0 (TreeRenderer)           ├── SL-f290ad (issueReadyCmd)
├── SL-0f6a9c (Render methods)         └── SL-ad1faa (runIssueReady)
├── SL-26f469 (Cycle detection)
└── SL-62f51c (--tree in list) ────────┐
    │                                  │
    ▼                                  │
US3: Single Issue Tree                 │
└── SL-eed51a (show --tree) ◄──────────┘
    │
    ▼
US4: Implement Integration
└── SL-687a6e (Update implement.md)
    │
    ▼
Polish
├── SL-e01a49 (Validation)
└── SL-15bfe3 (Documentation)
```

## Definition of Done Summary

| Issue ID | Title | DoD Items |
|----------|-------|-----------|
| SL-b75157 | IsReady() method | Open+no blockers→true, Open+closed blockers→true, Open+open blockers→false, Closed→false |
| SL-ab3620 | GetBlockers() method | Returns Blocker structs, Skips missing refs, Empty for no blockers |
| SL-f06199 | ListReady() method | Returns IsReady issues, Filter works, Empty slice ok |
| SL-e05db0 | TreeRenderer | Created with defaults, Options configurable |
| SL-0f6a9c | Render methods | Correct indentation, Connecting lines, Title truncation |
| SL-26f469 | Cycle detection | No crash, Warning displayed, Path shown |
| SL-62f51c | --tree in list | Hierarchical output, --all groups by spec, Filters work |
| SL-f290ad | issueReadyCmd | Command exists, Help text, Registered |
| SL-ad1faa | runIssueReady | Lists ready only, --all works, --json works, Blocked message |
| SL-eed51a | show --tree | Dependency context, Blocks/blocked_by, Standalone for no deps |
| SL-687a6e | Implement update | Uses ready command, Only unblocked shown, Blocked handling |
| SL-e01a49 | Validation | All 7 scenarios pass |
| SL-15bfe3 | Documentation | CLAUDE.md updated |

## MVP Scope

**Recommended MVP**: Complete Foundational + US1 + US2

This delivers:
- Ready state computation (IsReady, GetBlockers, ListReady)
- Tree view visualization (`sl issue list --tree`)
- Ready command (`sl issue ready`)

Total MVP tasks: 7 (SL-b75157, SL-ab3620, SL-f06199, SL-e05db0, SL-0f6a9c, SL-26f469, SL-62f51c, SL-f290ad, SL-ad1faa)

## Execution Strategy

1. **Foundational first** - Complete all 3 foundational tasks
2. **US1 and US2 in parallel** - Can be worked simultaneously after foundational
3. **US3 after US1** - Requires tree renderer
4. **US4 after US2** - Requires ready command
5. **Polish last** - After all user stories complete
