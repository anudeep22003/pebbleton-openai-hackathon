import type { WorldResponse } from "./types";

export async function requestWorld(action = "world"): Promise<WorldResponse> {
  const response = await fetch(`/api/${action}`, {
    method: action === "world" ? "GET" : "POST",
  });
  const result = await response.json();
  if (!response.ok)
    throw new Error(result.error || "The town hall could not be reached.");
  return result;
}
