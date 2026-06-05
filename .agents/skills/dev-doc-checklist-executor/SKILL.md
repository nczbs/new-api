---
name: dev-doc-checklist-executor
description: Manual-only coordinated development execution with design docs, checklists, commits, cleanup, and reports.
---

# Dev Doc Checklist Executor

Guide development by using the relevant design document and checklist together.

## Flow

1. Locate the design document and checklist for the requested work.
2. Read the design document first to understand goals, constraints, data structures, API contracts, and implementation
   boundaries.
3. Read the checklist to identify independent feature slices, phases, task items, and verification items.
4. Execute tasks in checklist order, one independent feature slice at a time.
5. For each feature slice:

    * implement the required tasks
    * run required tests, validation, and regression checks
    * review the changes
    * update or confirm matching checklist items
    * commit the completed slice
    * push the commit
    * clean temporary or intermediate artifacts
    * generate or update the completion report
6. Do not move to the next feature slice until the current slice is implemented, verified, reviewed, committed, pushed,
   cleaned up, and reported.

## Rules

* Follow both the design document and checklist.
* Do not execute mechanically from the checklist alone.
* If the checklist conflicts with the design document, follow the design document and point out the conflict.
* If the checklist lacks detail, use the design document to complete the implementation basis.
* Keep each commit scoped to one completed independent feature slice.
* Do not commit unrelated changes.
* Do not push if tests or required verification fail, unless the user explicitly approves.
* Clean only safe intermediate artifacts such as caches, generated scratch files, temporary logs, and unused drafts.
* Do not delete source files, docs, fixtures, migrations, or generated assets required by the project unless explicitly
  intended.
* Do not output unrelated background, long explanations, architecture essays, or filler.
* If required documents are missing, paths do not exist, git state is unsafe, push target is unclear, or requirements
  have blocking ambiguity, state the missing information or ask only the necessary questions.

## Review Checklist

Before each commit:

* Confirm the implementation matches the design document.
* Confirm the completed checklist items are actually done.
* Inspect changed files for unrelated edits.
* Run formatting, linting, tests, type checks, build checks, or regression checks as appropriate.
* Summarize any remaining risks or follow-ups.

## Report Format

Location: docs/dev/report
After each completed feature slice, create or update a report document using this format:

```markdown
# Development Completion Report

Date: <YYYY-MM-DD>
Feature Slice: <name>
Status: <completed | blocked | partially completed>
Branch: <branch>
Commit: <commit hash or pending>
Pushed: <yes | no>
Design Document: <path or title>
Checklist Document: <path or title>

## Task Summary

- <summary of completed task or checklist item>
- <summary of completed task or checklist item>

## Changed Files

- `<path>`: <summary of change>
- `<path>`: <summary of change>

## Verification

- Tests: <commands and results>
- Validation: <checks performed>
- Regression: <regression scope and result>

## Review Notes

- <design/checklist alignment notes>
- <conflicts found, if any>
- <risks or edge cases>

## Cleanup

- Removed: <temporary artifacts removed>
- Kept: <artifacts intentionally preserved>

## Follow-ups

- <remaining work, known limitations, or next feature slice>
```
