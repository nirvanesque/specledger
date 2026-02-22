# Tasks Index: Issue Create Fields Enhancement

Issue graph index into the tasks and phases for this feature implementation.
This index does **not contain tasks directly**—those are fully managed through `sl issue` CLI.

## Feature Tracking

* **Epic ID**: `SL-bb6f5b`
* **User Stories Source**: `specledger/597-issue-create-fields/spec.md`
* **Research Inputs**: `specledger/597-issue-create-fields/research.md`
* **Planning Details**: `specledger/597-issue-create-fields/plan.md`
* **Data Model**: `specledger/597-issue-create-fields/data-model.md`
* **CLI Contracts**: `specledger/597-issue-create-fields/contracts/cli-contract.md`

## Query Hints

Use the `sl issue` CLI to query and manipulate the issue graph:

```bash
# Find all open tasks for this feature
sl issue list --label spec:597-issue-create-fields --status open

# Find ready tasks to implement (not blocked)
sl issue ready --label spec:597-issue-create-fields

# See dependencies for issue
sl issue show --tree SL-xxxxxx

# View issues by story
sl issue list --label 'spec:597-issue-create-fields' --label 'story:US1'

# View issues by component
sl issue list --label 'spec:597-issue-create-fields' --label 'component:cli'

# Define dependencies
sl issue link SL-xxxxx blocks SL-yyyyy
```

## Tasks and Phases Structure

This feature follows a phase-based structure:

* **Epic**: SL-bb6f5b → Issue Create Fields Enhancement
* **Phases**: Issues of type `feature`, tracked by this epic
  * Phase = a user story group or technical milestone

## Convention Summary

| Type    | Description                  | Labels                                    |
| ------- | ---------------------------- | ----------------------------------------- |
| epic    | Full feature epic            | `spec:597-issue-create-fields`            |
| feature | Implementation phase / story | `phase:<name>`, `story:<US#>`             |
| task    | Implementation task          | `component:<area>`, `fr:<requirement-id>` |

## Phases Overview

### Phase 1: US1 - Create Issues with Complete Field Set (P1)

**Feature ID**: SL-05d112
**Goal**: Add --acceptance-criteria, --dod, --design, --notes flags to sl issue create

| Issue ID | Title | Blocked By |
|----------|-------|------------|
| SL-225aa7 | Add create flags to issueCreateCmd | - |
| SL-7b3e32 | Populate issue fields from create flags | SL-225aa7 |
| SL-cf082a | Update issue show display for new fields | - |

### Phase 2: US2 - Update Definition of Done (P2)

**Feature ID**: SL-0dd751
**Goal**: Add --dod, --check-dod, --uncheck-dod flags to sl issue update
**Blocked By**: US1 (SL-05d112)

| Issue ID | Title | Blocked By |
|----------|-------|------------|
| SL-3c9872 | Add DoD update flags to issueUpdateCmd | - |
| SL-5ba439 | Handle DoD operations in runIssueUpdate | SL-3c9872 |

### Phase 3: US3+US4 - Tasks Prompt Updates (P2)

**Feature ID**: SL-d42fee
**Goal**: Update tasks prompt to use new CLI flags and improve blocking instructions
**Blocked By**: US2 (SL-0dd751)

| Issue ID | Title | Blocked By |
|----------|-------|------------|
| SL-fe19e2 | Update CLI examples in tasks prompt | - |
| SL-de1a02 | Add design field population instruction | - |
| SL-c141f0 | Improve blocking relationship instructions | - |

### Phase 4: US5 - Implement Prompt Updates (P3)

**Feature ID**: SL-a594b9
**Goal**: Update implement prompt to utilize design, AC, and DoD fields
**Blocked By**: US2 (SL-0dd751)

| Issue ID | Title | Blocked By |
|----------|-------|------------|
| SL-dfa042 | Add field reading instructions to implement prompt | - |
| SL-909448 | Add AC verification and DoD check instructions | SL-dfa042 |

### Phase 5: Polish - Tests (P3)

**Feature ID**: SL-f17a81
**Goal**: Add unit tests for CLI flag changes
**Blocked By**: US1, US2, US3+US4, US5

| Issue ID | Title | Blocked By |
|----------|-------|------------|
| SL-9f721e | Add unit tests for CLI flag changes | All implementation tasks |

## Dependency Graph

```
Epic: SL-bb6f5b
├── US1: SL-05d112 (Create fields)
│   ├── SL-225aa7 → SL-7b3e32
│   └── SL-cf082a (parallel)
│
├── US2: SL-0dd751 (Update DoD) [blocked by US1]
│   └── SL-3c9872 → SL-5ba439
│
├── US3+US4: SL-d42fee (Tasks prompt) [blocked by US2]
│   ├── SL-fe19e2 (parallel)
│   ├── SL-de1a02 (parallel)
│   └── SL-c141f0 (parallel)
│
├── US5: SL-a594b9 (Implement prompt) [blocked by US2]
│   └── SL-dfa042 → SL-909448
│
└── Polish: SL-f17a81 (Tests) [blocked by all above]
    └── SL-9f721e
```

## Parallel Execution Opportunities

**Within US1**: SL-225aa7 and SL-cf082a can run in parallel (different concerns)

**Within US3+US4**: All 3 tasks can run in parallel (different prompt sections)

**US3+US4 and US5**: Can run in parallel after US2 completes

## Definition of Done Summary

| Issue ID | DoD Items |
|----------|-----------|
| SL-225aa7 | 5 items: flags declared, StringArrayVar used, go build succeeds |
| SL-7b3e32 | 5 items: all fields populated, DoD items unchecked, manual test passes |
| SL-cf082a | 5 items: all sections display, empty fields omitted, manual test passes |
| SL-3c9872 | 4 items: all flags declared, go build succeeds |
| SL-5ba439 | 5 items: DoD replace/check/uncheck work, error format correct, manual tests pass |
| SL-fe19e2 | 5 items: all flags in examples, description simplified, all 3 files updated |
| SL-de1a02 | 3 items: instruction added, references plan.md, all files updated |
| SL-c141f0 | 6 items: all blocking rules documented, examples added, all files updated |
| SL-dfa042 | 3 items: read design/AC instructions, both files updated |
| SL-909448 | 4 items: verify AC instruction, --check-dod instruction, both files updated |
| SL-9f721e | 8 items: all test cases pass, make test passes |

## Implementation Strategy

### MVP Scope (Recommended)

Implement US1 only for minimum viable functionality:
- Create issues with all 4 new flags
- Display fields in issue show

This enables immediate value: users can create fully-specified issues in one command.

### Incremental Delivery

1. **MVP (US1)**: CLI create flags + issue show display
2. **Iteration 2 (US2)**: DoD update operations
3. **Iteration 3 (US3+US4)**: Tasks prompt improvements
4. **Iteration 4 (US5)**: Implement prompt improvements
5. **Final (Polish)**: Test coverage

## Status Tracking

Status is tracked only in issues:

* **Open** → default
* **In Progress** → task being worked on
* **Closed** → complete

Use `sl issue ready` and `sl issue list --status open` to query progress.

---

> This file is intentionally index-only. Implementation data lives in issue tracker.
