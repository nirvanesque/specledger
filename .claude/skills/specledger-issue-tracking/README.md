# SpecLedger Issue Tracking Skill

A comprehensive Claude Code skill for tracking complex, multi-session work with dependency graphs using the SpecLedger issue tracker (bd).

## What is This?

This is a [Claude Code](https://claude.com/claude-code) skill that teaches Claude how to use SpecLedger's issue tracking system effectively for multi-session coding workflows.

## What Does It Provide?

**Main skill file (`SKILL.md`):**
- Core workflow patterns (discovery, execution, planning phases)
- Decision criteria for when to use issue tracking vs TodoWrite
- Session start protocols and ready work checks
- Issue lifecycle management with self-check checklists
- Integration patterns with other tools

**Reference documentation:**
- `references/BOUNDARIES.md` - Detailed decision criteria with examples
- `references/CLI_REFERENCE.md` - Complete command reference with all flags
- `references/DEPENDENCIES.md` - Deep dive into dependency types and relationships
- `references/WORKFLOWS.md` - Step-by-step workflows with checklists
- `references/ISSUE_CREATION.md` - When to ask vs create issues, quality guidelines
- `references/RESUMABILITY.md` - Making issues resumable across sessions

## Why is This Useful?

The skill helps Claude understand:

1. **When to use issue tracking** - Not every task needs it. Learn when issue tracking helps vs when TodoWrite is better.

2. **How to structure issues** - Proper use of dependency types, issue metadata, and relationship patterns.

3. **Workflow patterns** - Proactive issue creation during discovery, status maintenance during execution.

4. **Integration with other tools** - How issue tracking and TodoWrite can coexist.

## Usage

Once installed, Claude will automatically:
- Check for ready work at session start (if `.beads/` exists)
- Suggest creating issues for multi-session work
- Use appropriate dependency types when linking issues
- Maintain proper issue lifecycle (create → in_progress → close)

You can also explicitly ask Claude to use issue tracking:

```
Let's track this work in the issue tracker since it spans multiple sessions
```

```
Create an issue for this bug we discovered
```

```
Show me what's ready to work on
```

## See Also

- **[specledger-deps](../specledger-deps/README.md)** - Manage specification dependencies
