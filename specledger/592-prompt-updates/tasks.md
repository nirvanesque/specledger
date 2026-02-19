# Tasks: Improve SpecLedger Command Prompts

**Feature**: 592-prompt-updates
**Epic**: SL-e353a6
**Generated**: 2026-02-20

## Overview

Update three core SpecLedger command prompts to improve dependency handling, issue quality, and Definition of Done verification.

## Issue Summary

| Type | Count |
|------|-------|
| Epic | 1 |
| Features | 4 |
| Tasks | 10 |
| **Total** | **15** |

## Epic

**SL-e353a6** - Improve SpecLedger Command Prompts

Update three core SpecLedger command prompts (specify, tasks, implement) to improve dependency handling, issue quality, and Definition of Done verification. Changes apply to both `.claude/commands/` and `pkg/embedded/skills/commands/`.

## Phase 1: US1 - Reference Dependencies During Specification

**Feature**: SL-b7212a - US1: Reference Dependencies During Specification

**Goal**: Update specledger.specify.md to detect and load dependency references using explicit syntax.

**Independent Test**: Create a spec with `deps:api-contracts` and verify content loads from cache.

### Tasks

| ID | Title | Status | Labels |
|----|-------|--------|--------|
| SL-1f8167 | Add dependency detection to specledger.specify.md | open | story:US1, fr:FR-001, fr:FR-002, fr:FR-003 |
| SL-f39984 | Sync specledger.specify.md to embedded templates | open | story:US1, fr:FR-011 |

**Dependencies**: SL-1f8167 → SL-f39984

---

## Phase 2: US2 - Generate Descriptive, Complete Issues

**Feature**: SL-1629bc - US2: Generate Descriptive, Complete Issues

**Goal**: Update specledger.tasks.md to generate issues with structured content, DoD items, and error handling.

**Independent Test**: Generate tasks from a plan and verify each issue has complete structure with DoD items.

### Tasks

| ID | Title | Status | Labels |
|----|-------|--------|--------|
| SL-337c9b | Add issue content structure to specledger.tasks.md | open | story:US2, fr:FR-004 |
| SL-eaef67 | Add DoD population from acceptance criteria | open | story:US2, fr:FR-005 |
| SL-76e443 | Add DoD summary section to tasks.md template | open | story:US2, fr:FR-006 |
| SL-a2574d | Add error handling section to specledger.tasks.md | open | story:US2, fr:FR-007 |
| SL-0aeb93 | Sync specledger.tasks.md to embedded templates | open | story:US2, fr:FR-011 |

**Dependencies**:
- SL-337c9b → SL-0aeb93
- SL-eaef67 → SL-0aeb93
- SL-76e443 → SL-0aeb93
- SL-a2574d → SL-0aeb93

**Parallel Execution**: SL-337c9b, SL-eaef67, SL-76e443, SL-a2574d can run in parallel

---

## Phase 3: US3 - Utilize Definition of Done During Implementation

**Feature**: SL-d9921b - US3: Utilize Definition of Done During Implementation

**Goal**: Update specledger.implement.md to verify DoD items before closing issues.

**Independent Test**: Implement a task with DoD items and verify the system checks each item.

### Tasks

| ID | Title | Status | Labels |
|----|-------|--------|--------|
| SL-b7a777 | Add DoD verification section to specledger.implement.md | open | story:US3, fr:FR-008, fr:FR-009, fr:FR-010 |
| SL-7edfe9 | Sync specledger.implement.md to embedded templates | open | story:US3, fr:FR-011 |

**Dependencies**: SL-b7a777 → SL-7edfe9

---

## Phase 4: Polish - Verification and Consistency

**Feature**: SL-7b8790 - Polish: Verification and Consistency

**Goal**: Verify prompt consistency and run manual tests.

### Tasks

| ID | Title | Status | Labels |
|----|-------|--------|--------|
| SL-1571dc | Verify prompt file consistency | open | phase:polish, fr:FR-011 |
| SL-610a8c | Run manual tests from quickstart.md | open | phase:polish |

**Dependencies**: SL-1571dc → SL-610a8c

---

## Definition of Done Summary

| Issue ID | DoD Items |
|----------|-----------|
| SL-1f8167 | Pattern matching works, sl deps list called, content loads from cache, error message displayed |
| SL-f39984 | Files identical, no content left behind |
| SL-337c9b | Structure section defines all fields, examples show format, purposes explained |
| SL-eaef67 | DoD derived from acceptance criteria, Then clauses become items, DoD included in creation |
| SL-76e443 | DoD Summary section added, table format, items referenceable by ID |
| SL-a2574d | Error handling section added, sanitization documented, retry logic explained |
| SL-0aeb93 | Files identical, all US2 changes included |
| SL-b7a777 | DoD verification section added, automated patterns documented, fallback described, failed handling defined |
| SL-7edfe9 | Files identical, all US3 changes included |
| SL-1571dc | All diffs show no differences, report generated |
| SL-610a8c | All 6 scenarios pass, checklist completed, issues documented |

## Query Commands

```bash
# View all issues for this feature
sl issue list --label "spec:592-prompt-updates"

# View open issues
sl issue list --status open

# View by story
sl issue list --label "story:US1"
sl issue list --label "story:US2"
sl issue list --label "story:US3"

# View by phase
sl issue list --label "phase:polish"
```

## MVP Scope

**Recommended MVP**: US1 only (SL-1f8167, SL-f39984)

This delivers immediate value by enabling dependency references in specifications.

## Links

- [spec.md](./spec.md) - Feature specification
- [plan.md](./plan.md) - Implementation plan
- [quickstart.md](./quickstart.md) - Test scenarios
