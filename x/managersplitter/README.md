# Manager Splitter Module

A manager splitter is a module-derived address that acts as the manager of one
or more tokenization collections. Its admin delegates individual manager
actions (metadata, token ids, approvals, standards, archive, ...) to approved
addresses through `MsgExecuteUniversalUpdateCollection`; every other field of
the forwarded message (collection creation, mint-escrow coin transfers,
invariants, default balances, collection permissions) stays admin-only.

## Notes for admins

- `canUpdateManager` is full-control delegation: an approved address can set
  itself as the collection manager and then holds every manager power directly,
  so grant it only to addresses you would trust with the whole collection.
- `MsgDeleteManagerSplitter` does not reassign the collections the splitter
  manages; the derived address is never recreated, so move each collection to a
  new manager with `UpdateManager` before deleting the splitter.
