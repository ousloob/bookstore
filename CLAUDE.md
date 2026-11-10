# Working on this repo

- Never run `kind` or the cluster-mutating make targets (`dev-up`, `dev-down`,
  `dev-load`, `dev-secrets`, `dev-apply`, `dev-update`). Two worktrees share
  this machine and the kind cluster; the user runs those commands manually.
  Read-only checks (`kubectl kustomize`, `make -n`) are fine.
- Dev/prod parity is a goal of this project: the kind dev environment must
  behave exactly like prod. Don't accept loose shortcuts "because it's only
  dev" — prefer the production-grade solution (e.g. an init container over
  letting the API crash-loop until the database is up).
