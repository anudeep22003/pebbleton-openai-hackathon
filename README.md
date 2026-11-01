# Pebbleton — A Tiny Society with Big Problems

Forty citizens. Three ambitious neighbours. One daily newspaper. A 30-day society simulator built for OpenAI DevDay, with Go, React, TypeScript and TanStack Query.

## Run locally

Requires Go 1.24+ and Node.js 22.12+ (or 20.19+).

```sh
make install
make start
```

Open **http://127.0.0.1:8080**. Mock mode is the default and needs no API key. `make start` builds the frontend and serves everything from one Go process.

For frontend hot reload, run `make backend` and `make frontend` in separate terminals, then open **http://127.0.0.1:5173**. Vite proxies API and image requests to port 8080.

```sh
make test       # focused backend tests
make simulate   # complete reproducible 30-day mock story
make format     # gofmt and Prettier
```

## The rules

Pip, Momo and Zig take turns presenting one of three proposals. Every proposal moves one coin of daily income per citizen between two factions. All forty citizens vote; ties fail. Protagonist reactions do not count as votes. A passed policy replaces the previous policy and affects income starting the following day. Citizens remember attempted policies even when they fail.

Only Go changes the world. The frontend displays server facts, including effective income and historical snapshots. Random sampling uses the configured seed and day number, so retries and restarts preserve mock outcomes. Reset restores the configured seed and gives the new world a distinct identity.

## Configuration

Copy `.env.example` to `.env` to change settings. Make loads this file; direct `go run` reads exported environment variables. `DATA_DIR` defaults to `./data`, `PORT` to `8080`, and `SEED` to `42`. The backend binds to localhost.

Snapshots are atomically written to `data/world.json`. Images belong in `data/images/`. This is a single-process local app: ephemeral disks and multiple backend instances are not supported. Invalid snapshots are reported in the server log and start a fresh Day 0.

## Live OpenAI mode

Set `MOCK_AI=false` and `OPENAI_API_KEY` in your local `.env`, then restart `make backend` (or `make start`). Keep the key out of frontend files and Git. Each day makes two Responses requests and two Decisions requests, with one extra proposal request only if semantic validation fails. A live failure never silently changes to mock decisions; newspaper failure alone uses the factual template.

| Variable | Default | Purpose |
| --- | --- | --- |
| `MOCK_AI` | `true` | Use repeatable local proposals and decisions |
| `OPENAI_API_KEY` | unset | Required only in live mode |
| `OPENAI_TEXT_MODEL` | `gpt-5.4-mini` | Structured proposals and newspaper prose |
| `OPENAI_DECISIONS_MODEL` | `gpt-6-luna` | Choice distribution and six batched predicates |
| `OPENAI_IMAGE_MODEL` | `gpt-image-2.5-flare` | Optional editorial illustrations |
| `GENERATE_IMAGES` | `false` | Optional asynchronous illustration generation |
| `DATA_DIR` | `./data` | Snapshot and image storage |
| `PORT` | `8080` | Backend port; update the Vite proxy if changing it |
| `SEED` | `42` | Sampling seed restored on reset |

Live prompts contain structured current state and at most seven recent factual summaries. The engine samples the `choice` answer's `probabilities` array; it does not use `confidence` or blindly accept the API's preferred `choice`. Probability completeness and bounds are checked before any state is saved.

Request contracts were checked against the official [Decisions guide](https://developers.openai.com/api/docs/guides/decisions), [Structured Outputs guide](https://developers.openai.com/api/docs/guides/structured-outputs), and [image generation guide](https://developers.openai.com/api/docs/guides/image-generation). Tests use local HTTP fixtures and never require API credits.

Set `GENERATE_IMAGES=true` in live mode for one GPT Image illustration per new edition. The day commits before image generation starts; the UI polls while artwork is pending. Images are stored locally and survive restart. Image failures keep the newspaper playable and show the local storybook illustration. Mock mode always uses local artwork, even if the image flag is enabled. Restart marks interrupted image jobs failed rather than charging for automatic retries. Reset clears history but leaves old image files on disk.

## Verification and demo notes

Verified locally: production build; full 30-day mock simulation and HTTP run; browser Next Day → changed state → newspaper → historical snapshot → reset; desktop and 390px layouts; two real live days with Responses/Decisions; and a real GPT Image illustration saved and displayed in the second live edition. The full 30-day live run was not exercised.

- `go test -race ./...` covers a repeatable 30-day run, policy replacement/timing, debt, validation, ties, probability extremes, relationships, immutable history, concurrency, persistence and image/reset races.
- Provider fixtures cover the exact Responses/Decisions request flow, refusals, missing probabilities, HTTP errors and newspaper fallback. Live model output is probabilistic even with a fixed sampling seed.
- The initial world, Next Day, voting results, historical snapshots and reset are available in mock mode. Inspect the candidate probabilities under each newspaper and expand locals/factions for relationships and loyalties.
- Newspaper prose is grounded through prompts and supplied resolved facts; there is no separate semantic fact-checking model. The recorded policy, vote tally and dashboard remain authoritative.
- Fonts use Google Fonts with local serif/sans-serif fallbacks. The mock simulation and local vector artwork work without network access after installing dependencies.
- No multiplayer, database, fast-forward, account system or additional economics are included.

## API

- `GET /api/world` returns `{world, running, max_days, mode}`.
- `POST /api/advance` completes one day and returns the same response shape. Concurrent advancement/reset returns 409; Day 31 returns 422. Failed essential stages or persistence leave the world unchanged.
- `POST /api/reset` restores Day 0 and returns the same shape.

The newspaper, votes, candidates and probabilities remain in every day record. Historical navigation displays that day's economic and relationship snapshot.
