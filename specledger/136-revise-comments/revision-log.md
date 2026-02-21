# Revision Log: 136-revise-comments

Tracks all revision options proposed and user choices during the `/specledger.revise` workflow.

## Comment 1

- **File**: `specledger/136-revise-comments/plan.md`
- **Target**: "New sl resolve command:" (Future Enhancements section, lines 279-344)
- **Feedback**: "I think a command group to work with comments `sl comments pull`, `sl comment resolve` will allow a better agentic workflow than the proposed `sl revise` workflow"
- **Author Guidance**: Clarify that future work will focus on dedicated command group to work with comments left on specs (summarise them by artifact, get all unresolved comments on a specific artifact, resolve a comment with a --reason, ...)

### Options Presented

| Option | Description |
|--------|-------------|
| **A: Full rewrite (Recommended)** | Replace the entire 'Agent-driven resolve — sl resolve command' block with a new 'sl comments command group' section listing subcommands (pull, list, resolve, summarize, reply), describing the agentic workflow, and removing old sl resolve details. |
| **B: Reframe under umbrella** | Keep most existing sl resolve details (bash examples, UX flow) but reframe under a broader sl comments group. Adds other subcommands as bullet points alongside existing resolve description. |
| **C: Minimal annotation** | Keep existing section mostly intact, add a prominent direction-update note at the top, rename sl resolve to sl comments resolve throughout. |

### User Choice: **Option A — Full rewrite**

### Changes Applied
- Replaced lines 280-344 in `plan.md` (the entire "Agent-driven resolve — `sl resolve` command" subsection) with a concise "Agent-driven comment management — `sl comments` command group" section
- New section lists 5 planned subcommands: `pull`, `list`, `summarize`, `resolve`, `reply`
- Removed detailed bash examples, UX flow walkthrough, and prompt template update subsections (these belong in the future spec for the `sl comments` feature itself)
- Updated estimated effort from 3-5 days to 5-8 days (broader scope)
