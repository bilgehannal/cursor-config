# Get Conflict Responsible

Find who is responsible for git merge conflicts between two remote branches.

The user will provide two branch names (bare names, e.g. `main`, `feature/foo`). Ask for them if not provided:
- **branch**: The target branch (the command will use `origin/<branch>`)
- **conflicted_branch**: The branch to merge into the target (the command will use `origin/<conflicted_branch>`)

> **Note:** This command automatically runs `git fetch origin` and operates entirely on remote refs (`origin/...`) via detached HEAD. Local branches are never checked out or modified. This ensures the analysis always reflects the latest remote state.

Execute the following steps exactly in order using the Shell tool for all git commands. If any step fails unexpectedly, skip immediately to Step 6 (Cleanup) to restore the repo.

## Step 1 — Save Current State & Fetch

1. Record the current branch:
   ```
   git rev-parse --abbrev-ref HEAD
   ```
   Store this as `ORIGINAL_BRANCH`.

2. Check for uncommitted changes:
   ```
   git status --porcelain
   ```
   If there is any output, stash the changes:
   ```
   git stash push -m "get-conflict-responsible auto-stash"
   ```
   Remember whether a stash was created (STASH_CREATED = true or false).

3. Fetch the latest remote refs:
   ```
   git fetch origin
   ```

## Step 2 — Checkout & Merge (Remote Refs)

1. Checkout the target branch as a detached HEAD using the remote ref:
   ```
   git checkout origin/<branch>
   ```

2. Attempt the merge without committing, using the remote ref:
   ```
   git merge --no-commit --no-ff origin/<conflicted_branch>
   ```

3. If the merge succeeds with no conflicts, report "No conflicts found between origin/<branch> and origin/<conflicted_branch>." and skip to Step 6.

## Step 3 — Identify Conflicted Files

1. List all conflicted files:
   ```
   git diff --name-only --diff-filter=U
   ```

2. For each conflicted file, read its contents and locate all conflict markers:
   - `<<<<<<<` marks the start of the current change (from the target branch)
   - `=======` separates current from incoming
   - `>>>>>>>` marks the end of the incoming change (from the conflicted branch)

3. For each conflict hunk, note the content of both the current and incoming sides.

## Step 4 — Blame Analysis

For each conflict hunk found in Step 3:

1. **Current changes author** (this is the PRIMARY output — most important):
   Run git blame on the remote target branch version of the file to find who last modified the lines in the current-change block:
   ```
   git blame origin/<branch> -- <file>
   ```
   Look at the lines that correspond to the current-change content. Extract the author name and email.

2. **Incoming changes author** (secondary info):
   Run git blame on the remote conflicted branch version of the file to find who last modified the lines in the incoming-change block:
   ```
   git blame origin/<conflicted_branch> -- <file>
   ```
   Look at the lines that correspond to the incoming-change content. Extract the author name and email.

Tip: Use `git blame -L <start>,<end> origin/<branch> -- <file>` with line ranges when you can determine the exact lines from the pre-merge version. If exact line mapping is difficult, blame the full file and match by content.

## Step 5 — Report

Present a clear, structured report. The current changes user is the most important.

Format the output like this:

```
## Conflict Analysis: origin/<branch> ← origin/<conflicted_branch>

### File: <filepath>

#### Conflict #N (lines X-Y)

Current changes (from origin/<branch>):
  Author: <name> (<email>)
  Lines:
    <show the conflicting lines briefly>

Incoming changes (from origin/<conflicted_branch>):
  Author: <name> (<email>)
  Lines:
    <show the conflicting lines briefly>

---

### Summary

| File | Current Changes By | Incoming Changes By |
|------|-------------------|---------------------|
| ...  | ...               | ...                 |

Primary responsible person for current conflicts: <most frequent current-changes author>
```

## Step 6 — Cleanup / Restore (ALWAYS RUN THIS)

This step MUST always run, even if previous steps failed.

1. Abort the merge (ignore errors if no merge is in progress):
   ```
   git merge --abort
   ```

2. Return to the original branch:
   ```
   git checkout <ORIGINAL_BRANCH>
   ```

3. If STASH_CREATED was true, restore the stash:
   ```
   git stash pop
   ```

4. Confirm the repo is restored:
   ```
   git status
   ```
   Report: "Repository restored to original state on branch <ORIGINAL_BRANCH>."
