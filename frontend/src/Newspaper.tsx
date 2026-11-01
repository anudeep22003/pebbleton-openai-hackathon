import { useState } from "react";
import type { DayRecord, Snapshot } from "./types";
import { TownArt } from "./TownArt";

const names: Record<string, string> = { pip: "Pip", momo: "Momo", zig: "Zig" };
export function Newspaper({
  edition,
  snapshot,
}: {
  edition?: DayRecord;
  snapshot: Snapshot;
}) {
  const proposal = edition?.selected_proposal;
  const [imageFailed, setImageFailed] = useState(false);
  const hasImage =
    edition?.illustration_status === "ready" &&
    edition.illustration_url &&
    !imageFailed;
  return (
    <article className="newspaper">
      <div className="paper-topline">
        <span>INDEPENDENT. MOSTLY.</span>
        <span>ONE COIN. MANY OPINIONS.</span>
      </div>
      <h1 className="masthead">The Pebbleton Post</h1>
      <div className="edition-line">
        <span>
          VOL. I ·{" "}
          {edition
            ? `EDITION ${String(edition.day).padStart(2, "0")}`
            : "THE FOUNDING EDITION"}
        </span>
        <span>ALL THE NEWS THAT FITS IN A TINY TOWN</span>
        <span>PRICE: YOUR ATTENTION</span>
      </div>
      <div className="story-kicker">
        <span>
          ◆ {edition ? "POLITICS & POCKET CHANGE" : "A NEW CHAPTER BEGINS"}
        </span>
        {edition && (
          <span className={`verdict ${edition.passed ? "passed" : "rejected"}`}>
            {edition.passed ? "MOTION PASSED" : "MOTION REJECTED"}
          </span>
        )}
      </div>
      <h2 className="headline">
        {edition?.headline ??
          "A tiny town. A fresh start. What could possibly go wrong?"}
      </h2>
      <p className="story-deck">
        {proposal
          ? proposal.reason
          : "Forty citizens, three ambitious neighbours, and an alarming number of opinions about one coin."}
      </p>
      <figure className="editorial">
        <div className="editorial-art">
          {hasImage ? (
            <img
              className="generated-art"
              src={edition.illustration_url}
              alt={`Editorial illustration for day ${edition.day}: ${edition.headline}`}
              onError={() => setImageFailed(true)}
            />
          ) : (
            <TownArt day={edition?.day} />
          )}
          {edition?.illustration_status === "pending" && (
            <span className="illustration-badge">
              <span className="spinner" /> OUR ARTIST IS AT WORK
            </span>
          )}
        </div>
        <figcaption>
          {hasImage
            ? "AN EDITORIAL ILLUSTRATION FROM THE PEBBLETON POST · GPT IMAGE"
            : "PEBBLETON, AS SEEN BY OUR RESIDENT OPTIMIST."}{" "}
          {(edition?.illustration_status === "failed" || imageFailed) &&
            "THE ARTIST MISSED DEADLINE. THE NEWS DID NOT."}
        </figcaption>
      </figure>
      <div className="byline">
        <span>
          BY THE PEBBLETON POST <i>· Our town affairs desk</i>
        </span>
        <span>DISPATCH № {String(edition?.day ?? 0).padStart(2, "0")}</span>
      </div>
      <div className="article-body">
        {(
          edition?.article ??
          "Welcome to Pebbleton, a place small enough for everyone to know your business and ambitious enough to call that an economy. Pip wants prosperity, preferably in his own pocket. Momo wants a fair deal. Zig would quite like everybody to remain calm.\n\nEach day, one neighbour will put a proposal to forty citizens: move one coin of daily income from one faction to another. The town votes. Friendships wobble. Tomorrow's pay might look a little different.\n\nThe Post will be here to document every triumph, every grievance, and every magnificently overblown dispute. Turn the page. History is waiting to happen."
        )
          .split("\n\n")
          .map((paragraph, i) => (
            <p key={i}>{paragraph}</p>
          ))}
      </div>
      {edition ? (
        <>
          <section className="ballot">
            <div>
              <span className="eyebrow">THE PEOPLE HAVE SPOKEN</span>
              <h3>
                {edition.total_yes} for. {edition.total_no} against.
              </h3>
            </div>
            <div
              className="vote-bar"
              aria-label={`${edition.total_yes} yes votes, ${edition.total_no} no votes`}
            >
              <span style={{ width: `${(edition.total_yes / 40) * 100}%` }} />
            </div>
            <div className="faction-votes">
              {edition.faction_votes.map((vote) => (
                <div key={vote.faction_id}>
                  <b>{vote.faction_id}</b>
                  <span>
                    {vote.yes_votes} yes <i> / {vote.no_votes} no</i>
                  </span>
                </div>
              ))}
            </div>
          </section>
          <div className="reactions">
            {edition.reactions.map((reaction) => (
              <span key={reaction.protagonist_id}>
                <b>{names[reaction.protagonist_id]}</b>{" "}
                {reaction.supports ? "backs the plan ↗" : "is unconvinced ↘"}
              </span>
            ))}
          </div>
          <details className="decision-details">
            <summary>Inside the decision · candidates & probabilities</summary>
            <p className="muted">
              Estimated relative choices, followed by sampled votes.
            </p>
            {edition.candidates.map((candidate) => (
              <div key={candidate.id}>
                <b>
                  {candidate.id === proposal?.id ? "✓ " : ""}
                  {candidate.title}
                </b>
                <span>
                  {Math.round(
                    edition.selection_probabilities[candidate.id] * 100,
                  )}
                  %
                </span>
                <small>
                  {candidate.from_faction} → {candidate.to_faction}
                </small>
              </div>
            ))}
            <p>
              {edition.faction_votes
                .map(
                  (vote) =>
                    `${vote.faction_id}: ${Math.round(vote.support_probability * 100)}% support`,
                )
                .join(" · ")}
            </p>
            <p>
              {edition.reactions
                .map(
                  (reaction) =>
                    `${names[reaction.protagonist_id]}: ${Math.round(reaction.support_probability * 100)}% likely to support`,
                )
                .join(" · ")}
            </p>
          </details>
        </>
      ) : (
        <div className="welcome-note">
          <span>✳</span>
          <div>
            <b>The presses are ready.</b>
            <p>Click Next Day to see what your little society does first.</p>
          </div>
        </div>
      )}
      <div className="paper-bottom">
        {snapshot.active_policy
          ? `INCOME POLICY: ${snapshot.active_policy.from_faction} → ${snapshot.active_policy.to_faction} · 1 COIN / PERSON / DAY`
          : "NO POLICY YET. ENJOY THE PEACE WHILE IT LASTS."}
      </div>
    </article>
  );
}
