---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "pkg/billingexpr/expr.md"
  - "pkg/billingexpr/"
  - "relay/helper/price.go"
  - "relay/helper/billing_expr_request.go"
  - "service/billing_session.go"
  - "service/pre_consume_quota.go"
  - "service/tiered_settle.go"
  - "service/text_quota.go"
  - "setting/billing_setting/tiered_billing.go"
---

# Billing Module

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `service/billing_session.go`
- `service/pre_consume_quota.go`
- `service/quota.go`
- `service/text_quota.go`
- `service/tiered_settle.go`
- `relay/helper/price.go`
- `relay/helper/billing_expr_request.go`
- `pkg/billingexpr/`
- `pkg/billingexpr/expr.md`
- `setting/billing_setting/tiered_billing.go`

## Responsibility

Billing estimates and reserves quota before an upstream relay call, then settles against actual usage after the response. It supports wallet quota, subscriptions, token quota, trust-quota bypass, task billing, violation fees, and expression-based tiered billing.

## Billing Session Lifecycle

`BillingSession` wraps pre-consume, reserve, settle, and refund behavior. It is the main stateful guard for a single billable request.

- Pre-consume can reserve token quota and a funding source.
- Trust quota can bypass pre-consumption when configured and the user/token has enough balance.
- Settlement adjusts the funding source and token quota by the actual delta.
- Refund is idempotent and skips funding once funding has already been settled.

## Tiered Expression Billing

Read `pkg/billingexpr/expr.md` before modifying this area.

Important facts:

- One expression is the billing contract for a model.
- Coefficients are real prices in dollars per 1M tokens.
- `p` and `c` are fallback token variables. Cache, image, and audio subcategory tokens are subtracted only when the expression references their specific variables.
- `len` is the total input context length for tier conditions and is not reduced by subcategory exclusion.
- Expressions can use request-aware helpers such as `param()` and `header()`.
- Expressions carry versions such as `v1:`; no prefix currently means v1.

Pre-consume path:

1. `relay/helper/price.go` detects `tiered_expr` billing.
2. `relay/helper/billing_expr_request.go` builds request input for expression helpers.
3. `pkg/billingexpr` compiles and runs the expression against estimated tokens.
4. A frozen `BillingSnapshot` and request input are stored on `RelayInfo`.

Settlement path:

1. Usage is normalized by `service.BuildTieredTokenParams()`.
2. `service.TryTieredSettle()` reuses the frozen pre-consume snapshot.
3. On expression errors, settlement falls back to pre-consumed or estimated quota.

## Invariants

- Do not split expression behavior across hidden tables or implicit multipliers.
- Do not use `p` for long-context tier checks; use `len`.
- Keep token normalization upstream-agnostic: GPT/OpenAI usage and Claude/Anthropic usage report cache tokens differently.
- Keep request input frozen across pre-consume and settlement.

## Tests

Important tests:

- `pkg/billingexpr/billingexpr_test.go`
- `service/tiered_settle_test.go`
- `service/text_quota_test.go`
- `relay/helper/billing_expr_request_test.go`
- `relay/helper/price_test.go`
- `controller/channel_test_internal_test.go`
