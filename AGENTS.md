# Pebbleton — Agent Guidelines

Pebbleton is an **OpenAI DevDay hackathon project**. Optimize for a fun, working demo—not production-grade engineering.

## Source of truth
- Follow the product spec. Use conversation context only to clarify intent.
- Keep V1 tightly scoped: 3 protagonists, 4 citizen factions, 1 economic proposal and vote per day, persistent consequences, and a daily newspaper.

## Stack
- Backend: Go, preferably `net/http` and the standard library.
- Frontend: React, TypeScript, Vite, TanStack Query.
- State: in memory with simple JSON snapshot persistence. No database.
- AI: Responses for proposals/newspapers, Decisions for choices/votes, GPT Image for illustrations.
- Mock mode must work without API keys.

## Working style
- Inspect the repo before editing; preserve unrelated work.
- Build vertically: simulation → API/UI → live AI → visual polish.
- **Commit each major working chunk** with a descriptive Git message.
- Run the app and verify the core flow after each chunk.
- Write a few focused simulation tests; avoid extensive test suites and premature abstractions.
- Make routine decisions independently. Consult current API docs instead of guessing.

## Invariants
- Only the deterministic Go simulation mutates world state; models propose and evaluate.
- Day advancement is atomic: failures must not leave partial state.
- Keep an append-only day history; newspaper content must reflect committed facts.
- One active income-transfer policy at a time; new policies replace old ones.
- Do not add multiplayer, individual NPC agents, complex economics, or infrastructure not in V1.

## Definition of done
A user can open the app, advance a day, see votes and changed state, read the day's newspaper, browse past days, and reset. Mock mode works end-to-end. Report run commands, commits, tested flows, and any gaps.
