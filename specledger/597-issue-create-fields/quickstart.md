# Quickstart: Issue Create Fields Enhancement

**Feature**: 597-issue-create-fields
**Date**: 2026-02-22

## Prerequisites

- Built `sl` binary: `make build`
- In a SpecLedger project directory

## Test Scenarios

### Scenario 1: Create Issue with All New Fields

```bash
# Create an issue with all 4 new fields
./bin/sl issue create \
  --title "Implement user authentication" \
  --type task \
  --priority 1 \
  --acceptance-criteria "User can log in with email/password. Invalid credentials show error." \
  --dod "Write unit tests" \
  --dod "Add integration test" \
  --dod "Update documentation" \
  --design "Use JWT tokens with 24h expiry. Store hashed passwords using bcrypt." \
  --notes "Consider adding OAuth support in future"

# Verify the issue was created
./bin/sl issue show <issue-id>
```

**Expected**: Issue displays with dedicated sections for Acceptance Criteria, Design, Definition of Done, and Notes.

### Scenario 2: Create Issue with Repeated --dod Flags

```bash
# Create issue with multiple DoD items
./bin/sl issue create \
  --title "Add API endpoint" \
  --type task \
  --dod "Create route handler" \
  --dod "Add request validation" \
  --dod "Write tests" \
  --dod "Update API docs"

# Show the issue to verify DoD items
./bin/sl issue show <issue-id>
```

**Expected**: Definition of Done shows 4 unchecked items.

### Scenario 3: Update DoD on Existing Issue

```bash
# Replace entire DoD
./bin/sl issue update <issue-id> \
  --dod "New requirement 1" \
  --dod "New requirement 2"

# Verify replacement
./bin/sl issue show <issue-id>
```

**Expected**: Previous DoD items replaced with 2 new unchecked items.

### Scenario 4: Check DoD Item

```bash
# Check a specific DoD item
./bin/sl issue update <issue-id> --check-dod "New requirement 1"

# Verify item is checked
./bin/sl issue show <issue-id>
```

**Expected**: "New requirement 1" shows `[x]` with a verified_at timestamp.

### Scenario 5: Check Non-existent DoD Item (Error Case)

```bash
# Try to check an item that doesn't exist
./bin/sl issue update <issue-id> --check-dod "Nonexistent item"
```

**Expected**: Error message: `DoD item not found: 'Nonexistent item'`

### Scenario 6: Uncheck DoD Item

```bash
# Uncheck a previously checked item
./bin/sl issue update <issue-id> --uncheck-dod "New requirement 1"

# Verify item is unchecked
./bin/sl issue show <issue-id>
```

**Expected**: "New requirement 1" shows `[ ]` and verified_at is cleared.

### Scenario 7: Exact Text Matching

```bash
# Create issue with DoD item
./bin/sl issue create --title "Test" --type task --dod "Write Tests"

# These should all FAIL (exact match required):
./bin/sl issue update <issue-id> --check-dod "write tests"    # Wrong case
./bin/sl issue update <issue-id> --check-dod "Write Tests "   # Trailing space
./bin/sl issue update <issue-id> --check-dod "Write  Tests"   # Double space

# This should SUCCEED:
./bin/sl issue update <issue-id> --check-dod "Write Tests"    # Exact match
```

**Expected**: Only exact match succeeds.

### Scenario 8: Issue Show Display

```bash
# Create a fully-specified issue
./bin/sl issue create \
  --title "Complex task" \
  --description "This is the description" \
  --type task \
  --acceptance-criteria "AC1: Success case. AC2: Error case." \
  --design "Use factory pattern with dependency injection" \
  --dod "Implement core logic" \
  --dod "Add error handling" \
  --notes "Reference: https://example.com/docs"

# Display the issue
./bin/sl issue show <issue-id>
```

**Expected Output Format**:
```
Issue: SL-xxxxxx
  Title: Complex task
  Type: task
  Status: open
  Priority: 2
  Spec: <spec-context>

Description:
  This is the description

Acceptance Criteria:
  AC1: Success case. AC2: Error case.

Design:
  Use factory pattern with dependency injection

Definition of Done:
  [ ] Implement core logic
  [ ] Add error handling

Notes:
  Reference: https://example.com/docs

Created: 2026-02-22 10:30:00
Updated: 2026-02-22 10:30:00
```

## Verification Checklist

- [ ] `--acceptance-criteria` flag works on create
- [ ] `--dod` flag works with repeated values on create
- [ ] `--design` flag works on create
- [ ] `--notes` flag works on create
- [ ] `--dod` flag replaces entire DoD on update
- [ ] `--check-dod` marks item as checked with timestamp
- [ ] `--uncheck-dod` marks item as unchecked
- [ ] Error returned when checking non-existent DoD item
- [ ] Exact text matching (case-sensitive, no normalization)
- [ ] `sl issue show` displays all 4 fields in dedicated sections
