# Recommendation engine

The engine turns **average utilization** and **declared requests** into **conservative** recommendations. It is deliberately small and testable: heuristics propose a draft, policy caps growth, **finalize** enforces hard invariants, and **scoring** summarizes risk.

## Pipeline (per workload)

1. **Heuristics (CPU)** — `round(cpu_usage_avg × cpu_safety_buffer)` then floor to `min_cpu_millicores`. The effective buffer is never below **1.0** even if configuration were corrupted (defense in depth after startup normalization).
2. **Heuristics (memory)** — same pattern with `memory_safety_buffer` and `min_memory_mb`.
3. **Downsize-only cap** — if the draft exceeds the submitted `cpu_request` / `memory_request`, it is clamped back to the current request (the service does not recommend **increasing** requests from this endpoint).
4. **Finalize (guardrails)**  
   - Never return **zero** when the corresponding request is positive.  
   - Never recommend **below** observed average usage.  
   - Apply configured **minimum floors** and a minimum of **1** unit when requests are positive.  
   - Clamp back to **≤ current request** after raising floors.  
   - If the draft was **zero or below usage** before repair, the item is marked **degraded** for scoring.
5. **Scoring** — derives `confidence_score` and `severity`.

## Severity

| Value | Meaning |
|-------|---------|
| `low` | Small modeled waste (low reduction percentages). |
| `moderate` | Meaningful but not extreme reduction. |
| `high` | Large reduction opportunity. |
| `critical` | Very large reduction opportunity (**only when not degraded**). |
| `degraded` | Emergency clamps fired; **do not** interpret reduction magnitude as a confident FinOps signal until metrics/config are reviewed. |

## Confidence heuristic

- Baseline **0.85** when both CPU and memory usage averages are **> 0**.
- If both `recommended / usage` ratios exceed **2.0**, a small boost applies.
- If either ratio falls **below 1.2**, a penalty applies (tighter headroom vs averages).
- If either average is **zero**, a penalty applies and ratios are **not** divided (avoids NaN/Inf).
- If the workload is **degraded**, confidence is pinned low (**0.25**) so automation does not over-trust the number.

## Workload archetypes (expected behavior)

| Profile | Typical signal | Expected outcome |
|---------|----------------|------------------|
| Heavily over-provisioned | Usage ≪ request | Meaningful **reduction** with moderate/high severity (if not degraded) |
| Well provisioned | Usage near request | **Little or no reduction**, low severity |
| Nearly saturated | Usage approaches request | **No reduction** after caps; may show **low confidence** if ratios sit near buffers |
| Invalid metrics | Usage > request | **400** from validation (by design today) |

## Related code

- `internal/engine/heuristics` — draft math  
- `internal/engine/recommendation` — orchestration + finalize  
- `internal/engine/scoring` — confidence + severity  
- `internal/config` — buffers, floors, validation  
