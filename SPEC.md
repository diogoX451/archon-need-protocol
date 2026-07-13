# Archon Need Protocol v1

## Overview

The **need protocol** is how an orchestrator offloads work to external executors
without blocking the core runtime.

```
orchestrator --NeedEvent--> bus subject (e.g. need.http)
executor    --WebhookBody--> POST {api}/webhooks/needs/{correlation_id}
orchestrator resumes paused agent by correlation_id
```

## NeedEvent

| Field | Type | Notes |
|-------|------|-------|
| schema_version | int | `1` current; `0`/omitted = legacy |
| net_id | string | optional workflow/net id |
| workflow_id | string | optional |
| need.type | string | need kind (`http`, `planner`, `mcp`, …) |
| need.correlation_id | string | **required** resume key |
| need.payload | object | tool-specific JSON |
| need.created_at | RFC3339 | optional |

## Executor rules

1. **Ack vs redeliver**: permanent errors (bad JSON, 4xx) should not redeliver forever.
2. **Idempotency**: handlers may see redelivery; make side effects safe.
3. **Webhook**: POST JSON `{"payload": ...}` to the orchestrator; include delivery auth header if configured.
4. **Schema**: prefer `schema_version: 1` on publish; dual-read legacy.

## Error classification (recommended)

| Class | Action |
|-------|--------|
| Transient (timeout, 5xx, 429) | return error → bus redelivers |
| Permanent (invalid, 4xx auth) | ack and drop / report failure payload |

## Out of scope

Tenant billing, channel delivery, and product UI are **not** part of this protocol.

## Compliance suite

This repository includes golden JSON under `testdata/golden/` and tests in
`compliance_test.go`. Language ports should pass the same fixtures:

| Fixture | Expectation |
|---------|-------------|
| `need_v1.json` | Parse schema_version=1, correlation_id, nested payload |
| `need_legacy_toplevel.json` | Dual-read legacy top-level need fields |
| `webhook_body.json` | `{ "payload": ... }` shape for orchestrator resume |

### Executor checklist (community implementations)

- [ ] Refuse empty `correlation_id`
- [ ] Classify permanent vs transient errors
- [ ] Idempotent side effects under redelivery
- [ ] POST webhook before ack on success
- [ ] Do not ack on transient webhook/network failure

## Example

See [`examples/echo-executor`](examples/echo-executor) for a full NATS → parse → webhook loop.
