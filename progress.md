The {{aka}} command will be responsible to send the jira card to "In Progress" column, change the assignee to user and create a new branch worktree using a branch pattern `feature/TEC-1600`.

# Flags

- `s jira progress 1600 nojira` :: On that situation, `nojira` is the flag and ignore the jira assignee and move to progress column part, just create the git branch.

- `s jira progress 1600 nogit` :: On that situation, `nogit` is the flag and ignore the git branch creation part, just assign and move jira card to progress column.
