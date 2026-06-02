---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "router/relay-router.go"
  - "controller/relay.go"
  - "middleware/distributor.go"
  - "relay/relay_adaptor.go"
  - "relay/common/relay_info.go"
  - "relay/channel/adapter.go"
  - "dto/openai_request.go"
  - "types/relay_format.go"
---

# Relay Module

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `router/relay-router.go`
- `controller/relay.go`
- `middleware/distributor.go`
- `relay/relay_adaptor.go`
- `relay/common/relay_info.go`
- `relay/channel/`
- `relay/helper/`
- `relay/common/`
- `dto/`
- `types/relay_format.go`

## Responsibility

The relay module receives OpenAI-compatible, Claude, Gemini, rerank, image, audio, realtime, Midjourney, Suno, and video/task-style requests, selects an eligible channel, converts request/response formats, handles streaming, retries failed channels, and settles billing.

## Route Surface

`SetRelayRouter()` applies CORS, decompression, body storage cleanup, and stats middleware globally for relay routing. Important groups include:

- `/v1/models`, `/v1beta/models`, `/v1beta/openai/models`: token-authenticated model listing/retrieval.
- `/pg/chat/completions`: playground route using user auth.
- `/v1/realtime`: WebSocket realtime relay.
- `/v1/messages`, `/v1/chat/completions`, `/v1/responses`, `/v1/images/*`, `/v1/audio/*`, `/v1/embeddings`, `/v1/rerank`, `/v1/models/*path`: main HTTP relay formats.
- `/mj`, `/:mode/mj`, `/suno`, `/v1beta/models/*path`: task and provider-specific relay routes.

## Request Flow

1. `middleware.TokenAuth()` or `middleware.UserAuth()` establishes user/token context.
2. `middleware.ModelRequestRateLimit()` applies model-level rate limits for token relay paths.
3. `middleware.Distribute()` extracts the model, checks token model restrictions, resolves affinity, selects a channel, and stores channel metadata in Gin context.
4. `controller.Relay()` validates the request DTO, builds `relay/common.RelayInfo`, checks sensitive text when enabled, estimates tokens, computes price data, and pre-consumes billing.
5. `controller.Relay()` retries channels through `getChannel()` until success, no retry budget remains, or an error is marked skip-retry.
6. Format-specific helpers in `relay/` call the provider adaptor from `relay/relay_adaptor.go`.
7. On success, usage and billing are settled. On failure after pre-consume, `BillingSession.Refund()` or violation fee handling runs.

## Provider Adaptors

`relay/relay_adaptor.go` maps channel/API types to concrete `relay/channel/*` adaptors. When adding or changing a provider:

- Add or update the constant mapping in `constant/` and adaptor switch logic as needed.
- Keep provider-specific DTOs and conversion inside `relay/channel/<provider>/`.
- Confirm `StreamOptions` support and update `streamSupportedChannels` in `relay/common/relay_info.go` only when supported.
- Preserve explicit client zero values in DTOs that are parsed and re-marshaled upstream.
- Add focused tests beside the provider or helper code.

## Invariants

- `RelayInfo.ChannelMeta` is initialized from Gin context after channel selection.
- `RelayInfo.RequestConversionChain` records format conversion order and should remain accurate for logging/billing/debugging.
- `common.GetBodyStorage(c)` is used so retry paths can replay the request body.
- Streaming helpers must maintain normal end/error status for billing and diagnostics.

## Tests

Useful test areas include:

- `dto/openai_request_zero_value_test.go`
- `relay/common/*_test.go`
- `relay/helper/*_test.go`
- `relay/channel/*/*_test.go`
- `service/text_quota_test.go`
- `service/tiered_settle_test.go`
