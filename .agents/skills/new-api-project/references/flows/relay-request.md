---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "router/relay-router.go"
  - "middleware/distributor.go"
  - "controller/relay.go"
  - "relay/common/relay_info.go"
  - "relay/relay_adaptor.go"
  - "service/billing_session.go"
  - "service/tiered_settle.go"
---

# Relay Request Flow

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Trigger

A client sends a relay request through an OpenAI-compatible, Claude, Gemini, realtime, image, audio, rerank, playground, Midjourney, Suno, or task endpoint.

## Actors

- Client
- Gin router and middleware
- User/token auth middleware
- Distributor middleware
- Controller relay orchestration
- Relay format helper
- Provider adaptor
- Billing session
- Model/log persistence

## Sequence

1. The route group in `router/relay-router.go` selects auth and relay format.
2. Auth middleware stores token/user/group/quota context.
3. `middleware.Distribute()` parses the model from JSON, form data, query, or provider-specific path. It checks token model restrictions and group access.
4. Distributor resolves a preferred affinity channel when possible; otherwise it asks `service.CacheGetRandomSatisfiedChannel()`.
5. `SetupContextForSelectedChannel()` stores channel metadata, selected key, overrides, mappings, base URL, and provider-specific `Other` fields in context.
6. `controller.Relay()` validates the request DTO with `helper.GetAndValidateRequest()`.
7. `relay/common.GenRelayInfo()` builds `RelayInfo`, including token context, model, relay format, stream status, request headers, and conversion chain.
8. Sensitive text checks and token estimation run when enabled.
9. `helper.ModelPriceHelper()` computes price/pre-consume data. Tiered billing freezes a `BillingSnapshot` and request input.
10. `service.PreConsumeBilling()` creates the billing session unless the model is free.
11. The controller retries channel execution up to `common.RetryTimes`, replaying the stored body for each attempt.
12. The format-specific relay helper selects a provider adaptor via `relay.GetAdaptor()` or task adaptor via `relay.GetTaskAdaptor()`.
13. On success, usage is settled and metrics/logs are recorded.
14. On failure after pre-consume, billing is refunded or violation fees are charged according to the error.

## Failure Modes

- Invalid or oversized body maps to request-body errors, with oversized bodies normalized to HTTP 413.
- No eligible channel returns a model/channel availability error.
- Provider errors can disable channels or skip retry depending on error metadata and status-code rules.
- Billing errors before upstream call skip retry.
- Tiered expression settlement errors fall back to the frozen pre-consume estimate.

## Side Effects

- Context stores selected channel, key, base URL, mappings, and billing metadata.
- User/token/subscription quota can be pre-consumed, settled, or refunded.
- Channel affinity can be recorded after successful relay.
- Logs and performance metrics can be recorded asynchronously.

## Tests

Relevant tests live in `middleware/`, `controller/`, `relay/`, `service/`, `dto/`, and `pkg/billingexpr/`.
