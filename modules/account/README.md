# FTN Account Engine

Production account, ledger, mail-event and reconciliation module for FTN-AI.

## Current boundaries

- Account HTTP API
- Mail event processing API
- Deterministic reconciliation matcher
- Financial actions remain behind the account service boundary
- AI is not permitted to write directly to the ledger

## Safety invariant

Financial posting must be performed by the authoritative ledger service and must reject unbalanced journals. Mail-derived events are evidence and require deterministic matching or the configured review policy before posting.
