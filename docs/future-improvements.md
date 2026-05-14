# Future improvements

This section tracks **directionally accurate** enhancements that fit the current architecture without rewriting it.

## Metrics & statistics

- Ingest **percentile** CPU/memory series (P90/P99) alongside averages.
- Weight recommendations by **seasonality** or **deployment class** (batch vs serving).

## Policy & multi-tenancy

- Pluggable **policy engine** (OPA/Cedar) to enforce org-specific floors per namespace.
- Per-tenant **buffer overrides** sourced from CRDs instead of global env.

## Control plane integration

- Optional **admission webhook** or **VPA integration** that reads recommendations as **hints** rather than blind patches.
- **Dry-run** mode returning a **diff** against live Deployment specs.

## Data plane

- Optional **PostgreSQL** or **S3** artifact store for historical runs and audit trails.
- **Idempotent** batch job mode (CLI + CronJob) writing structured reports.

## API surface

- Optional **JWT** or **mTLS** hooks documented as deployment patterns.
- Pagination / streaming for **very large** batches.

## Non-goals (recap)

- Becoming a full **cost optimization** suite (node pricing, reservations) without additional data sources.
- Replacing **HPA** or **cluster autoscaler** decision loops.

Contributions should keep the **finalize** guardrails and **config validation** paths strong — they are what keep recommendations trustworthy when the world is messy.
