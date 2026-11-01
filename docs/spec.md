# Pebbleton — Product & Implementation Spec (V1)

**Project:** OpenAI DevDay Hackathon  
**Product:** Pebbleton — A Tiny Society with Big Problems  
**Objective:** Build a playable, AI-driven miniature society that evolves over 30 days and tells its own history through a daily illustrated newspaper.

## 0. Instructions to the coding agent

Implement this application end-to-end.

This is a hackathon project. Prioritize a working, polished, demonstrable product over architectural sophistication.

- Inspect the existing repository before making changes. Follow any existing project conventions where reasonable.
- Use the specified tech stack unless an existing repository structure strongly justifies a small adjustment.
- Implement incrementally, keeping the application runnable at every checkpoint.
- Start with deterministic mock AI behavior before integrating live model calls.
- Do not introduce functionality outside the explicit V1 scope.
- If implementation details are unspecified, choose the simplest reasonable solution and document your assumption.
- Do not stop at scaffolding or a design document. Implement, test and run the application.
- Update the README with setup instructions, configuration and what is actually implemented.
- Keep the frontend visually appealing. The newspaper is the primary product experience, not an afterthought.

**Most important success criterion:** A user can open the app, click "Next Day," watch the simulation update, and read an entertaining newspaper explaining what happened and why.

---

# 1. Product overview

Pebbleton is an AI-driven miniature society consisting of:

- **Three protagonists:** Named characters with personal ambitions, wealth, occupations and feelings toward one another.
- **Four citizen factions:** Farmers, Merchants, Workers and Clergy. Each faction has a population, income, wealth and loyalty toward the protagonists.
- **One world:** A shared economic and political state that changes over time.
- **One decision per day:** A protagonist proposes an economic policy, and the society votes.
- **One newspaper per day:** The newspaper reports the proposal, reactions, results and consequences with a humorous headline and an AI-generated illustration.

The simulation is neither a scripted story nor an unconstrained AI chatbot.

It is a small state machine in which AI proposes actions and evaluates choices, while deterministic code updates the world.

### Product personality

The tone should be humorous, whimsical and slightly satirical.

Characters should feel like ambitious, flawed little citizens taking themselves far too seriously.

Political disputes over one coin should be reported like major historical crises.

Think of a tiny fantasy civilization with a newspaper that behaves like a sensationalist national publication.

### Target demonstration

During a live demonstration:

1. The user opens Pebbleton and sees Day 0 with three protagonists and four citizen factions.
2. They click **Next Day**.
3. A protagonist proposes one of three AI-generated economic policies.
4. The Decisions API selects which proposal that character would pursue.
5. The remaining protagonists react, and the citizens vote.
6. The policy passes or fails.
7. Money, loyalties and relationships change.
8. A newspaper appears describing the day's events.
9. The user advances through subsequent days and sees how previous events influence new decisions.

The experience should support a coherent 30-day history.

---

# 2. Scope

## In scope for V1

- One world, named Pebbleton.
- Three predefined protagonists.
- Four predefined citizen factions.
- In-memory simulation state.
- JSON snapshot persistence to local disk.
- One economic policy proposal per day.
- Three AI-generated candidate proposals.
- Decisions API evaluation of the protagonist's preferred proposal.
- Decisions API evaluation of protagonist reactions and citizen support.
- Probabilistic sampling of choices and votes.
- Deterministic economic and relationship updates.
- Append-only day history.
- AI-generated newspaper headlines and articles.
- One optional asynchronous AI-generated illustration per newspaper edition.
- React dashboard and newspaper reader.
- Next Day and Reset controls.
- Historical day navigation.
- Mock AI mode, requiring no OpenAI API key.
- Unit tests for core simulation rules.

## Explicitly out of scope

Do not implement:

- User accounts, authentication or multiplayer.
- Crowd participation or real human voting.
- Multiple worlds or branching timelines.
- Free-form protagonist conversations.
- Individual NPC agents or NPC memories.
- Family trees, reproduction, births or deaths.
- Resource types beyond money.
- Complex taxation, markets, banking or trading.
- AI agents with autonomous tool access.
- Arbitrary generated actions or general-purpose law engines.
- WebSockets, message queues or microservices.
- 3D rendering or Three.js.
- Automatic multi-day fast-forward.
- A database or ORM.

The architecture should allow these to be added later, but V1 must not implement them.

---

# 3. Technology stack

## Backend

- Go, using `net/http`.
- Go standard library wherever practical.
- Official OpenAI Go SDK or direct REST requests for AI integration.
- One backend process containing all simulation state.
- JSON file persistence.
- No database.

## Frontend

- React with TypeScript.
- Vite.
- React Router only if needed; a single-page app is sufficient.
- TanStack React Query for API state and mutations.
- CSS or lightweight styling utilities.
- Responsive layout, primarily optimized for desktop presentation.

## AI

- OpenAI Responses API for proposal and newspaper generation.
- OpenAI Decisions API for proposal selection and support probabilities.
- GPT Image for newspaper illustrations.

Configuration should be environment-driven.

Use the current documented API shapes rather than assuming old interfaces.

Relevant documentation:

- https://developers.openai.com/api/docs/guides/decisions
- https://developers.openai.com/api/docs/guides/structured-outputs
- https://developers.openai.com/api/docs/guides/image-generation

---

# 4. Core simulation concepts

The simulation has three temporal concepts.

### Past

The append-only history of committed daily events.

Every protagonist and faction has access to the same public history in V1.

There are no private memories yet.

### Present

One proposed economic policy, its evaluation, reactions, voting and consequences.

This is transient during simulation and becomes part of history once committed.

### Future

Only protagonists have persistent goals.

Their goals influence which proposal they choose today.

Citizen factions do not have generated goals or autonomous plans. They simply respond to proposals based on their current circumstances.

### Primary causal loop

A protagonist proposes a policy.

The citizens react based on their wealth, occupation and loyalty.

The result changes income, wealth, relationships and loyalties.

Those changed circumstances influence tomorrow's proposal and voting decisions.

The application must maintain this feedback loop.

---

# 5. Initial world configuration

Start with three predefined protagonists.

| Name | Occupation | Goal | Initial money |
|---|---|---|---:|
| Pip | Merchant | Become the wealthiest and most influential person in Pebbleton | 15 |
| Momo | Farmer | Build a fair society where ordinary citizens prosper | 15 |
| Zig | Worker | Maintain stability, security and order | 15 |

Initial feelings toward other protagonists can be set to zero (neutral).

The three protagonists should have consistent visual identities:

- Pip: Ambitious orange fox, merchant clothing.
- Momo: Friendly blue penguin, practical farmer clothing.
- Zig: Clever green frog, engineer/inventor aesthetic.

These identities should be reused when generating newspaper illustrations.

## Citizen factions

Each faction contains 10 NPC citizens.

Faction financial values represent **money and income per NPC citizen**, not aggregate faction totals.

| Faction | Population | Initial money/person | Base income/day | Upkeep/day |
|---|---:|---:|---:|---:|
| Farmers | 10 | 10 | 3 | 2 |
| Merchants | 10 | 10 | 5 | 2 |
| Workers | 10 | 10 | 2 | 2 |
| Clergy | 10 | 10 | 2 | 2 |

Protagonists are separate individuals and are not included in these NPC counts.

A protagonist earns the effective income and pays the upkeep associated with their occupation.

Each faction also has a loyalty score toward each protagonist, ranging from 0 to 100.

Suggested initial loyalties:

| Faction | Pip | Momo | Zig |
|---|---:|---:|---:|
| Farmers | 40 | 70 | 50 |
| Merchants | 70 | 40 | 50 |
| Workers | 40 | 60 | 60 |
| Clergy | 50 | 50 | 60 |

These values create initial differentiation.

---

# 6. Core Go data structures

These are conceptual structures. Refine naming, JSON tags and helper types as needed, but preserve the model.

```go
type Protagonist struct {
    ID           string         `json:"id"`
    Name         string         `json:"name"`
    OccupationID string         `json:"occupation_id"`
    Goal         string         `json:"goal"`
    Money        int            `json:"money"`

    // Directed feelings about other protagonists.
    // Range: -100 to +100.
    Feelings map[string]int `json:"feelings"`
}

type Faction struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    Population int    `json:"population"`

    // Per NPC citizen.
    Money      int `json:"money"`
    BaseIncome int `json:"base_income"`
    Upkeep     int `json:"upkeep"`

    // Protagonist ID -> loyalty (0..100)
    Loyalty map[string]int `json:"loyalty"`
}

type Proposal struct {
    ID          string `json:"id"`
    AuthorID    string `json:"author_id"`
    Title       string `json:"title"`
    Reason      string `json:"reason"`
    FromFaction string `json:"from_faction"`
    ToFaction   string `json:"to_faction"`
}

type FactionVote struct {
    FactionID          string  `json:"faction_id"`
    SupportProbability float64 `json:"support_probability"`
    YesVotes           int     `json:"yes_votes"`
    NoVotes            int     `json:"no_votes"`
}

type ProtagonistReaction struct {
    ProtagonistID      string  `json:"protagonist_id"`
    SupportProbability float64 `json:"support_probability"`
    Supports           bool    `json:"supports"`
}

type DayRecord struct {
    Day       int    `json:"day"`
    AuthorID  string `json:"author_id"`

    Candidates             []Proposal          `json:"candidates"`
    SelectionProbabilities map[string]float64  `json:"selection_probabilities"`
    SelectedProposal       Proposal            `json:"selected_proposal"`

    FactionVotes  []FactionVote         `json:"faction_votes"`
    Reactions     []ProtagonistReaction `json:"reactions"`

    TotalYes int  `json:"total_yes"`
    TotalNo  int  `json:"total_no"`
    Passed   bool `json:"passed"`

    Headline string `json:"headline"`
    Article  string `json:"article"`

    IllustrationStatus string `json:"illustration_status"`
    IllustrationURL    string `json:"illustration_url"`

    Snapshot WorldSnapshot `json:"snapshot"`
}

type WorldSnapshot struct {
    Protagonists []Protagonist `json:"protagonists"`
    Factions     []Faction     `json:"factions"`
    ActivePolicy *Proposal     `json:"active_policy"`
}

type World struct {
    ID           string        `json:"id"`
    Day          int           `json:"day"`
    Protagonists []Protagonist `json:"protagonists"`
    Factions     []Faction     `json:"factions"`

    ActivePolicy *Proposal    `json:"active_policy"`
    History      []DayRecord  `json:"history"`

    Seed int64 `json:"seed"`
}
```

`World.Day` means the number of **completed days**.

An initialized world starts at Day 0. Its first advance generates the Day 1 record.

`WorldSnapshot` holds the end-of-day state without recursively including history. This lets the frontend inspect historical economics and loyalties alongside a historical newspaper.

The history array is append-only during normal gameplay.

The selected policy must refer to valid factions by stable IDs.

---

# 7. Proposal schema and economic rules

V1 supports exactly one kind of economic proposal:

**Transfer one coin of daily income per citizen from one faction to another.**

Example:

```json
{
  "id": "proposal-7-a",
  "author_id": "pip",
  "title": "The Merchant Prosperity Act",
  "reason": "Commerce must thrive if Pebbleton is to prosper.",
  "from_faction": "farmers",
  "to_faction": "merchants"
}
```

If enacted:

- Farmers' effective daily income decreases by 1.
- Merchants' effective daily income increases by 1.
- Other factions remain unchanged.

Because all four factions contain 10 NPC citizens, this income transfer is balanced at the NPC faction level.

### One active policy only

The world can have at most one active economic policy.

When a new policy passes, it **replaces** the previous policy.

Policies do not stack.

Example:

- Base farmer income: 3.
- Policy A takes from farmers: effective income becomes 2.
- Policy B later benefits farmers: effective income becomes 4, because Policy B replaces Policy A.

Always calculate income from the baseline and the active policy.

Do not permanently mutate `BaseIncome` when a policy passes.

### Daily economics

For each faction:

`money_per_person += effective_income - upkeep`

For each protagonist:

`money += effective_income_for_occupation - upkeep_for_occupation`

Negative money is allowed and represents debt or economic distress.

Do not implement starvation, bankruptcy or death.

### Proposal validation

The backend must validate all generated proposals.

Rules:

1. Author must be a valid protagonist.
2. Source and destination must be valid factions.
3. Source and destination must differ.
4. The proposed policy must not produce negative effective daily income for any faction.
5. Three generated candidates must have distinct `(from_faction, to_faction)` pairs.
6. The model must not specify arbitrary money amounts or state mutations.
7. IDs are assigned or verified by the backend, not trusted blindly from the model.

If a candidate is invalid, request a corrected set once. If still invalid, return an error and do not advance the day.

---

# 8. One-day simulation algorithm

Implement a single orchestrator such as:

`Simulator.AdvanceDay(ctx)`

The order of operations is important.

## Step 1 — Prepare the next day

Clone the current world state.

Calculate the next day number:

`nextDay = world.Day + 1`

Apply daily economics using the **currently active policy**.

An approved new policy will begin affecting income on the following day, not retroactively.

## Step 2 — Select today's author

Use round-robin selection.

- Day 1: Pip
- Day 2: Momo
- Day 3: Zig
- Day 4: Pip
- And so on.

No AI call is required for scheduling.

## Step 3 — Generate three possible proposals

Call the Responses API with Structured Outputs.

Supply:

- Current day.
- Current protagonist and goal.
- Current money and economic circumstances.
- Current feelings about other protagonists.
- All faction economic states.
- Current active policy.
- A compact recent history, such as the latest 5–7 day summaries.

Request exactly three distinct valid economic proposals.

The generator should create proposals that are plausible for the author but varied in strategy.

For example, Pip should not exclusively benefit merchants. Occasionally he may support workers or farmers to gain popularity, form alliances, or achieve an indirect objective.

However, the model can only choose economic effects within the permitted proposal schema.

The generator must not decide which proposal is selected.

## Step 4 — Decision point 1: Choose a proposal

Call the Decisions API using a `choice` question.

Provide the character's circumstances as input and the three generated proposals as choices.

Question intent:

"Given this protagonist's established goal, personality, resources, relationships and recent public history, which proposal would they most plausibly put forward today?"

Read the returned probability distribution.

Do not always use the highest-probability choice.

**Use weighted random sampling** to choose the proposal, so unlikely but plausible actions can occur.

Store all three choice probabilities and the selected proposal in the DayRecord.

The probabilities are model-estimated relative choices, not mathematically calibrated real-world frequencies.

## Step 5 — Decision point 2: Evaluate reactions and votes

Using the selected proposal, make one batched Decisions API request containing six `predicate` questions:

- Two questions for the remaining protagonists.
- Four questions for the citizen factions.

### Protagonist reaction questions

For each non-author protagonist:

"Given this character's goal, occupation, current wealth, feelings toward the author, and recent history, how likely are they to publicly support this proposal?"

Sample support or opposition using the returned probability.

These reactions do not count as additional votes in V1. They inform the newspaper and update interpersonal feelings.

### Faction voting questions

For each faction:

"How likely is a randomly selected citizen of this faction to vote YES on this policy, given the faction's economic position, loyalty toward the author, and the proposal's consequences?"

Use the returned support probability.

For a faction with 10 citizens, sample 10 independent yes/no votes.

A vote is either YES or NO. There are no abstentions in V1.

The sampled votes are totaled across the four factions.

### Voting rule

A proposal passes if:

`total_yes > total_no`

A 20–20 tie fails.

The protagonist's personal political reaction does not override the vote.

Store the probabilities, sampled votes and final result.

## Step 6 — Resolve consequences

If the proposal passes:

- Replace `ActivePolicy` with the selected proposal.

If it fails:

- Keep the existing policy.

In both cases, update political loyalties and interpersonal feelings based on what was proposed and how the characters reacted.

Use these simple deterministic rules initially.

### Faction loyalty changes

The target faction receiving a proposed benefit gains +5 loyalty toward the author.

The faction paying for the proposed benefit loses −5 loyalty toward the author.

These changes occur even if the proposal fails, because citizens remember what was proposed.

If the proposal passes, apply an additional +3 and −3 respectively.

Clamp loyalty values to the range 0–100.

### Protagonist feelings

For each non-author protagonist:

- If they supported the proposal, feelings toward the author increase by 2.
- If they opposed it, feelings toward the author decrease by 2.
- If their occupation benefits from the proposal, add another +5.
- If their occupation is harmed, add another −5.

Apply the occupation effect regardless of whether the law passed, representing their reaction to the attempt.

Clamp feelings to −100 through +100.

Feelings are directional. Changing Momo's feelings toward Pip does not automatically change Pip's feelings toward Momo.

These are game-design rules, not model predictions. Keep them isolated in one place so they can be tuned easily.

## Step 7 — Generate the newspaper

The Responses API receives only committed facts:

- Day number.
- Author and selected proposal.
- Voting result.
- Other protagonists' reactions.
- Consequent economic or political changes.
- Relevant previous history.

Generate:

- One dramatic headline.
- One short newspaper article, approximately 100–180 words.
- Optionally one short image prompt describing the main scene.

The newspaper must not invent uncommitted laws, deaths, characters or economic outcomes.

Humorous dialogue, jokes and colorful descriptions are welcome, but material events must remain faithful to the simulation.

If newspaper text generation fails after a valid simulation result, use a deterministic fallback headline and article. The core simulation must remain playable.

## Step 8 — Commit state

Commit the new world and append a DayRecord only after the essential proposal and Decisions stages succeed.

Store an end-of-day snapshot in the DayRecord.

Increment the completed day count.

Persist the world to JSON.

A failed core AI request must not leave the world partially advanced.

## Step 9 — Generate illustration

Generate one illustration for the day's newspaper asynchronously after committing the day.

The image is presentation-only; it must not affect world simulation state.

Set the day illustration status to `pending`, then update it to `ready` or `failed`.

Store generated image files locally, with their URLs referenced in the day record.

Use a stable art direction: whimsical illustrated newspaper editorial art featuring the same recurring characters.

Do not ask the image model to render newspaper headlines or long text. The frontend renders typography separately.

If image generation fails, show a visually appealing placeholder and preserve the newspaper text.

---

# 9. AI integration requirements

Build AI providers behind small, testable interfaces.

Suggested conceptual contracts:

```go
type ProposalGenerator interface {
    GenerateProposals(
        ctx context.Context,
        world World,
        actor Protagonist,
    ) ([]Proposal, error)

    GenerateNewspaper(
        ctx context.Context,
        world World,
        event DayRecord,
    ) (NewspaperContent, error)
}

type DecisionMaker interface {
    ChooseProposal(
        ctx context.Context,
        world World,
        actor Protagonist,
        candidates []Proposal,
    ) (ChoiceDistribution, error)

    EvaluateReactions(
        ctx context.Context,
        world World,
        proposal Proposal,
    ) (ReactionProbabilities, error)
}

type ImageGenerator interface {
    GenerateIllustration(
        ctx context.Context,
        prompt string,
    ) ([]byte, error)
}
```

Use appropriate concrete structs for `NewspaperContent`, `ChoiceDistribution` and `ReactionProbabilities`.

The simulation engine should not care whether these interfaces are implemented by real AI APIs or mocks.

### OpenAI Decisions API

Use:

`POST https://api.openai.com/v1/decisions`

Model:

`gpt-6-luna`

The proposal selection should be a `choice` question with three supplied options.

Conceptual request shape:

```json
{
  "model": "gpt-6-luna",
  "input": "Structured description of Pip, his goal, the world and recent history",
  "questions": [
    {
      "type": "choice",
      "name": "proposal_selection",
      "instructions": "Which proposal is this character most likely to introduce?",
      "choices": [
        {
          "value": "proposal-a",
          "description": "Transfer income from farmers to merchants"
        },
        {
          "value": "proposal-b",
          "description": "Transfer income from clergy to workers"
        },
        {
          "value": "proposal-c",
          "description": "Transfer income from merchants to farmers"
        }
      ]
    }
  ]
}
```

Read the `probabilities` array from the `choice` answer.

For reactions, use six separate named `predicate` questions in one request.

Each predicate returns the estimated probability that the stated condition is true.

Do not use the returned `confidence` field as the probability of a particular choice.

Handle model refusals, missing answers and malformed responses explicitly.

### Structured generation

Use Responses API Structured Outputs with an appropriate JSON schema for the three proposals.

Validate the returned candidates again in Go.

Use a configurable text model.

### Image generation

Use a currently supported GPT Image model.

Suggested default:

`gpt-image-2.5-flare`

Generate a single editorial illustration, preferably landscape-oriented.

Maintain a consistent style across editions.

### Model context management

Do not pass the complete raw historical transcript on every request.

For V1, provide:

- Current structured state.
- Protagonist goals and feelings.
- Current policy.
- Summaries of the most recent 5–7 days.

The complete historical log remains stored locally.

This is sufficient to establish recent continuity without growing prompts indefinitely.

---

# 10. Mock AI mode

**Mock mode is mandatory.**

The application must work without an OpenAI API key.

Use hard-coded example proposals and plausible deterministic support probabilities, with seed-controlled sampling.

Mock mode should use exactly the same simulation engine, data structures and endpoints as live mode.

No separate fake frontend.

Requirements:

- Produce three valid proposal candidates.
- Produce a probability distribution over candidates.
- Produce support probabilities for two protagonists and four factions.
- Generate a template-based headline and article.
- Use placeholder illustrations.
- Allow 30 days of simulation.
- Produce reproducible outcomes for identical starting seed and mock inputs.

It should be possible to demonstrate the entire product using mock mode.

Do not silently switch a live session to mock decisions after an OpenAI API failure. Surface the error instead.

---

# 11. Backend API

Only three application endpoints are required initially.

## `GET /api/world`

Returns the current world state, including completed history and simulation status.

Example response shape:

```json
{
  "world": {
    "id": "pebbleton-1",
    "day": 7,
    "protagonists": [],
    "factions": [],
    "active_policy": {},
    "history": []
  },
  "running": false,
  "max_days": 30
}
```

Use the actual structs for nested fields.

## `POST /api/advance`

Simulates exactly one day.

Request body can be empty.

Returns the updated world state or the newly completed day plus current state. Choose one consistent shape and document it.

Expected behavior:

- Only one advancement can run at a time.
- A second concurrent advance returns HTTP 409.
- Attempting to advance beyond Day 30 returns HTTP 422.
- Core AI failure returns an appropriate 5xx response and leaves the world unchanged.
- A successful request commits exactly one day.

Model calls can be slow. The frontend must show an explicit progress/loading state.

## `POST /api/reset`

Restores the predefined Day 0 world.

Clears active policy and history.

Reset the random seed to the configured default for reproducible mock runs.

Do not allow reset while an advance is actively running; return HTTP 409.

## Static assets

Serve generated illustrations from a local assets directory.

No external asset storage is required for V1.

## Concurrency

Protect the world with a mutex.

Do not hold the world mutex during slow external model requests.

Use a `running` guard, create a working world copy, execute simulation, then atomically commit the result.

Image generation must not overwrite or attach an illustration to an incorrect day after reset. Check world identity and day identity before applying an asynchronous image result.

---

# 12. Persistence

Use an in-memory `World` as the runtime source of truth.

After each committed day, write a JSON snapshot to disk.

Suggested location:

`data/world.json`

Use a temporary file plus atomic rename to avoid partially written snapshots.

On startup:

- Load an existing snapshot if valid.
- Otherwise initialize the predefined Day 0 world.

Images can be stored in:

`data/images/`

This is a single-instance local application. No database, migrations or persistence abstraction is required.

The README should explain that local filesystem persistence is not suitable for deployments with ephemeral storage or multiple independent backend instances.

---

# 13. Frontend experience

The frontend should feel like a playful miniature civilization dashboard, not an admin tool.

The newspaper is the visual centerpiece.

## Layout

### Header

- Pebbleton title/logo.
- Tagline: "A Tiny Society with Big Problems."
- Current day / 30.
- Next Day button.
- Reset button.

### Main region — newspaper

Show the selected day's newspaper with:

- Newspaper masthead: **The Pebbleton Post**.
- Day/edition number.
- Dramatic headline.
- Large illustrated image.
- Short article.
- Compact vote summary.
- Pass/fail result.

Use a vintage editorial aesthetic with warm paper tones, readable serif headlines, attractive spacing and well-designed typography.

Use CSS for newspaper text. Do not rely on AI to render text inside images.

### Side region — world dashboard

Show three protagonist cards.

Each card displays:

- Character name.
- Visual identity (simple icon, emoji or static portrait).
- Occupation.
- Persistent goal.
- Current money.
- Feelings toward the other protagonists.

Also show four faction cards containing:

- Occupation.
- Population.
- Money per person.
- Effective daily income.
- Loyalty toward protagonists.

Keep the initial dashboard compact. Detailed data can be shown in expandable sections or tooltips.

### Bottom region — history

Show a horizontal timeline of completed days.

Clicking an earlier day displays that day's newspaper and its recorded snapshot.

The user should be able to inspect previous policy decisions and see how the state has changed.

### During simulation

When the user clicks Next Day:

- Disable additional advance clicks.
- Show that the society is deliberating.
- Show simple phase messages, if possible, without implementing real-time streaming.
- When the response completes, show the new day.
- Display an illustration placeholder while image generation is pending.

React Query should manage API fetching and mutation state.

The frontend is a renderer. It must not implement simulation logic, calculate votes or locally mutate world state.

### Visual direction

Aim for:

- Storybook/fantasy editorial newspaper aesthetics.
- Expressive character identities.
- A slightly absurd sense of drama.
- Smooth transitions between days.
- Polished typography.
- A readable desktop layout suitable for projecting on a large screen.

Avoid generic dashboard styling.

The newspaper should be attractive enough to screenshot and share.

---

# 14. Suggested project structure

Use a structure similar to:

```text
pebbleton/
├── backend/
│   ├── cmd/server/
│   │   └── main.go
│   ├── internal/
│   │   ├── world/
│   │   │   ├── types.go
│   │   │   ├── seed.go
│   │   │   └── economy.go
│   │   ├── simulation/
│   │   │   ├── simulator.go
│   │   │   ├── voting.go
│   │   │   └── effects.go
│   │   ├── ai/
│   │   │   ├── interfaces.go
│   │   │   ├── openai.go
│   │   │   └── mock.go
│   │   ├── store/
│   │   │   └── json.go
│   │   └── httpapi/
│   │       └── handlers.go
│   ├── data/
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Newspaper.tsx
│   │   │   ├── ProtagonistCards.tsx
│   │   │   ├── FactionDashboard.tsx
│   │   │   └── DayTimeline.tsx
│   │   ├── api/
│   │   ├── types/
│   │   ├── App.tsx
│   │   └── main.tsx
│   └── package.json
├── .env.example
└── README.md
```

This is a suggested organization, not a mandate to create a complex architecture.

Avoid unnecessary wrappers, services and dependency injection frameworks.

---

# 15. Environment configuration

Provide `.env.example` with:

```env
OPENAI_API_KEY=
OPENAI_TEXT_MODEL=gpt-5.4-mini
OPENAI_DECISIONS_MODEL=gpt-6-luna
OPENAI_IMAGE_MODEL=gpt-image-2.5-flare

MOCK_AI=true
GENERATE_IMAGES=false

DATA_DIR=./data
PORT=8080
```

Make text and image model choices configurable.

When `MOCK_AI=true`, the backend should not require an API key.

When `MOCK_AI=false`, validate the OpenAI configuration on startup or return a clear configuration error.

Never expose the API key in frontend code or API responses.

Provide convenient local development commands and configure the Vite development proxy to reach the Go backend.

---

# 16. Tests

Write unit tests for core simulation logic.

At minimum cover:

### Economic calculations

- Income and upkeep are applied correctly.
- Active policy adjusts income.
- Replacement policy does not stack with old policy.
- Policy changes take effect on the following day.
- Negative balances are handled consistently.

### Proposal validation

- Cannot transfer income to the same faction.
- Cannot reference an unknown faction.
- Cannot create negative effective income.
- Reject malformed or duplicate proposal candidates.

### Voting

- A majority passes.
- A tie fails.
- Faction vote totals equal faction population.
- Probabilities of 0 and 1 produce expected votes.
- Mock mode is reproducible with the same seed.

### Relationships

- Loyalty changes are correctly applied.
- Feelings are directional.
- Values remain within their permitted ranges.

### Day advancement

- Exactly one DayRecord is added.
- Day increments exactly once.
- Failed model calls do not partially commit state.
- Concurrent advance requests cannot both execute.
- Day 30 is the maximum.
- Reset restores the initial world.

### History and persistence

- Records remain in chronological order.
- Historical snapshots do not change when future days advance.
- JSON serialization and reload preserve world state.

### Smoke test

Run a complete 30-day mock simulation.

Verify:

- All 30 days complete.
- Every day has a valid proposal.
- Every day has voting results and a newspaper.
- There are no invalid faction references.
- No income becomes negative.
- No state invariants are violated.
- The run is repeatable for a fixed seed.

Live model integration can be tested separately, but automated tests must not depend on live API calls.

---

# 17. Implementation milestones

Implement in this order.

## Milestone 1 — Deterministic core

Build:

- All Go structs.
- Predefined world seed.
- Economy calculations.
- Proposal validation.
- Voting and weighted sampling.
- Relationship updates.
- Mock AI.
- Append-only history.
- Simulator unit tests.

**Definition of done:** Run a 30-day simulation in Go with mock choices and no OpenAI dependency.

## Milestone 2 — HTTP backend and frontend

Build:

- Three API endpoints.
- React application.
- World dashboard.
- Newspaper component with mock content.
- Day navigation.
- Next Day and Reset controls.
- Loading and error states.

**Definition of done:** A person can play the simulation through their browser.

## Milestone 3 — Live OpenAI integration

Build:

- Structured proposal generation using Responses.
- Proposal selection using Decisions `choice`.
- Batched reaction and voting using Decisions `predicate`.
- Newspaper text generation using Responses.
- Explicit error handling and logging.

**Definition of done:** Live model calls drive the same simulator without modifications to core business logic.

## Milestone 4 — Illustration and visual polish

Build:

- GPT Image editorial illustrations.
- Asynchronous image lifecycle.
- Stable character art descriptions.
- Newspaper visual refinements.
- Historical newspaper navigation.
- Responsive desktop UI.

**Definition of done:** The product looks presentable during a live hackathon demonstration.

At the end of each milestone, verify the application still runs.

If time is constrained, preserve core functionality first. Illustration generation must not compromise the simulation.

---

# 18. Acceptance criteria

The project is considered V1-complete when:

- [ ] The application starts locally with a documented command.
- [ ] Mock mode works without an OpenAI API key.
- [ ] The world initializes with exactly 3 protagonists and 4 citizen factions.
- [ ] The user can advance one day through the UI.
- [ ] One protagonist proposes 3 candidate economic policies.
- [ ] The selected proposal comes from the Decisions API in live mode.
- [ ] The other 2 protagonists react.
- [ ] All 4 factions vote using model-estimated probabilities.
- [ ] Votes are sampled and tallied deterministically by Go logic.
- [ ] Policies pass or fail according to the voting rules.
- [ ] Policy results persist into subsequent days.
- [ ] Money, loyalties and feelings change over time.
- [ ] Every completed day generates a newspaper.
- [ ] The newspaper accurately reflects recorded simulation events.
- [ ] Newspaper illustrations work when enabled.
- [ ] Historical days can be viewed.
- [ ] The simulation can complete all 30 days.
- [ ] Reset restores the original world.
- [ ] The simulation is protected against partial commits.
- [ ] Backend unit tests pass.
- [ ] Frontend builds successfully.
- [ ] The README explains how to run mock mode and live AI mode.

---

# 19. Product principles for future extensibility

Do not implement these future features yet, but preserve the following boundaries.

### Generation and decisions are different

Generation produces candidate actions.

Decisions evaluates supplied possibilities.

Neither directly updates the world.

### The world engine owns facts

Only validated, deterministic state transitions create economic or political consequences.

The newspaper reports the results; it does not invent new facts for the world.

### History is append-only

Every day produces a persistent record.

Future features may include private memories, secrets, institutions and factions with their own histories.

V1 requires only one public history.

### Human interventions should reuse existing contracts

Eventually the crowd can select one of the generated proposals or inject a valid proposal.

The voting and state-transition engine should operate identically regardless of whether the proposal was selected by AI or humans.

Do not implement crowd participation in V1.

### Interesting behavior comes from feedback

The purpose is not to generate random dramatic headlines.

A character's previous actions should alter faction loyalties, relationships and economic conditions, which then affect future decisions.

A successful simulation produces histories that are surprising but causally understandable.

---

# 20. Expected handoff

When implementation is complete:

1. Report what was implemented.
2. Provide exact local setup and run commands.
3. Explain the environment variables.
4. State whether mock mode and live AI mode were both tested.
5. Summarize tests run and their results.
6. Identify any incomplete features or limitations.
7. State any assumptions or deviations from the spec.

Do not describe functionality as complete unless it has actually been implemented.

**Final product goal:** A small, charming, autonomous civilization in which three ambitious characters make economic proposals, ordinary citizens respond, and history unfolds one newspaper edition at a time.

The simulation can be simple. Its consequences must persist. Its characters should have recognizable motivations. And its newspaper should make the whole experience feel alive.
