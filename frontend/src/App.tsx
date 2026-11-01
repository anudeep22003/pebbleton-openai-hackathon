import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { requestWorld } from "./api";
import { Newspaper } from "./Newspaper";
import { Dashboard } from "./Dashboard";
import type { WorldResponse } from "./types";

export default function App() {
  const client = useQueryClient();
  const [selectedDay, setSelectedDay] = useState<number | null>(null);
  const worldQuery = useQuery({
    queryKey: ["world"],
    queryFn: () => requestWorld(),
    refetchInterval: (query) =>
      query.state.data?.running ||
      query.state.data?.world.history.some(
        (d) => d.illustration_status === "pending",
      )
        ? 1500
        : 5000,
  });
  function received(data: WorldResponse) {
    client.setQueryData(["world"], data);
    setSelectedDay(null);
  }
  const advance = useMutation({
    mutationFn: () => requestWorld("advance"),
    onSuccess: received,
    onSettled: () => client.invalidateQueries({ queryKey: ["world"] }),
  });
  const reset = useMutation({
    mutationFn: () => requestWorld("reset"),
    onSuccess: received,
  });
  const data = worldQuery.data;
  if (!data)
    return (
      <main className="loading">
        <span className="brand-icon">✳</span>
        <h1>Opening the town gates…</h1>
        {worldQuery.error && (
          <>
            <p role="alert">
              Could not reach Pebbleton. Check that the Go server is running.
            </p>
            <button onClick={() => worldQuery.refetch()}>Try again</button>
          </>
        )}
      </main>
    );
  const { world, max_days, mode } = data;
  const busy = advance.isPending || reset.isPending || data.running;
  const day =
    selectedDay === null ? world.day : Math.min(selectedDay, world.day);
  const edition = world.history.find((record) => record.day === day);
  // Day 0 is displayed only before the first advance. All past views use recorded snapshots.
  const snapshot = edition?.snapshot ?? world;
  const error = advance.error || reset.error;
  return (
    <div className="app-shell">
      <header className="app-header">
        <a className="brand" href="/" aria-label="Pebbleton home">
          <span className="brand-icon">✳</span>
          <span>
            <strong>
              Pebbleton<span className="brand-period">.</span>
            </strong>
            <small>A Tiny Society with Big Problems.</small>
          </span>
        </a>
        <div className="controls">
          <span className="day-counter">
            DAY <b>{String(world.day).padStart(2, "0")}</b>
            <span> / {max_days}</span>
          </span>
          <button
            className="reset-button"
            disabled={busy}
            onClick={() => {
              if (
                window.confirm(
                  "Begin a new 30-day story? This clears the current history.",
                )
              ) {
                advance.reset();
                reset.mutate();
              }
            }}
          >
            ↺ <span>Reset</span>
          </button>
          <button
            className="next-button"
            disabled={busy || world.day >= max_days}
            onClick={() => {
              reset.reset();
              advance.mutate();
            }}
          >
            {busy
              ? "Deliberating…"
              : world.day >= max_days
                ? "Story complete ✓"
                : "Next Day"}{" "}
            <span aria-hidden="true">↗</span>
          </button>
        </div>
      </header>
      <div className="town-strip">
        <span>
          <i className="live-dot" />{" "}
          {mode === "mock"
            ? "A LITTLE WORLD, RUNNING ON IMAGINATION"
            : "A LITTLE WORLD, POWERED BY OPENAI"}
        </span>
        <span>40 CITIZENS · 3 BIG PERSONALITIES · 1 SHARED FUTURE</span>
      </div>
      {error && (
        <div className="error-banner" role="alert">
          {error.message}
        </div>
      )}
      {busy && (
        <div className="progress-banner" role="status">
          <span className="spinner" /> The town is deliberating. Proposals,
          opinions, forty little votes…
        </div>
      )}
      {world.day === max_days && (
        <div className="complete-banner">
          Thirty days. A town changed forever. Browse the archives, or reset to
          tell the story again.
        </div>
      )}
      <main className="main-grid">
        <section className="edition-region">
          <div className="section-heading">
            <span>THE DAILY DISPATCH</span>
            <span>
              {day < world.day
                ? `FROM THE ARCHIVES · DAY ${day}`
                : "FRESH OFF THE PRESS"}{" "}
              <span aria-hidden="true">↙</span>
            </span>
          </div>
          <Newspaper
            key={`${world.id}-${day}`}
            edition={edition}
            snapshot={snapshot}
          />
        </section>
        <aside>
          <div className="section-heading">
            <span>MEET THE LOCALS</span>
            <span>DAY {day}</span>
          </div>
          <Dashboard snapshot={snapshot} />
        </aside>
      </main>
      <nav className="history" aria-label="Newspaper archive">
        <div className="history-title">
          <span className="eyebrow">THE TOWN ARCHIVES</span>
          <h2>Every day leaves a story.</h2>
        </div>
        <div className="timeline">
          {world.history.length === 0 ? (
            <p>Your first edition is one click away.</p>
          ) : (
            world.history.map((record) => (
              <button
                key={record.day}
                className={day === record.day ? "selected" : ""}
                aria-label={`Read day ${record.day}`}
                aria-pressed={day === record.day}
                onClick={() => setSelectedDay(record.day)}
              >
                <small>DAY</small>
                <b>{String(record.day).padStart(2, "0")}</b>
                <i className={record.passed ? "passed-dot" : "rejected-dot"} />
              </button>
            ))
          )}
        </div>
      </nav>
      <footer>
        <span>
          PEBBLETON <span aria-hidden="true">✳</span> EST. THIS MORNING
        </span>
        <span>Small town. Long memories.</span>
      </footer>
    </div>
  );
}
