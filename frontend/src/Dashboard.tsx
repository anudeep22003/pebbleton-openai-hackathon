import type { Snapshot } from "./types";
import { CharacterArt } from "./TownArt";

const portraits: Record<string, string> = { pip: "🦊", momo: "🐧", zig: "🐸" };
const occupations: Record<string, string> = {
  merchants: "Merchant",
  farmers: "Farmer",
  workers: "Worker",
};
const factionIcons: Record<string, string> = {
  farmers: "✿",
  merchants: "◇",
  workers: "⚒",
  clergy: "✧",
};
export function Dashboard({ snapshot }: { snapshot: Snapshot }) {
  return (
    <>
      <div className="characters">
        {snapshot.protagonists.map((person) => (
          <details className={`character ${person.id}`} key={person.id}>
            <summary>
              <span className="portrait">
                <CharacterArt id={person.id} size={47} />
              </span>
              <span className="character-name">
                <b>{person.name}</b>
                <small>
                  {occupations[person.occupation_id]} <span>·</span>{" "}
                  {person.id === "pip"
                    ? "The opportunist"
                    : person.id === "momo"
                      ? "The idealist"
                      : "The pragmatist"}
                </small>
              </span>
              <span className="money">
                {person.money}
                <small>COINS</small>
              </span>
            </summary>
            <p className="goal">“{person.goal}.”</p>
            <div className="feelings">
              {Object.entries(person.feelings).map(([id, feeling]) => (
                <span key={id}>
                  {portraits[id]} {feeling > 0 ? "+" : ""}
                  {feeling} <small>toward {id}</small>
                </span>
              ))}
            </div>
          </details>
        ))}
      </div>
      <div className="sidebar-note">
        Three ambitions.
        <br />
        <i>One very small town.</i>
      </div>
      <section className="factions">
        <div className="section-heading">
          <span>THE CITIZENRY</span>
          <span>40 STRONG</span>
        </div>
        {snapshot.factions.map((faction) => (
          <details className="faction" key={faction.id}>
            <summary>
              <span className="faction-icon">{factionIcons[faction.id]}</span>
              <span>
                <b>{faction.name}</b>
                <small>{faction.population} citizens</small>
              </span>
              <span className="faction-wealth">
                <b>
                  {faction.money} <small>¢ / person</small>
                </b>
                <small>
                  {faction.effective_income} income − {faction.upkeep} upkeep
                </small>
              </span>
            </summary>
            <div className="loyalty">
              <small>LOYALTY TO THE LOCALS</small>
              {snapshot.protagonists.map((person) => (
                <div key={person.id}>
                  <span>{person.name}</span>
                  <meter
                    min={0}
                    max={100}
                    value={faction.loyalty[person.id]}
                    aria-label={`${faction.name} loyalty to ${person.name}`}
                  />
                  <b>{faction.loyalty[person.id]}</b>
                </div>
              ))}
            </div>
          </details>
        ))}
      </section>
      <section className="policy-note">
        <span className="eyebrow">ON THE BOOKS</span>
        <h3>{snapshot.active_policy?.title ?? "A clean slate."}</h3>
        <p>
          {snapshot.active_policy
            ? `One coin of daily income moves from ${snapshot.active_policy.from_faction} to ${snapshot.active_policy.to_faction}.`
            : "Everyone earns their base income. For now, the only thing being redistributed is gossip."}
        </p>
        <small>
          Income shown is the rate for the next day.
          <br />
          Click a local or faction to inspect their loyalties.
        </small>
      </section>
    </>
  );
}
