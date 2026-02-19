# Quickstart: Improve SpecLedger Command Prompts

**Feature**: 592-prompt-updates
**Date**: 2026-02-20

## Overview

This quickstart provides test scenarios for validating the updated SpecLedger command prompts. Each scenario tests a specific functional requirement.

## Prerequisites

- SpecLedger CLI installed
- Test project initialized with `sl init`
- At least one dependency added: `sl deps add --alias test-deps https://example.com/spec.md`

## Test Scenarios

### Scenario 1: Dependency Detection in /specledger.specify

**Tests**: FR-001, FR-002, FR-003

**Steps**:
1. Run `/specledger.specify Create feature integrating with deps:test-deps for user authentication`
2. **Expected**: Spec generation loads content from test-deps cache
3. Verify spec.md references the dependency context

**Negative Case**:
1. Run `/specledger.specify Create feature using deps:nonexistent-deps`
2. **Expected**: Error message "Dependency 'nonexistent-deps' not found. Use 'sl deps add --alias nonexistent-deps <source>' to add it."

### Scenario 2: Issue Structure in /specledger.tasks

**Tests**: FR-004, FR-005, FR-006

**Steps**:
1. Complete a spec with `/specledger.specify`
2. Complete a plan with `/specledger.plan`
3. Run `/specledger.tasks`
4. **Expected**: Each issue has:
   - Title under 80 characters
   - Problem statement (WHY)
   - Implementation details (HOW/WHERE)
   - Acceptance criteria (WHAT)
   - definition_of_done items
5. Verify tasks.md contains DoD Summary section

### Scenario 3: Error Handling in /specledger.tasks

**Tests**: FR-007

**Steps**:
1. Create a plan with special characters in component names (e.g., `component: "test's-feature"`)
2. Run `/specledger.tasks`
3. **Expected**: Commands succeed without errors (special characters auto-escaped)

### Scenario 4: Automated DoD Verification in /specledger.implement

**Tests**: FR-008, FR-009, FR-010

**Steps**:
1. Create a task with DoD item "file exists: /tmp/test.txt"
2. Run `/specledger.implement`
3. Implement the task without creating the file
4. **Expected**: Verification fails with clear message about missing file
5. Create the file: `touch /tmp/test.txt`
6. Re-run verification
7. **Expected**: Verification passes

### Scenario 5: Interactive DoD Confirmation

**Tests**: FR-009

**Steps**:
1. Create a task with DoD item "User interface is intuitive" (cannot be automated)
2. Run `/specledger.implement`
3. **Expected**: System prompts "Is 'User interface is intuitive' complete? (y/n)"
4. Answer 'y'
5. **Expected**: Issue closes successfully

### Scenario 6: Prompt Consistency

**Tests**: FR-011

**Steps**:
1. Diff `.claude/commands/specledger.specify.md` vs `pkg/embedded/skills/commands/specledger.specify.md`
2. Diff `.claude/commands/specledger.tasks.md` vs `pkg/embedded/skills/commands/specledger.tasks.md`
3. Diff `.claude/commands/specledger.implement.md` vs `pkg/embedded/skills/commands/specledger.implement.md`
4. **Expected**: No differences (files are identical)

## Validation Checklist

After completing all scenarios, verify:

- [ ] Dependency references with explicit syntax are detected
- [ ] Dependency content loads from cache when available
- [ ] Missing dependencies show helpful error message
- [ ] Generated issues have complete structure
- [ ] DoD items are populated from acceptance criteria
- [ ] tasks.md contains DoD Summary section
- [ ] Special characters are handled gracefully
- [ ] Automated DoD verification works for supported patterns
- [ ] Interactive confirmation works for non-automatable items
- [ ] Failed verification displays clear messages
- [ ] Dev and embedded prompts are identical
