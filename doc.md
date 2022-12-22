The {{aka}} branch should be responsible to deal with jira workflows, it relies heavily on the [jira-cli](https://github.com/ankitpokhrel/jira-cli) binary to manipulate jira cards.

## Features

- `progress 1642` :: Should assign the card **1642** to me and move it to progress, also should create a new git branch using worktrees matching our pattern(feature/TEC-1642).

- `review 1642` :: Shold move the card to review and insert to system clipboard a template that can be inserted into slack, look something like this:

```
:github: PR: [TEC-$1]($prurl) :github:
CARD: [https://lamimed.atlassian.net/browse/TEC-$1]
URL: [homologation_url]
NOTE/OBS:
cc @here
[status guide](https://lamimed.atlassian.net/wiki/spaces/TEC/pages/210436097/Messages+template+on+Slack)
```

- `homol 1642` :: Should move the card to "In Homol" column.

- `done 1642` :: Should move the card to "Done" column.
