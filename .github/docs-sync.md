# Documentation ownership

The private `BitBadges/bitbadges-monorepo` owns the documentation site. It pins
this public repository at `public/chain` and generates the chain API and proto
reference after an accepted pin update, with scheduled reconciliation as backup.
Public chain source and release workflows stay in this repository.

Activate and verify the parent replacement before merging removal of the old
notifier. Retire `DOCS_DISPATCH_PAT` after that handoff; do not grant a public
workflow credentials for writing to the private parent. Parent docs generation
opens a separate PR and does not publish this chain or change its default branch.
