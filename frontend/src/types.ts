export type Protagonist = {
  id: string;
  name: string;
  occupation_id: string;
  goal: string;
  money: number;
  feelings: Record<string, number>;
};
export type Faction = {
  id: string;
  name: string;
  population: number;
  money: number;
  base_income: number;
  effective_income: number;
  upkeep: number;
  loyalty: Record<string, number>;
};
export type Proposal = {
  id: string;
  author_id: string;
  title: string;
  reason: string;
  from_faction: string;
  to_faction: string;
};
export type Snapshot = {
  protagonists: Protagonist[];
  factions: Faction[];
  active_policy: Proposal | null;
};
export type DayRecord = {
  day: number;
  author_id: string;
  candidates: Proposal[];
  selection_probabilities: Record<string, number>;
  selected_proposal: Proposal;
  faction_votes: {
    faction_id: string;
    support_probability: number;
    yes_votes: number;
    no_votes: number;
  }[];
  reactions: {
    protagonist_id: string;
    support_probability: number;
    supports: boolean;
  }[];
  total_yes: number;
  total_no: number;
  passed: boolean;
  headline: string;
  article: string;
  illustration_status: string;
  illustration_url: string;
  snapshot: Snapshot;
};
export type World = Snapshot & {
  id: string;
  day: number;
  seed: number;
  history: DayRecord[];
};
export type WorldResponse = {
  world: World;
  running: boolean;
  max_days: number;
  mode: string;
};
