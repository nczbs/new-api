---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "AGENTS.md"
  - "main.go"
  - "router/"
  - "relay/"
  - "service/"
  - "model/"
  - "web/default/"
---

# Glossary

Last verified: 2026-06-02
Verified against commit: 7aaa5332

- API type: Internal category mapping from channel/provider type to adaptor behavior.
- Channel: An upstream provider configuration containing type, keys, base URL, model mappings, settings, and availability.
- Channel affinity: Logic that prefers a previously successful channel for matching model/group/request attributes.
- Distributor: `middleware.Distribute()`, which parses model intent, validates access, and selects a channel before relay execution.
- Funding source: Wallet quota or subscription quota used by `BillingSession`.
- Relay format: Incoming or outgoing API contract such as OpenAI, Claude, Gemini, OpenAI Responses, audio, image, embedding, rerank, realtime, task, or MJ proxy.
- RelayInfo: Shared per-request state passed through relay helpers and adaptors.
- StreamOptions: Optional upstream streaming usage option; only supported channels should be marked in `streamSupportedChannels`.
- Tiered expression billing: Billing mode where a self-contained expression calculates quota from token counts and request context.
- Token group: Group context used for token access, model restrictions, pricing, and channel selection.
- Trust quota: Threshold that can bypass pre-consume for sufficiently funded users/tokens.
- web/default: Primary React 19 frontend.
- web/classic: Separate classic frontend theme.
